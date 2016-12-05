package jobs

import (
	"fmt"
	"strings"

	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/util"
)

type Jobs struct {
	Account   string
	Jobs      []*Job `mapstructure:"jobs" json:"jobs" yaml:"jobs" toml:"jobs"`
	JobMap    map[string]JobResults
}

func EmptyJobs() *Jobs {
	return &Jobs{}
}

func (jobs *Jobs) RunJobs(do *definitions.Do) (string, error) {
	var jobNames []string
	if len(do.DefaultSets) >= 1 {
		jobs.defaultSetJobs()
	}
	if do.DefaultAddr != "" {
		jobs.defaultAddrJob(do.DefaultOutput)
	}
	for i, job := range jobs {
		found, overwrite, at := checkForDuplicateQueryOverwrite(job.Name, jobNames, do.Overwrite)
		if found && !overwrite {
			continue
		} else if found && overwrite {
			//overwrite the name
			jobNames = append(jobNames[:at], jobNames[at+1:]...)
		}
		jobNames = append(jobNames, job.Name)
		jobs.announce(job, i)
		job.Preprocess(jobs)
		job.Execute()
	}
}

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
			log.WithField("Overwriting job name", job.JobName)
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

/* [rj] commenting out for the sake of unit testing loading
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
				for _, theJob := range job.JobVars {
					log.WithField("=>", fmt.Sprintf("%s,%s", theJob.Name, theJob.Value)).Info("Job Vars")
				}
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
				for _, theJob := range job.JobVars {
					log.WithField("=>", fmt.Sprintf("%s,%s", theJob.Name, theJob.Value)).Info("Job Vars")
				}
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
}*/

func (jobs *Jobs) announce(job JobsCommon, index int) {
	log.Warn("\n*****Executing Job*****\n")
	log.WithField("=>", jobs.Jobs[index].Name).Warn("Job Name")
	typ := fmt.Sprintf("%T", job)	
	log.WithField("=>", typ).Info("Type")
}

func (jobs *Jobs) defaultAddrJob(do *definitions.Do) {
	oldJobs := jobs

	newJob := &definitions.Jobs{
		JobName: "defaultAddr",
		Job: &definitions.Job{
			Account: &definitions.Account{
				Address: do.DefaultAddr,
			},
		},
	}

	jobs = append([]*definitions.Jobs{newJob}, oldJobs...)
}

func (jobs *Jobs) defaultSetJobs() {
	oldJobs := jobs

	newJobs := []*Jobs{}

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

	jobs = append(newJobs, oldJobs...)
}

/*func (jobs *Jobs) postProcess(defaultOutput string) error {
	switch defaultOutput {
	case "csv":
		log.Info("Writing [epm.csv] to current directory")
		for _, job := range jobs {
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
}*/