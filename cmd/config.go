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
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		enablePlugin, _ := cmd.Flags().GetBool("enable-as-plugin")
		_, execFile := path.Split(os.Args[0])
		execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
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
		color.Yellow("Oc Plugin already installed.")
		found = true
	}

	execKubectlLink := execPath + "/kubectl-olm"
	err = os.Symlink(execFile, execKubectlLink)
	if err != nil {
		color.Yellow("Kubectl Plugin already installed.")
		found = true
	}
	if found {
		// TODO: check if tasty is in a valid $PATH and if not, add the symlink to one in $PATH
		color.Cyan("Tasty already installed as oc and kubectl plugin, you can try oc olm --help or kubeclt olm --help")
	}
}
