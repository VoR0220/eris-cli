package compilers

import (
	"os"
	"testing"

	"github.com/eris-ltd/eris/log"
	"github.com/eris-ltd/eris/util"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.InfoLevel)
	util.DockerConnect(false, "eris")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSolcCompilerNormal(t *testing.T) {

	var solFile string = `pragma solidity >= 0.0.0;
	contract main {
		uint a;
		function f() {
			a = 1;
		}
	}`
	file, err := os.Create("simpleContract.sol")
	defer os.Remove("simpleContract.sol")
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString(solFile)
	template := &SolcTemplate{
		CombinedOutput: []string{"bin", "abi"},
	}

	solReturn, err := template.Compile([]string{"simpleContract.sol"}, "stable")
	if err != nil {
		t.Fatal(err)
	}

	if solReturn.Error != nil || solReturn.Warning != "" || len(solReturn.Contracts) != 1 {
		t.Fatalf("Expected no errors or warnings and expected contract items. Got %v for errors, %v for warnings, and %v for contract items", solReturn.Error, solReturn.Warning, solReturn.Contracts)
	}
}

func TestSolcCompilerError(t *testing.T) {
	var solFile string = `pragma solidity >= 0.0.0;
	contract main {
		uint a;
		function f() {
			a = 1;
		}
	`
	file, err := os.Create("faultyContract.sol")
	defer os.Remove("faultyContract.sol")
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString(solFile)
	template := &SolcTemplate{
		CombinedOutput: []string{"bin", "abi"},
	}

	solReturn, err := template.Compile([]string{"faultyContract.sol"}, "stable")
	if err != nil {
		t.Fatal(err)
	}
	if solReturn.Error == nil {
		t.Fatal("Expected an error, got nil.")
	}
}

func TestSolcCompilerWarning(t *testing.T) {
	var solFile string = `contract main {
		uint a;
		function f() {
			a = 1;
		}
	}`
	file, err := os.Create("simpleContract.sol")
	defer os.Remove("simpleContract.sol")
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString(solFile)
	template := &SolcTemplate{
		CombinedOutput: []string{"bin", "abi"},
	}

	solReturn, err := template.Compile([]string{"simpleContract.sol"}, "stable")
	if err != nil {
		t.Fatal(err)
	}
	if solReturn.Warning == "" {
		t.Fatal("Expected a warning.")
	}
}

func TestLinkingBinaries(t *testing.T) {

}

func TestLinkingBinariesAndNormalCompileMixed(t *testing.T) {

}

func TestDefaultCompilerUnmarshalling(t *testing.T) {

}

func TestPullingDifferentVersions(t *testing.T) {

}

func TestPullingInvalidVersions(t *testing.T) {

}

func TestDefaultCompiling(t *testing.T) {

}

func TestDefinedCompiling(t *testing.T) {

}
