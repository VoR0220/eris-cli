package maker

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"

	keys "github.com/eris-ltd/eris-keys/eris-keys"
)

func MakeAccounts(name, chainType string, accountTypes []*definitions.ErisDBAccountType) ([]*definitions.ErisDBAccount, error) {
	accounts := []*definitions.ErisDBAccount{}

	for _, accountT := range accountTypes {
		log.WithField("type", accountT.Name).Info("Making Account Type")

		perms := &definitions.ErisDBAccountPermissions{}
		var err error
		if chainType == "mint" {
			perms, err = ErisDBAccountPermissions(accountT.Perms, []string{}) // TODO: expose roles
			if err != nil {
				return nil, err
			}
		}

		for i := 0; i < accountT.Number; i++ {
			thisAct := &definitions.ErisDBAccount{}
			thisAct.Name = fmt.Sprintf("%s_%s_%03d", name, accountT.Name, i)
			thisAct.Name = strings.ToLower(thisAct.Name)

			log.WithField("name", thisAct.Name).Debug("Making Account")

			thisAct.Amount = accountT.Tokens
			thisAct.ToBond = accountT.ToBond

			thisAct.PermissionsMap = accountT.Perms
			thisAct.Validator = false

			if thisAct.ToBond != 0 {
				thisAct.Validator = true
			}

			if chainType == "mint" {
				thisAct.ErisDBPermissions = &definitions.ErisDBAccountPermissions{}
				thisAct.ErisDBPermissions.ErisDBBase = &definitions.ErisDBBasePermissions{}
				thisAct.ErisDBPermissions.ErisDBBase.ErisDBPerms = perms.ErisDBBase.ErisDBPerms
				thisAct.ErisDBPermissions.ErisDBBase.ErisDBSetBit = perms.ErisDBBase.ErisDBSetBit
				thisAct.ErisDBPermissions.ErisDBRoles = perms.ErisDBRoles
				log.WithField("perms", thisAct.ErisDBPermissions.ErisDBBase.ErisDBPerms).Debug()

				if err := makeKey("ed25519,ripemd160", thisAct); err != nil {
					return nil, err
				}
			}

			accounts = append(accounts, thisAct)
		}
	}

	return accounts, nil
}

func makeKey(keyType string, account *definitions.ErisDBAccount) error {
	log.WithFields(log.Fields{
		"path": keys.DaemonAddr,
		"type": keyType,
	}).Debug("Sending Call to eris-keys server")

	var err error
	log.WithField("endpoint", "gen").Debug()
	account.Address, err = keys.Call("gen", map[string]string{"auth": "", "type": keyType, "name": account.Name}) // note, for now we use not password to lock/unlock keys
	if _, ok := err.(keys.ErrConnectionRefused); ok {
		return fmt.Errorf("Could not connect to eris-keys server. Start it with `eris services start keys`. Error: %v", err)
	}
	if err != nil {
		return err
	}

	log.WithField("endpoint", "pub").Debug()
	account.PubKey, err = keys.Call("pub", map[string]string{"addr": account.Address, "name": account.Name})
	if _, ok := err.(keys.ErrConnectionRefused); ok {
		return fmt.Errorf("Could not connect to eris-keys server. Start it with `eris services start keys`. Error: %v", err)
	}
	if err != nil {
		return err
	}

	// log.WithField("endpoint", "to-mint").Debug()
	// mint, err := keys.Call("to-mint", map[string]string{"addr": account.Address, "name": account.Name})

	log.WithField("endpoint", "mint").Debug()
	mint, err := keys.Call("mint", map[string]string{"addr": account.Address, "name": account.Name})
	if _, ok := err.(keys.ErrConnectionRefused); ok {
		return fmt.Errorf("Could not connect to eris-keys server. Start it with `eris services start keys`. Error: %v", err)
	}
	if err != nil {
		return err
	}
	// [zr] leave MintKey / MintPrivValidator
	account.MintKey = &definitions.MintPrivValidator{}
	err = json.Unmarshal([]byte(mint), account.MintKey)
	if err != nil {
		log.Error(string(mint))
		log.Error(account.MintKey)
		return err
	}

	account.MintKey.Address = account.Address
	return nil
}

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
