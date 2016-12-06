package loaders

import (
	"fmt"
	"path/filepath"

	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/pkgs/jobs"

	"github.com/spf13/viper"
)

func LoadJobs(fileName string) (*jobs.Jobs, error) {
	log.Info("Loading Eris Run File...")
	var jobset = jobs.EmptyJobs()
	var epmJobs = viper.New()

	// setup file
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to find the absolute path to the eris-pm jobs file.")
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

	// load file
	if err := epmJobs.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to load the eris jobs file. Please check your path.\nERROR =>\t\t\t%v", err)
	}

	// marshall file
	if err := epmJobs.Unmarshal(jobset); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots could not figure that eris jobs file out.\nPlease check your epm.yaml is properly formatted.\n")
	}

	return jobset, nil
}
