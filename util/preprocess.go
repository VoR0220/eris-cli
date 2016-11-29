package util

import (
	"fmt"
	"strings"
)

func StringPreProcess(val string, do *Do) (string, error) {
	switch {
	/*case strings.HasPrefix(val, "$block"):
		return replaceBlockVariable(val, do)*/
	case strings.HasPrefix(val, "$"):
		key := strings.TrimPrefix(val, "$")
		if results, ok := do.Jobs.JobMap[key]; ok {
			index := strings.Index(key, ".")
			if index == -1 {
				return "", fmt.Errorf("Could not find results for job %v", index)
			} else {
				return results.JobVars[key[index:]], nil
			}
			return results.JobResult, nil
		}
		return "", fmt.Errorf("Could not find results for job %v", key)
	default: 
		return val, nil
	}
}