// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// *** Local Build ***

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	_ = os.RemoveAll("roller")
}

// Build for current platform
func Build() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "roller", ".")
	return cmd.Run()
}

func Run() error {
	mg.Deps(Build)
	return sh.Run("./roller")
}

// Build for Lambda or Linux
func BuildLinux() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "roller", ".")
	cmd.Env = append(os.Environ(), []string{
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
