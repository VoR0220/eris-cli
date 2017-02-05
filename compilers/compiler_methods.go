package compilers

import (
	"os"
	"os/exec"

	"github.com/eris-ltd/eris/util"

	docker "github.com/fsouza/go-dockerclient"
)

func Install(language, version string) error {
	return util.DockerClient.PullImage(
		docker.PullImageOptions{
			Repository:   language,
			Tag:          version,
			OutputStream: os.Stdout,
		},
		docker.AuthConfiguration{},
	)
}

func List(image string) ([]string, error) {
	return []string{}, nil
}

func executeCompilerCommand(image, command string) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return exec.Command("docker", "run", "-v", pwd+":/toCompile", "-w", "/toCompile", "--rm", image, command).Output()
}
