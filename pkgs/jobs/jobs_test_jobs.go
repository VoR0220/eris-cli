package jobs

import (
//"encoding/hex"
//"fmt"
//"strconv"

/*"github.com/eris-ltd/eris-cli/log"
"github.com/eris-ltd/eris-cli/pkgs/abi"
"github.com/eris-ltd/eris-cli/util"

"github.com/eris-ltd/eris-db/client"*/
)

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

/*func QueryContractJob(query *definitions.QueryContract, do *definitions.Do) (string, []*definitions.Variable, error) {
	// Preprocess variables. We don't preprocess data as it is processed by ReadAbiFormulateCall
	query.Source, _ = util.PreProcess(query.Source, do)
	query.Destination, _ = util.PreProcess(query.Destination, do)
	query.ABI, _ = util.PreProcess(query.ABI, do)

	var queryDataArray []string
	var err error
	query.Function, queryDataArray, err = util.PreProcessInputData(query.Function, query.Data, do, false)
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}
	// Set the from and the to
	fromAddrBytes, err := hex.DecodeString(query.Source)
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}
	toAddrBytes, err := hex.DecodeString(query.Destination)
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}

	// Get the packed data from the ABI functions
	var data string
	var packedBytes []byte
	if query.ABI == "" {
		packedBytes, err = abi.ReadAbiFormulateCall(query.Destination, query.Function, queryDataArray, do)
		data = hex.EncodeToString(packedBytes)
	} else {
		packedBytes, err = abi.ReadAbiFormulateCall(query.ABI, query.Function, queryDataArray, do)
		data = hex.EncodeToString(packedBytes)
	}
	if err != nil {
		var str, err = util.ABIErrorHandler(do, err, nil, query)
		return str, make([]*definitions.Variable, 0), err
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}

	// Call the client
	nodeClient := client.NewErisNodeClient(do.ChainName)
	result, _, err := nodeClient.QueryContract(fromAddrBytes, toAddrBytes, dataBytes)
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}

	// Formally process the return
	log.WithField("res", result).Debug("Decoding Raw Result")
	if query.ABI == "" {
		log.WithField("abi", query.Destination).Debug()
		query.Variables, err = abi.ReadAndDecodeContractReturn(query.Destination, query.Function, result, do)
	} else {
		log.WithField("abi", query.ABI).Debug()
		query.Variables, err = abi.ReadAndDecodeContractReturn(query.ABI, query.Function, result, do)
	}
	if err != nil {
		return "", make([]*definitions.Variable, 0), err
	}

	result2 := util.GetReturnValue(query.Variables)
	// Finalize
	if result2 != "" {
		log.WithField("=>", result2).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result2, query.Variables, nil
}

func QueryAccountJob(query *definitions.QueryAccount, do *definitions.Do) (string, error) {
	// Preprocess variables
	query.Account, _ = util.PreProcess(query.Account, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Perform Query
	arg := fmt.Sprintf("%s:%s", query.Account, query.Field)
	log.WithField("=>", arg).Info("Querying Account")

	result, err := util.AccountsInfo(query.Account, query.Field, do)
	if err != nil {
		return "", err
	}

	// Result
	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryNameJob(query *definitions.QueryName, do *definitions.Do) (string, error) {
	// Preprocess variables
	query.Name, _ = util.PreProcess(query.Name, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithFields(log.Fields{
		"name":  query.Name,
		"field": query.Field,
	}).Info("Querying")
	result, err := util.NamesInfo(query.Name, query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryValsJob(query *definitions.QueryVals, do *definitions.Do) (string, error) {
	var result string

	// Preprocess variables
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithField("=>", query.Field).Info("Querying Vals")
	result, err := util.ValidatorsInfo(query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}*/

/*func AssertJob(assertion *definitions.Assert, do *definitions.Do) (string, error) {
	var result string
	// Preprocess variables
	assertion.Key, _ = util.PreProcess(assertion.Key, do)
	assertion.Relation, _ = util.PreProcess(assertion.Relation, do)
	assertion.Value, _ = util.PreProcess(assertion.Value, do)

	// Switch on relation
	log.WithFields(log.Fields{
		"key":      assertion.Key,
		"relation": assertion.Relation,
		"value":    assertion.Value,
	}).Info("Assertion =>")

	switch assertion.Relation {
	case "==", "eq":
		if assertion.Key == assertion.Value {
			return assertPass("==", assertion.Key, assertion.Value)
		} else {
			return assertFail("==", assertion.Key, assertion.Value)
		}
	case "!=", "ne":
		if assertion.Key != assertion.Value {
			return assertPass("!=", assertion.Key, assertion.Value)
		} else {
			return assertFail("!=", assertion.Key, assertion.Value)
		}
	case ">", "gt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k > v {
			return assertPass(">", assertion.Key, assertion.Value)
		} else {
			return assertFail(">", assertion.Key, assertion.Value)
		}
	case ">=", "ge":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k >= v {
			return assertPass(">=", assertion.Key, assertion.Value)
		} else {
			return assertFail(">=", assertion.Key, assertion.Value)
		}
	case "<", "lt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k < v {
			return assertPass("<", assertion.Key, assertion.Value)
		} else {
			return assertFail("<", assertion.Key, assertion.Value)
		}
	case "<=", "le":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k <= v {
			return assertPass("<=", assertion.Key, assertion.Value)
		} else {
			return assertFail("<=", assertion.Key, assertion.Value)
		}
	}

	return result, nil
}

func bulkConvert(key, value string) (int, int, error) {
	k, err := strconv.Atoi(key)
	if err != nil {
		return 0, 0, err
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, err
	}
	return k, v, nil
}

func assertPass(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Succeeded")
	return "passed", nil
}

func assertFail(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Failed")
	return "failed", fmt.Errorf("assertion failed")
}

func convFail() (string, error) {
	return "", fmt.Errorf("The Key of your assertion cannot be converted into an integer.\nFor string conversions please use the equal or not equal relations.")
}*/
