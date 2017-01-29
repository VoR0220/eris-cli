package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	//"github.com/eris-ltd/eris/log"
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
		return Type{}, fmt.Errorf("Could not get epm type.")
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
