package compilers

import (
	"os"

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

}
