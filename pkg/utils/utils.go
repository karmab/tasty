package utils

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type Operator struct {
	Name           string
	Source         string
	DefaultChannel string
	Csv            string
	Namespace      string
}

var OperatorTemplate = `{{ if ne .Namespace "openshift-operators" }}
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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetClient() (client dynamic.Interface) {
	kubeconfig, _ := os.LookupEnv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Check(err)
	client, err = dynamic.NewForConfig(config)
	Check(err)
	return client
}

func WaitCrd(crd string, timeout int) {
	client := GetClient()
	i := 0
	for i < timeout {
		crds := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
		list, err := client.Resource(crds).Namespace("").List(context.TODO(), metav1.ListOptions{})
		Check(err)
		for _, d := range list.Items {
			if d.GetName() == crd {
				return
			}
		}
		color.Cyan("Waiting for CRD %s to be created\n", crd)
		time.Sleep(5 * time.Second)
		i += 5
	}
	color.Red("Timeout waiting ffor CRD %s\n", crd)
}

func GetOperator(operator string) (source string, defaultchannel string, csv string, description string, target_namespace string, channels []string, crd string) {
	target_namespace = strings.Split(operator, "-operator")[0]
	kubeconfig, _ := os.LookupEnv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Check(err)
	client, err := dynamic.NewForConfig(config)
	Check(err)
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	operatorinfo, err := client.Resource(packagemanifests).Namespace("openshift-marketplace").Get(context.TODO(), operator, metav1.GetOptions{})
	Check(err)
	source, _, err = unstructured.NestedString(operatorinfo.Object, "status", "catalogSource")
	Check(err)
	defaultchannel, _, err = unstructured.NestedString(operatorinfo.Object, "status", "defaultChannel")
	Check(err)
	allchannels, _, err := unstructured.NestedSlice(operatorinfo.Object, "status", "channels")
	Check(err)
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
				}
			}
		}
	}
	return source, defaultchannel, csv, description, target_namespace, channels, crd
}
