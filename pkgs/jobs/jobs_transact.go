package jobs

import (
	//"encoding/csv"
	"fmt"
	//"io"
	//"os"

	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/util"
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

func (send *Send) PreProcess(jobs *Jobs) error {
	var err error
	send.Source, err = stringPreProcess(send.Source, jobs)
	send.Destination, err = stringPreProcess(send.Destination, jobs)
	send.Amount, err = stringPreProcess(send.Amount, jobs)
	send.Source = useDefault(send.Source, jobs.Package.Account)
}

func (send *Send) Execute(jobs *Jobs) (*JobResults, error) {

	// Formulate tx
	log.WithFields(log.Fields{
		"source":      send.Source,
		"destination": send.Destination,
		"amount":      send.Amount,
	}).Info("Sending Transaction")

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Send(erisNodeClient, erisKeyClient, jobs.PublicKey, send.Source, send.Destination, send.Amount, send.Nonce)
	if err != nil {
		return util.MintChainErrorHandler(jobs, err)
	}

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
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

/*func (name *RegisterName) PreProcess(jobs *Jobs) error {
	name.DataFile, err := util.StringPreProcess(name.DataFile, do)
	if err != nil {
		return err
	}
	name.Amount, err = util.StringPreProcess(name.Amount, do)
	if err != nil {
		return err
	}
	name.Fee, err = util.StringPreProcess(name.Fee, do)
	if err != nil {
		return err
	}

	// Set Defaults
	name.Source = useDefault(name.Source, jobs.Package.Account)
	name.Fee = useDefault(name.Fee, jobs.DefaultFee)
	name.Amount = useDefault(name.Amount, jobs.DefaultAmount)
	// If a data file is given it should be in csv format and
	// it will be read first. Once the file is parsed and sent
	// to the chain then a single nameRegTx will be sent if that
	// has been populated.
	if name.DataFile != "" {
		// open the file and use a reader
		fileReader, err := os.Open(name.DataFile)
		if err != nil {
			log.Error("ERROR =>")
			return "", err
		}

		defer fileReader.Close()
		r := csv.NewReader(fileReader)

		// loop through the records
		for {
			// Read the record
			record, err := r.Read()

			// Catch the errors
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error("ERROR =>")
				return "", err
			}

			// Sink the Amount into the third slot in the record if
			// it doesn't exist
			if len(record) <= 2 {
				record = append(record, name.Amount)
			}
		}
	}
}

func RegisterNameJob(name *definitions.RegisterName, jobs *Jobs) (string, error) {



	if name.DataFile != "" {
		// open the file and use a reader
		fileReader, err := os.Open(name.DataFile)
		if err != nil {
			log.Error("ERROR =>")
			return "", err
		}

		defer fileReader.Close()
		r := csv.NewReader(fileReader)

		// loop through the records
		for {
			// Read the record
			record, err := r.Read()

			// Catch the errors
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error("ERROR =>")
				return "", err
			}

			// Sink the Amount into the third slot in the record if
			// it doesn't exist
			if len(record) <= 2 {
				record = append(record, name.Amount)
			}

			// Send an individual Tx for the record
			// [TODO]: move these to async using goroutines?
			r, err := registerNameTx(&definitions.RegisterName{
				Source: name.Source,
				Name:   record[0],
				Data:   record[1],
				Amount: record[2],
				Fee:    name.Fee,
				Nonce:  name.Nonce,
			}, do)

			if err != nil {
				log.Error("ERROR =>")
				return "", err
			}

			n := fmt.Sprintf("%s:%s", record[0], record[1])
			// TODO: fix this... simple and naive result just now.
			if err = util.WriteJobResultCSV(n, r); err != nil {
				log.Error("ERROR =>")
				return "", err
			}
		}
	}

	// If the data field is populated then there is a single
	// nameRegTx to send. So do that *now*.
	if name.Data != "" {
		return registerNameTx(name, do)
	} else {
		return "data_file_parsed", nil
	}
}

// Runs an individual nametx.
func registerNameTx(name *definitions.RegisterName, jobs *Jobs) (string, error) {
	// Process Variables
	name.Source, _ = util.PreProcess(name.Source, do)
	name.Name, _ = util.PreProcess(name.Name, do)
	name.Data, _ = util.PreProcess(name.Data, do)


	// Don't use pubKey if account override
	var oldKey string
	if name.Source != jobs.Package.Account {
		oldKey = jobs.PublicKey
		jobs.PublicKey = ""
	}

	// Formulate tx
	log.WithFields(log.Fields{
		"name":   name.Name,
		"data":   name.Data,
		"amount": name.Amount,
	}).Info("NameReg Transaction")

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Name(erisNodeClient, erisKeyClient, jobs.PublicKey, name.Source, name.Amount, name.Nonce, name.Fee, name.Name, name.Data)
	if err != nil {
		return util.MintChainErrorHandler(jobs, err)
	}

	// Don't use pubKey if account override
	if name.Source != jobs.Package.Account {
		jobs.PublicKey = oldKey
	}

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
}
*/

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
	perm.Source, err = util.StringPreProcess(perm.Source, jobs)
	if err != nil {
		return err
	}
	perm.Action, err = util.StringPreProcess(perm.Action, jobs)
	if err != nil {
		return err
	}
	perm.PermissionFlag, err = util.StringPreProcess(perm.PermissionFlag, jobs)
	if err != nil {
		return err
	}
	perm.Value, err = util.StringPreProcess(perm.Value, jobs)
	if err != nil {
		return err
	}
	perm.Target, err = util.StringPreProcess(perm.Target, jobs)
	if err != nil {
		return err
	}
	perm.Role, err = util.StringPreProcess(perm.Role, jobs)
	if err != nil {
		return err
	}
	// Set defaults
	perm.Source = useDefault(perm.Source, jobs.Account)
}

func (perm *Permission) Execute(jobs *Jobs) (*JobResults, error) {
	log.Debug("Target: ", perm.Target)
	log.Debug("Marmots Deny: ", perm.Role)
	log.Debug("Action: ", perm.Action)
	// Populate the transaction appropriately
	var args []string
	switch perm.Action {
	case "set_global":
		args = []string{perm.PermissionFlag, perm.Value}
	case "set_base":
		args = []string{perm.Target, perm.PermissionFlag, perm.Value}
	case "unset_base":
		args = []string{perm.Target, perm.PermissionFlag}
	case "add_role", "rm_role":
		args = []string{perm.Target, perm.Role}
	}

	// Don't use pubKey if account override
	var oldKey string
	if perm.Source != jobs.Account {
		oldKey = jobs.PublicKey
		jobs.PublicKey = ""
	}

	// Formulate tx
	arg := fmt.Sprintf("%s:%s", args[0], args[1])
	log.WithField(perm.Action, arg).Info("Setting Permissions")

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Permissions(erisNodeClient, erisKeyClient, jobs.PublicKey, perm.Source, perm.Nonce, perm.Action, args)
	if err != nil {
		return util.MintChainErrorHandler(jobs, err)
	}

	log.Debug("What are the args returned in transaction: ", tx.PermArgs)

	// Don't use pubKey if account override
	if perm.Source != jobs.Account {
		jobs.PublicKey = oldKey
	}

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
}

type Bond struct {
	// (Required) public key of the address which will be bonded
	PublicKey string `mapstructure:"pub_key" json:"pub_key" yaml:"pub_key" toml:"pub_key"`
	// (Required) address of the account which will be bonded
	Source string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) amount of tokens which will be bonded
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

func (bond *Bond) PreProcess(jobs *Jobs) (err error) {
	// Process Variables
	bond.Source, err = util.StringPreProcess(bond.Source, do)
	if err != nil {
		return err
	}
	bond.Amount, err = util.StringPreProcess(bond.Amount, do)
	if err != nil {
		return err
	}
	bond.PublicKey, err = util.StringPreProcess(bond.PublicKey, do)
	if err != nil {
		return err
	}
	// Use Defaults
	bond.Source = useDefault(bond.Source, jobs.Package.Account)
	jobs.PublicKey = useDefault(jobs.PublicKey, bond.PublicKey)
}

func (bond *Bond) Execute(jobs *Jobs) (*JobResults, error) {
	// Formulate tx
	log.WithFields(log.Fields{
		"public key": jobs.PublicKey,
		"amount":     bond.Amount,
	}).Infof("Bond Transaction")

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Bond(erisNodeClient, erisKeyClient, jobs.PublicKey, bond.Account, bond.Amount, bond.Nonce)
	if err != nil {
		return util.MintChainErrorHandler(jobs, err)
	}

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
}

type Unbond struct {
	// (Required) address of the account which to unbond
	Source string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the unbonding will take place (users may unbond at any
	// time >= currentBlock)
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

func (unbond *Unbond) PreProcess(jobs *Jobs) (err error) {
	unbond.Source, err = util.PreProcess(unbond.Source, do)
	if err != nil {
		return err
	}
	unbond.Height, err = util.PreProcess(unbond.Height, do)
	if err != nil {
		return err
	}
	// Use defaults
	unbond.Source = useDefault(unbond.Source, jobs.Package.Account)
}

func (unbond *Unbond) Execute(jobs *Jobs) (*JobResults, error) {
	// Formulate tx
	log.WithFields(log.Fields{
		"account": unbond.Source,
		"height":  unbond.Height,
	}).Info("Unbond Transaction")

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
}

type Rebond struct {
	// (Required) address of the account which to rebond
	Source string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the rebonding will take place (users may rebond at any
	// time >= (unbondBlock || currentBlock))
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

func (rebond *Rebond) PreProcess(jobs *Jobs) error {
	// Process Variables
	var err error
	rebond.Source, err = util.PreProcess(rebond.Source, do)
	rebond.Height, err = util.PreProcess(rebond.Height, do)
	if err != nil {
		return "", err
	}

	// Use defaults
	rebond.Source = useDefault(rebond.Source, jobs.Account)
}

func (rebond *Rebond) Execute(jobs *Jobs) (*JobResults, error) {

	// Formulate tx
	log.WithFields(log.Fields{
		"account": rebond.Source,
		"height":  rebond.Height,
	}).Info("Rebond Transaction")

	// Sign, broadcast, display
	return txFinalize(jobs, tx)
}

func txFinalize(jobs *Jobs, tx interface{}) (*JobResults, error) {
	var result *JobResults

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	res, err := core.SignAndBroadcast(jobs.ChainID, erisNodeClient, erisKeyClient, tx.(txs.Tx), true, true, true)
	if err != nil {
		return util.MintChainErrorHandler(jobs, err)
	}

	if err := util.ReadTxSignAndBroadcast(res, err); err != nil {
		log.Error("ERROR =>")
		return "", err
	}

	result.JobResult = fmt.Sprintf("%X", res.Hash)
	return result, nil
}
