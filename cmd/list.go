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
	"context"
	"fmt"
	"sort"
	"tasty/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/syohex/go-texttable"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var result map[string]interface{}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List operators",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var operators []string
		dynamic := utils.GetDynamicClient()
		installed, _ := cmd.Flags().GetBool("installed")
		if installed {
			subscriptionsGVR := schema.GroupVersionResource{
				Group:    "operators.coreos.com",
				Version:  "v1alpha1",
				Resource: "subscriptions",
			}
			list, err := dynamic.Resource(subscriptionsGVR).Namespace("").List(context.TODO(), metav1.ListOptions{})
			utils.Check(err)
			for _, d := range list.Items {
				operators = append(operators, d.GetName())
			}
		} else {
			packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
			list, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
			utils.Check(err)
			for _, d := range list.Items {
				operators = append(operators, d.GetName())
			}
		}
		sort.Strings(operators)
		operatortable := &texttable.TextTable{}
		operatortable.SetHeader("Name")
		for _, operator := range operators {
			operatortable.AddRow(operator)
		}
		fmt.Println(operatortable.Draw())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("installed", "i", false, "Display installed operators")
}
