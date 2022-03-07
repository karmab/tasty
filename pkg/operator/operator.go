package operator

type Operator struct {
	Name           string
	Source         string
	DefaultChannel string
	Description    string
	Csv            string
	Namespace      string
	Channels       []string
	Crd            string
	ConfigExecFile string
	ConfigExecPath string
}

func NewOperator() *Operator {
	return &Operator{
		Name:           "",
		Source:         "",
		DefaultChannel: "",
		Description:    "",
		Csv:            "",
		Namespace:      "",
		Channels:       []string{},
		Crd:            "",
		ConfigExecFile: "",
		ConfigExecPath: "",
	}
}

func NewOperatorWithOptions(name, source, defaultChannel, description, csv, namespace, crd, configExecFile, configExecPath string) *Operator {
	return &Operator{
		Name:           name,
		Source:         source,
		DefaultChannel: defaultChannel,
		Description:    description,
		Csv:            csv,
		Namespace:      namespace,
		Channels:       []string{},
		Crd:            crd,
		ConfigExecFile: configExecFile,
		ConfigExecPath: configExecPath,
	}
}

var OperatorTemplate = `{{ if ne .Namespace "openshift-operators" -}}
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Namespace }}
  labels:
    openshift.io/cluster-monitoring: "true"
  annotations:
    workload.openshift.io/allowed: management
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: {{ .Name }}-operatorgroup
  namespace: {{ .Namespace }}
spec:
  targetNamespaces:
  - {{ .Namespace }}
---
{{ end -}}
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  channel: "{{ .DefaultChannel }}"
  name: {{ .Name }}
  source: {{ .Source }}
  sourceNamespace: openshift-marketplace
  startingCSV: {{ .Csv }}
  installPlanApproval: Automatic
`

func GetOperatorTemplate() string {
	return OperatorTemplate
}
