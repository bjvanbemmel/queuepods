{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs
  ];

  name = "web";

  shellHook = ''
    exec zsh
  '';
}
