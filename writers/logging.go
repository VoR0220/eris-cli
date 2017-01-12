package writers

import (
	"os"

	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
)

// ------------------------------------------------------------------------
// Logging
// ------------------------------------------------------------------------

func ClearJobResults() error {
	if err := os.Remove(setJsonPath()); err != nil {
		return err
	}

	return os.Remove(setCsvPath())
}

func PrintPathPackage(do *definitions.Do) {
	log.WithField("=>", do.Compiler).Info("Using Compiler at")
	log.WithField("=>", do.ChainName).Info("Using Chain at")
	log.WithField("=>", do.ChainID).Debug("With ChainID")
	log.WithField("=>", do.Signer).Info("Using Signer at")
}
