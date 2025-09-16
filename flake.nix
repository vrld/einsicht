{
  description = "Einsicht, the mail viewer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        einsicht = pkgs.buildGoModule {
          pname = "einsicht";
          version = "0.0.0.2";
          src = ./.;
          # vendorHash = pkgs.lib.fakeHash;
          vendorHash = "sha256-T1X/od0lEN/7tfG+mNSVQ8KUZwFIMgXq5U7d0fI4nO0=";

          meta = {
            license = pkgs.lib.licenses.mit;
          };

        };

      in
      {
        packages = {
          inherit einsicht;
          default = einsicht;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            cobra-cli
            go-tools # linter (`staticcheck`)
            delve # debugger
          ];
        };
      }
    );
}
