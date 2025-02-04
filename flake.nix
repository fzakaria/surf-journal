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
  in {
    formatter = eachSystem ({pkgs, ...}: pkgs.alejandra);

    devShells = eachSystem ({
      pkgs,
      system,
    }: {
      default = let
        rubyNix = ruby-nix.lib pkgs;
        gemset =
          if builtins.pathExists ./gemset.nix
          then import ./gemset.nix
          else {};
        ruby = pkgs.ruby_3_4;
        bundixcli = bundix.packages.${system}.default;
        inherit
          (rubyNix {
            name = "surf-journal-gems";
            inherit gemset ruby;
          })
          env;
      in
        pkgs.mkShell {
          buildInputs = [
            pkgs.ruby_3_4
            env
            bundixcli
          ];
        };
    });
  };
}
