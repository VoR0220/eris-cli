package writers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"

	"github.com/BurntSushi/toml"
)

// XXX: this is temporary until eris-keys.js is more tightly integrated with eris-contracts.js
type accountInfo struct {
	Address string `mapstructure:"address" json:"address" yaml:"address" toml:"address"`
	PubKey  string `mapstructure:"pubKey" json:"pubKey" yaml:"pubKey" toml:"pubKey"`
	PrivKey string `mapstructure:"privKey" json:"privKey" yaml:"privKey" toml:"privKey"`
}

func WritePrivVals(name string, account *definitions.ErisDBAccount) error {
	return writer(account.MintKey, name, account.Name, "priv_validator.json")
}

func SaveAccountResults(do *definitions.Do) error {
	addrFile, err := os.Create(filepath.Join(config.ChainsPath, do.Name, "addresses.csv"))
	if err != nil {
		return fmt.Errorf("Error creating addresses file. This usually means that there was a problem with the chain making process.")
	}
	defer addrFile.Close()

	log.WithField("name", do.Name).Debug("Creating file")
	actFile, err := os.Create(filepath.Join(config.ChainsPath, do.Name, "accounts.csv"))
	if err != nil {
		return fmt.Errorf("Error creating accounts file.")
	}
	log.WithField("path", filepath.Join(config.ChainsPath, do.Name, "accounts.csv")).Debug("File successfully created")
	defer actFile.Close()

	log.WithField("name", do.Name).Debug("Creating file")
	actJSONFile, err := os.Create(filepath.Join(config.ChainsPath, do.Name, "accounts.json"))
	if err != nil {
		return fmt.Errorf("Error creating accounts file.")
	}
	log.WithField("path", filepath.Join(config.ChainsPath, do.Name, "accounts.json")).Debug("File successfully created")
	defer actJSONFile.Close()

	valFile, err := os.Create(filepath.Join(config.ChainsPath, do.Name, "validators.csv"))
	if err != nil {
		return fmt.Errorf("Error creating validators file.")
	}
	defer valFile.Close()

	accountJsons := make(map[string]*accountInfo)

	for _, account := range do.Accounts {
		accountJsons[account.Name] = &accountInfo{
			Address: account.Address,
			PubKey:  account.PubKey,
			PrivKey: account.MintKey.PrivKey[1].(string),
		}

		_, err := addrFile.WriteString(fmt.Sprintf("%s,%s\n", account.Address, account.Name))
		if err != nil {
			log.Error("Error writing addresses file.")
			return err
		}
		_, err = actFile.WriteString(fmt.Sprintf("%s,%d,%s,%d,%d\n", account.PubKey, account.Amount, account.Name, account.ErisDBPermissions.ErisDBBase.ErisDBPerms, account.ErisDBPermissions.ErisDBBase.ErisDBSetBit))
		if err != nil {
			log.Error("Error writing accounts file.")
			return err
		}
		if account.Validator {
			_, err = valFile.WriteString(fmt.Sprintf("%s,%d,%s,%d,%d\n", account.PubKey, account.ToBond, account.Name, account.ErisDBPermissions.ErisDBBase.ErisDBPerms, account.ErisDBPermissions.ErisDBBase.ErisDBSetBit))
			if err != nil {
				log.Error("Error writing validators file.")
				return err
			}
		}
	}
	addrFile.Sync()
	actFile.Sync()
	valFile.Sync()

	j, err := json.MarshalIndent(accountJsons, "", "  ")
	if err != nil {
		return err
	}

	_, err = actJSONFile.Write(j)
	if err != nil {
		return err
	}

	log.WithField("path", actJSONFile.Name()).Debug("Saving File.")
	log.WithField("path", addrFile.Name()).Debug("Saving File.")
	log.WithField("path", actFile.Name()).Debug("Saving File.")
	log.WithField("path", valFile.Name()).Debug("Saving File.")

	return nil
}

func SaveAccountType(thisActT *definitions.ErisDBAccountType) error {
	writer, err := os.Create(filepath.Join(config.AccountsTypePath, fmt.Sprintf("%s.toml", thisActT.Name)))
	defer writer.Close()
	if err != nil {
		return err
	}

	enc := toml.NewEncoder(writer)
	enc.Indent = ""
	err = enc.Encode(thisActT)
	if err != nil {
		return err
	}
	return nil
}
