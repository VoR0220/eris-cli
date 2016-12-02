package jobs

// ------------------------------------------------------------------------
// Testing Jobs
// ------------------------------------------------------------------------

// aka. Simulated Call.
type QueryContract struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the contract which should be called
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required) data which should be called. will use the eris-abi tooling under the hood to formalize the
	// transaction. QueryContract will usually be used with "accessor" functions in contracts
	Function string `mapstructure:"function" json:"function" yaml:"function" toml:"function"`
	// (Optional) data to be used in the function arguments. Will use the eris-abi tooling under the hood to formalize the
	// transaction.
	Data interface{} `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) location of the abi file to use (can be relative path or in abi path)
	// deployed contracts save ABI artifacts in the abi folder as *both* the name of the contract
	// and the address where the contract was deployed to
	ABI string `mapstructure:"abi" json:"abi" yaml:"abi" toml:"abi"`

	Variables []*Variable
}

type QueryAccount struct {
	// (Required) address of the account which should be queried
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) field which should be queried. If users are trying to query the permissions of the
	// account one can get either the `permissions.base` which will return the base permission of the
	// account, or one can get the `permissions.set` which will return the setBit of the account.
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type QueryName struct {
	// (Required) name which should be queried
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// (Required) field which should be quiried (generally will be "data" to get the registered "name")
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type QueryVals struct {
	// (Required) should be of the set ["bonded_validators" or "unbonding_validators"] and it will
	// return a comma separated listing of the addresses which fall into one of those categories
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type Assert struct {
	// (Required) key which should be used for the assertion. This is usually known as the "expected"
	// value in most testing suites
	Key string `mapstructure:"key" json:"key" yaml:"key" toml:"key"`
	// (Required) must be of the set ["eq", "ne", "ge", "gt", "le", "lt", "==", "!=", ">=", ">", "<=", "<"]
	// establishes the relation to be tested by the assertion. If a strings key:value pair is being used
	// only the equals or not-equals relations may be used as the key:value will try to be converted to
	// ints for the remainder of the relations. if strings are passed to them then eris:pm will return an
	// error
	Relation string `mapstructure:"relation" json:"relation" yaml:"relation" toml:"relation"`
	// (Required) value which should be used for the assertion. This is usually known as the "given"
	// value in most testing suites. Generally it will be a variable expansion from one of the query
	// jobs.
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}
