package operator

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func (o *Operator) NewConfiguration(cmd *cobra.Command, args []string) error {
	enablePlugin, _ := cmd.Flags().GetBool("enable-as-plugin")

	_, execFile := path.Split(os.Args[0])

	filePath, err := exec.LookPath(execFile)
	if err != nil {
		color.Red("it is required to install tasty within a path that is in your $PATH environment variable")
		log.Fatalf("%s", err)
		return err
	}

	execPath, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		log.Fatal(err)
		return err
	}

	o.ConfigExecPath = execPath
	o.ConfigExecFile = execFile

	if enablePlugin {
		o.enableAsPlugin()
	}
	return nil
}

func (o *Operator) enableAsPlugin() {

	var found bool
	color.Cyan("Installing tasty as kubectl and oc CLI plugin")

	execOcLink := o.ConfigExecPath + "/oc-olm"
	err := os.Symlink(o.ConfigExecFile, execOcLink)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			color.Yellow("Oc Plugin already installed.")
			found = true
		} else {
			log.Fatal(err)
		}
	}

	execKubectlLink := o.ConfigExecPath + "/kubectl-olm"
	err = os.Symlink(o.ConfigExecFile, execKubectlLink)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			color.Yellow("Kubectl Plugin already installed.")
			found = true
		} else {
			log.Fatal(err)
		}
	}

	if !found {
		color.Green("Tasty installed successfully as oc and kubectl plugin")
	}

}
