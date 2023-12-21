package operator

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/karmab/tasty/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (o *Operator) RemoveOperator(ns string, remove, rmns, rmgroup bool, args []string) error {
	if remove {
		var confirmation string
		color.Green("Are you sure? [y/N]:")
		if _, err := fmt.Scanln(&confirmation); err != nil {
			fmt.Printf("Failed to scan confirmation with %s\n", err)
		}
		if strings.ToLower(confirmation) != "y" {
			color.Red("Leaving..")
			return errors.New("leaving")
		}
	}
	for _, operator := range args {
		color.Cyan("Removing operator %s", operator)
		err := o.GetOperator(operator)
		if err != nil {
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

		olmClient := utils.GetOlmClient()
		subscriptions, err := olmClient.OperatorsV1alpha1().Subscriptions(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing subscriptions: %v\n", err)
			return err
		}

		for _, s := range subscriptions.Items {
			if s.Spec.Package == operator {
				err := olmClient.OperatorsV1alpha1().Subscriptions(ns).Delete(context.TODO(), s.Name, metav1.DeleteOptions{})
				if err != nil {
					return err
				}
			}
		}

		if o.Namespace != "openshift-operators" {
			color.Cyan("Removing all operator groups")

			if rmgroup {
				operatorGroups, err := olmClient.OperatorsV1().OperatorGroups(ns).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					fmt.Printf("Error listing OperatorGroups: %v\n", err)
					return err
				}

				fmt.Println("Deleting OperatorGroups:")
				for _, operatorGroup := range operatorGroups.Items {
					color.Cyan("Removing %s", operatorGroup.GetName())
					err := olmClient.OperatorsV1().OperatorGroups(ns).Delete(context.TODO(), operatorGroup.GetName(), metav1.DeleteOptions{})
					if err != nil {
						return err
					}
				}
			}
			if rmns {
				color.Cyan("Removing namespace group %s", o.Namespace)
				k8sclient := utils.GetK8sClient()
				err = k8sclient.CoreV1().Namespaces().Delete(context.TODO(), o.Namespace, metav1.DeleteOptions{})
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}
