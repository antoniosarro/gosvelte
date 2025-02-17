{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    nixpkgs-unstable,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem
    (
      system: let
        pkgs = import nixpkgs {
          system = system;
        };
        pkgs-unstable = import nixpkgs-unstable {
          system = system;
        };

        common = with pkgs; [
          go
          nodejs_22
          corepack_22
        ];

        unstable = with pkgs-unstable; [
          gopls
          go-tools
          golangci-lint
          air
        ];

        # runtime Deps
        libraries = with pkgs;
          [
          ]
          ++ common;

        # compile-time deps
        packages = with pkgs;
          [
          ]
          ++ common ++ unstable;
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = packages;
          buildInputs = libraries;
          shellHook = ''
            echo "🚀 Development environment loaded"
          '';
          env = {
            GOPATH = "${placeholder "out"}/go";
            CGO_ENABLED = "1";
          };
        };

        formatter = pkgs.alejandra;
      }
    );
}
