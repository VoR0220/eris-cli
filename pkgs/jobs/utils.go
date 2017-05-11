package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/monax/cli/log"
	"github.com/monax/cli/pkgs/abi"

	"github.com/hyperledger/burrow/client/rpc"
	"github.com/hyperledger/burrow/txs"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
)

//preprocesses an interface type into a type type
func preProcessInterface(toProcess interface{}, jobs *Jobs) (Type, error) {
	switch typ := toProcess.(type) {
	case string:
		typString, typActual, err := preProcessString(typ, jobs)
		return Type{StringResult: typString, ActualResult: typActual}, err
	case bool:
		return Type{StringResult: fmt.Sprintf("%v", typ), ActualResult: typ}, nil
	case int:
		return Type{StringResult: fmt.Sprintf("%v", typ), ActualResult: typ}, nil
	case []byte:
		typString, typActual, err := preProcessString(string(typ), jobs)
		return Type{StringResult: typString, ActualResult: typActual}, err
	case []interface{}:
		var wrangledTypes []interface{}
		for _, toWrangle := range typ {
			if wrangled, err := preProcessInterface(toWrangle, jobs); err == nil {
				wrangledTypes = append(wrangledTypes, wrangled.ActualResult)
			} else {
				return Type{}, err
			}
		}
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(wrangledTypes)
		return Type{StringResult: strings.Trim(buf.String(), " \n\t"), ActualResult: wrangledTypes}, err
	case Type:
		return typ, nil
	default:
		return Type{}, fmt.Errorf("Could not get job type.")
	}
}

// preprocesses a string for $ references that indicate a job result and returns them if found
func preProcessString(key string, jobs *Jobs) (string, interface{}, error) {
	unfound := "Could not find results for job %v"

	switch {
	/*case strings.HasPrefix(val, "$block"): //todo: [rj] include this when we get to bond jobs
	return replaceBlockVariable(val, do)*/
	case strings.HasPrefix(key, "$"):
		var name string
		key = strings.TrimPrefix(key, "$")
		if index := strings.Index(key, "."); index != -1 {
			name = key[index+1:]
			key = key[:index]
		}
		if typeType, ok := jobs.JobMap[key]; ok {
			if len(name) > 1 {
				if namedResults, ok := typeType.NamedResults[name]; !ok {
					return "", nil, fmt.Errorf(unfound, name)
				} else {
					return namedResults.StringResult, namedResults.ActualResult, nil
				}
			}
			return typeType.FullResult.StringResult, typeType.FullResult.ActualResult, nil
		}
		return "", nil, fmt.Errorf(unfound, key)
	default:
		return key, key, nil
	}
}

// preprocesses for a job itself that has already been run.
func preProcessPluginJob(plugin interface{}, jobs *Jobs) (JobsRunner, error) {
	unfound := "Cannot deduce valid plugin type from %v"
	switch plugin := plugin.(type) {
	case string:
		if strings.HasPrefix(plugin, "$") {
			key := strings.TrimPrefix(plugin, "$")
			for _, job := range jobs.Jobs {
				if job.Name == key {
					return job.getType()
				}
			}
		}
		return nil, fmt.Errorf(unfound, plugin)
	default:
		return nil, fmt.Errorf(unfound, plugin)
	}
}

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}

// This is a utility function for signing, broadcasting and gathering a return,
func txFinalize(tx txs.Tx, jobs *Jobs, request TxResult) (*JobResults, error) {
	result, err := rpc.SignAndBroadcast(jobs.ChainID, jobs.NodeClient, jobs.KeyClient, tx, true, true, true)
	if err != nil {
		return MintChainErrorHandler(jobs, err)
	}
	// if there is nothing to unpack then just return.
	if result == nil {
		return nil, nil
	}

	// Unpack and display for the user.
	addr := fmt.Sprintf("%X", result.Hash)
	hash := fmt.Sprintf("%X", result.Hash)
	blkHash := fmt.Sprintf("%X", result.BlockHash)
	ret := fmt.Sprintf("%X", result.Return)

	if result.Address != nil {
		log.WithField("=>", addr).Warn("Address")
	}
	log.WithField("=>", hash).Info("Transaction Hash")
	log.WithField("=>", blkHash).Debug("Block Hash")
	if ret != "" {
		log.WithField("=>", ret).Warn("Return Value")
	} else {
		log.Warn("No return.")
	}
	log.WithField("=>", result.Exception).Debug("Exception")

	switch request {
	case TxHash:
		return &JobResults{Type{hash, result.Hash}, nil}, nil
	case Address:
		return &JobResults{Type{addr, result.Address}, nil}, nil
	case Return:
		return &JobResults{Type{ret, result.Return}, nil}, nil
	case BlockHash:
		return &JobResults{Type{blkHash, result.BlockHash}, nil}, nil
	default:
		return &JobResults{Type{ret, result.Exception}, nil}, fmt.Errorf(result.Exception)
	}
}

func Unpacker(contractAbi ethAbi.ABI, funcName string, output []byte) (*JobResults, error) {
	if output == nil {
		return &JobResults{}, nil
	}
	namedResults := make(map[string]Type)
	// Create blank slate of whether we are looking to use an interface or an []interface{}
	toUnpackInto, method, err := abi.CreateBlankSlate(contractAbi, funcName)
	if err != nil {
		return &JobResults{}, err
	}
	err = contractAbi.Unpack(&toUnpackInto, funcName, output)
	if err != nil {
		return &JobResults{}, err
	}
	// get names of the types, get string results, get actual results, return them.
	fullStringResults := []string{}
	for i, methodOutput := range method.Outputs {
		if methodOutput.Name == "" {
			methodOutput.Name = strconv.FormatInt(int64(i), 10)
		}
		var strResult string
		if unpackedValues, ok := toUnpackInto.([]interface{}); ok {
			strResult, err = abi.GetStringValue(unpackedValues[i], methodOutput.Type)
			log.WithField("=>", strResult).Debug("Result")
			namedResults[methodOutput.Name] = Type{ActualResult: unpackedValues, StringResult: strResult}
		} else {
			if unpackedValues, ok := toUnpackInto.(interface{}); !ok {
				return &JobResults{}, fmt.Errorf("Unexpected error while converting unpacked values to readable format")
			} else {
				strResult, err = abi.GetStringValue(unpackedValues, methodOutput.Type)
				log.WithField("=>", strResult).Debug("Result")
				namedResults[methodOutput.Name] = Type{ActualResult: unpackedValues, StringResult: strResult}
			}
		}
		if err != nil {
			return &JobResults{}, err
		}
		fullStringResults = append(fullStringResults, strResult)

	}
	return &JobResults{
		FullResult: Type{
			StringResult: "(" + strings.Join(fullStringResults, ", ") + ")",
			ActualResult: toUnpackInto},
		NamedResults: namedResults,
	}, nil
}

func WriteJobResultCSV(name, result string) error {

	pwd, _ := os.Getwd()
	logFile := filepath.Join(pwd, "jobs_output.csv")

	var file *os.File
	var err error

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err = os.Create(logFile)
	} else {
		file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0600)
	}

	if err != nil {
		return err
	}

	defer file.Close()

	text := fmt.Sprintf("%s,%s\n", name, result)
	if _, err = file.WriteString(text); err != nil {
		return err
	}

	return nil
}
