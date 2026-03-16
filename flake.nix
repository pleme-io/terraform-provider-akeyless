{
  description = "Akeyless Terraform Provider - manage Akeyless resources via Terraform";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    substrate = {
      url = "github:pleme-io/substrate";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, substrate, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs { inherit system; };
      mkGoTool = (import "${substrate}/lib/go-tool.nix").mkGoTool;
    in {
      packages.default = mkGoTool pkgs {
        pname = "terraform-provider-akeyless";
        version = "0.0.0-dev";
        src = self;
        proxyVendor = true;
        vendorHash = "sha256-ItO75dB0d4zH09IxNyoG1gJ0NmLge7ml/hFRT0hXsJE=";
        description = "Akeyless Terraform Provider - manage Akeyless resources via Terraform";
        homepage = "https://github.com/pleme-io/terraform-provider-akeyless";
        license = pkgs.lib.licenses.mpl20;
      };

      devShells.default = pkgs.mkShellNoCC {
        packages = with pkgs; [ go gopls gotools ];
      };
    });
}
