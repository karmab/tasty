package operator

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"strings"
	"tasty/pkg/utils"
)

func (o *Operator) SearchOperator(args []string) error {
	var currentoperator string

	if len(args) != 1 {
		log.Printf("Usage: tasty search OPERATOR_NAME")
		return errors.New("Invalid number of arguments. Usage: tasty search OPERATOR_NAME")
	}

	dynamic := utils.GetDynamicClient()
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	list, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
	utils.Check(err)

	for _, d := range list.Items {
		currentoperator = d.GetName()
		if strings.Contains(currentoperator, args[0]) {
			fmt.Println(currentoperator)
		}
	}
	return nil
}
