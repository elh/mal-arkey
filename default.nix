# latest channel, 23.05, does not have go 1.21
{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go_1_21
    # CI dependencies
    golangci-lint
    golint
  ];
}
