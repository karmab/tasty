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
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove operators",
	Long:  `Remove operators`,
	Run: func(cmd *cobra.Command, args []string) {
		yes, _ := cmd.Flags().GetBool("stdout")
		if yes != true {
			var confirmation string
			color.Green("Are you sure? [y/N]:")
			fmt.Scanln(&confirmation)
			if strings.ToLower(confirmation) != "y" {
				color.Red("Leaving..")
				os.Exit(0)
			}
		}
		for _, operator := range args {
			color.Cyan("Removing operator %s", operator)
			source, defaultchannel, csv, _, target_namespace, _, _, _ := get_operator(operator)
			t := template.New("Template")
			tpl, err := t.Parse(operatordata)
			check(err)
			operatordata := Operator{
				Name:           operator,
				Source:         source,
				DefaultChannel: defaultchannel,
				Csv:            csv,
				Namespace:      target_namespace,
			}
			buf := &bytes.Buffer{}
			err = tpl.Execute(buf, operatordata)
			check(err)
			tmpfile, err := os.CreateTemp("", "tasty")
			check(err)
			_, err = tmpfile.Write(buf.Bytes())
			check(err)
			tmpfile.Close()
			applyout, _ := exec.Command("oc", "delete", "-f", tmpfile.Name()).Output()
			fmt.Println(string(applyout))
			os.Remove(tmpfile.Name())
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("yes", "y", false, "Confirm")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
