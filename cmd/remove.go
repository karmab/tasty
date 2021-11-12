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
	"os"
	"strings"
	"tasty/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
		subscriptionsGVR := schema.GroupVersionResource{
			Group:    "operators.coreos.com",
			Version:  "v1alpha1",
			Resource: "subscriptions",
		}
		operatorgroupsGVR := schema.GroupVersionResource{
			Group:    "operators.coreos.com",
			Version:  "v1",
			Resource: "operatorgroups",
		}
		for _, operator := range args {
			color.Cyan("Removing operator %s", operator)
			_, _, _, _, target_namespace, _, _ := utils.GetOperator(operator)
			dynamic := utils.GetDynamicClient()
			color.Cyan("Removing subscription %s", operator)
			err := dynamic.Resource(subscriptionsGVR).Namespace(target_namespace).Delete(context.TODO(), operator, metav1.DeleteOptions{})
			utils.Check(err)
			if target_namespace != "openshift-operators" {
				color.Cyan("Removing operator group %s-operatorgroup", operator)
				k8sclient := utils.GetK8sClient()
				err := dynamic.Resource(operatorgroupsGVR).Namespace(target_namespace).Delete(context.TODO(), operator+"-operatorgroup", metav1.DeleteOptions{})
				utils.Check(err)
				color.Cyan("Removing namespace group %s", target_namespace)
				err = k8sclient.CoreV1().Namespaces().Delete(context.TODO(), target_namespace, metav1.DeleteOptions{})
				utils.Check(err)
			}
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
