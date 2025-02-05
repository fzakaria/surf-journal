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
        buildInputs = with pkgs; [
          gemsets.${system}.env
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
        tag = "latest";
        # This removes reproducibility
        created = "now";
        runAsRoot = ''
          # Expects this folder to exist for checking jemalloc
          mkdir -p /usr/lib
        '';
        copyToRoot = with pkgs.dockerTools; [
          (pkgs.buildEnv {
            name = "image-root";
            pathsToLink = ["/bin"];
            paths = with pkgs; [
              coreutils
              findutils
            ];
          })
          # This provides the env utility at /usr/bin/env.
          usrBinEnv
          # This provides bashInteractive at /bin/sh.
          binSh
          # This sets up /etc/ssl/certs/ca-certificates.crt.
          caCertificates
          # Provides /etc/passwd and /etc/group that contain root and nobody.
          # Useful when packaging binaries that insist on using nss to look up username/groups
          fakeNss
          # Everything seems to assume that the application is in the root lots of relative
          # paths
          default
        ];
        config = {
          Env = [
            "RAILS_ENV=production"
            # Enable jemalloc for reduced memory usage and latency.
            "LD_PRELOAD=${pkgs.jemalloc}/lib/libjemalloc.so"
          ];
          Cmd = ["./bin/thrust" "./bin/rails" "server"];
          Entrypoint = ["${default}/bin/docker-entrypoint"];
          ExposedPorts = {
            "3000/tcp" = {};
          };
        };
      };
    });

    apps = eachSystem ({system, ...}: {
      default = {
        type = "app";
        program = "${packages.${system}.default}/bin/dev";
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
