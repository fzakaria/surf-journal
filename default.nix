{
  version,
  buildGoApplication,
  sqlite,
  lib,
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
          ./tailwind.config.js
        ]);
    };

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
