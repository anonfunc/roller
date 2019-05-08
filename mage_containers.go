// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/sh"
)

// *** Docker Build ***

// Docker builds the docker image
func BuildDocker() error {
	cmd := exec.Command("docker", "build", "-t", "anonfunc/roller", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DockerRun runs the docker image
func Docker() error {
	cmd := exec.Command("docker", "run",
		"--publish", "3000:3000",
		"--name", "roller",
		"--rm", "anonfunc/roller")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// *** K8s ***

func BuildMinikube() error {
	mkEnv, err := sh.Output("minikube", "docker-env")
	if err != nil {
		return err
	}
	lines := strings.Split(mkEnv, "\n")
	var env []string
	for _, l := range lines {
		if strings.HasPrefix(l, "#") {
			continue
		}
		l = strings.TrimPrefix(l, "export ")
		l = strings.Replace(l, "\"", "", -1)
		env = append(env, l)
	}
	cmd := exec.Command("docker", "build", "-t", "anonfunc/roller", ".")
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Minikube() error {
	cmd := exec.Command("kubectl", "create", "-f", "k8s/deployment.yml", "--context", "minikube")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	url, err := sh.Output("minikube", "service", "roller", "--url")
	if err != nil {
		return err
	}
	fmt.Printf("%s/roll/2d6\n", url)
	return nil
}
