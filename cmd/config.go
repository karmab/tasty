/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "This options allow to perform configuration of tasty itself",
	Long: `This options allow to perform configuration of tasty itself. For example: you can install 
	tasty as kubectl and oc plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		enablePlugin, _ := cmd.Flags().GetBool("enable-as-plugin")

		_, execFile := path.Split(os.Args[0])

		filePath, err := exec.LookPath(execFile)
		if err != nil {
			color.Red("it is required to install tasty within a path that is in your $PATH environment variable")
			log.Fatalf("%s", err)
		}

		execPath, err := filepath.Abs(filepath.Dir(filePath))
		if err != nil {
			log.Fatal(err)
		}

		if enablePlugin {
			enableAsPlugin(execPath, execFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolP("enable-as-plugin", "p", false, "Install as kubeclt and oc plugin")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func enableAsPlugin(execPath, execFile string) {

	var found bool
	color.Cyan("Installing tasty as kubectl and oc CLI plugin")

	execOcLink := execPath + "/oc-olm"
	err := os.Symlink(execFile, execOcLink)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			color.Yellow("Oc Plugin already installed.")
			found = true
		} else {
			log.Fatal(err)
		}
	}

	execKubectlLink := execPath + "/kubectl-olm"
	err = os.Symlink(execFile, execKubectlLink)
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
