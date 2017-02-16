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

func TestCompilerNormal(t *testing.T) {

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

	t.Log("Contracts: ", solReturn.Contracts)
	t.Log("Warning: ", solReturn.Warning)
	t.Log("Error: ", solReturn.Error)
}

func TestCompilerError(t *testing.T) {

}

func TestDefaultCompilerUnmarshalling(t *testing.T) {

}

func TestLinkingBinaries(t *testing.T) {

}

func TestLinkingBinariesAndNormalCompileMixed(t *testing.T) {

}

func TestPullingDifferentVersions(t *testing.T) {

}

func TestPullingInvalidVersions(t *testing.T) {

}

func TestDefaultCompiling(t *testing.T) {

}

func TestDefinedCompiling(t *testing.T) {

}
