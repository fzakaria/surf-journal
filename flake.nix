{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    self,
    systems,
    nixpkgs,
    ...
  }: let
    eachSystem = f:
      nixpkgs.lib.genAttrs (import systems) (system:
        f {
          pkgs = nixpkgs.legacyPackages.${system};
          inherit system;
        });
  in {
    formatter = eachSystem ({pkgs, ...}: pkgs.alejandra);

    devShells = eachSystem ({pkgs, ...}: {
      default = let
        ruby = pkgs.ruby_3_4;
        gems = pkgs.bundlerEnv {
          name = "surf-journal-gems";
          inherit ruby;
          gemdir = ./.;
        };
      in
        pkgs.mkShell {
          buildInputs = with pkgs; [
            ruby_3_4
            gems
            bundix
          ];
        };
    });
  };
}
