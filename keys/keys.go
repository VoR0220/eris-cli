package keys

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/data"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/services"

	eKeys "github.com/eris-ltd/eris-keys/eris-keys"
)

type KeyClient struct {
	IpAddr string
}

func InitKeyClient() (*KeyClient, error) {
	keys := &KeyClient{}
	err := keys.ensureRunning()
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "darwin" {
		eKeys.DaemonAddr = "http://127.0.0.1:4767"
		keys.IpAddr = "http://127.0.0.1:4767"
	} else {
		eKeys.DaemonAddr = "http://172.17.0.2:4767"
		keys.IpAddr = "http://172.17.0.2:4767"
	}
	return keys, nil
}

func (keys *KeyClient) ListKeys(host, container, quiet bool) ([]string, error) {
	var result []string
	if host {
		keysPath := filepath.Join(config.KeysPath, "data")
		addrs, err := ioutil.ReadDir(keysPath)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			result = append(result, addr.Name())
		}
		if !quiet {
			if len(addrs) == 0 {
				log.Warn("No keys found on host")
			} else {
				// First key.
				log.WithField("=>", result[0]).Warn("The keys on your host kind marmot")
				// Subsequent keys.
				if len(result) > 1 {
					for _, addr := range result[1:] {
						log.WithField("=>", addr).Warn()
					}
				}
			}
		}
	}

	if container {
		err := keys.ensureRunning()
		if err != nil {
			return nil, err
		}

		keysOut, err := services.ExecHandler("keys", []string{"ls", "/home/eris/.eris/keys/data"})
		if err != nil {
			return nil, err
		}
		result = strings.Fields(keysOut.String())
		if !quiet {
			if len(result) == 0 || result[0] == "" {
				log.Warn("No keys found in container")
			} else {
				// First key.
				log.WithField("=>", result[0]).Warn("The keys in your container kind marmot")
				// Subsequent keys.
				if len(result) > 1 {
					for _, addr := range result[1:] {
						log.WithField("=>", addr).Warn()
					}
				}
			}
		}
	}
	return result, nil
}

func (keys *KeyClient) GenerateKey(save bool, password string) error {
	err := keys.ensureRunning()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if password != "" {
		return fmt.Errorf("Password currently unimplemented. Marmots are confused at how you got here.")
	} else {
		buf, err = services.ExecHandler("keys", []string{"eris-keys", "gen", "--no-pass"})
		if err != nil {
			return err
		}
	}

	if save {
		addr := new(bytes.Buffer)
		addr.ReadFrom(buf)

		exportAddress := strings.TrimSpace(addr.String())

		log.WithField("=>", exportAddress).Warn("Saving key to host")
		if err := keys.ExportKey(exportAddress, false); err != nil {
			return err
		}
	}

	io.Copy(config.Global.Writer, buf)

	return nil
}

func (keys *KeyClient) ExportKey(address string, all bool) error {
	err := keys.ensureRunning()
	if err != nil {
		return err
	}
	do := definitions.NowDo()
	if all && address == "" {
		do.Destination = config.KeysPath
		do.Source = path.Join(config.KeysContainerPath)
	} else {
		do.Destination = config.KeysDataPath
		do.Source = path.Join(config.KeysContainerPath, do.Address)
	}
	return data.ExportData(do)
}

func (keys *KeyClient) ImportKey(address string, all bool) error {
	err := keys.ensureRunning()
	if err != nil {
		return err
	}

	do := definitions.NowDo()
	if all && address == "" {
		// get all keys from host
		result, err := keys.ListKeys(true, false, true)
		if err != nil {
			return err
		}
		// flip them for the import
		do.Container = true
		do.Host = false
		do.Quiet = false
		for _, addr := range result {
			do.Source = filepath.Join(config.KeysDataPath, addr)
			do.Destination = path.Join(config.KeysContainerPath, addr)
			if err := data.ImportData(do); err != nil {
				return err
			}
		}
	} else {
		do.Source = filepath.Join(config.KeysDataPath, do.Address)
		do.Destination = path.Join(config.KeysContainerPath, do.Address)
		if err := data.ImportData(do); err != nil {
			return err
		}
	}

	return nil
}

func (keys *KeyClient) ensureRunning() (err error) {
	doKeys := definitions.NowDo()
	doKeys.Name = "keys"
	err = services.EnsureRunning(doKeys)
	return
}

func (keys *KeyClient) PubKey(address string) (string, error) {
	err := keys.ensureRunning()
	if err != nil {
		return "", err
	}

	addr := strings.TrimSpace(address)
	buf, err := services.ExecHandler("keys", []string{"eris-keys", "pub", "--addr", addr, "--name", ""})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
