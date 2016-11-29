package jobs

import (
	"strings"

	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/util"
)

type Jobs struct {
	Account   string
	Jobs      []*Job `mapstructure:"jobs" json:"jobs" yaml:"jobs" toml:"jobs"`
	JobMap    map[string]JobResults
	Libraries map[string]string
}

func BlankJobs() *Jobs {
	return &Jobs{}
}

type JobResults struct {
	// Full Result
	JobResult string
	// Map of job name to results
	JobVars map[string]string
}

type Job struct {
	// Name of the job
	JobName string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
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
}

func (job *Job) Announce(inter JobsCommon) {
	log.Warn("\n*****Executing Job*****\n")
	log.WithField("=>", job.JobName).Warn("Job Name")
	typ := fmt.Sprintf("%T", inter)	
	log.WithField("=>", typ).Info("Type")
}

type JobsCommon interface {
	PreProcess(*Do) error
	Execute(Do) (*JobResults, error)
}

func RunJobs(do *definitions.Do) error {
	var err error
	var dup bool = false
	// ADD DefaultAddr and DefaultSet to jobs array....
	// These work in reverse order and the addendums to the
	// the ordering from the loading process is lifo
	if len(do.DefaultSets) >= 1 {
		defaultSetJobs(do)
	}

	if do.DefaultAddr != "" {
		defaultAddrJob(do)
	}

	for index, job := range do.Package.Jobs {
		for _, checkForDup := range do.Package.Jobs[0:index] {
			if checkForDup.JobName == job.JobName {
				dup = true
				break
			}
		}
		if do.Overwrite == true && dup == true {
			log.WithField("Overwriting job name", job.JobName)
		} else if do.Overwrite == false && dup == true {
			overwriteWarning := "You are about to overwrite a previous job name, continue?"

			if util.QueryYesOrNo(overwriteWarning, []int{}...) == util.No {
				continue
			}
		}

		switch {
		// Util jobs
		case job.Job.Account != nil:
			announce(job.JobName, "Account")
			job.JobResult, err = SetAccountJob(job.Job.Account, do)
		case job.Job.Set != nil:
			announce(job.JobName, "Set")
			job.JobResult, err = SetValJob(job.Job.Set, do)

		// Transaction jobs
		case job.Job.Send != nil:
			announce(job.JobName, "Sent")
			job.JobResult, err = SendJob(job.Job.Send, do)
		case job.Job.RegisterName != nil:
			announce(job.JobName, "RegisterName")
			job.JobResult, err = RegisterNameJob(job.Job.RegisterName, do)
		case job.Job.Permission != nil:
			announce(job.JobName, "Permission")
			job.JobResult, err = PermissionJob(job.Job.Permission, do)
		case job.Job.Bond != nil:
			announce(job.JobName, "Bond")
			job.JobResult, err = BondJob(job.Job.Bond, do)
		case job.Job.Unbond != nil:
			announce(job.JobName, "Unbond")
			job.JobResult, err = UnbondJob(job.Job.Unbond, do)
		case job.Job.Rebond != nil:
			announce(job.JobName, "Rebond")
			job.JobResult, err = RebondJob(job.Job.Rebond, do)

		// Contracts jobs
		case job.Job.Deploy != nil:
			announce(job.JobName, "Deploy")
			job.JobResult, err = DeployJob(job.Job.Deploy, do)
		case job.Job.Call != nil:
			announce(job.JobName, "Call")
			job.JobResult, job.JobVars, err = CallJob(job.Job.Call, do)
			if len(job.JobVars) != 0 {
				log.WithField("=>", job.JobVars).Info("Job Vars")
			}
		// State jobs
		case job.Job.RestoreState != nil:
			announce(job.JobName, "RestoreState")
			job.JobResult, err = RestoreStateJob(job.Job.RestoreState, do)
		case job.Job.DumpState != nil:
			announce(job.JobName, "DumpState")
			job.JobResult, err = DumpStateJob(job.Job.DumpState, do)

		// Test jobs
		case job.Job.QueryAccount != nil:
			announce(job.JobName, "QueryAccount")
			job.JobResult, err = QueryAccountJob(job.Job.QueryAccount, do)
		case job.Job.QueryContract != nil:
			announce(job.JobName, "QueryContract")
			job.JobResult, job.JobVars, err = QueryContractJob(job.Job.QueryContract, do)
			if len(job.JobVars) != 0 {
				log.WithField("=>", job.JobVars).Info("Job Vars")
			}
		case job.Job.QueryName != nil:
			announce(job.JobName, "QueryName")
			job.JobResult, err = QueryNameJob(job.Job.QueryName, do)
		case job.Job.QueryVals != nil:
			announce(job.JobName, "QueryVals")
			job.JobResult, err = QueryValsJob(job.Job.QueryVals, do)
		case job.Job.Assert != nil:
			announce(job.JobName, "Assert")
			job.JobResult, err = AssertJob(job.Job.Assert, do)
		}

		if err != nil {
			return err
		}

	}

	postProcess(do)
	return nil
}

func announce(job, typ string) {
	log.Warn("\n*****Executing Job*****\n")
	log.WithField("=>", job).Warn("Job Name")
	log.WithField("=>", typ).Info("Type")
}

func defaultAddrJob(do *definitions.Do) {
	oldJobs := do.Package.Jobs

	newJob := &definitions.Jobs{
		JobName: "defaultAddr",
		Job: &definitions.Job{
			Account: &definitions.Account{
				Address: do.DefaultAddr,
			},
		},
	}

	do.Package.Jobs = append([]*definitions.Jobs{newJob}, oldJobs..)
}

func defaultSetJobs(do *definitions.Do) {
	oldJobs := do.Package.Jobs

	newJobs := []*definitions.Jobs{}

	for _, setr := range do.DefaultSets {
		blowdUp := strings.Split(setr, "=")
		if blowdUp[0] != "" {
			newJobs = append(newJobs, &definitions.Jobs{
				JobName: blowdUp[0],
				Job: &definitions.Job{
					Set: &definitions.Set{
						Value: blowdUp[1],
					},
				},
			})
		}
	}

	do.Package.Jobs = append(newJobs, oldJobs...)
}

func postProcess(do *definitions.Do) error {
	switch do.DefaultOutput {
	case "csv":
		log.Info("Writing [epm.csv] to current directory")
		for _, job := range do.Package.Jobs {
			if err := util.WriteJobResultCSV(job.JobName, job.JobResult); err != nil {
				return err
			}
		}
	case "json":
		log.Info("Writing [jobs_output.json] to current directory")
		results := make(map[string]string)
		for _, job := range do.Package.Jobs {
			results[job.JobName] = job.JobResult
		}
		return util.WriteJobResultJSON(results)
	}

	return nil
}
