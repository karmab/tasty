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
package main

import (
	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"tasty/cmd"
)

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf(color.InRed("[ERROR] %s"), err.Error())
	}
}

func newCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "tasty",
		Short: "This application allows you to interact with olm operators\n using a yum-like workflow",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	c.AddCommand(cmd.NewConfig())
	c.AddCommand(cmd.NewInfo())
	c.AddCommand(cmd.NewSearch())

	return c
}
