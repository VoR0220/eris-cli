package compilers

import (
	"os"
	"testing"

	"github.com/eris-ltd/eris/log"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCompilerNormal(t *testing.T) {

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
