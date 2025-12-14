{
  description = "Agent Smith installer and development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        agents = pkgs.buildGoModule rec {
          pname = "agent-smith";
          version = builtins.replaceStrings ["\n" "\r"] ["" ""] (builtins.readFile ./VERSION);
          src = ./.;
          vendorHash = "sha256-hicKdQauPnxag4DFS+xoYxwwXFSmPUsTm3y3Cxuw8UM=";

          ldflags = [ "-X agent-smith/cmd.Version=${version}" ];


          nativeBuildInputs = [ pkgs.pandoc ];

          postInstall = ''
            mv $out/bin/agent-smith $out/bin/agents
            mkdir -p $out/share/agent-smith/agents
            mkdir -p $out/etc/agent-smith/agents

            mkdir -p $out/share/man/man1
            pandoc -s -t man README.md -o $out/share/man/man1/agents.1 \
              -V title=AGENTS \
              -V section=1 \
              -V date="$(date +%Y-%m-%d)" \
              -V header="Agent Smith Manual"
          '';

          meta.mainProgram = "agents";
        };
      in
      {
        packages.default = agents;

        packages.fhs = pkgs.buildFHSEnv {
          name = "agents-fhs";
          targetPkgs = pkgs: [ agents ];
          runScript = "agents";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            gcc
            go
            mage
            gopls
            gotools
            go-tools
            pandoc
          ];

          shellHook = ''
            echo "Welcome to the Agent Smith dev shell!"
            echo "Run 'mage' to build the project."
          '';
        };
      }
    );
}
