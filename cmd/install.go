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
	"html/template"
	"os"
	"tasty/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [operator]",
	Short: "Install operators",
	Long:  `Install operators`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, operator := range args {
			stdout, _ := cmd.Flags().GetBool("stdout")
			wait, _ := cmd.Flags().GetBool("wait")
			// sno, _ := cmd.Flags().GetBool("sno")
			target_channel, _ := cmd.Flags().GetString("channel")
			target_namespace, _ := cmd.Flags().GetString("namespace")
			source, default_channel, csv, _, default_namespace, channels, crd := utils.GetOperator(operator)
			if target_channel != "" {
				if utils.Contains(channels, target_channel) {
					default_channel = target_channel
				} else {
					color.Red("Target channel %s not found in %s", target_channel, channels)
					os.Exit(1)
				}
			}
			if target_namespace == "" {
				target_namespace = default_namespace
			}
			if stdout == true {
				t := template.New("Template")
				tpl, err := t.Parse(utils.OperatorTemplate)
				utils.Check(err)
				operatordata := utils.Operator{
					Name:           operator,
					Source:         source,
					DefaultChannel: default_channel,
					Csv:            csv,
					Namespace:      target_namespace,
				}
				err = tpl.Execute(os.Stdout, operatordata)
				utils.Check(err)
			} else {
				color.Cyan("Installing operator %s", operator)
				dynamic := utils.GetDynamicClient()
				if target_namespace != "openshift-operators" {
					color.Cyan("Creating namespace %s", target_namespace)
					k8sclient := utils.GetK8sClient()
					nsmetaSpec := metav1.ObjectMeta{Name: target_namespace, Annotations: map[string]string{"workload.openshift.io/allowed": "management"}}
					nsSpec := &v1.Namespace{ObjectMeta: nsmetaSpec}
					_, err := k8sclient.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
					utils.Check(err)
					color.Cyan("Creating operator group %s-operatorgroup", operator)
					operatorgroupsGVR := schema.GroupVersionResource{
						Group:    "operators.coreos.com",
						Version:  "v1",
						Resource: "operatorgroups",
					}

					operatorgroupspec := &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "OperatorGroup",
							"apiVersion": "operators.coreos.com/v1",
							"metadata": map[string]interface{}{
								"name":      operator + "-operatorgroup",
								"namespace": target_namespace,
							},
							"spec": map[string]interface{}{
								"targetNamespaces": []string{target_namespace},
							},
						},
					}
					_, err = dynamic.Resource(operatorgroupsGVR).Namespace(target_namespace).Create(context.TODO(), operatorgroupspec, metav1.CreateOptions{})
					utils.Check(err)
				}
				color.Cyan("Creating subscription %s", operator)
				subscriptionsGVR := schema.GroupVersionResource{
					Group:    "operators.coreos.com",
					Version:  "v1alpha1",
					Resource: "subscriptions",
				}

				subspec := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Subscription",
						"apiVersion": "operators.coreos.com/v1alpha1",
						"metadata": map[string]interface{}{
							"name":      operator,
							"namespace": target_namespace,
						},
						"spec": map[string]interface{}{
							"channel":         target_channel,
							"name":            operator,
							"source":          source,
							"sourceNamespace": "openshift-marketplace",
						},
					},
				}
				_, err := dynamic.Resource(subscriptionsGVR).Namespace(target_namespace).Create(context.TODO(), subspec, metav1.CreateOptions{})
				utils.Check(err)
				if wait == true && crd != "" {
					utils.WaitCrd(crd, 60)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("channel", "c", "", "Target channel")
	installCmd.Flags().StringP("namespace", "n", "", "Target namespace")
	installCmd.Flags().BoolP("stdout", "s", false, "Print to stdout")
	installCmd.Flags().BoolP("wait", "w", false, "Wait for crd to show up")
}
