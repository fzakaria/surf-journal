{
  version,
  buildGoApplication,
  sqlite,
  lib,
  importNpmLock,
  nodejs,
}: let
  fs = lib.fileset;
in
  buildGoApplication {
    pname = "surf-journal";
    inherit version;
    pwd = ./.;
    src = fs.toSource {
      root = ./.;
      fileset =
        fs.intersection
        (fs.gitTracked ./.)
        (fs.unions [
          ./database
          ./handlers
          ./static
          ./templates
          ./passwords
          ./embed.go
          ./cmd
          ./go.mod
          ./go.sum
        ]);
    };

    preBuild = ''
      go generate
    '';

    npmDeps = importNpmLock.buildNodeModules {
      npmRoot = ./.;
      inherit nodejs;
    };

    nativeBuildInputs = [
      nodejs
      importNpmLock.hooks.linkNodeModulesHook
    ];

    buildInputs = [
      sqlite
    ];

    modules = ./gomod2nix.toml;

    tags = ["libsqlite3"];

    meta = {
      description = "A simple site to log surf sessions.";
      homepage = "https://github.com/fzakaria/surf-journal";
      changelog = "https://github.com/lxc/incus/releases/tag/v${version}";
      mainProgram = "server";
    };
  }
