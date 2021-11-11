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
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func get_client() (client dynamic.Interface) {
	kubeconfig, _ := os.LookupEnv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	check(err)
	client, err = dynamic.NewForConfig(config)
	check(err)
	return client
}

func wait_crd(crd string, version string, timeout int) {
	crdsplit := strings.SplitN(crd, ".", 2)
	resource := crdsplit[0]
	group := crdsplit[1]
	client := get_client()
	crdGVR := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	i := 0
	for i < timeout {
		_, err := client.Resource(crdGVR).Namespace("default").List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			fmt.Printf("CRD %s ready\n", crd)
			break
		} else {
			fmt.Printf("Waiting for CRD %s to be created\n", crd)
			time.Sleep(5 * time.Second)
			i += 5
		}
	}
}

func get_operator(operator string) (source string, defaultchannel string, csv string, description string, target_namespace string, channels []string, crd string, crdversion string) {
	kubeconfig, _ := os.LookupEnv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	check(err)
	client, err := dynamic.NewForConfig(config)
	check(err)
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	operatorinfo, err := client.Resource(packagemanifests).Namespace("openshift-marketplace").Get(context.TODO(), operator, metav1.GetOptions{})
	check(err)
	source, _, err = unstructured.NestedString(operatorinfo.Object, "status", "catalogSource")
	check(err)
	defaultchannel, _, err = unstructured.NestedString(operatorinfo.Object, "status", "defaultChannel")
	check(err)
	allchannels, _, err := unstructured.NestedSlice(operatorinfo.Object, "status", "channels")
	check(err)
	for _, channel := range allchannels {
		channelmap, _ := channel.(map[string]interface{})
		channelname := channelmap["name"]
		channels = append(channels, channelname.(string))
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
				ownedlistmap := customresourcedefinitions["owned"]
				if ownedlistmap == nil {
					crd = ""
				} else {
					ownedlist := ownedlistmap.([]interface{})
					owned := ownedlist[0].(map[string]interface{})
					crd = owned["name"].(string)
					crdversion = owned["version"].(string)
				}
			}
		}
	}
	return source, defaultchannel, csv, description, target_namespace, channels, crd, crdversion
}

type Operator struct {
	Name           string
	Source         string
	DefaultChannel string
	Csv            string
	Namespace      string
}

var operatordata = `{{ if ne .Namespace "openshift-operators" }}
apiVersion: v1
kind: Namespace
metadata:
  labels:
    openshift.io/cluster-monitoring: "true"
  name: {{ .Namespace }}
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: {{ .Name }}-operatorgroup
  namespace: {{ .Namespace }}
spec:
  targetNamespaces:
  - {{ .Namespace }}
---
{{ end }}
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  channel: "{{ .DefaultChannel }}"
  name: {{ .Name }}
  source: {{ .Source }}
  sourceNamespace: openshift-marketplace
`

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
