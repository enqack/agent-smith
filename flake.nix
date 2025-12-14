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

          ldflags = [ "-X agent-smith/internal/cli.Version=${version}" ];


          nativeBuildInputs = [ pkgs.pandoc ];

          postInstall = ''
            # Binary is already named 'agents' because it comes from cmd/agents
            mkdir -p $out/share/agent-smith/agents
            mkdir -p $out/etc/agent-smith/agents

            mkdir -p $out/share/man/man1
            mkdir -p $out/share/man/man1
            pandoc -s -t man docs/man/agents.1.md -o $out/share/man/man1/agents.1

            mkdir -p $out/share/man/man5
            pandoc -s -t man docs/man/agents-config.5.md -o $out/share/man/man5/agents-config.5

            mkdir -p $out/share/man/man7
            pandoc -s -t man docs/man/agents-format.7.md -o $out/share/man/man7/agents-format.7

            pandoc -s -t man docs/man/agents-status.5.md -o $out/share/man/man5/agents-status.5
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
            python3Packages.sphinx
            python3Packages.myst-parser
            python3Packages.furo
            python3Packages.sphinx-copybutton
          ];

          shellHook = ''
            echo "Welcome to the Agent Smith dev shell!"
            echo "Run 'mage' to build the project."
          '';
        };
      }
    );
}
