package jobs

import (
	"fmt"

	"github.com/eris-ltd/eris/log"

	"github.com/eris-ltd/eris-db/client/rpc"
)

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

func (send *Send) PreProcess(jobs *Jobs) (err error) {
	send.Source, _, err = preProcessString(send.Source, jobs)
	if err != nil {
		return err
	}
	send.Destination, _, err = preProcessString(send.Destination, jobs)
	if err != nil {
		return err
	}
	send.Amount, _, err = preProcessString(send.Amount, jobs)
	if err != nil {
		return err
	}
	send.Nonce, _, err = preProcessString(send.Nonce, jobs)
	if err != nil {
		return err
	}
	send.Source = useDefault(send.Source, jobs.Account)
	send.Amount = useDefault(send.Amount, jobs.DefaultAmount)
	return nil
}

func (send *Send) Execute(jobs *Jobs) (*JobResults, error) {
	// Use Default

	// Don't use pubKey if account override
	var oldKey string
	if send.Source != jobs.Account {
		oldKey = jobs.PublicKey
		jobs.PublicKey = ""
	}

	// Formulate tx
	log.WithFields(log.Fields{
		"source":      send.Source,
		"destination": send.Destination,
		"amount":      send.Amount,
	}).Info("Sending Transaction")

	tx, err := rpc.Send(jobs.NodeClient, jobs.KeyClient, jobs.PublicKey, send.Source, send.Destination, send.Amount, send.Nonce)
	if err != nil {
		return MintChainErrorHandler(jobs, err)
	}

	// Don't use pubKey if account override
	if send.Source != jobs.Account {
		jobs.PublicKey = oldKey
	}

	// Sign, broadcast, display
	return txFinalize(tx, jobs)
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

func (perm *Permission) PreProcess(jobs *Jobs) (err error) {
	perm.Source, _, err = preProcessString(perm.Source, jobs)
	if err != nil {
		return err
	}
	perm.Action, _, err = preProcessString(perm.Action, jobs)
	if err != nil {
		return err
	}
	perm.PermissionFlag, _, err = preProcessString(perm.PermissionFlag, jobs)
	if err != nil {
		return err
	}
	perm.Value, _, err = preProcessString(perm.Value, jobs)
	if err != nil {
		return err
	}
	perm.Target, _, err = preProcessString(perm.Target, jobs)
	if err != nil {
		return err
	}
	perm.Role, _, err = preProcessString(perm.Role, jobs)
	if err != nil {
		return err
	}
	// Set defaults
	perm.Source = useDefault(perm.Source, jobs.Account)
	return nil
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

func (bond *Bond) PreProcess(jobs *Jobs) (err error) {
	// Process Variables
	bond.Account, _, err = preProcessString(bond.Account, jobs)
	if err != nil {
		return err
	}
	bond.Amount, _, err = preProcessString(bond.Amount, jobs)
	if err != nil {
		return err
	}
	bond.PublicKey, _, err = preProcessString(bond.PublicKey, jobs)
	if err != nil {
		return err
	}
	bond.Nonce, _, err = preProcessString(bond.Nonce, jobs)
	if err != nil {
		return err
	}
	// Use Defaults
	bond.Account = useDefault(bond.Account, jobs.Account)
	bond.Amount = useDefault(bond.Amount, jobs.DefaultAmount)
	jobs.PublicKey = useDefault(jobs.PublicKey, bond.PublicKey)
	return nil
}

func (bond *Bond) Execute(jobs *Jobs) (*JobResults, error) {
	return nil, fmt.Errorf("Job bond currently unimplemented.")
}

type Unbond struct {
	// (Required) address of the account which to unbond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the unbonding will take place (users may unbond at any
	// time >= currentBlock)
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

func (unbond *Unbond) PreProcess(jobs *Jobs) (err error) {
	unbond.Account, _, err = preProcessString(unbond.Account, jobs)
	if err != nil {
		return err
	}
	unbond.Height, _, err = preProcessString(unbond.Height, jobs)
	if err != nil {
		return err
	}
	// Use defaults
	unbond.Account = useDefault(unbond.Account, jobs.Account)
	return nil
}

func (unbond *Unbond) Execute(jobs *Jobs) (*JobResults, error) {
	return nil, fmt.Errorf("Job unbond currently unimplemented.")
}

type Rebond struct {
	// (Required) address of the account which to rebond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the rebonding will take place (users may rebond at any
	// time >= (unbondBlock || currentBlock))
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

func (rebond *Rebond) PreProcess(jobs *Jobs) error {
	// Process Variables
	var err error
	rebond.Account, _, err = preProcessString(rebond.Account, jobs)
	rebond.Height, _, err = preProcessString(rebond.Height, jobs)
	if err != nil {
		return err
	}

	// Use defaults
	rebond.Account = useDefault(rebond.Account, jobs.Account)
	return nil
}

func (rebond *Rebond) Execute(jobs *Jobs) (*JobResults, error) {
	return nil, fmt.Errorf("Job rebond currently unimplemented.")
}
