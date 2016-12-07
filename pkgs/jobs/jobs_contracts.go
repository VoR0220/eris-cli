package jobs

import (
	//"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	//"github.com/eris-ltd/eris-cli/pkgs/abi"

	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/util"

	"github.com/eris-ltd/eris-db/client"
	"github.com/eris-ltd/eris-db/client/core"
	"github.com/eris-ltd/eris-db/keys"
	"github.com/eris-ltd/eris-db/txs"
)

type Deploy struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required/Optional if Compiler used) the filepath to the contract file. this should be relative to the current path **or**
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
	// (Optional/Required if Contract not specified) Choose which compiler settings to use. Must be prefixed with $ and added to jobname.
	// If instance is specified, instance takes priority over compiler files and only instance will be deployed.
	// If instance is not specified, all instances from contract and instances in files from compile will deploy.
	// If contract is not specified, files of compiler will by default be deployed.
	// Will not take gas estimates.
	Compiler string `mapstructure:"compiler" json:"compiler" yaml:"compiler" toml:"compiler"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
}

func (deploy *Deploy) PreProcess(jobs *Jobs) (err error) {
	// Preprocess variables
	deploy.Source, err = stringPreProcess(deploy.Source, jobs)
	deploy.Contract, err = stringPreProcess(deploy.Contract, jobs)
	deploy.Instance, err = stringPreProcess(deploy.Instance, jobs)
	deploy.Amount, err = stringPreProcess(deploy.Amount, jobs)
	deploy.Nonce, err = stringPreProcess(deploy.Nonce, jobs)
	deploy.Fee, err = stringPreProcess(deploy.Fee, jobs)
	deploy.Gas, err = stringPreProcess(deploy.Gas, jobs)

	deploy.Libraries, err = SplitAndPreProcessStringPairs(deploy.Libraries, ":", jobs)

	// Use defaults
	deploy.Source = useDefault(deploy.Source, jobs.Account)
	deploy.Instance = useDefault(deploy.Instance, contractName)
	deploy.Amount = useDefault(deploy.Amount, jobs.DefaultAmount)
	deploy.Fee = useDefault(deploy.Fee, jobs.DefaultFee)
	deploy.Gas = useDefault(deploy.Gas, jobs.DefaultGas)

	// Compile contracts before deploying
	if deploy.Compiler != "" {
		if !strings.HasPrefix(deploy.Compiler, "$"){
			return fmt.Errorf("Could not use compiler %v. Could not link properly to pre-run job.", deploy.Compiler)
		}
	}
	// additional data may be sent along with the contract
	// these are naively added to the end of the contract code using standard
	// mint packing

	if deploy.Data != nil {
		_, callDataArray, err = util.PreProcessInputData(compilersResponse.Objectname, deploy.Data, jobs, true)
		packedBytes, err = abi.ReadAbiFormulateCall(compilersResponse.Objectname, "", callDataArray, jobs)
		callData := hex.EncodeToString(packedBytes)
		contractCode = contractCode + callData
	}
	return err
}

func (deploy *Deploy) Execute(jobs *Jobs) (*JobResults, error) {
	// trim the extension
	contractName := strings.TrimSuffix(deploy.Contract, filepath.Ext(deploy.Contract))

	// assemble contract
	var contractPath string
	if _, err := os.Stat(deploy.Contract); err == nil {
		contractPath = deploy.Contract
	} else {
		contractPath = filepath.Join(jobs.ContractsPath, deploy.Contract)
	}
	log.WithField("=>", contractPath).Info("Contract path")

	// use the proper compiler
	if jobs.Compiler != "" {
		log.WithField("=>", jobs.Compiler).Info("Setting compiler path")
	}

	// Don't use pubKey if account override
	var oldKey string
	if deploy.Source != jobs.Package.Account {
		oldKey = jobs.PublicKey
		jobs.PublicKey = ""
	}

	// compile
	if filepath.Ext(deploy.Contract) == ".bin" {
		log.Info("Binary file detected. Using binary deploy sequence.")
		// binary deploy sequence
		contractCode, err := ioutil.ReadFile(contractPath)
		if err != nil {
			return "could not read binary file", err
		}
		tx, err := deployRaw(jobs, deploy, contractName, string(contractCode))
		if err != nil {
			return "could not deploy binary contract", err
		}
		result, err := deployFinalize(jobs, tx)
		if err != nil {
			return "", fmt.Errorf("Error finalizing contract deploy from path %s: %v", contractPath, err)
		}
		return result, err
	} else {
		// normal compilation/deploy sequence
		resp, err := compilers.RequestCompile(jobs.Compiler, contractPath, false, deploy.Libraries)

		if err != nil {
			log.Errorln("Error compiling contracts: Compilers error:")
			return "", err
		} else if resp.Error != "" {
			log.Errorln("Error compiling contracts: Language error:")
			return "", fmt.Errorf("%v", resp.Error)
		}
		// loop through objects returned from compiler
		switch {
		case len(resp.Objects) == 1:
			log.WithField("path", contractPath).Info("Deploying the only contract in file")
			response := resp.Objects[0]
			if response.Bytecode != "" {
				result, err = deployContract(deploy, jobs, response, contractPath)
				if err != nil {
					return "", err
				}
			}
		case deploy.Instance == "all":
			log.WithField("path", contractPath).Info("Deploying all contracts")
			var baseObj string
			for _, response := range resp.Objects {
				if response.Bytecode == "" {
					continue
				}
				result, err = deployContract(deploy, jobs, response, contractPath)
				if err != nil {
					return "", err
				}
				if strings.ToLower(response.Objectname) == strings.ToLower(strings.TrimSuffix(filepath.Base(deploy.Contract), filepath.Ext(filepath.Base(deploy.Contract)))) {
					baseObj = result
				}
			}
			if baseObj != "" {
				result = baseObj
			}
		default:
			log.WithField("contract", deploy.Instance).Info("Deploying a single contract")
			for _, response := range resp.Objects {
				if response.Bytecode == "" {
					continue
				}
				if strings.ToLower(response.Objectname) == strings.ToLower(deploy.Instance) {
					result, err = deployContract(deploy, jobs, response, contractPath)
					if err != nil {
						return "", err
					}
				}
			}
		}
	}

	// Don't use pubKey if account override
	if deploy.Source != jobs.Package.Account {
		jobs.PublicKey = oldKey
	}

	return result, nil
}

// TODO [rj] refactor to remove [contractPath] from functions signature => only used in a single error throw.
func deployContract(deploy Deploy, jobs *Jobs, compilersResponse compilers.ResponseItem, contractPath string) (string, error) {
/*	log.WithField("=>", string(compilersResponse.ABI)).Debug("ABI Specification (From Compilers)")
	contractCode := compilersResponse.Bytecode

	// Save ABI
	if _, err := os.Stat(jobs.ABIPath); os.IsNotExist(err) {
		if err := os.Mkdir(jobs.ABIPath, 0775); err != nil {
			return "", err
		}
	}

	// saving contract/library abi
	var abiLocation string
	if compilersResponse.Objectname != "" {
		abiLocation = filepath.Join(jobs.ABIPath, compilersResponse.Objectname)
		log.WithField("=>", abiLocation).Warn("Saving ABI")
		if err := ioutil.WriteFile(abiLocation, []byte(compilersResponse.ABI), 0664); err != nil {
			return "", err
		}
	} else {
		log.Debug("Objectname from compilers is blank. Not saving abi.")
	}

	// saving binary
	if deploy.SaveBinary {
		contractDir := filepath.Dir(deploy.Contract)
		contractName := filepath.Join(contractDir, fmt.Sprintf("%s.bin", strings.TrimSuffix(deploy.Contract, filepath.Ext(deploy.Contract))))
		log.WithField("=>", contractName).Warn("Saving Binary")
		if err := ioutil.WriteFile(contractName, []byte(contractCode), 0664); err != nil {
			return "", err
		}
	} else {
		log.Debug("Not saving binary.")
	}

	tx, err := deployRaw(jobs, deploy, compilersResponse.Objectname, contractCode)
	if err != nil {
		return "", err
	}

	// Sign, broadcast, display
	result, err := deployFinalize(jobs, tx)
	if err != nil {
		return "", fmt.Errorf("Error finalizing contract deploy %s: %v", contractPath, err)
	}

	// saving contract/library abi at abi/address
	if result != "" {
		abiLocation := filepath.Join(jobs.ABIPath, result)
		log.WithField("=>", abiLocation).Debug("Saving ABI")
		if err := ioutil.WriteFile(abiLocation, []byte(compilersResponse.ABI), 0664); err != nil {
			return "", err
		}
	} else {
		// we shouldn't reach this point because we should have an error before this.
		log.Error("The contract did not deploy. Unable to save abi to abi/contractAddress.")
	}

	return result, err
}

func deployRaw(jobs *Jobs, deploy Deploy, contractName, contractCode string) (*txs.CallTx, error) {

	// Deploy contract
	log.WithFields(log.Fields{
		"name": contractName,
	}).Warn("Deploying Contract")

	log.WithFields(log.Fields{
		"source": deploy.Source,
		"code":   contractCode,
	}).Info()

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Call(erisNodeClient, erisKeyClient, jobs.PublicKey, deploy.Source, "", deploy.Amount, deploy.Nonce, deploy.Gas, deploy.Fee, contractCode)
	if err != nil {
		return &txs.CallTx{}, fmt.Errorf("Error deploying contract %s: %v", contractName, err)
	}

	return tx, err
}*/

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
}

/*func CallJob(call Call, jobs *Jobs) (string, []Variable, error) {
	var err error
	var callData string
	var callDataArray []string
	// Preprocess variables
	call.Source, _ = util.PreProcess(call.Source, do)
	call.Destination, _ = util.PreProcess(call.Destination, do)
	//todo: find a way to call the fallback function here
	call.Function, callDataArray, err = util.PreProcessInputData(call.Function, call.Data, jobs, false)
	if err != nil {
		return "", make([]Variable, 0), err
	}
	call.Function, _ = util.PreProcess(call.Function, do)
	call.Amount, _ = util.PreProcess(call.Amount, do)
	call.Nonce, _ = util.PreProcess(call.Nonce, do)
	call.Fee, _ = util.PreProcess(call.Fee, do)
	call.Gas, _ = util.PreProcess(call.Gas, do)
	call.ABI, _ = util.PreProcess(call.ABI, do)

	// Use default
	call.Source = useDefault(call.Source, jobs.Package.Account)
	call.Amount = useDefault(call.Amount, jobs.DefaultAmount)
	call.Fee = useDefault(call.Fee, jobs.DefaultFee)
	call.Gas = useDefault(call.Gas, jobs.DefaultGas)

	// formulate call
	var packedBytes []byte
	if call.ABI == "" {
		packedBytes, err = abi.ReadAbiFormulateCall(call.Destination, call.Function, callDataArray, do)
		callData = hex.EncodeToString(packedBytes)
	} else {
		packedBytes, err = abi.ReadAbiFormulateCall(call.ABI, call.Function, callDataArray, do)
		callData = hex.EncodeToString(packedBytes)
	}
	if err != nil {
		if call.Function == "()" {
			log.Warn("Calling the fallback function")
		} else {
			var str, err = util.ABIErrorHandler(jobs, err, call, nil)
			return str, make([]Variable, 0), err
		}
	}

	// Don't use pubKey if account override
	var oldKey string
	if call.Source != jobs.Account {
		oldKey = jobs.PublicKey
		jobs.PublicKey = ""
	}

	log.WithFields(log.Fields{
		"destination": call.Destination,
		"function":    call.Function,
		"data":        callData,
	}).Info("Calling")

	erisNodeClient := client.NewErisNodeClient(jobs.ChainName)
	erisKeyClient := keys.NewErisKeyClient(jobs.Signer)
	tx, err := core.Call(erisNodeClient, erisKeyClient, jobs.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Fee, callData)
	if err != nil {
		return "", make([]Variable, 0), err
	}

	// Don't use pubKey if account override
	if call.Source != jobs.Package.Account {
		jobs.PublicKey = oldKey
	}

	// Sign, broadcast, display

	res, err := core.SignAndBroadcast(jobs.ChainID, erisNodeClient, erisKeyClient, tx, true, true, true)
	if err != nil {
		var str, err = util.MintChainErrorHandler(jobs, err)
		return str, make([]Variable, 0), err
	}

	txResult := res.Return
	var result string
	log.Debug(txResult)

	// Formally process the return
	if txResult != nil {
		log.WithField("=>", result).Debug("Decoding Raw Result")
		if call.ABI == "" {
			call.Variables, err = abi.ReadAndDecodeContractReturn(call.Destination, call.Function, txResult, do)
		} else {
			call.Variables, err = abi.ReadAndDecodeContractReturn(call.ABI, call.Function, txResult, do)
		}
		if err != nil {
			return "", make([]Variable, 0), err
		}
		log.WithField("=>", call.Variables).Debug("call variables:")
		result = util.GetReturnValue(call.Variables)
		if result != "" {
			log.WithField("=>", result).Warn("Return Value")
		} else {
			log.Debug("No return.")
		}
	} else {
		log.Debug("No return from contract.")
	}

	if call.Save == "tx" {
		log.Info("Saving tx hash instead of contract return")
		result = fmt.Sprintf("%X", res.Hash)
	}

	return result, call.Variables, nil
}

func deployFinalize(jobs *Jobs, tx interface{}) (string, error) {
	var result string

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

	result = fmt.Sprintf("%X", res.Address)
	return result, nil
}*/
