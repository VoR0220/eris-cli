package compilers

import (
	"strings"

	"github.com/eris-ltd/eris/util"

	docker "github.com/fsouza/go-dockerclient"
)

//The following represents solidity outputs
type SolcReturn struct {
	Warning   string
	Error     error
	Version   string                `json:"version"`
	Contracts map[string]*SolcItems `json:"contracts"`
}

type SolcItems struct {
	//Note: There will be more fields than this in the final version, this is just a base
	Bin string `json:"bin"`
	Abi string `json:"abi"`
}

//This is a template to define our inputs for the compiler image
type SolcTemplate struct {
	// (Optional) what to get in the output, can be any combination of [abi, bin, clone-bin, bin-runtime, userdoc, devdoc, asm]
	// abi: application binary interface. Necessary for interaction with contracts.
	// bin: binary bytecode. Necessary for creating and deploying and interacting with contracts.
	// clone-bin: Cloned contracts binary.
	// bin-runtime: Runtime binaries for contracts.
	// userdoc: natspec for users.
	// devdoc: natspec for devs.
	// asm: assembly opcodes.
	CombinedOutput []string `mapstructure:"combined-json" yaml:"combined-json"`
	// (Optional) Direct string of library address mappings.
	//  Syntax: <libraryName>:<address>
	//  Address is interpreted as a hex string optionally prefixed by 0x.
	Libraries []string `mapstructure:"libraries" yaml:"libraries"`
	// (Optional) Remappings, see https://solidity.readthedocs.io/en/latest/layout-of-source-files.html#use-in-actual-compilers
	// Syntax: <remoteName>=<localName>
	Remappings []string `mapstructure:"remappings" yaml:"remappings"`
	// (Optional) Whether or not to get a gas estimation. Default false.
	// Cannot get abi, binaries and documentations while enabled.
	GasEstimate bool `mapstructure:"gas-estimate" yaml:"gas-estimate"`
	// (Optional) if true, optimizes solidity code
	Optimize bool `mapstructure:"optimize" yaml:"optimize"`
	// (Optional) the number of optimization runs to run on solidity
	OptimizeRuns uint `mapstructure:"optimize-runs" yaml:"optimize-runs"`
	// (Optional) For anything else we may have missed
	Exec string `mapstructure:"exec" yaml:"exec"`
}

func (s *SolcTemplate) Compile(files []string) (Return, error) {
	solcExecute := []string{"solc"}

	switch {
	case s.Exec != "":
		solcExecute = append(solcExecute, s.Exec)
	case s.GasEstimate:
		solcExecute = append(solcExecute, "--gas")
	default:
		if len(s.CombinedOutput) > 0 {
			solcExecute = append(solcExecute, "--combined-json", strings.Join(s.CombinedOutput, ","))
		}
		if len(s.Libraries) > 0 {
			solcExecute = append(solcExecute, "--libraries", strings.Join(s.Libraries, ","))
		}
		if len(s.Remappings) > 0 {
			solcExecute = append(solcExecute, strings.Join(s.Remappings, " "))
		}
		if s.Optimize {
			solcExecute = append(solcExecute, "--optimize")
		}
		if s.OptimizeRuns != nil {
			solcExecute = append(solcExecute, "--optimize-runs", strconv.FormatUint(s.OptimizeRuns, 10))
		}
	}
	solcExecute = append(solcExecute, files...)
	finalCommand := strings.Join(solcExecute, " ")

	volume, err := util.DockerClient.CreateVolume(opts)
	if err != nil {
		return nil, err
	}
	util.DockerClient.CreateContainer(opts)
}
