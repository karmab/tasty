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
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [operator]",
	Short: "Install operator",
	Long: `Install operator
	Examples needed here`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operator := args[0]
		outputdir, _ := cmd.Flags().GetString("output")
		fmt.Println("install called with " + operator)
		fmt.Println("output set to " + outputdir)
		namespace, source, defaultchannel, csv, _, target_namespace, crd := get_operator(operator)
		fmt.Println("namespace: " + namespace)
		fmt.Println("source: " + source)
		fmt.Println("defaultchannel: " + defaultchannel)
		fmt.Println("target_namespace: " + target_namespace)
		fmt.Println("csv: " + csv)
		fmt.Println("crd: " + crd)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("output", "o", "", "Output directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
