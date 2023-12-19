package operator

import (
	"context"
	"fmt"
	"sort"

	"github.com/karmab/tasty/pkg/utils"
	"github.com/syohex/go-texttable"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Operator) GetList(installed bool) error {
	operators := make(map[string]*Operator)

	dynamic := utils.GetDynamicClient()
	if installed {
		olmClient := utils.GetOlmClient()
		subscriptions, err := olmClient.OperatorsV1alpha1().Subscriptions("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing subscriptions: %v\n", err)
			return err
		}

		for _, s := range subscriptions.Items {
			op := NewOperator()
			op.Name = s.Spec.Package
			op.Namespace = s.Namespace
			op.InstalledChannel = s.Spec.Channel
			op.InstalledSource = s.Spec.CatalogSource
			op.InstalledSourceNS = s.Spec.CatalogSourceNamespace
			op.InstalledCsv = s.Spec.StartingCSV
			operators[op.Name+"."+op.Namespace] = op
		}
	}
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	packages, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if !installed {
		for _, d := range packages.Items {
			op := NewOperator()
			err = op.ParseOperator(&d)
			if err != nil {
				fmt.Printf("Failed parsing operator manifest")
				continue
			}
			if _, ok := operators[op.Name]; ok {
				operators[op.Name].SupportedInstallModes = op.SupportedInstallModes
				operators[op.Name].SuggestedNamespace = op.SuggestedNamespace
			} else {
				operators[op.Name] = op
			}
		}
	}
	sortedOperators := []*Operator{}
	for _, op := range operators {
		sortedOperators = append(sortedOperators, op)
	}

	sort.Slice(sortedOperators, func(i, j int) bool {
		return sortedOperators[i].Name < sortedOperators[j].Name
	})

	operatortable := &texttable.TextTable{}
	if installed {
		if err := operatortable.SetHeader("Name", "Namespace", "Source", "Source Namespace", "Channel", "Csv"); err != nil {
			fmt.Printf("Failed to set header with %s\n", err)
		}
	} else {
		if err := operatortable.SetHeader("Name", "Suggested Namespace", "Install Modes"); err != nil {
			fmt.Printf("Failed to set header with %s\n", err)
		}
	}

	for _, operator := range sortedOperators {
		if installed {
			if err := operatortable.AddRow(operator.Name, operator.Namespace, operator.InstalledSource, operator.InstalledSourceNS, operator.InstalledChannel, operator.InstalledCsv); err != nil {
				fmt.Printf("Failed to add a row with %s\n", err)
			}
		} else {
			if err := operatortable.AddRow(operator.Name, operator.SuggestedNamespace, operator.GetInstallModes()); err != nil {
				fmt.Printf("Failed to add a row with %s\n", err)
			}
		}
	}
	fmt.Println(operatortable.Draw())
	return nil
}
