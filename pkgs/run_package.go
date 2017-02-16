package pkgs

import (
	"fmt"

	"github.com/eris-ltd/eris/definitions"
	"github.com/eris-ltd/eris/loaders"
	"github.com/eris-ltd/eris/util"
)

func RunPackage(do *definitions.Do) error {
	// sets do.ChainIP and do.ChainPort
	if err := setChainIPandPort(do); err != nil {
		return err
	}

	do.ChainURL = fmt.Sprintf("tcp://%s:%s", do.ChainIP, do.ChainPort)

	loadedJobs, err := loaders.LoadJobs(do)
	if err != nil {
		return err
	}

	return loadedJobs.RunJobs()
}

func setChainIPandPort(do *definitions.Do) error {

	if !util.IsChain(do.ChainName, true) {
		return fmt.Errorf("chain (%s) is not running", do.ChainName)
	}

	containerName := util.ContainerName(definitions.TypeChain, do.ChainName)

	cont, err := util.DockerClient.InspectContainer(containerName)
	if err != nil {
		return util.DockerError(err)
	}

	do.ChainIP = cont.NetworkSettings.IPAddress
	do.ChainPort = "46657" // [zr] this can be hardcoded even if [--publish] is used

	return nil
}
