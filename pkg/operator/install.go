package operator

import (
	"context"
	"errors"
	"html/template"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/karmab/tasty/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Operator) InstallOperator(wait, out bool, ns, channel, csv, src, srcNS, installPlan string, args []string) error {
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
		o.Source = src
		o.SourceNS = srcNS
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
				SourceNS:       o.SourceNS,
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
				namespace := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: ns,
						Annotations: map[string]string{
							"workload.openshift.io/allowed": "management",
						},
					},
				}

				_, err := k8sclient.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
				if !apierrors.IsAlreadyExists(err) {
					utils.Check(err)
				} else {
					color.Yellow("Namespace %s already exists, continuing...", namespace.Name)
				}

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
					if !apierrors.IsAlreadyExists(err) {
						log.Printf("Error creating operator group %s: %s", operator+"-operatorgroup", err)
						return err
					} else {
						color.Yellow("OperatorGroup %s already exists, continuing...", operatorgroupspec.GetName())
					}
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
						"sourceNamespace":     o.SourceNS,
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
