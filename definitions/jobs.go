package definitions


// ------------------------------------------------------------------------
// Contracts Jobs
// ------------------------------------------------------------------------

type PackageDeploy struct {
	// TODO
}

type Variable struct {
	Name  string
	Value string
}

type Deploy struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) the filepath to the contract file. this should be relative to the current path **or**
	// relative to the contracts path established via the --contracts-path flag or the $EPM_CONTRACTS_PATH
	// environment variable. If contract has a "bin" file extension then it will not be sent to the
	// compilers but rather will just be sent to the chain. Note, if you use a "call" job after deploying
	// a binary contract then you will be **required** to utilize an abi field in the call job.
	Contract string `mapstructure:"contract" json:"contract" yaml:"contract" toml:"contract"`
	// (Optional) the name of contract to instantiate (it has to be one of the contracts present)
	// in the file defined in Contract above.
	// When none is provided, the system will choose the contract with the same name as that file.
	// use "all" to override and deploy all contracts in order. if "all" is selected the result
	// of the job will default to the address of the contract which was deployed that matches
	// the name of the file (or the last one deployed if there are no matching names; not the "last"
	// one deployed" strategy is non-deterministic and should not be used).
	Instance string `mapstructure:"instance" json:"instance" yaml:"instance" toml:"instance"`
	// (Optional) list of Name:Address separated by commas of libraries (see solc --help)
	Libraries string `mapstructure:"libraries" json:"libraries" yaml:"libraries" toml:"libraries"`
	// (Optional) TODO: additional arguments to send along with the contract code
	Data interface{} `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) amount of tokens to send to the contract which will (after deployment) reside in the
	// contract's account
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional) amount of gas which should be sent along with the contract deployment transaction
	Gas string `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	// (Optional) after compiling the contract save the binary in filename.bin in same directory
	// where the *.sol or *.se file is located. This will speed up subsequent installs
	SaveBinary bool `mapstructure:"save" json:"save" yaml:"save" toml:"save"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

type Call struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the contract which should be called
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required unless testing fallback function) function inside the contract to be called
	Function string `mapstructure:"function" json:"function" yaml:"function" toml:"function"`
	// (Optional) data which should be called. will use the eris-abi tooling under the hood to formalize the
	// transaction
	Data interface{} `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) amount of tokens to send to the contract
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional) amount of gas which should be sent along with the call transaction
	Gas string `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) location of the abi file to use (can be relative path or in abi path)
	// deployed contracts save ABI artifacts in the abi folder as *both* the name of the contract
	// and the address where the contract was deployed to
	ABI string `mapstructure:"abi" json:"abi" yaml:"abi" toml:"abi"`
	// (Optional) by default the call job will "store" the return from the contract as the
	// result of the job. If you would like to store the transaction hash instead of the
	// return from the call job as the result of the call job then select "tx" on the save
	// variable. Anything other than "tx" in this field will use the default.
	Save string `mapstructure:"save" json:"save" yaml:"save" toml:"save"`
	// (Optional) the call job's returned variables
	Variables []*Variable
}
