{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = builtins.substring 0 8 self.lastModifiedDate;
      in
      rec {
        packages = rec {
          default = pkgs.buildGo122Module rec {
            pname = "run";
            inherit version;
            src = ./.;
            subPackages = [ "cmd/run" ];
            vendorHash = "sha256-hCeOLtylQ1Jb1u43pB+2dHU3dB2CsvlOfdusmYNSxLw=";
          };
        };

        apps.default = {
          type = "app";
          program = "${packages.default}/bin/run";
        };

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            go-tools
            gotools
            gopls
          ];
        };
      });
}
