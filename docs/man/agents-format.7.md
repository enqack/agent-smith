---
title: AGENTS-FORMAT
section: 7
header: Agent Smith Manual
date: 2025-12-14
---

# agents-format(7)

## NAME

agents-format - specifications and concepts for Agent Smith

## DESCRIPTION

Agent Smith implements a dynamic version of the **AGENTS.md** pattern.

### The AGENTS.md Standard

Refer to **https://agents.md/**. An `AGENTS.md` file defines the context, behavior, and capabilities of an AI agent. It acts as a "system prompt" or "instruction set" for LLMs interacting with a repository.

### Dynamic Context

Standard `AGENTS.md` is static. Agent Smith makes it dynamic by introducing **Personas** and a **Canonical Target**.

#### Personas

A **Persona** is a specialized `AGENTS.md` file, named using the convention `AGENTS.<name>.md`.

* `AGENTS.coder.md`
* `AGENTS.writer.md`
* `AGENTS.architect.md`

#### Canonical Target

The **Canonical Target** is a specific symbolic link (defaults to `$XDG_CONFIG_HOME/agents/AGENTS.md`).

This link is the **Source of Truth**.
* If it points to `AGENTS.coder.md`, the system is in "Coder" mode.
* If it points to `AGENTS.writer.md`, the system is in "Writer" mode.

### Drift and Reconciliation

Because the Canonical Target is a standard filesystem symlink, you can modify it using standard tools (`ln -sf`).

**Source of Truth Pattern:**
1. **agents**(1) inspects the Canonical Target to determine the Active Persona.
2. If the Canonical Target has changed manually (Drift), **agents**(1) accepts this as the new valid state.
3. **agents**(1) (via `reconcile`) forces all other configured targets (copies, other links) to match the Active Persona.

## SEE ALSO

**agents**(1), **agents-config**(5)
