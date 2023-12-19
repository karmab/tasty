package operator

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/karmab/tasty/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (o *Operator) GetInfo(args []string) error {
	if len(args) != 1 {
		color.Set(color.FgRed)
		return errors.New("invalid number of arguments. Usage: tasty info OPERATOR_NAME")
	}

	err := o.GetOperator(args[0])
	if err != nil {
		color.Set(color.FgRed)
		return err
	}

	color.Cyan("Providing information on operator %s", args[0])
	fmt.Println("source: ", o.Source)
	fmt.Println("channels: ", o.Channels)
	fmt.Println("defaultchannel: ", o.DefaultChannel)
	fmt.Printf("supported install modes: %v\n", o.SupportedInstallModes)
	fmt.Println("suggested namespace: ", o.Namespace)
	fmt.Println("csv: ", o.Csv)
	fmt.Println("description: ", o.Description)
	return nil
}

func (o *Operator) GetOperator(operator string) error {
	dynamic := utils.GetDynamicClient()
	packagemanifests := schema.GroupVersionResource{Group: "packages.operators.coreos.com", Version: "v1", Resource: "packagemanifests"}
	operatorinfo, err := dynamic.Resource(packagemanifests).Namespace("openshift-marketplace").Get(context.TODO(), operator, metav1.GetOptions{})
	if err != nil {
		color.Set(color.FgRed)
		return err
	}

	return o.ParseOperator(operatorinfo)
}

func (o *Operator) ParseOperator(operatorinfo *unstructured.Unstructured) error {
	var err error
	o.Name = operatorinfo.GetName()
	o.Source, _, err = unstructured.NestedString(operatorinfo.Object, "status", "catalogSource")
	if err != nil {
		color.Set(color.FgRed)
		return err
	}

	o.DefaultChannel, _, err = unstructured.NestedString(operatorinfo.Object, "status", "defaultChannel")
	if err != nil {
		color.Set(color.FgRed)
		return err
	}

	allchannels, _, err := unstructured.NestedSlice(operatorinfo.Object, "status", "channels")
	if err != nil {
		color.Set(color.FgRed)
		return err
	}

	for _, channel := range allchannels {
		channelmap, _ := channel.(map[string]interface{})
		channelname := channelmap["name"]
		o.Channels = append(o.Channels, channelname.(string))
		if channelname == o.DefaultChannel {
			o.Csv = channelmap["currentCSV"].(string)
			csvdescmap, _ := channelmap["currentCSVDesc"].(map[string]interface{})
			o.Description = csvdescmap["description"].(string)
			installmodes := csvdescmap["installModes"].([]interface{})
			for _, mode := range installmodes {
				modemap, _ := mode.(map[string]interface{})
				if modemap["supported"] == true {
					if str, ok := modemap["type"].(string); ok {
						if channelString, ok := channelname.(string); ok {
							o.SupportedInstallModes[channelString] = append(o.SupportedInstallModes[channelString], str)
						}
					}
				}
			}
			if _, ok := csvdescmap["annotations"]; ok {
				csvdescannotations := csvdescmap["annotations"].(map[string]interface{})
				if suggestedNamespace, ok := csvdescannotations["operatorframework.io/suggested-namespace"].(string); ok {
					o.SuggestedNamespace = suggestedNamespace
				}
			}

			if customresourcedefinitionsmap, ok := csvdescmap["customresourcedefinitions"]; ok {
				customresourcedefinitions, _ := customresourcedefinitionsmap.(map[string]interface{})
				ownedlistmap := customresourcedefinitions["owned"]
				if ownedlistmap != nil {
					ownedlist := ownedlistmap.([]interface{})
					owned := ownedlist[0].(map[string]interface{})
					o.Crd = owned["name"].(string)
				}
			}
		}
	}
	return nil
}
