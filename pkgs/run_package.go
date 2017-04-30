package pkgs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/monax/cli/definitions"
	"github.com/monax/cli/loaders"
	"github.com/monax/cli/log"
	"github.com/monax/cli/pkgs/jobs"
	"github.com/monax/cli/services"
	"github.com/monax/cli/util"
)

func RunPackage(do *definitions.Do) error {
	var gotwd string
	if do.Path == "" {
		var err error
		gotwd, err = os.Getwd()
		if err != nil {
			return err
		}
		do.Path = gotwd
	}

	// useful for debugging
	printPathPackage(do)

	// sets do.ChainIP and do.ChainPort
	if err := setChainIPandPort(do); err != nil {
		return err
	}

	do.ChainURL = fmt.Sprintf("tcp://%s:%s", do.ChainIP, do.ChainPort)
	if err := util.GetChainID(do); err != nil {
		return err
	}

	// if --dir is used and --file is left default, concat
	// so that the job will run
	// note: [zr] this could be problematic with a combo of
	// other flags, however, at least the --dir flag isn't
	// completely broken now
	if do.Path != gotwd {
		if do.YAMLPath == "epm.yaml" {
			do.YAMLPath = filepath.Join(do.Path, do.YAMLPath)
		}
		if do.BinPath == "./bin" {
			do.BinPath = filepath.Join(do.Path, "bin")
		}
		if do.ABIPath == "./abi" {
			do.ABIPath = filepath.Join(do.Path, "abi")
		}
		// TODO enable this feature
		// if do.ContractsPath == "./contracts" {
		//do.ContractsPath = filepath.Join(do.Path, "contracts")
		//}
	}

	var err error
	// Load the package if it doesn't exist
	loadedJobs, err := loaders.LoadJobs(do)
	if err != nil {
		return err
	}

	if do.Path != gotwd {
		for _, job := range loadedJobs.Jobs {
			if job.Deploy != nil {
				job.Deploy.Contract = filepath.Join(do.Path, job.Job.Deploy.Contract)
			}
		}
	}

	if len(loadedJobs.DefaultSets) >= 1 {
		loadedJobs.DefaultSetJobs()
	}
	if loadedJobs.DefaultAddr != "" {
		loadedJobs.DefaultAddrJob()
	}

	return jobs.RunJobs()
}

func setChainIPandPort(do *definitions.Do) error {

	if !util.IsChain(do.ChainName, true) {
		return fmt.Errorf("chain (%s) is not running", do.ChainName)
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		var err error
		do.ChainIP, err = util.DockerWindowsAndMacIP(do)
		if err != nil {
			return err
		}
	} else {
		containerName := util.ContainerName(definitions.TypeChain, do.ChainName)

		cont, err := util.DockerClient.InspectContainer(containerName)
		if err != nil {
			return util.DockerError(err)
		}

		do.ChainIP = cont.NetworkSettings.IPAddress
	}
	do.ChainPort = "46657" // [zr] this can be hardcoded even if [--publish] is used

	return nil
}

func printPathPackage(do *definitions.Do) {
	log.WithField("=>", do.ChainName).Info("Using Chain at")
	log.WithField("=>", do.ChainID).Debug("With ChainID")
	log.WithField("=>", do.Signer).Info("Using Signer at")
}
