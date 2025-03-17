{
  version,
  buildGoApplication,
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
    modules = ./gomod2nix.toml;
  }
