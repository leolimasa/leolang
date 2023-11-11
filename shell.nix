{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    rustc
    llvmPackages_13.libllvm
    llvmPackages_13.clangUseLLVM
    rustup
  ];

  # The export LD_LIBRARY_PATH is needed because jupyter-lab
  # needs to access libstdc++ which the shell doesn't export
  shellHook = ''
    export ENV_NAME=leolang
  '';
}
