{
  description = "A simple site to log surf sessions.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    systems,
    nixpkgs,
    nixpkgs-unstable,
    ...
  } @ inputs: let
    eachSystem = f:
      nixpkgs.lib.genAttrs (import systems) (system:
        f {
          pkgs = import nixpkgs {
            overlays = [
              (final: _prev: {
                surf-journal = final.callPackage ./default.nix {
                  # https://github.com/nix-community/nixos-facter/blob/906098c600609d95a475449272d59b68bda2ef83/nix/packages/nixos-facter/default.nix#L18
                  # there's no good way of tying in the version to a git tag or branch
                  # so for simplicity's sake we set the version as the commit revision hash
                  # we remove the `-dirty` suffix to avoid a lot of unnecessary rebuilds in local dev
                  version = final.lib.removeSuffix "-dirty" (self.shortRev or self.dirtyShortRev);
                };
                unstable = import nixpkgs-unstable {
                  inherit system;
                };
              })
              (inputs.gomod2nix.overlays.default)
            ];
            inherit system;
          };
          inherit system;
        });
  in {
    formatter = eachSystem ({pkgs, ...}: pkgs.alejandra);

    packages = eachSystem ({pkgs, ...}: {
      default = pkgs.surf-journal;
    });

    devShells = eachSystem ({pkgs, ...}: {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [
          (mkGoEnv {pwd = ./.;})
          gomod2nix
          unstable.tailwindcss_4
        ];
      };
    });
  };
}
