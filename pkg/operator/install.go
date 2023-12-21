package operator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/karmab/tasty/pkg/utils"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Operator) InstallOperator(wait, out bool, targetNSs, ns, channel, csv, src, srcNS, installPlan string, args []string) error {
	const (
		nsDelimiter = ","
		allNSString = "*"
	)

	// convert comma separated list of target namespaces to a slice of target namespace
	targetNamespaces := strings.Split(targetNSs, ",")

	// if any namespaces in the list is "*", assume all namespaces install mode. All other namespaces are discarded
	IsAllNamespacesMode := StringInSlice(targetNamespaces, allNSString)

	for _, operator := range args {
		err := o.GetOperator(operator)
		if err != nil {
			log.Printf("Error getting operator %s: %s", operator, err)
			return err
		}
		o.Namespace = ns
		if ns == "" {
			if o.SuggestedNamespace == "" {
				// if not suggested namespace and not namespaces, create one
				o.Namespace = "openshift-" + strings.Split(operator, "-operator")[0]
			} else {
				// by default use the suggested namespace
				o.Namespace = o.SuggestedNamespace
			}
		}

		if channel != "" {
			if utils.Contains(o.Channels, channel) {
				o.DefaultChannel = channel
			} else {
				color.Red("Target channel %s not found in %s", channel, o.Channels)
				return errors.New("target channel not found")
			}
		}
		if csv != "" {
			o.Csv = csv
		}
		if installPlan == "" {
			installPlan = "Automatic"
		} else if installPlan != "Manual" && installPlan != "Automatic" {
			color.Red("Invalid installplan %s", installPlan)
			return errors.New("invalid installplan")
		}
		o.Source = src
		o.SourceNS = srcNS
		yamlOutput := ""
		if !out {
			color.Cyan("Installing operator %s", operator)
		}
		dynamic := utils.GetDynamicClient()
		if o.Namespace != "openshift-operators" {
			if !out {
				color.Cyan("Creating namespace %s", o.Namespace)
			}

			namespaceGVR := schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "namespaces",
			}

			namespacespec := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Namespace",
					"apiVersion": "v1",
					"metadata": map[string]interface{}{
						"name": o.Namespace,
						"annotations": map[string]string{
							"workload.openshift.io/allowed": "management",
						},
					},
				},
			}

			if out {
				yamlOutput += fmt.Sprintf("---\n%s", printObjectYAML(namespacespec.Object))
			} else {
				_, err = dynamic.Resource(namespaceGVR).Create(context.TODO(), namespacespec, metav1.CreateOptions{})
				if err != nil {
					if !apierrors.IsAlreadyExists(err) {
						log.Printf("Error creating namespace %s: %s", o.Namespace, err)
						return err
					} else {
						color.Yellow("Namespace %s already exists, continuing...", namespacespec.GetName())
					}
				}
			}

			if !out {
				color.Cyan("Creating operator group %s-operatorgroup", operator)
			}
			operatorgroupsGVR := schema.GroupVersionResource{
				Group:    "operators.coreos.com",
				Version:  "v1",
				Resource: "operatorgroups",
			}

			// If all namespaces install mode is detected, configure an empty spec
			aSpec := map[string]interface{}{}

			// If one or more target namespaces are passed, use them as target namespaces
			if !IsAllNamespacesMode {
				aSpec = map[string]interface{}{
					"targetNamespaces": targetNamespaces,
				}
			}

			operatorgroupspec := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "OperatorGroup",
					"apiVersion": "operators.coreos.com/v1",
					"metadata": map[string]interface{}{
						"name":      operator + "-operatorgroup",
						"namespace": o.Namespace,
					},
					"spec": aSpec,
				},
			}
			if out {
				yamlOutput += fmt.Sprintf("---\n%s", printObjectYAML(operatorgroupspec.Object))
			} else {
				_, err = dynamic.Resource(operatorgroupsGVR).Namespace(o.Namespace).Create(context.TODO(), operatorgroupspec, metav1.CreateOptions{})
				if err != nil {
					if !apierrors.IsAlreadyExists(err) {
						log.Printf("Error creating operator group %s: %s", operator+"-operatorgroup", err)
						return err
					} else {
						color.Yellow("OperatorGroup %s already exists, continuing...", operatorgroupspec.GetName())
					}
				}
			}

		}
		if !out {
			color.Cyan("Creating subscription %s", operator)
		}
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
					"namespace": o.Namespace,
				},
				"spec": map[string]interface{}{
					"channel":             o.DefaultChannel,
					"name":                operator,
					"source":              o.Source,
					"sourceNamespace":     o.SourceNS,
					"startingCSV":         o.Csv,
					"installPlanApproval": installPlan,
				},
			},
		}
		if out {
			yamlOutput += fmt.Sprintf("---\n%s", printObjectYAML(subspec.Object))
			fmt.Print(yamlOutput)
		} else {
			_, err = dynamic.Resource(subscriptionsGVR).Namespace(o.Namespace).Create(context.TODO(), subspec, metav1.CreateOptions{})
			if err != nil {
				log.Printf("Error creating subscription %s: %s", operator, err)
				return err
			}
			if wait && o.Crd != "" {
				utils.WaitCrd(o.Crd, 60)
			}
		}
	}
	return nil
}

// StringInSlice checks a slice for a given string.
func StringInSlice(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(strings.TrimSpace(string(v)), string(str)) {
			return true
		}
	}
	return false
}

// Prints the yaml for the k8s object
func printObjectYAML(obj map[string]interface{}) string {
	yamlData, err := yaml.Marshal(obj)
	if err != nil {
		log.Printf("Error marshaling YAML: %s", err)
		return ""
	}
	return string(yamlData)
}
