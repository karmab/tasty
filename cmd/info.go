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
	"log"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
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
		kubeconfig, _ := os.LookupEnv("KUBECONFIG")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Panicln("failed to create K8s config")
		}
		client, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Panicln("Failed to create K8s clientset")
		}
		packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
		operatorinfo, err := client.Resource(packagemanifests).Namespace("openshift-marketplace").Get(context.TODO(), operator, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		namespace, _, err := unstructured.NestedString(operatorinfo.Object, "metadata", "namespace")
		if err != nil {
			log.Printf("Error getting namespace %v", err)
		}
		source, _, err := unstructured.NestedString(operatorinfo.Object, "status", "catalogSource")
		if err != nil {
			log.Printf("Error getting source %v", err)
		}
		channel, _, err := unstructured.NestedString(operatorinfo.Object, "status", "defaultChannel")
		if err != nil {
			log.Printf("Error getting channel %v", err)
		}
		fmt.Println("Providing information for app", operator)
		fmt.Println("source: ", source)
		fmt.Println("channel: ", channel)
		fmt.Println("target namespace: ", operator)
		fmt.Println("csv: ", namespace)
		fmt.Println("description: ", namespace)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
