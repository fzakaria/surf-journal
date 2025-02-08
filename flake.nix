{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    systems.url = "github:nix-systems/default";
    bundix = {
      url = "github:inscapist/bundix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    ruby-nix = {
      url = "github:inscapist/ruby-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    systems,
    nixpkgs,
    bundix,
    ruby-nix,
    ...
  }: let
    eachSystem = f:
      nixpkgs.lib.genAttrs (import systems) (system:
        f {
          pkgs = nixpkgs.legacyPackages.${system};
          inherit system;
        });
  in rec {
    formatter = eachSystem ({pkgs, ...}: pkgs.alejandra);

    gemsets = eachSystem (
      {pkgs, ...}: let
        ruby = pkgs.ruby_3_4;
        rubyNix = ruby-nix.lib pkgs;
        gemset =
          if builtins.pathExists ./gemset.nix
          then import ./gemset.nix
          else {};
      in (rubyNix
        {
          name = "surf-journal-gems";
          inherit gemset ruby;
        })
    );

    packages = eachSystem ({
      pkgs,
      system,
    }: let
      fs = pkgs.lib.fileset;
    in rec {
      default = pkgs.stdenv.mkDerivation {
        name = "surf-journal";
        version = "0.1.0";
        src = fs.toSource {
          root = ./.;
          fileset =
            fs.difference (fs.unions [
              ./bin
              ./app
              ./config
              ./db
              ./lib
              ./public
              ./config.ru
              ./Gemfile
              ./Gemfile.lock
              ./Rakefile
            ]) (fs.unions [
              (fs.maybeMissing ./config/master.key)
            ]);
        };
        env = {
          RAILS_ENV = "production";
        };
        buildInputs = [
          gemsets.${system}.envMinimal
          gemsets.${system}.ruby
        ];
        buildPhase = ''
          patchShebangs --build ./bin

          bundle exec bootsnap precompile --gemfile
          # Precompile bootsnap code for faster boot times
          bundle exec bootsnap precompile app/ lib/
          # Precompiling assets for production without requiring secret RAILS_MASTER_KEY
          SECRET_KEY_BASE_DUMMY=1 ./bin/rails assets:precompile
        '';
        installPhase = ''
          mkdir -p $out
          cp -r . $out
        '';
      };

      container = pkgs.dockerTools.buildImage {
        name = "registry.fly.io/surf-journal";
        runAsRoot = ''
        '';
        copyToRoot = with pkgs.dockerTools; [
          (pkgs.buildEnv {
            name = "image-root";
            pathsToLink = ["/bin"];
            paths = with pkgs; [
              coreutils
              findutils
              bashInteractive
              binSh
            ];
          })
          (pkgs.buildEnv {
            name = "journal-env";
            # The cache is pretty large so let's
            pathsToLink = ["/"];
            extraPrefix = "/surf-journal";
            paths = [
              # Everything seems to assume that the application is in the root lots of relative
              # paths
              default
            ];
          })
          # This provides the env utility at /usr/bin/env.
          usrBinEnv
          # This sets up /etc/ssl/certs/ca-certificates.crt.
          caCertificates
          # Provides /etc/passwd and /etc/group that contain root and nobody.
          # Useful when packaging binaries that insist on using nss to look up username/groups
          fakeNss
        ];
        # This ensures symlinks to directories are preserved in the image
        keepContentsDirlinks = true;
        config = {
          WorkingDir = "/surf-journal";
          Env = [
            "RAILS_ROOT=/surf-journal"
            "RAILS_ENV=production"
            # Enable jemalloc for reduced memory usage and latency.
            "LD_PRELOAD=${pkgs.jemalloc}/lib/libjemalloc.so"
          ];
          Cmd = ["./bin/thrust" "./bin/rails" "server"];
          Entrypoint = ["${default}/bin/docker-entrypoint"];
          ExposedPorts = {
            "80/tcp" = {};
          };
        };
      };
    });

    devShells = eachSystem ({
      pkgs,
      system,
    }: {
      default = with pkgs;
        mkShell {
          buildInputs = [
            bundix.packages.${system}.default
            flyctl
            skopeo
          ];
          inputsFrom = [
            packages.${system}.default
          ];
        };
    });
  };
}
