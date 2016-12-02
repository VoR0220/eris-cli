package jobs

/*import (
	"testing"
)

var testConfigFile string = `
jobs:

- name: setStorageBase
  set:
    val: 5

- name: deployStorageK
  deploy:
    contract: storage.sol

- name: setStorage
  call:
    destination: $deployStorageK
    function: set
    data:
      - [1, 2, 3]

- name: queryStorage
  query-contract:
    destination: $deployStorageK
    function: get

- name: assertStorage
  assert:
    key: $queryStorage
    relation: eq
    val: $setStorageBase
`

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPreProcessingStrings(t *testing.T) {

}*/