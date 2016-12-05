
import (
	"fmt"
	"strings"
)

func stringPreProcess(val string, jobs *Jobs) (string, error) {
	switch {
	/*case strings.HasPrefix(val, "$block"):
		return replaceBlockVariable(val, do)*/
	case strings.HasPrefix(val, "$"):
		key := strings.TrimPrefix(val, "$")
		if results, ok := jobs.JobMap[key]; ok {
			index := strings.Index(key, ".")
			if index == -1 {
				return "", fmt.Errorf("Could not find results for job %v", index)
			} else {
				return results.JobVars[key[index:]], nil
			}
			return results.JobResult, nil
		}
		return "", fmt.Errorf("Could not find results for job %v", key)
	default : 
		return val, nil
	}
}

func replaceBlockVariable(toReplace string, do *definitions.Do) (string, error) {
	log.WithFields(log.Fields{
		"chain": do.ChainName,
		"var":   toReplace,
	}).Debug("Correcting $block variable")
	blockHeight, err := GetBlockHeight(do)
	block := strconv.Itoa(blockHeight)
	log.WithField("=>", block).Debug("Current height is")
	if err != nil {
		return "", err
	}

	if toReplace == "$block" {
		log.WithField("=>", block).Debug("Replacement (=)")
		return block, nil
	}

	catchEr := regexp.MustCompile("\\$block\\+(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 + h2)
		log.WithField("=>", height).Debug("Replacement (+)")
		return height, nil
	}

	catchEr = regexp.MustCompile("\\$block\\-(\\d*)")
	if catchEr.MatchString(toReplace) {
		height := catchEr.FindStringSubmatch(toReplace)[1]
		h1, err := strconv.Atoi(height)
		if err != nil {
			return "", err
		}
		h2, err := strconv.Atoi(block)
		if err != nil {
			return "", err
		}
		height = strconv.Itoa(h1 - h2)
		log.WithField("=>", height).Debug("Replacement (-)")
		return height, nil
	}

	log.WithField("=>", toReplace).Debug("Replacement (unknown)")
	return toReplace, nil
}

func PreProcessInputData(function string, data interface{}, do *definitions.Do, constructor bool) (string, []string, error) {
	var callDataArray []string
	var callArray []string
	if function == "" && !constructor {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			return "", []string{""}, fmt.Errorf("Incorrect formatting of epm run file. Please update your epm run file to include a function field.")
		}
		log.Warn("Deprecation Warning: The use of the 'data' field to specify the name of the contract function has been deprecated. Please update your epm jobs file to utilize a combination of 'function' and 'data' fields instead. See documentation for further details.")
		function = strings.Split(data.(string), " ")[0]
		callArray = strings.Split(data.(string), " ")[1:]
		for _, val := range callArray {
			output, _ := PreProcess(val, do)
			callDataArray = append(callDataArray, output)
		}
	} else if data != nil {
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			if constructor {
				log.Warn("Deprecation Warning: Your deploy job is currently using a soon to be deprecated way of declaring constructor values. Please remember to update your run file to store them as a array rather than a string. See documentation for further details.")
				callArray = strings.Split(data.(string), " ")
				for _, val := range callArray {
					output, _ := PreProcess(val, do)
					callDataArray = append(callDataArray, output)
				}
				return function, callDataArray, nil
			} else {
				return "", make([]string, 0), fmt.Errorf("Incorrect formatting of epm run file. Please update your epm run file to include a function field.")
			}
		}
		val := reflect.ValueOf(data)
		for i := 0; i < val.Len(); i++ {
			s := val.Index(i)
			var newString string
			switch s.Interface().(type) {
			case bool:
				newString = strconv.FormatBool(s.Interface().(bool))
			case int, int32, int64:
				newString = strconv.FormatInt(int64(s.Interface().(int)), 10)
			case []interface{}:
				var args []string
				for _, index := range s.Interface().([]interface{}) {
					value := reflect.ValueOf(index)
					var stringified string
					switch value.Kind() {
					case reflect.Int:
						stringified = strconv.FormatInt(value.Int(), 10)
					case reflect.String:
						stringified = value.String()
					}
					index, _ = PreProcess(stringified, do)
					args = append(args, stringified)
				}
				newString = "[" + strings.Join(args, ",") + "]"
				log.Debug(newString)
			default:
				newString = s.Interface().(string)
			}
			newString, _ = PreProcess(newString, do)
			callDataArray = append(callDataArray, newString)
		}
	}
	return function, callDataArray, nil
}

func SplitAndPreProcessStringPairs(value string, operator string, do *definitions.Do) (string, error) {
	if strings.HasPrefix(value, "$") {
		return StringPreProcess(value, do)
	} else if value == "" {
		return "", nil
	} else {
		var result string
		pairs := strings.Split(value, operator)
		switch operator {
		case "=":
			log.Info("Preprocessing remappings")
		case ":":
			log.Info("Preprocessing libraries")
			strings.ToLower(pairs[1])
		default : 
			return "", fmt.Errorf("Got invalid operator: %v", operator)
		}	
		pairs[0] = StringPreProcess(pairs[0], do)
		pairs[1] = StringPreProcess(pairs[1], do)
		result = strings.Join(pairs, operator)
		return result, nil
	}
}


func GetReturnValue(vars []*definitions.Variable) string {
	var result []string

	if len(vars) > 1 {
		for _, value := range vars {
			log.WithField("=>", []byte(value.Value)).Debug("Value")
			result = append(result, value.Value)
		}
		return "(" + strings.Join(result, ", ") + ")"
	} else if len(vars) == 1 {
		log.Debug("Debugging: ", vars[0].Value)
		return vars[0].Value
	} else {
		return ""
	}
}

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}