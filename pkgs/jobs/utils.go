package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eris-ltd/eris/log"

	"github.com/eris-ltd/eris-db/client/rpc"
	"github.com/eris-ltd/eris-db/txs"
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

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}

// This is a closer function which is called by most of the tx_run functions
func txFinalize(tx txs.Tx, jobs *Jobs) (*JobResults, error) {
	result, err := rpc.SignAndBroadcast(jobs.ChainID, jobs.NodeClient, jobs.KeyClient, tx, true, true, true)
	if err != nil {
		return MintChainErrorHandler(jobs, err)
	}
	// if there is nothing to unpack then just return.
	if result == nil {
		return nil, nil
	}

	// Unpack and display for the user.
	addr := fmt.Sprintf("%X", result.Address)
	hash := fmt.Sprintf("%X", result.Hash)
	blkHash := fmt.Sprintf("%X", result.BlockHash)
	ret := fmt.Sprintf("%X", result.Return)

	if result.Address != nil {
		log.WithField("=>", addr).Warn("Address")
		log.WithField("=>", hash).Info("Transaction Hash")
	} else {
		log.WithField("=>", hash).Warn("Transaction Hash")
		log.WithField("=>", blkHash).Debug("Block Hash")
		if len(result.Return) != 0 {
			if ret != "" {
				log.WithField("=>", ret).Warn("Return Value")
			} else {
				log.Debug("No return.")
			}
			log.WithField("=>", result.Exception).Debug("Exception")
		}
	}

	return &JobResults{Type{hash, result.Hash}, nil}, nil
}
