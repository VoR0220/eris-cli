package orchestration

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eris/log"
	docker "github.com/fsouza/go-dockerclient"
)

var DockerClient *docker.Client

type DockerFunc func(...interface{}) (interface{}, error)

type ContainerOrchestrator interface {
	ConfigureContainer()
	Create()
	Start()
	Stop(id string, timeout uint)
	Remove()
	Rename()
	Exec()
	ListContainers()
	Inspect(id string)
	Upload(id string)
	Download(id string)
	Wait(id string)
	Attach()
}

type ImageOrchestrator interface {
	ConfigureImage()
	Pull()
	ListImages()
}

type DockerBase struct {
	docker.PullImageOptions
	docker.ListImagesOptions
	docker.ListContainersOptions
	docker.CreateContainerOptions
	docker.LogsOptions
	docker.DownloadFromContainerOptions
	docker.UploadToContainerOptions
	docker.AttachToContainerOptions
	*docker.HostConfig
	docker.RemoveContainerOptions
	docker.RenameContainerOptions
	docker.AuthConfiguration
}

func CreateBase() (*DockerBase, error) {
	log.Debug("Creating docker base")
	base := &DockerBase{}
	if err := base.ConfigureContainer(); err != nil {
		return nil, err
	}
	if err = base.ConfigureImage(); err != nil {
		return nil, err
	}
	return base, nil
}

func (base *DockerBase) ConfigureContainer() error {
	//Default container configuration
	log.Debug("Configuring Container from Docker Base.")
	base.ListContainersOptions = docker.ListContainersOptions{All: true}
	base.HostConfig = &docker.HostConfig{
		ReadonlyRootfs: false,
		RestartPolicy:  docker.NeverRestart(),
	}
	base.CreateContainerOptions = docker.CreateContainerOptions{
		Config: &docker.Config{
			AttachStderr:    false,
			AttachStdin:     false,
			AttachStdout:    false,
			Tty:             false,
			OpenStdin:       false,
			NetworkDisabled: false,
		},
		HostConfig: base.HostConfig,
	}
	base.LogsOptions = docker.LogsOptions{
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Follow:       true,
		Since:        0,
		Timestamps:   false,
		Tail:         true,
	}
	base.AttachToContainerOptions = docker.AttachToContainerOptions{
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Logs:         false,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
		RawTerminal:  true,
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	base.DownloadFromContainerOptions = docker.DownloadFromContainerOptions{
		OutputStream: os.Stdout,
		Path:         pwd,
	}
	base.UploadToContainerOptions = docker.UploadToContainerOptions{
		Path:                 pwd,
		NoOverwriteDirNonDir: true,
	}
	return nil
}

func (base *DockerBase) ConfigureImage() {
	base.PullImageOptions = docker.PullImageOptions{
		RawJSONStream: true,
		OutputStream:  os.Stdout,
	}
}

func DockerError(err error) error {
	if _, ok := err.(*docker.Error); ok {
		return fmt.Errorf("Docker: %v", err.(*docker.Error).Message)
	}
	return err
}

func (base *DockerBase) Pull() error {
	log.Debug("Docker base pulling image.")
	if err := DockerClient.PullImage(base.PullImageOptions, base.AuthConfiguration); err != nil {
		return nil, fmt.Errorf("Error in pulling image: %v", DockerError(err))
	}
	return after, nil
}

func (base *DockerBase) ListImages() (interface{}, error) { //actual signature ([]APIImages, error)
	log.Debug("Docker base listing image.")
	if images, err := DockerClient.ListImages(base.ListImagesOptions); err != nil {
		return nil, fmt.Errorf("Error in listing images: %v", DockerError(err))
	} else {
		return images, nil
	}
}

func (base *DockerBase) ListContainers() (interface{}, error) { //actual signature ([]APIContainers, error)
	log.Debug("Docker base listing containers.")
	if containers, err := DockerClient.ListContainers(base.ListContainersOptions); err != nil {
		return nil, fmt.Errorf("Error in listing containers: %v", DockerError(err))
	} else {
		return containers, nil
	}
}

func (base *DockerBase) Create() (interface{}, error) { //(*Container, error)
	log.Debug("Docker base creating container.")
	if container, err := DockerClient.CreateContainer(base.CreateContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in creating container: %v", DockerError(err))
	} else {
		return container, nil
	}
}

func (base *DockerBase) Logs() (interface{}, error) { // error
	log.Debug("Docker base grabbing logs.")
	if err := DockerClient.Logs(base.LogsOptions); err != nil {
		return nil, fmt.Errorf("Error in getting logs: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Download(id string) (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base downloading from container.")
	if err := DockerClient.DownloadFromContainer(id, base.DownloadFromContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in downloading from container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Upload(id string) (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base uploading to container")
	if err := DockerClient.UploadToContainer(id, base.UploadToContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in uploading to container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Inspect(id string) (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base inspecting container")
	if container, err := DockerClient.InspectContainer(id); err != nil {
		return nil, fmt.Errorf("Error in inspecting container: %v", DockerError(err))
	} else {
		return container, nil
	}
}

func (base *DockerBase) Stop(id string, timeout uint) (interface{}, error) {
	log.WithFields(log.Fields{
		"id":      id,
		"timeout": timeout,
	}).Debug("Docker base stopping container")
	if err := DockerClient.StopContainer(id, timeout); err != nil {
		return nil, fmt.Errorf("Error in stopping container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Remove() (interface{}, error) {
	log.Debug("Docker base removing container")
	if err := DockerClient.RemoveContainer(base.RemoveContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in removing container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Wait(id string) (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base waiting on container")
	if exitCode, err := DockerClient.WaitContainer(id); err != nil {
		return nil, fmt.Errorf("Error in waiting on container: %v", DockerError(err))
	} else {
		return exitCode, nil
	}
}

func (base *DockerBase) Attach() (interface{}, error) {
	log.Debug("Docker base attaching to container")
	if err := DockerClient.AttachToContainer(base.AttachToContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in attaching to container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Rename() (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base waiting on container")
	if err := DockerClient.RenameContainer(base.RenameContainerOptions); err != nil {
		return nil, fmt.Errorf("Error in renaming container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}

func (base *DockerBase) Start(id string) (interface{}, error) {
	log.WithField("=>", id).Debug("Docker base starting container")
	if err := DockerClient.StartContainer(id, base.HostConfig); err != nil {
		return nil, fmt.Errorf("Error in starting container: %v", DockerError(err))
	} else {
		return nil, nil
	}
}
