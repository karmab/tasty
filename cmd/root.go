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
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func get_operator(operator string) (namespace string, source string, defaultchannel string, csv string, description string, target_namespace string, crd string) {
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
	namespace, _, err = unstructured.NestedString(operatorinfo.Object, "metadata", "namespace")
	if err != nil {
		log.Printf("Error getting namespace %v", err)
	}
	source, _, err = unstructured.NestedString(operatorinfo.Object, "status", "catalogSource")
	if err != nil {
		log.Printf("Error getting source %v", err)
	}
	defaultchannel, _, err = unstructured.NestedString(operatorinfo.Object, "status", "defaultChannel")
	if err != nil {
		log.Printf("Error getting channel %v", err)
	}
	channels, _, err := unstructured.NestedSlice(operatorinfo.Object, "status", "channels")
	if err != nil {
		log.Printf("Error getting channel %v", err)
	}
	for _, channel := range channels {
		channelmap, _ := channel.(map[string]interface{})
		channelname := channelmap["name"]
		if channelname == defaultchannel {
			csv = channelmap["currentCSV"].(string)
			csvdescmap, _ := channelmap["currentCSVDesc"].(map[string]interface{})
			description = csvdescmap["description"].(string)
			installmodes := csvdescmap["installModes"].([]interface{})
			for _, mode := range installmodes {
				modemap, _ := mode.(map[string]interface{})
				if modemap["type"] == "OwnNamespace" && modemap["supported"] == false {
					target_namespace = "openshift-operators"
				}
			}
			csvdescannotations := csvdescmap["annotations"].(map[string]interface{})
			if suggested_namespace, ok := csvdescannotations["operatorframework.io/suggested-namespace"].(string); ok {
				target_namespace = suggested_namespace
			}
			if customresourcedefinitionsmap, ok := csvdescmap["customresourcedefinitions"]; ok {
				customresourcedefinitions, _ := customresourcedefinitionsmap.(map[string]interface{})
				ownedlist := customresourcedefinitions["owned"].([]interface{})
				owned := ownedlist[0].(map[string]interface{})
				crd = owned["name"].(string)
			}
		}
	}

	return namespace, source, defaultchannel, csv, description, target_namespace, crd
}

var cfgFile string
var kubeconfig string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tasty",
	Short: "Handles OLM operators",
	Long: `This application allows you to interact with olm operators
using a yum-like workflow`,

	//Run: func(cmd *cobra.Command, args []string) { fmt.Println("Coucou c est Karim") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	kubeconfig = os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		color.Red("KUBECONFIG env variable needs to be set")
		os.Exit(1)
	}
}
