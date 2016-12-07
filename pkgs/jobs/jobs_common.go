package jobs

import (
	"fmt"
	"reflect"

	"github.com/eris-ltd/eris-cli/log"
)

type JobResults struct {
	// Full Result
	JobResult string
	// Map of job name to results
	JobVars map[string]string
}

type JobsCommon interface {
	PreProcess(*Jobs) error
	Execute(*Jobs) (*JobResults, error)
}

type Job struct {
	// Name of the job
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// Sets/Resets the primary account to use
	Account *Account `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// Set an arbitrary value
	Set *Set `mapstructure:"set" json:"set" yaml:"set" toml:"set"`
	// Contract compile and send to the chain functions
	Deploy *Deploy `mapstructure:"deploy" json:"deploy" yaml:"deploy" toml:"deploy"`
	// Send tokens from one account to another
	Send *Send `mapstructure:"send" json:"send" yaml:"send" toml:"send"`
	// Utilize eris:db's native name registry to register a name
	RegisterName *RegisterName `mapstructure:"register" json:"register" yaml:"register" toml:"register"`
	// Sends a transaction which will update the permissions of an account. Must be sent from an account which
	// has root permissions on the blockchain (as set by either the genesis.json or in a subsequence transaction)
	Permission *Permission `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	// Sends a bond transaction
	Bond *Bond `mapstructure:"bond" json:"bond" yaml:"bond" toml:"bond"`
	// Sends an unbond transaction
	Unbond *Unbond `mapstructure:"unbond" json:"unbond" yaml:"unbond" toml:"unbond"`
	// Sends a rebond transaction
	Rebond *Rebond `mapstructure:"rebond" json:"rebond" yaml:"rebond" toml:"rebond"`
	// Sends a transaction to a contract. Will utilize eris-abi under the hood to perform all of the heavy lifting
	Call *Call `mapstructure:"call" json:"call" yaml:"call" toml:"call"`
	// Wrapper for mintdump dump. WIP
	DumpState *DumpState `mapstructure:"dump-state" json:"dump-state" yaml:"dump-state" toml:"dump-state"`
	// Wrapper for mintdum restore. WIP
	RestoreState *RestoreState `mapstructure:"restore-state" json:"restore-state" yaml:"restore-state" toml:"restore-state"`
	// Sends a "simulated call" to a contract. Predominantly used for accessor functions ("Getters" within contracts)
	QueryContract *QueryContract `mapstructure:"query-contract" json:"query-contract" yaml:"query-contract" toml:"query-contract"`
	// Queries information from an account.
	QueryAccount *QueryAccount `mapstructure:"query-account" json:"query-account" yaml:"query-account" toml:"query-account"`
	// Queries information about a name registered with eris:db's native name registry
	QueryName *QueryName `mapstructure:"query-name" json:"query-name" yaml:"query-name" toml:"query-name"`
	// Queries information about the validator set
	QueryVals *QueryVals `mapstructure:"query-vals" json:"query-vals" yaml:"query-vals" toml:"query-vals"`
	// Makes and assertion (useful for testing purposes)
	Assert *Assert `mapstructure:"assert" json:"assert" yaml:"assert" toml:"assert"`
	// Results of the job
	Results JobResults
}

func (job *Job) getType() (JobsCommon, error) {
	jobType := reflect.ValueOf(job)
	jobValue := jobType.Elem()
	//iterate through struct field and execute non nil and non name fields
	//break after executing so that we don't run into other fields
	for i := 1; i < jobValue.NumField(); i++ {
		field := jobValue.Field(i)
		if ptr := field.Pointer(); ptr != 0 {
			return field.Interface().(JobsCommon), nil
		}
	}
	return nil, fmt.Errorf("Could not find a job to execute.")
}

func (job *Job) announce(jobType JobsCommon) {
	log.Warn("\n*****Executing Job*****\n")
	log.WithField("=>", job.Name).Warn("Job Name")
	typ := fmt.Sprintf("%T", jobType)
	log.WithField("=>", typ).Info("Type")
}

func (job *Job) beginJob(jobType JobsCommon, jobs *Jobs) (*JobResults, error) {
	err := jobType.PreProcess(jobs)
	if err != nil {
		return &JobResults{}, err
	}
	results, err := jobType.Execute(jobs)
	if err != nil {
		return &JobResults{}, err
	}
	return results, nil
}
