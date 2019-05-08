package main

import (
	"os"
	"os/exec"
)

// *** Docker Build ***

// Docker builds the docker image
func Docker() error {
	cmd := exec.Command("docker", "build", "-t", "anonfunc/roller", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DockerRun runs the docker image
func DockerRun() error {
	cmd := exec.Command("docker", "run",
		"--publish", "3000:3000",
		"--name", "roller",
		"--rm", "anonfunc/roller")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
