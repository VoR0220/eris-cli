package jobs

import (
	"github.com/eris-ltd/eris-cli/log"
)

type Jobs struct {
	Account   string
	Jobs      []*Job `mapstructure:"jobs" json:"jobs" yaml:"jobs" toml:"jobs"`
	JobMap    map[string]JobResults
	Libraries map[string]string
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

type JobResults struct {
	// Full Result
	JobResult string
	// Map of job name to results
	JobVars map[string]string
}

func BlankJobs() *Jobs {
	return &Jobs{}
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

func dbCall(job JobsCommon, do *definitions.Do) (*definitions.JobResults, error) {
	erisNodeClient := client.NewErisNodeClient(do.ChainName)
	erisKeyClient := keys.NewErisKeyClient(do.Signer)
	switch job.(type) {
	case *Send:
		tx, err := core.Send(erisNodeClient, erisKeyClient, do.PublicKey, send.Source, send.Destination, send.Amount, send.Nonce)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
	case *BondJob:
		tx, err := core.Bond(erisNodeClient, erisKeyClient, do.PublicKey, bond.Account, bond.Amount, bond.Nonce)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
	case *Permission:
		tx, err := core.Permissions(erisNodeClient, erisKeyClient, do.PublicKey, perm.Source, perm.Nonce, perm.Action, args)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
	case *Rebond:
		tx, err := core.Rebond(rebond.Account, rebond.Height)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
	case *Name:
		tx, err := core.Name(erisNodeClient, erisKeyClient, do.PublicKey, name.Source, name.Amount, name.Nonce, name.Fee, name.Name, name.Data)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
	case *Call:
		tx, err := core.Call(erisNodeClient, erisKeyClient, do.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Fee, callData)
		if err != nil {
			return "", make([]*definitions.Variable, 0), err
		}
	case *Deploy:
		tx, err := core.Call(erisNodeClient, erisKeyClient, do.PublicKey, deploy.Source, "", deploy.Amount, deploy.Nonce, deploy.Gas, deploy.Fee, contractCode)
		if err != nil {
			return &txs.CallTx{}, fmt.Errorf("Error deploying contract %s: %v", contractName, err)
		}
	default :
		return nil, fmt.Errorf("Error, invalid job")
	}

	res, err := core.SignAndBroadcast(do.ChainID, erisNodeClient, erisKeyClient, tx.(txs.Tx), true, true, true)
	if err != nil {
		return util.MintChainErrorHandler(do, err)
	}

	if err := util.ReadTxSignAndBroadcast(res, err); err != nil {
		log.Error("ERROR =>")
		return "", err
	}

}