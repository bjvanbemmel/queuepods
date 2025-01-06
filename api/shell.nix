{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    air
  ];

  name = "api";

  shellHook = ''
    exec zsh
  '';
}
