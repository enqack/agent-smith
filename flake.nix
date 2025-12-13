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
        agents = pkgs.buildGoModule {
          pname = "agents";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-hicKdQauPnxag4DFS+xoYxwwXFSmPUsTm3y3Cxuw8UM=";


          postInstall = ''
            mv $out/bin/agent-smith $out/bin/agents
            mkdir -p $out/share/agent-smith/agents
            mkdir -p $out/etc/agent-smith/agents
          '';
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
            go
            mage
            gopls
            gotools
            go-tools
          ];

          shellHook = ''
            echo "Welcome to the Agent Smith dev shell!"
            echo "Run 'mage' to build the project."
          '';
        };
      }
    );
}
