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

type Set struct {
	// (Required) value which should be saved along with the jobName (which will be the key)
	// this is useful to set variables which can be used throughout the epm definition file.
	// It should be noted that arrays and bools must be defined using strings as such "[1,2,3]"
	// if they are intended to be used further in a assert job.
	Value interface{} `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}

