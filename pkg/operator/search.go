package operator

import (
	"context"
	"errors"
	"fmt"
	"github.com/karmab/tasty/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

func (o *Operator) SearchOperator(args []string) error {
	var currentoperator string

	if len(args) != 1 {
		return errors.New("invalid number of arguments. Usage: tasty search OPERATOR_NAME")
	}

	dynamic := utils.GetDynamicClient()
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	list, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range list.Items {
		currentoperator = d.GetName()
		if strings.Contains(currentoperator, args[0]) {
			fmt.Println(currentoperator)
		}
	}
	return nil
}
