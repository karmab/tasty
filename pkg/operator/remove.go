package operator

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"tasty/pkg/utils"
)

func (o *Operator) RemoveOperator(remove bool, args []string) error {

	if remove != true {
		var confirmation string
		color.Green("Are you sure? [y/N]:")
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			color.Red("Leaving..")
			return errors.New("Leaving..")
		}
	}
	subscriptionsGVR := schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}
	operatorgroupsGVR := schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1",
		Resource: "operatorgroups",
	}
	for _, operator := range args {
		color.Cyan("Removing operator %s", operator)
		err := o.GetOperator(operator)
		if err != nil {
			return err
		}
		d := utils.GetDynamicClient()
		color.Cyan("Removing subscription %s", operator)
		err = d.Resource(subscriptionsGVR).Namespace(o.Namespace).Delete(context.TODO(), operator, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		if o.Namespace != "openshift-operators" {
			color.Cyan("Removing operator group %s-operatorgroup", operator)
			k8sclient := utils.GetK8sClient()
			err := d.Resource(operatorgroupsGVR).Namespace(o.Namespace).Delete(context.TODO(), operator+"-operatorgroup", metav1.DeleteOptions{})
			if err != nil {
				return err
			}
			color.Cyan("Removing namespace group %s", o.Namespace)
			err = k8sclient.CoreV1().Namespaces().Delete(context.TODO(), o.Namespace, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
