package operator

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func Test_Search(t *testing.T) {
	t.Log("Test search Operator found")
	operatorOK := NewOperator()
	err := operatorOK.SearchOperator([]string{"web-terminal"})
	assert.NoError(t, err)

	t.Log("Test search without Operator")
	operatorEmpty := NewOperator()
	err = operatorEmpty.SearchOperator([]string{""})
	assert.NoError(t, err)

	t.Log("Test search without params")
	operatorWithOutParam := NewOperator()
	err = operatorWithOutParam.SearchOperator([]string{})
	assert.Error(t, err)

	t.Log("Test search Operator not found")
	operatorNotFound := NewOperator()
	err = operatorNotFound.SearchOperator([]string{"xxx"})
	assert.NoError(t, err)

}
