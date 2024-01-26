let
  unstable = import (fetchTarball https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz) { };
in
{ nixpkgs ? import <nixpkgs> {} }:
with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    unstable.go_1_21
    unstable.gopls
    unstable.golangci-lint
    sqlite
  ];
}
