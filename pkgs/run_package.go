package pkgs

import (
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/loaders"
	"github.com/eris-ltd/eris-cli/pkgs/jobs"
	"github.com/eris-ltd/eris-cli/services"
	"github.com/eris-ltd/eris-cli/util"
)

func RunPackage(do *definitions.Do) error {
	if err := getChainIPandURL(do); err != nil {
		return err
	}

	if err := util.GetChainID(do); err != nil {
		return err
	}

	var err error
	// Load the package if it doesn't exist
	if do.Package == nil {
		do.Package, err = loaders.LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
	}

	if do.LocalCompiler {
		if err := bootCompiler(); err != nil {
			return err
		}
		getLocalCompilerData(do)
	}

	return jobs.RunJobs(do)
}

func getChainIPandURL(do *definitions.Do) error {

	if !util.IsChain(do.ChainName, true) {
		return fmt.Errorf("chain (%s) is not running", do.ChainName)
	}

	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		dockerHost := os.Getenv("DOCKER_HOST")
		port := strings.LastIndex(dockerHost, ":")
		do.ChainIP = dockerHost[:port]
		chainURL, err := url.Parse(do.ChainIP)
		if err != nil {
			return err
		}
		chainURL.Scheme = "tcp"
		do.ChainURL = fmt.Sprintf("%s:%s", chainURL.String(), "46657")
	} else {
		containerName := util.ContainerName(definitions.TypeChain, do.ChainName)

		cont, err := util.DockerClient.InspectContainer(containerName)
		if err != nil {
			return util.DockerError(err)
		}

		do.ChainIP = cont.NetworkSettings.IPAddress

		// TODO flexible port
		do.ChainURL = fmt.Sprintf("tcp://%s:%s", do.ChainIP, "46657")
	}

	return nil
}

func bootCompiler() error {

	// add the compilers to the local services if the flag is pushed
	// [csk] note - when we move to default local compilers we'll remove
	// the compilers service completely and this will need to get
	// reworked to utilize DockerRun with a populated service def.
	doComp := definitions.NowDo()
	doComp.Name = "compilers"
	return services.StartService(doComp)
}

// getLocalCompilerData populates the IP:port combo for the compilers.
func getLocalCompilerData(do *definitions.Do) {
	// [csk]: note this is brittle we should only expose one port in the
	// docker file by default for the compilers service we can expose more
	// forcibly

	do.Compiler = "http://compilers:9099"
}
