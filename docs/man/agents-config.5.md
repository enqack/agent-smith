---
title: AGENTS-CONFIG
section: 5
header: Agent Smith Manual
date: 2025-12-14
---

# NAME

agents-config - configuration file for Agent Smith

# SYNOPSIS

**$XDG_CONFIG_HOME/agent-smith/config.yaml**

# DESCRIPTION

The `config.yaml` file controls the behavior of the **agents**(1) CLI tool, defining where personas are found and how targets are managed.

# CONFIGURATION OPTIONS

## agents_dir (list of strings)

A list of directories to search for agent persona files (e.g., `AGENTS.coder.md`).

**Default:**
* `$XDG_DATA_HOME/agent-smith/personas`
* `/usr/share/agent-smith/personas`

**Example:**
```yaml
agents_dir:
  - "/home/user/my-personas"
  - "/shared/team-personas"
```

## target_file (string)

The path to the **Canonical Target** symlink. This link acts as the Source of Truth for the active persona.

**Default:** `$XDG_CONFIG_HOME/agents/AGENTS.md`

**Example:**
```yaml
target_file: "/home/user/.config/agents/AGENTS.md"
```

## targets (list of objects)

A list of additional targets to manage automatically. These will be updated by **agents**(1) (specifically the `reconcile` and `use` commands).

Each target has:
* **path**: The file path to update.
* **mode**: `link` (symlink) or `copy` (file copy).

**Example:**
```yaml
targets:
  - path: "./docs/AGENTS.md"
    mode: "copy"
  - path: "./.github/AGENTS.md"
    mode: "link"
```

# PRECEDENCE

Configuration is resolved in the following order (highest priority first):
1. CLI Flags (`--target-file`)
2. Environment Variables (`AGENTS_TARGET_FILE`)
3. Config File (`config.yaml`)
4. Defaults

# SEE ALSO

**agents**(1), **agents-format**(7)
