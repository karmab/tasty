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
	"log"

	"tasty/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Provides information about specified operator",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var operator string
		if len(args) != 1 {
			log.Printf("Usage: tasty info OPERATOR_NAME")
		} else {
			operator = args[0]
		}
		source, defaultchannel, csv, description, target_namespace, channels, _ := utils.GetOperator(operator)
		color.Cyan("Providing information for app %s", operator)
		fmt.Println("source: ", source)
		fmt.Println("channels: ", channels)
		fmt.Println("defaultchannel: ", defaultchannel)
		fmt.Println("target namespace: ", target_namespace)
		fmt.Println("csv: ", csv)
		fmt.Println("description: ", description)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
