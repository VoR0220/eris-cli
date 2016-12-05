package jobs 

import (
	"fmt"

	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/util"

	"github.com/eris-ltd/eris-db/client"
	"github.com/eris-ltd/eris-db/client/core"
	"github.com/eris-ltd/eris-db/keys"
	"github.com/eris-ltd/eris-db/txs"
)

func dbCall(job JobsCommon, do *definitions.Do, additionalInput interface{}) (*definitions.JobResults, error) {
	var txGroup string
	var tx txs.Tx
	var err error
	
	erisNodeClient := client.NewErisNodeClient(do.ChainName)
	erisKeyClient := keys.NewErisKeyClient(do.Signer)

	switch jobType := job.(type) {
	case *Send, *BondJob, *Rebond, *Unbond, *Name, *Permission:
		// Don't use pubKey if account override
		var oldKey string
		if jobType.Source != do.Package.Account {
			oldKey = do.PublicKey
			do.PublicKey = ""
		}
		//clean way to get the tx
		tx, err = getTransactJobTx(job, do, additionalInput)
		if err != nil {
			return util.MintChainErrorHandler(do, err)
		}
		// Don't use pubKey if Source override
		if job.Source != do.Package.Source {
			do.PublicKey = oldKey
		}
	case *Call:
		callData := additionalInput.(string)
		tx, err = core.Call(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, jobType.Destination, jobType.Amount, jobType.Nonce, jobType.Gas, jobType.Fee, callData)
		if err != nil {
			return "", make([]*definitions.Variable, 0), err
		}
	case *Deploy:
		contractCode := additionalInput.(string)
		tx, err = core.Call(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, "", jobType.Amount, jobType.Nonce, jobType.Gas, jobType.Fee, contractCode)
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

func getTransactJobTx(job JobsCommon, do *definitions.Do, additionalInput interface{}) (tx txs.Tx, err error) {
	switch jobType := jobType.(type) {
	case *Send:
		tx, err = core.Send(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, jobType.Destination, jobType.Amount, jobType.Nonce)
	case *BondJob:
		tx, err = core.Bond(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Account, jobType.Amount, jobType.Nonce)
	case *Permission:
		args := additionalInput.([]string)
		tx, err = core.Permissions(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, jobType.Nonce, jobType.Action, args)
	case *Rebond:
		tx, err = core.Rebond(jobType.Account, jobType.Height)
	case *Name:
		tx, err = core.Name(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, jobType.Amount, jobType.Nonce, jobType.Fee, jobType.Name, jobType.Data)
	}
	return
}