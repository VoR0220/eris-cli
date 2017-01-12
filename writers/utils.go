package writers

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/log"
)

func writer(toWrangle interface{}, chain_name, account_name, fileBase string) error {
	fileBytes, err := json.MarshalIndent(toWrangle, "", "  ")
	if err != nil {
		return err
	}

	file := filepath.Join(config.ChainsPath, chain_name, account_name, fileBase)

	log.WithField("path", file).Debug("Saving File.")
	err = config.WriteFile(string(fileBytes), file)
	if err != nil {
		return err
	}
	return nil
}

func convertExportPortsSliceToString(exportPorts []string) string {
	if len(exportPorts) == 0 {
		return ""
	}
	return `[ "` + strings.Join(exportPorts[:], `", "`) + `" ]`
}
