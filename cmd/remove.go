// Package cmd
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
	"github.com/karmab/tasty/pkg/operator"

	"github.com/spf13/cobra"
)

func NewRemover() *cobra.Command {
	var o *operator.Operator
	var (
		removed bool
		rmns    bool
		rmgroup bool
		ns      string
	)
	cmd := &cobra.Command{
		Use:          "remove",
		Short:        "Remove Operator",
		Long:         `Remove Operators`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			o = operator.NewOperator()
			return o.RemoveOperator(ns, removed, rmns, rmgroup, args)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&removed, "yes", "y", false, "Confirm")
	flags.BoolVarP(&rmns, "rmns", "r", false, "Remove namespace")
	flags.BoolVarP(&rmgroup, "rmgroups", "g", false, "Remove all operator groups in namespace")
	flags.StringVarP(&ns, "namespace", "n", "", "namespace")
	return cmd
}
