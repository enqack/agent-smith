//go:build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build compiles the binary
func Build() error {
	fmt.Println("Building...")
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		return fmt.Errorf("failed to read VERSION file: %w", err)
	}
	version := strings.TrimSpace(string(versionBytes))

	ldflags := fmt.Sprintf("-X agent-smith/cmd.Version=%s", version)
	return sh.RunV("go", "build", "-ldflags", ldflags, "-o", "bin/agents", ".")
}

// Clean removes build artifacts
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("bin")
}

// Install installs the binary (example)
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return sh.RunV("go", "install", ".")
}

// Test runs tests
func Test() error {
	fmt.Println("Testing...")
	return sh.RunV("go", "test", "./...")
}

// Docs generates the man page
func Docs() error {
	fmt.Println("Generating man page...")
	err := sh.RunV("nix", "run", "nixpkgs#pandoc", "--",
		"-s", "-t", "man", "docs/man/agents.1.md", "-o", "agents.1",
	)
	if err != nil {
		return err
	}

	err = sh.RunV("nix", "run", "nixpkgs#pandoc", "--",
		"-s", "-t", "man", "docs/man/agents-config.5.md", "-o", "agents-config.5",
	)
	if err != nil {
		return err
	}

	err = sh.RunV("nix", "run", "nixpkgs#pandoc", "--",
		"-s", "-t", "man", "docs/man/agents-format.7.md", "-o", "agents-format.7",
	)
	if err != nil {
		return err
	}

	return sh.RunV("nix", "run", "nixpkgs#pandoc", "--",
		"-s", "-t", "man", "docs/man/agents-status.5.md", "-o", "agents-status.5",
	)
}

// Site generates the static documentation site
func Site() error {
	fmt.Println("Generating documentation site...")
	return sh.RunV("sphinx-build", "-b", "html", "docs", "docs/_build")
}
