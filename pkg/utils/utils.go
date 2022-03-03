package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Check(e error) {
	if e != nil {
		color.Red("%s", e)
		os.Exit(1)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetK8sClient() (client *kubernetes.Clientset) {
	var config *rest.Config
	var err error
	homedir := os.Getenv("HOME")
	kubeport, _ := os.LookupEnv("KUBERNETES_PORT")
	if kubeport != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			color.Red("Couldn't load in-cluster config")
			os.Exit(1)
		}
	} else {
		kubeconfigenv, _ := os.LookupEnv("KUBECONFIG")
		kubeconfig := strings.Replace(kubeconfigenv, "~", homedir, 1)
		if kubeconfig == "" {
			kubeconfig = filepath.Join(homedir, ".kube", "config")
		}
		_, err = os.Stat(kubeconfig)
		if err != nil {
			color.Red("KUBECONFIG file %s not found", kubeconfig)
			os.Exit(1)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		Check(err)
	}
	client, err = kubernetes.NewForConfig(config)
	Check(err)
	return client
}

func GetDynamicClient() (client dynamic.Interface) {
	var config *rest.Config
	var err error
	homedir := os.Getenv("HOME")
	kubeport, _ := os.LookupEnv("KUBERNETES_PORT")
	if kubeport != "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			color.Red("Couldn't load in-cluster config")
			os.Exit(1)
		}
	} else {
		kubeconfigenv, _ := os.LookupEnv("KUBECONFIG")
		kubeconfig := strings.Replace(kubeconfigenv, "~", homedir, 1)
		if kubeconfig == "" {
			kubeconfig = filepath.Join(homedir, ".kube", "config")
		}
		_, err = os.Stat(kubeconfig)
		if err != nil {
			color.Red("KUBECONFIG file %s not found", kubeconfig)
			os.Exit(1)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		Check(err)
	}
	client, err = dynamic.NewForConfig(config)
	Check(err)
	return client
}

func WaitCrd(crd string, timeout int) {
	d := GetDynamicClient()
	i := 0
	for i < timeout {
		crds := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
		list, err := d.Resource(crds).Namespace("").List(context.TODO(), metav1.ListOptions{})
		Check(err)
		for _, d := range list.Items {
			if d.GetName() == crd {
				return
			}
		}
		color.Cyan("Waiting for CRD %s to be created\n", crd)
		time.Sleep(5 * time.Second)
		i += 5
	}
	color.Red("Timeout waiting for CRD %s\n", crd)
}
