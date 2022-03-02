package operator

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestInfo(t *testing.T) {
	t.Log("TODO: Implement TestInfo")
	prueba := Operator{
		Name:           "",
		Source:         "",
		DefaultChannel: "",
		Description:    "",
		Csv:            "",
		Namespace:      "",
		Channels:       nil,
		Crd:            "",
		ConfigExecFile: "",
		ConfigExecPath: "",
	}
	err := prueba.GetInfo([]string{""})
	assert.Error(t, err)

	prueba2 := Operator{
		Name:           "",
		Source:         "",
		DefaultChannel: "",
		Description:    "",
		Csv:            "",
		Namespace:      "",
		Channels:       nil,
		Crd:            "",
		ConfigExecFile: "",
		ConfigExecPath: "",
	}
	prueba2.GetInfo([]string{""})

}
