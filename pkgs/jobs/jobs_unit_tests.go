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

func addJobResultsToMap(jobs *Jobs) (*Jobs, error) {
	for i, job := range jobs {
		job = job.PreProcess(jobs)
		job, err = job.Execute(jobs)
		jobs.JobMap[job.Name] = job.JobResults
	}
}

var jobMap map[string]*JobResults {
	"testFalse": {
		StringResult: "false",
		Result: false,
		JobVars: nil,
	},
	"testTrue": {
		StringResult: "true",
		Result: true,
		JobVars: nil,
	},
	"testInt": {
		StringResult: "1",
		Result: 1,
		JobVars: nil,
	},
	"testString": {
		StringResult: "Hello",
		Result: "Hello",
		JobVars: nil,
	},
	"testIntArray": {
		StringResult: "[1,2,3]",
		Result: []int{1,2,3},
		JobVars: nil,
	},
	"testBoolArray": {
		StringResult: "[true,false,true]",
		Result: []bool{true,false,true},
		JobVars: nil,
	},
	"testStringArray": {
		StringResult: "[hello,world,marmot]",
		Result: []string{"hello","world","marmot"},
		JobVars: nil,
	},
	"testLibraryPairs": {
		StringResult: "libName:01234567890123456789,otherLibName:0123456789123456789",
		Result: []string{"libName:01234567890123456789","otherLibName:0123456789123456789"},
		JobVars: nil,
	},
	"testRemappingPairs": {
		StringResult: "github.com/ethereum/dapp-bin/=/usr/local/dapp-bin/,github.com/eris-ltd/monaximus/=/usr/local/marmot/",
		Result: []string{"github.com/ethereum/dapp-bin/=/usr/local/dapp-bin/","github.com/eris-ltd/monaximus/=/usr/local/marmot/"},
		JobVars: nil,
	},
	"testTupleReturn": {
		StringResults: `(1, true, "hello")`,
		Result: []abi.Variable{
				{
					Name: "0",
					Value: "1",
				},
				{
					Name: "1",
					Value: "true",
				},
				{
					Name: "2",
					Value: "hello",
				},
		},
		JobVars: {
			"0": "1",
			"1": "true",
			"2": "hello",
		},
	},
	"testDeployFunctions": {
		StringResults: "01234567890123456789",
		Result: "01234567890123456789",
		JobVars: {
			"f": "012345678901234567891234",
			"g": "012345678901234567894321",
			"h": "012345678901234567896789",
		},
	},
}
// preprocess(toProcess interface{}, mapJobResults map[string]JobResults)

func TestPreProcessingTypes(t *testing.T) {
	jobs, err := loadJobsForTesting(testPreProcessingConfig)
	if err != nil {
		t.Fatalf("could not load jobs: %v", err)
	}
	for _, test := range []struct {
		Assign string
		toTest string
		err string
	}{
	}{
		switch test.Assign {
		case "int":
			var v int
			v, err := preprocess(test.toTest, jobs)
		case "bool":
			var v bool
			v, err := preprocess(test.toTest, jobs)
		case "string":
			var v string
			v, err := preprocess(test.toTest, jobs)
		case "intSlice":
			var v []int
			v, err := preprocess(test.toTest, jobs)
		case "boolSlice":
			var v []bool
			v, err := preprocess(test.toTest, jobs)
		case "stringSlice":
			var v []string
			v, err := preprocess(test.toTest, jobs)
		case "interface":
			var v interface{}
			v, err := preprocess(test.toTest, jobs)
		case "deploy":
			
		case "call":
		default :
			t.Errorf("unsupported assignment type.")
		}
		if err != nil && len(test.err) == 0 {
			t.Errorf("%d failed. Expected no err but got: %v", i, err)
			continue
		}
		if err == nil && len(test.err) != 0 {
			t.Errorf("%d failed. Expected err: %v but got none", i, test.err)
			continue
		}
		if err != nil && len(test.err) != 0 && err.Error() != test.err {
			t.Errorf("%d failed. Expected err: '%v' got err: '%v'", i, test.err, err)
			continue
		}
	}
}