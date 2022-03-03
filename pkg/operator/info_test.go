package operator

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestInfo(t *testing.T) {
	t.Log("Test info Operator found")
	operatorOK := NewOperator()
	err := operatorOK.GetInfo([]string{"web-terminal"})
	assert.NoError(t, err)

	t.Log("Test info without Operator")
	operatorEmpty := NewOperator()
	err = operatorEmpty.GetInfo([]string{""})
	assert.Error(t, err)

	t.Log("Test info without params")
	operatorWithOutParam := NewOperator()
	err = operatorWithOutParam.GetInfo([]string{})
	assert.Error(t, err)

	t.Log("Test info Operator not found")
	operatorNotFound := NewOperator()
	err = operatorNotFound.GetInfo([]string{"xxx"})
	assert.Error(t, err)

}
