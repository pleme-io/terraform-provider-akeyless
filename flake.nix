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

  outputs = inputs: (import "${inputs.substrate}/lib/repo-flake.nix" {
    inherit (inputs) nixpkgs flake-utils;
  }) {
    self = inputs.self;
    language = "go";
    builder = "tool";
    pname = "terraform-provider-akeyless";
    vendorHash = "sha256-ItO75dB0d4zH09IxNyoG1gJ0NmLge7ml/hFRT0hXsJE=";
    proxyVendor = true;
    description = "Akeyless Terraform Provider - manage Akeyless resources via Terraform";
    homepage = "https://github.com/pleme-io/terraform-provider-akeyless";
  };
}
