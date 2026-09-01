{
  description = "A personal Go web application and blog";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "go-weebsite";
          version = "0.1.0";
          src = ./.;

          subPackages = [ "cmd" ];

          # vendorHash for go modules
          vendorHash = "sha256-Pn6INucE1/bK8Zu95PcSbLI2wl/+ls8wSkAbPp7VdEo=";

          postInstall = ''
            mv $out/bin/cmd $out/bin/go-weebsite
          '';

          meta = with pkgs.lib; {
            description = "Go weebsite server";
            homepage = "https://github.com/arjunomray/go-weebsite";
            license = licenses.mit;
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            air
            golangci-lint
            gnumake
            wrangler
          ];

          shellHook = ''
            echo "🚀 Go Weebsite development environment loaded!"
            echo "Go version: $(go version)"
          '';
        };

        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };
      });
}
