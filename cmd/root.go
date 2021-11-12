/*
Copyright © 2021 karmab <EMAIL ADDRESS>

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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasty",
	Short: "Handles OLM operators",
	Long: `This application allows you to interact with olm operators
using a yum-like workflow`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		color.Red("KUBECONFIG env variable needs to be set")
		os.Exit(1)
	}
}
