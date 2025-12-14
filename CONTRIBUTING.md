# Contributing to Agent Smith

First off, thanks for taking the time to contribute!

## Development Environment

This project uses **Nix** to provide a reproducible development environment.

### Prerequisites

- [Nix](https://nixos.org/download.html) with Flakes enabled.

### Setting up

Enter the development shell:

```bash
nix develop
```

This will provide:
- Go 1.23+
- Mage (build tool)
- Code generators and linters

## Building

To build the project:

```bash
go build ./cmd/agents
```

Or using Nix:

```bash
nix build .
```

## Testing

We have two tiers of tests:

### 1. Unit Tests

Fast, isolated tests for internal logic.

```bash
go test -v ./internal/...
```

### 2. End-to-End (E2E) Tests

Integration tests that build the binary and run it in a sandboxed environment to verify flows, configuration precedence, and error recovery.

```bash
go test -v ./tests/e2e/...
```

## Pull Requests

1.  Ensure all tests pass (including E2E).
2.  If you added features, add coverage in `tests/e2e`.
3.  Update `CHANGELOG.md` if appropriate.

## Release Process

Releases are automated via GitHub Actions and GoReleaser.
1.  Update `CHANGELOG.md` with a new version section.
2.  Update `VERSION` file.
3.  Tag the commit with `vX.Y.Z`.
4.  Push the tag.
