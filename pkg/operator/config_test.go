package operator

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	// Configure test environment
	execFile := "tasty"
	execPath := "/tmp/tasty"
	if err := os.Mkdir(execPath, 0777); err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}

	if _, err := os.Create(execPath + "/" + execFile); err != nil {
		t.Fatalf("Error creating testFile: %s", err)
	}

	// Run enableAsPlugin function
	o := &Operator{
		ConfigExecPath: execPath,
		ConfigExecFile: execFile,
	}
	o.enableAsPlugin()

	contentTempPaths, err := ioutil.ReadDir(execPath)
	if err != nil {
		log.Fatalf("Error reading execPath: %s", err)
	}

	linksFiles := 0
	for _, contentTempPath := range contentTempPaths {
		if contentTempPath.Name() == "oc-olm" {
			t.Log("Found oc-olm link")
			linksFiles++

		}
		if contentTempPath.Name() == "kubectl-olm" {
			t.Log("Found kubectl-olm link")
			linksFiles++
		}
	}

	if linksFiles != 2 {
		t.Error("Not all links were created")
	}

	cleanUp(execPath)

}

func cleanUp(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalf("Error removing testFolder: %s", err)
	}
}
