//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build compiles the binary
func Build() error {
	fmt.Println("Building...")
	return sh.RunV("go", "build", "-o", "bin/agents", ".")
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
