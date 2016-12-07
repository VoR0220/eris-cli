package jobs

import (
	compilers "github.com/eris-ltd/eris-compilers/perform"
	"github.com/eris-ltd/eris-db/keys"
)
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
	log.WithField("=>", jobs.Account).Info("Setting Account")

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

type Compile struct {
	// (Required) the files pointed to by path to be compiled.
	// If imports are included in the file, they will be found recursively.
	Files []string `mapstructure:"files" json:"files" yaml:"files" toml:"files"`
	// (Optional) what to get in the output, can be any combination of [abi, bin, clone-bin, bin-runtime, userdoc, devdoc, asm]
    // abi: application binary interface. Necessary for interaction with contracts.
    // bin: binary bytecode. Necessary for creating and deploying and interacting with contracts.
    // clone-bin: Cloned contracts binary.
    // bin-runtime: Runtime binaries for contracts.
    // userdoc: natspec for users.
    // devdoc: natspec for devs.
    // asm: assembly opcodes.
    CombinedOutput []string  `mapstructure:"combined-output" json:"combined-json" yaml:"combined-json" toml:"combined-json"`
	// (Optional) Direct string of library address mappings.
    //  Syntax: <libraryName>:<address>
    //  Address is interpreted as a hex string optionally prefixed by 0x.
    Libraries []string `mapstructure:"libraries" json:"libraries" yaml:"libraries" toml:"libraries"`
    // (Optional) Remappings, see https://solidity.readthedocs.io/en/latest/layout-of-source-files.html#use-in-actual-compilers
    // Syntax: <remoteName>=<localName>
    Remappings []string `mapstructure:"remappings" json:"remappings" yaml:"remappings" toml:"remappings"`
	// (Optional) Whether or not to get a gas estimation. Default false.
    // Cannot get abi, binaries and documentations while enabled.
    GasEstimate bool `mapstructure:"gas-estimate" json:"gas-estimate" yaml:"gas-estimate" toml:"gas-estimate"`
	// (Optional) if true, optimizes solidity code
	Optimize bool `mapstructure:"optimize" json:"optimize" yaml:"optimize" toml:"optimize"`
	// (Optional) the number of optimization runs to run on solidity
	OptimizeRuns uint `mapstructure:"optimize-runs" json:"optimize-runs" yaml:"optimize-runs" toml:"optimize-runs"`
    // (Optional) Output file for compile job with path and directory.
    // Default no file. To get a file in the present working directory input "."
    OutputFile string `mapstructure:"output-dir" json:"output-dir" yaml:"output-dir" toml:"output-dir"`
}

func (compile *Compile) PreProcess(jobs *Jobs) (err error) {
	for i, file := range compile.Files {
		compile.Files[i] = StringPreProcess(file, jobs)
	}
	for i, output := range compile.CombinedOutput {
		compile.CombinedOutput[i] = StringPreProcess(output, jobs)
	}
	for i, library := range compile.Libraries {
		compile.Libraries[i] = SplitAndPreProcessStringPairs(library, ":", jobs)
	}
	for i, remapping := range compile.Remappings {
		compile.Remappings[i] = SplitAndPreProcessStringPairs(remapping, "=", jobs)
	}
	
}