package loaders

import (
	"fmt"
	"path/filepath"

	"github.com/eris-ltd/eris/definitions"
	"github.com/eris-ltd/eris/log"
	"github.com/eris-ltd/eris/pkgs/jobs"

	"github.com/eris-ltd/eris-db/client"
	"github.com/eris-ltd/eris-db/keys"
	"github.com/eris-ltd/eris-db/logging/loggers"
	"github.com/spf13/viper"
)

func LoadJobs(do *definitions.Do) (*jobs.Jobs, error) {
	log.Info("Loading Eris Run File...")
	var fileName = do.YAMLPath
	var jobset = jobs.EmptyJobs()

	erisClient := client.NewErisNodeClient(do.ChainURL, loggers.NewNoopInfoTraceLogger())
	_, chainID, _, err := erisClient.ChainId()

	jobset.NodeClient = erisClient
	jobset.KeyClient = keys.NewErisKeyClient(do.Signer, loggers.NewNoopInfoTraceLogger())
	jobset.PublicKey = do.PublicKey
	jobset.DefaultAddr = do.DefaultAddr
	jobset.DefaultOutput = do.DefaultOutput
	jobset.DefaultSets = do.DefaultSets
	jobset.Overwrite = do.Overwrite
	jobset.DefaultAmount = do.DefaultAmount
	jobset.DefaultFee = do.DefaultFee
	jobset.DefaultGas = do.DefaultGas
	jobset.JobMap = make(map[string]*jobs.JobResults)
	jobset.ChainID = chainID

	if err != nil {
		return nil, err
	}
	var epmJobs = viper.New()

	// setup file
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to find the absolute path to the eris jobs file.")
	}

	path := filepath.Dir(abs)
	file := filepath.Base(abs)
	extName := filepath.Ext(file)
	bName := file[:len(file)-len(extName)]
	log.WithFields(log.Fields{
		"path": path,
		"name": bName,
	}).Debug("Loading eris-pm file")

	epmJobs.AddConfigPath(path)
	epmJobs.SetConfigName(bName)

	epmJobs.AddConfigPath(path)
	epmJobs.SetConfigName(bName)

	// load file
	if err := epmJobs.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to load the eris jobs file. Please check your path: %v", err)
	}

	// marshall file
	if err := epmJobs.Unmarshal(jobset); err != nil {
		return nil, fmt.Errorf(`Sorry, the marmots could not figure that eris jobs file out. 
			Please check that your epm.yaml is properly formatted: %v`, err)
	}

	return jobset, nil
}
