package jobs

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/eris-cli/log"

	"github.com/spf13/viper"
)

var testPreProcessingConfig string = `
jobs:

- name: setStorageBase
  set:
    val: 5

- name: deployStorageK
  deploy:
    contract: storage.sol
    libraries: [$someName:$someAddress]

- name: setStorage
  call:
    destination: $deployStorageK
    function: set
    data: [$setStorage]

- name: queryStorage
  query-contract:
    destination: $deployStorageK
    function: get

- name: assertStorage
  assert:
    key: $queryStorage
    relation: eq
    val: $setStorageBase
`

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
	exitCode := m.Run()
	os.Exit(exitCode)
}

//copy of loaders for ease of testing
func loadJobsForTesting(config string) (*Jobs, error){
	log.Info("Loading Eris Run File...")
	var jobset = EmptyJobs()
	var epmJobs = viper.New()
	
	err := ioutil.WriteFile("epm.yaml", []byte(config), 0644)
	defer os.Remove("epm.yaml")
	if err != nil {
		return nil, fmt.Errorf("cannot write config file: %v", err)
	}
	
	epmJobs.AddConfigPath(".")
	epmJobs.SetConfigName("epm")

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

// Test duplicates

func addJobResultsToMap(jobs *Jobs) (*Jobs, error) {
	for i, job := range jobs {
		job = job.PreProcess(jobs)
		job, err = job.Execute(jobs)
		jobs.JobMap[job.Name] = job.JobResults
	}
}

func TestPreProcessingStrings(t *testing.T) {
	jobs, err := loadJobsForTesting(testPreProcessingConfig)
	if err != nil {
		t.Fatalf("could not load jobs: %v", err)
	}
}