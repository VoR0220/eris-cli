package loaders

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"

	"github.com/spf13/viper"
)

func LoadAccountTypes() ([]*definitions.ErisDBAccountType, error) {
	loadedAccounts := []*definitions.ErisDBAccountType{}
	accounts, err := AccountTypes(config.AccountsTypePath)
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		thisAct, err := LoadAccountType(account)
		if err != nil {
			return nil, err
		}
		log.WithField("=>", thisAct.Name).Debug("Loaded Account Named")
		loadedAccounts = append(loadedAccounts, thisAct)
	}
	return loadedAccounts, nil
}

func LoadAccountType(fileName string) (*definitions.ErisDBAccountType, error) {
	log.WithField("=>", fileName).Debug("Loading Account Type")
	var accountType = viper.New()
	typ := definitions.BlankAccountType()

	if err := getSetup(fileName, accountType); err != nil {
		return nil, err
	}

	// marshall file
	if err := accountType.Unmarshal(typ); err != nil {
		return nil, fmt.Errorf(`Sorry, the marmots could not figure that account types file out.
			Please check your account type definition file is properly formatted: 
			%v`, err)
	}

	return typ, nil
}

// returns a list of filenames which are the account_types files
// these *should be* absolute paths, but this is not a contract
// with calling functions.
func AccountTypes(erisPath string) ([]string, error) {
	haveTyps, err := filepath.Glob(filepath.Join(erisPath, "*.toml"))
	if err != nil {
		return []string{}, err
	}
	return haveTyps, nil
}

func getSetup(fileName string, cfg *viper.Viper) error {
	// setup file
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return fmt.Errorf(`Sorry, the marmots were unable to find the absolute path to the account types file.`)
	}

	path := filepath.Dir(abs)
	file := filepath.Base(abs)
	extName := filepath.Ext(file)
	bName := file[:len(file)-len(extName)]

	cfg.AddConfigPath(path)
	cfg.SetConfigName(bName)
	cfg.SetConfigType(strings.Replace(extName, ".", "", 1))

	// load file
	if err := cfg.ReadInConfig(); err != nil {
		return fmt.Errorf(`Sorry, the marmots were unable to load the file: (%s). Please check your path: %v`, fileName, err)
	}

	return nil
}

func AccountTypesNames(erisPath string, withExt bool) ([]string, error) {
	files, err := AccountTypes(erisPath)
	if err != nil {
		return []string{}, err
	}
	names := []string{}
	for _, file := range files {
		names = append(names, filepath.Base(file))
	}
	if !withExt {
		for e, name := range names {
			names[e] = strings.Replace(name, ".toml", "", 1)
		}
	}
	return names, nil
}
