package jobs

/*func getResultsFromJobTransaction(job JobsCommon, do *definitions.Do, additionalInput interface{}) (*definitions.JobResults, error) {
	var txGroup string
	var tx txs.Tx
	var err error

	//erisNodeClient := client.NewErisNodeClient(do.ChainName)
	//erisKeyClient := keys.NewErisKeyClient(do.Signer)
	oldKey := accountOverride(job, do)
	switch jobType := job.(type) {
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
	case *Call:
		callData := additionalInput.(string)
		tx, err = core.Call(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, jobType.Destination, jobType.Amount, jobType.Nonce, jobType.Gas, jobType.Fee, callData)
	case *Deploy:
		contractCode := additionalInput.(string)
		tx, err = core.Call(erisNodeClient, erisKeyClient, do.PublicKey, jobType.Source, "", jobType.Amount, jobType.Nonce, jobType.Gas, jobType.Fee, contractCode)
	default :
		return nil, fmt.Errorf("Error, invalid job")
	}
	if err != nil {
		return util.MintChainErrorHandler(do, err)
	}
	accountUnoverride(job, do, oldKey)
	res, err := core.SignAndBroadcast(do.ChainID, erisNodeClient, erisKeyClient, tx.(txs.Tx), true, true, true)
	if err != nil {
		return util.MintChainErrorHandler(do, err)
	}

	if err := util.ReadTxSignAndBroadcast(res, err); err != nil {
		log.Error("ERROR =>")
		return "", err
	}
}*/
