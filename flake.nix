{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    self,
    systems,
    nixpkgs,
    nixpkgs-unstable,
    ...
  }: let
    eachSystem = f:
      nixpkgs.lib.genAttrs (import systems) (system:
        f {
          pkgs = import nixpkgs {
            overlays = [
              (final: _prev: {
                unstable = import nixpkgs-unstable {
                  inherit system;
                };
              })
            ];
            inherit system;
          };
          inherit system;
        });
  in {
    formatter = eachSystem ({pkgs, ...}: pkgs.alejandra);

    devShells = eachSystem ({pkgs, ...}: {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          unstable.tailwindcss_4
        ];
      };
    });
  };
}
