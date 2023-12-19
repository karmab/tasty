package operator

type Operator struct {
	Name                  string
	Source                string
	SourceNS              string
	DefaultChannel        string
	Description           string
	Csv                   string
	Namespace             string
	TargetNamespaces      []string
	SuggestedNamespace    string
	Channels              []string
	Crd                   string
	ConfigExecFile        string
	ConfigExecPath        string
	SupportedInstallModes map[string][]string

	InstalledChannel  string
	InstalledSource   string
	InstalledSourceNS string
	InstalledCsv      string
}

func (op *Operator) GetInstallModes() (out string) {
	for channelName, channel := range op.SupportedInstallModes {
		out += channelName + "[ "
		for _, mode := range channel {
			out += mode + " "
		}
		out += "] "
	}
	return out
}

func NewOperator() (op *Operator) {

	op = &Operator{
		Name:               "",
		Source:             "",
		SourceNS:           "",
		DefaultChannel:     "",
		Description:        "",
		Csv:                "",
		Namespace:          "",
		SuggestedNamespace: "",
		TargetNamespaces:   []string{},
		Channels:           []string{},
		Crd:                "",
		ConfigExecFile:     "",
		ConfigExecPath:     "",
	}
	op.SupportedInstallModes = make(map[string][]string)
	return op
}

func NewOperatorWithOptions(name, source, sourceNS, defaultChannel, description, csv, namespace, crd, configExecFile, configExecPath string) *Operator {
	return &Operator{
		Name:           name,
		Source:         source,
		SourceNS:       sourceNS,
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
