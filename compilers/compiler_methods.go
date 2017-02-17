package compilers

import (
	"bytes"
	"fmt"
	"os"

	"github.com/eris-ltd/eris/log"
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

//if we want to we can move this into perform, just that for now I would like to have
//more fine grained control over the process and perform is unfortunately something of
//a clusterfuck to wade through atm.
func executeCompilerCommand(image string, command []string) ([]byte, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	//create container with volumes premounted
	opts := docker.CreateContainerOptions{
		Name: util.UniqueName("compiler"),
		Config: &docker.Config{
			Image:           image,
			User:            "root",
			AttachStdout:    true,
			AttachStderr:    true,
			AttachStdin:     true,
			Tty:             true,
			NetworkDisabled: false,
			WorkingDir:      "/home/",
			Cmd:             command,
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{pwd + ":" + "/home/"},
		},
	}
	if err != nil {
		return nil, util.DockerError(err)
	}
	container, err := util.DockerClient.CreateContainer(opts)
	if err != nil {
		return nil, util.DockerError(err)
	}
	removeOpts := docker.RemoveContainerOptions{
		ID:            container.ID,
		RemoveVolumes: true,
		Force:         true,
	}
	defer util.DockerClient.RemoveContainer(removeOpts)
	// Start the container.
	log.WithField("=>", opts.Name).Info("Starting data container")
	if err = util.DockerError(util.DockerClient.StartContainer(opts.Name, opts.HostConfig)); err != nil {
		return nil, err
	}

	log.WithField("=>", opts.Name).Info("Waiting for data container to exit")
	if exitCode, err := util.DockerClient.WaitContainer(container.ID); err != nil {
		if exitCode != 0 {
			err1 := fmt.Errorf("Container %s exited with status %d", container.ID, exitCode)
			if err != nil {
				err = fmt.Errorf("%s. Error: %v", err1.Error(), err)
			} else {
				err = err1
			}
		}
		return nil, err
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	logOpts := docker.LogsOptions{
		Container:    container.ID,
		OutputStream: &stdout,
		ErrorStream:  &stderr,
		RawTerminal:  true,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Since:        0,
		Timestamps:   false,
		Tail:         "all",
	}
	log.WithField("=>", opts.Name).Info("Getting logs from container")
	if err = util.DockerClient.Logs(logOpts); err != nil {
		log.Warn("Can't get logs")
		return nil, util.DockerError(err)
	}

	// Return the logs as a byte slice, if possible.
	if stdout.Len() != 0 {
		log.Warn("Hit normal output")
		return stdout.Bytes(), nil
	} else if stderr.Len() != 0 {
		log.Warn("Hit stderr")
		return stderr.Bytes(), fmt.Errorf("Compiler error.")
	} else {
		return nil, nil
	}
}
