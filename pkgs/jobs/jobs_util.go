package jobs

// ------------------------------------------------------------------------
// Util Jobs
// ------------------------------------------------------------------------

type Account struct {
	// (Required) address of the account which should be used as the default (if source) is
	// not given for future transactions. Will make sure the eris-keys has the public key
	// for the account. Generally account should be the first job called unless it is used
	// via a flag or environment variables to establish what default to use.
	Address string `mapstructure:"address" json:"address" yaml:"address" toml:"address"`
}

func (acc *Account) PreProcess(jobs *Jobs) (err error) {
	acc.Address, err = stringPreProcess(acc.Address, jobs)
}

func (acc *Account) Execute(jobs *Jobs) (*JobResults, error) {
	var result *JobResults
	var err error

	// Set the Account in the Package & Announce
	jobs.Account = acc.Address
	log.WithField("=>", jobs.Package.Account).Info("Setting Account")

	// Set the public key from eris-keys
	keys.DaemonAddr = jobs.Signer
	log.WithField("from", keys.DaemonAddr).Info("Getting Public Key")
	jobs.PublicKey, err = keys.Call("pub", map[string]string{"addr": jobs.Account, "name": ""})
	if _, ok := err.(keys.ErrConnectionRefused); ok {
		keys.ExitConnectErr(err)
	}

	if err != nil {
		return util.KeysErrorHandler(do, err)
	}

	// Set result and return
	result.JobResult = account.Address
	return result, nil
}

type Set struct {
	// (Required) value which should be saved along with the jobName (which will be the key)
	// this is useful to set variables which can be used throughout the epm definition file.
	// It should be noted that arrays and bools must be defined using strings as such "[1,2,3]"
	// if they are intended to be used further in a assert job.
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}

func (set *Set) PreProcess(jobs *Jobs) (err error) {
	set.Value, err = util.StringPreProcess(acc.Address, jobs)
}

func (set *Set) Execute(jobs *Jobs) (*JobResults, error) {
	var result *JobResults
	log.WithField("=>", set.Value).Info("Setting Variable")
	result.JobResult = set.Value
	return result, nil
}
