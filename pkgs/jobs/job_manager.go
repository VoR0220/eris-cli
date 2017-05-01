package jobs

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/monax/cli/log"
	"github.com/monax/cli/util"

	"github.com/hyperledger/burrow/client"
	"github.com/hyperledger/burrow/keys"
)

// This is the jobset, the manager of all the job runners. It holds onto essential information for interacting
// with the chain that is passed to it from the initial Do struct from the CLI during the loading stage in
// the loaders package (See LoadJobs method). The main purpose of it is to maintain the ordering of job execution,
// and maintain awareness of the jobs that have been run and the results of those job runs.
type Jobs struct {
	// Chain and key specific variables
	Account       string            `json:"-"`
	NodeClient    client.NodeClient `json:"-"`
	KeyClient     keys.KeyClient    `json:"-"`
	ChainID       string            `json:"chain_ID"`
	PublicKey     string            `json:"-"`
	DefaultAddr   string            `json:"-"`
	DefaultAmount string            `json:"-"`
	DefaultGas    string            `json:"-"`
	DefaultFee    string            `json:"-"`
	// UI specific variables
	OutputFormat string   `json:"-"`
	DefaultSets  []string `json:"-"`
	Overwrite    bool     `json:"-"`
	//Path variables
	BinPath      string `json:"-"`
	AbiPath      string `json:"-"`
	ContractPath string `json:"-"`
	//Job variables
	timestamp  time.Time              `json:"time_jobset_started"`
	Jobs       []*Job                 `mapstructure:"jobs" yaml:"jobs" json:"-"`
	JobMap     map[string]*JobResults `json:"-"`
	jobCounter map[int]string         `json:"-"`
	//abi map
	AbiMap map[string]string `json:"-"`
}

// Returns an initialized empty jobset
func EmptyJobs() *Jobs {
	return &Jobs{}
}

// The main function out of the jobset, runs the jobs from the jobs config file in a sequential order
// checks for overwriting of the results of an old jobset if there is a repeat
func (jobs *Jobs) RunJobs() (err error) {
	var jobNames []string
	/*
		defer jobs.postProcess(err)
	*/
	jobs.jobCounter = make(map[int]string)
	for i, job := range jobs.Jobs {
		// handle duplicate job names. Request user input for permission to overwrite.
		found, overwrite, at := checkForDuplicateQueryOverwrite(job.Name, jobNames, jobs.Overwrite)
		if found && !overwrite {
			continue
		} else if found && overwrite {
			//overwrite the name
			jobs.JobMap[jobNames[at]] = &JobResults{}
			jobNames = append(jobNames[:at], jobNames[at+1:]...)
		}
		jobs.jobCounter[i] = job.Name
		jobNames = append(jobNames, job.Name)
		job.swapLegacyJob()
		results, err := job.beginJob(jobs)
		if err != nil {
			return err
		}
		jobs.JobMap[job.Name] = results
		// stored for later writing to the jobs_output file
	}

	return nil
}

// The default address to work from with future jobs. Placed at the beginning of the jobset.
func (jobs *Jobs) AddDefaultAddrJob() {
	oldJobs := jobs.Jobs

	newJob := []*Job{
		{
			Name: "defaultAddr",
			Account: &Account{
				Address: jobs.DefaultAddr,
			},
		},
	}

	jobs.Jobs = append(newJob, oldJobs...)
}

func (jobs *Jobs) AddDefaultSetJobs() {
	oldJobs := jobs.Jobs

	newJobs := []*Job{}

	for _, setr := range jobs.DefaultSets {
		blowdUp := strings.Split(setr, "=")
		if blowdUp[0] != "" {
			newJobs = append(newJobs, &Job{
				Name: blowdUp[0],
				Set: &Set{
					Value: blowdUp[1],
				},
			})
		}
	}

	jobs.Jobs = append(newJobs, oldJobs...)
}

/*func (jobs *Jobs) marshalJSON() ([]byte, error) {

}*/

// this needs to change so that it isn't within the loop of the job functions and is rather gathered on a first round loop
// whereby the duplicate names are picked up and asked about prior to execution of the loop. This might make for a weird UI
// but there will be definite performance increases.
func checkForDuplicateQueryOverwrite(name string, jobNames []string, defaultOverwrite bool) (bool, bool, int) {
	var dup bool = false
	var index int = -1
	for i, checkForDup := range jobNames {
		if checkForDup == name {
			dup = true
			index = i
			break
		}
	}
	if dup {
		if defaultOverwrite {
			log.WithField("Overwriting job name", name)
		} else {
			overwriteWarning := "You are about to overwrite a previous job name, continue?"
			if util.QueryYesOrNo(overwriteWarning, []int{}...) == util.No {
				return true, false, index
			}
			return true, true, index
		}
	}
	return dup, defaultOverwrite, index
}

// This function handles post processing whereby the results are recorded.
// Post processing should be handled by taking in an error, if the error is nil and the current
// job counter == length of the contracts, then everything went off smoothly, record the entire job
// results based on the format that was requested. Otherwise there will
// be an error returned via this method, formatted and detailed, that will be returned but not
// before recording all of the job outputs up to this point.
func (jobs *Jobs) postProcess() error {
	log.Info("Writing [jobs_output.json] to current directory")
	file, err := os.Create("jobs_output.json")
	defer file.Close()

	res, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	if _, err = file.Write(res); err != nil {
		return err
	}

	return nil
}
