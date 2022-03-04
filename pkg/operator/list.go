package operator

import (
	"context"
	"fmt"
	"github.com/karmab/tasty/pkg/utils"
	"github.com/syohex/go-texttable"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sort"
)

func (o *Operator) GetList(installed bool) error {
	var operators []string
	dynamic := utils.GetDynamicClient()
	if installed {
		subscriptionsGVR := schema.GroupVersionResource{
			Group:    "operators.coreos.com",
			Version:  "v1alpha1",
			Resource: "subscriptions",
		}
		list, err := dynamic.Resource(subscriptionsGVR).Namespace("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, d := range list.Items {
			operators = append(operators, d.GetName())
		}
	} else {
		packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
		list, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, d := range list.Items {
			operators = append(operators, d.GetName())
		}
	}
	sort.Strings(operators)
	operatortable := &texttable.TextTable{}
	if err := operatortable.SetHeader("Name"); err != nil {
		fmt.Printf("Failed to set header with %s\n", err)
	}
	for _, operator := range operators {
		if err := operatortable.AddRow(operator); err != nil {
			fmt.Printf("Failed to add a row with %s\n", err)
		}
	}
	fmt.Println(operatortable.Draw())
	return nil
}
