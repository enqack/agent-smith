---
title: AGENTS-STATUS
section: 5
header: Agent Smith Manual
date: 2025-12-14
---

# agents-status(5)

## NAME

agents-status - internal state file for Agent Smith

## SYNOPSIS

**$XDG_STATE_HOME/agent-smith/status.yaml**

## DESCRIPTION

The `status.yaml` file persists the internal state of **agents**(1). It tracks the active persona and the list of files managed for each persona.

**WARNING**: This file is managed automatically by the `agents` CLI. Manual editing is discouraged and may lead to inconsistent state.

## FILE FORMAT

The file uses YAML format with the following structure:

### Top-Level Fields

* **canonical_target** (string):
  The absolute path to the currently active canonical symlink (Source of Truth).

* **agent_files** (list):
  A list of tracked personas and their associated targets.

### Agent File Object

Each entry in `agent_files` represents a known persona:

* **name** (string): The persona name (e.g., "coder").
* **path** (string): The absolute path to the source definition file (e.g., `.../AGENTS.coder.md`).
* **targets** (list): A list of target files managed for this persona.

### Target Object

Each entry in `targets` represents a specific file or symlink updated when this persona is active:

* **path** (string): The absolute path to the target.
* **mode** (string):
    * `link`: The target is a symbolic link to the source.
    * `copy`: The target is a copy of the source.

## EXAMPLE

```yaml
canonical_target: /home/user/.config/agents/AGENTS.md
agent_files:
  - name: coder
    path: /home/user/.local/share/agent-smith/personas/AGENTS.coder.md
    targets:
      - path: /home/user/.config/agents/AGENTS.md
        mode: link
      - path: /home/work/repo/AGENTS.md
        mode: copy
  - name: writer
    path: /home/user/.local/share/agent-smith/personas/AGENTS.writer.md
    targets:
      - path: /home/user/.config/agents/AGENTS.md
        mode: link
```

## SEE ALSO

**agents**(1), **agents-config**(5)
