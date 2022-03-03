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
	"tasty/pkg/operator"

	"github.com/spf13/cobra"
)

func NewConfigurer() *cobra.Command {
	var o *operator.Operator
	cmd := &cobra.Command{
		Use:   "config",
		Short: "This options allow to perform configuration of tasty itself",
		Long: `This options allow to perform configuration of tasty itself. For example: you can install 
	tasty as kubectl and oc plugin.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			o = operator.NewOperator()
			return o.NewConfiguration(cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.BoolP("enable-as-plugin", "p", false, "Install as kubeclt and oc plugin")
	return cmd
}
