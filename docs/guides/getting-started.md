# Getting Started

## 1. Installation

### Via Nix (Recommended)

```bash
nix run github:enqack/agent-smith
```

### From Source

```bash
git clone https://github.com/enqack/agent-smith.git
cd agent-smith
go install ./cmd/agents
```

## 2. Setting Up Your Personas

By default, Agent Smith looks for personas in `~/.local/share/agent-smith/personas/`.

Let's create two simple personas.

### The Coder

```markdown
<!-- ~/.local/share/agent-smith/personas/AGENTS.coder.md -->
# Role
You are an expert Go engineer.

# Style
Prefer composition over inheritance. Write table-driven tests.
```

### The Writer

```markdown
<!-- ~/.local/share/agent-smith/personas/AGENTS.writer.md -->
# Role
You are a technical writer.

# Style
Use active voice. Use simple vocabulary.
```

## 3. Switching Contexts

Now, you can switch between them in your project.

```bash
# Set up the active persona
agents use coder
```

Verify it:

```bash
agents status
```

You should see that **coder** is the active persona.
