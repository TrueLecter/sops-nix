{
  description = "Integrates sops into nixos";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  nixConfig.extra-substituters = ["https://cache.garnix.io"];
  nixConfig.extra-trusted-public-keys = ["cache.garnix.io:CTFPyKSLcx5RMJKfLo5EEPUObbA78b0YQ2DTCJXqr9g="];
  outputs = {
    self,
    nixpkgs,
  }: let
    systems = [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-darwin"
      "aarch64-linux"
      "armv6l-linux"
      "armv7l-linux"
    ];
    forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
  in {
    overlay = final: prev: let
      localPkgs = import ./default.nix {pkgs = final;};
    in {
      inherit (localPkgs) sops-install-secrets sops-init-gpg-key sops-pgp-hook sops-import-keys-hook sops-ssh-to-age;
      # backward compatibility
      inherit (prev) ssh-to-pgp;
    };
    nixosModules.sops = import ./modules/sops;
    nixosModule = self.nixosModules.sops;
    packages = forAllSystems (system:
      import ./default.nix {
        pkgs = import nixpkgs {inherit system;};
      });
    checks =
      nixpkgs.lib.genAttrs ["x86_64-linux" "aarch64-linux"]
      (system: self.packages.${system}.sops-install-secrets.tests);
    defaultPackage = forAllSystems (system: self.packages.${system}.sops-init-gpg-key);
    devShell = forAllSystems (
      system:
        nixpkgs.legacyPackages.${system}.callPackage ./shell.nix {}
    );
    devShells = forAllSystems (system: {
      unit-tests = nixpkgs.legacyPackages.${system}.callPackage ./unit-tests.nix {};
    });
  };
}
