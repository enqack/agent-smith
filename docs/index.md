# Agent Smith Documentation

**Agent Smith** (`agents`) is a CLI tool for managing AI context files. It allows you to define multiple "personas" (configurations of `AGENTS.md`) and instantly switch between them using symbolic links.

## Project Overview

The core philosophy of Agent Smith is simple: **The Canonical Target is the Source of Truth**.

Instead of a complex database, `agents` manages a single symbolic link (the "active persona"). All other files are secondary and can be reconciled to match this link.

*   **Canonical Target**: A symlink (default: `~/.config/agents/AGENTS.md`) that points to the currently active persona file.
*   **Personas**: Markdown files (e.g., `AGENTS.coder.md`) stored in your data directory.
*   **Reconciliation**: The process of updating all other tracked files to match the active persona.

```{toctree}
:maxdepth: 1
:caption: Guides

guides/getting-started.md
```

```{toctree}
:maxdepth: 1
:caption: Reference Manual

man/agents.1.md
man/agents-config.5.md
man/agents-status.5.md
man/agents-format.7.md
```

```{toctree}
:maxdepth: 1
:caption: Development

CONTRIBUTING.md
CHANGELOG.md
```
