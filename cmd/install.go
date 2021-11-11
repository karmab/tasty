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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [operator]",
	Short: "Install operators",
	Long:  `Install operators`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, operator := range args {
			color.Cyan("Installing operator %s", operator)
			stdout, _ := cmd.Flags().GetBool("stdout")
			wait, _ := cmd.Flags().GetBool("wait")
			targetchannel, _ := cmd.Flags().GetString("channel")
			source, defaultchannel, csv, _, target_namespace, channels, crd := get_operator(operator)
			if targetchannel != "" {
				if contains(channels, targetchannel) {
					defaultchannel = targetchannel
				} else {
					color.Red("Target channel %s not found in %s", targetchannel, channels)
					os.Exit(1)
				}
			}
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
			if stdout == true {
				err = tpl.Execute(os.Stdout, operatordata)
				check(err)
			} else {
				buf := &bytes.Buffer{}
				err = tpl.Execute(buf, operatordata)
				check(err)
				tmpfile, err := os.CreateTemp("", "tasty")
				check(err)
				_, err = tmpfile.Write(buf.Bytes())
				check(err)
				tmpfile.Close()
				applyout, err := exec.Command("oc", "apply", "-f", tmpfile.Name()).Output()
				check(err)
				fmt.Println(string(applyout))
				os.Remove(tmpfile.Name())
				if wait == true {
					wait_crd(crd, 60)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("channel", "c", "", "Target channel")
	installCmd.Flags().BoolP("stdout", "s", false, "Print to stdout")
	installCmd.Flags().BoolP("wait", "w", false, "Wait for crd to show up")
}
