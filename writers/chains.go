package writers

import (
	"fmt"
	"path/filepath"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
)

func WriteGenesisFile(name string, genesis *definitions.ErisDBGenesis, account *definitions.ErisDBAccount) error {
	return writer(genesis, name, account.Name, "genesis.json")
}

func WriteConfigurationFile(chain_name, account_name, seeds string, chainImageName string,
	useDataContainer bool, exportedPorts []string, containerEntrypoint string) error {
	if account_name == "" {
		account_name = "anonymous_marmot"
	}
	if chain_name == "" {
		return fmt.Errorf("No chain name provided.")
	}
	var fileBytes []byte
	var err error
	if fileBytes, err = config.GetConfigurationFileBytes(chain_name,
		account_name, seeds, chainImageName, useDataContainer,
		convertExportPortsSliceToString(exportedPorts), containerEntrypoint); err != nil {
		return err
	}

	file := filepath.Join(config.ChainsPath, chain_name, account_name, "config.toml")

	log.WithField("path", file).Debug("Saving File.")
	if err := config.WriteFile(string(fileBytes), file); err != nil {
		return err
	}
	return nil
}
