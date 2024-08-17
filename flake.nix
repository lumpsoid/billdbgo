{
  description = "Dev golang env";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in 
  {

    devShells.${system}.default = pkgs.mkShell {
      name = "billdbgo-devshell";

      buildInputs = with pkgs; [
        go
      ];
      shellHook = ''
        #!/bin/sh

        session="billdbgo"
        sessioExist=$(tmux list-sessions | grep $session)

        if [ "$sessioExist" != "" ]; then
            tmux kill-session -t $session
        fi
        window=0

        tmux -L # make it able to inherit the shell variable

        tmux new-session -d -s $session
        tmux set -g mouse on
        tmux set -g mouse-select-window on

        tmux split-window -h -t $session:$window.0

        tmux attach-session -t $session
      '';
    };
  };
}
