{
  description = "A Nix-flake-based Go development environment with Templ support";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    templ.url = "github:a-h/templ/v0.3.1001";
  };

  outputs =
    {
      self,
      nixpkgs,
      templ,
      ...
    }:
    let
      goVersion = 25;

      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forEachSupportedSystem =
        f:
        nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            # Directly access the templ package for this specific system
            templPkg = templ.packages.${system}.default;
            pkgs = import nixpkgs {
              inherit system;
              overlays = [ self.overlays.default ];
            };
          }
        );
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem (
        { pkgs, templPkg }:
        {
          default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              gcc
              go
              gopls
              templPkg # Use the verified package here
              gotools
              golangci-lint
              tailwindcss_4
              tailwindcss-language-server
              air
            ];
          };
        }
      );

      packages = forEachSupportedSystem (
        { pkgs, templPkg }:
        {
          default = pkgs.buildGoModule {
            pname = "my-go-templ-app";
            version = "0.1.0";
            src = ./.;
            vendorHash = null; # Set this to nixpkgs.lib.fakeHash and run 'nix build' to find the real hash

            preBuild = ''
              ${templPkg}/bin/templ generate
            '';
          };
        }
      );
    };
}
