package operator

import (
	"context"
	"errors"
	"html/template"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/karmab/tasty/pkg/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Operator) InstallOperator(wait bool, out bool, ns string, channel string, csv string, installPlan string, args []string) error {
	for _, operator := range args {

		err := o.GetOperator(operator)
		if err != nil {
			log.Printf("Error getting operator %s: %s", operator, err)
			return err
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
		if ns == "" {
			ns = o.Namespace
		}
		if installPlan == "" {
			installPlan = "Automatic"
		} else if installPlan != "Manual" && installPlan != "Automatic" {
			color.Red("Invalid installplan %s", installPlan)
			return errors.New("invalid installplan")
		}
		if out {
			t := template.New("Template")
			tpl, err := t.Parse(GetOperatorTemplate())
			if err != nil {
				log.Printf("Error parsing template: %s", err)
				return err
			}
			operatordata := Operator{
				Name:           operator,
				Source:         o.Source,
				DefaultChannel: o.DefaultChannel,
				Csv:            o.Csv,
				Namespace:      ns,
			}
			err = tpl.Execute(os.Stdout, operatordata)
			if err != nil {
				log.Printf("Error executing template: %s", err)
				return err
			}
		} else {
			color.Cyan("Installing operator %s", operator)
			dynamic := utils.GetDynamicClient()
			if ns != "openshift-operators" {
				color.Cyan("Creating namespace %s", ns)
				k8sclient := utils.GetK8sClient()
				nsmetaSpec := metav1.ObjectMeta{Name: ns, Annotations: map[string]string{"workload.openshift.io/allowed": "management"}}
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
							"namespace": ns,
						},
						"spec": map[string]interface{}{
							"targetNamespaces": []string{ns},
						},
					},
				}
				_, err = dynamic.Resource(operatorgroupsGVR).Namespace(ns).Create(context.TODO(), operatorgroupspec, metav1.CreateOptions{})
				if err != nil {
					log.Printf("Error creating operator group %s: %s", operator+"-operatorgroup", err)
					return err
				}

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
						"namespace": ns,
					},
					"spec": map[string]interface{}{
						"channel":             o.DefaultChannel,
						"name":                operator,
						"source":              o.Source,
						"sourceNamespace":     "openshift-marketplace",
						"startingCSV":         o.Csv,
						"installPlanApproval": installPlan,
					},
				},
			}
			_, err := dynamic.Resource(subscriptionsGVR).Namespace(ns).Create(context.TODO(), subspec, metav1.CreateOptions{})
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
