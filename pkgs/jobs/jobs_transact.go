package jobs

// ------------------------------------------------------------------------
// Transaction Jobs
// ------------------------------------------------------------------------

type Send struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the account to send the tokens
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required) amount of tokens to send from the `source` to the `destination`
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

type RegisterName struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) name which will be registered
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// (Optional, if data_file is used; otherwise required) data which will be stored at the `name` key
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) csv file in the form (name,data[,amount]) which can be used to bulk register names
	DataFile string `mapstructure:"data_file" json:"data_file" yaml:"data_file" toml:"data_file"`
	// (Optional) amount of blocks which the name entry will be reserved for the registering user
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

type Permission struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) actions must be in the set ["set_base", "unset_base", "set_global", "add_role" "rm_role"]
	Action string `mapstructure:"action" json:"action" yaml:"action" toml:"action"`
	// (Required, unless add_role or rm_role action selected) the name of the permission flag which is to
	// be updated
	PermissionFlag string `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	// (Required) the value of the permission or role which is to be updated
	Value string `mapstructure:"value" json:"value" yaml:"value" toml:"value"`
	// (Required) the target account which is to be updated
	Target string `mapstructure:"target" json:"target" yaml:"target" toml:"target"`
	// (Required, if add_role or rm_role action selected) the role which should be given to the account
	Role string `mapstructure:"role" json:"role" yaml:"role" toml:"role"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

type Bond struct {
	// (Required) public key of the address which will be bonded
	PublicKey string `mapstructure:"pub_key" json:"pub_key" yaml:"pub_key" toml:"pub_key"`
	// (Required) address of the account which will be bonded
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) amount of tokens which will be bonded
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

type Unbond struct {
	// (Required) address of the account which to unbond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the unbonding will take place (users may unbond at any
	// time >= currentBlock)
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

type Rebond struct {
	// (Required) address of the account which to rebond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the rebonding will take place (users may rebond at any
	// time >= (unbondBlock || currentBlock))
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}
