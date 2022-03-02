package operator

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetList(t *testing.T) {
	t.Log("Test list Operator found")
	operatorOK := NewOperator()
	err := operatorOK.GetList(false)
	assert.NoError(t, err)

	t.Log("Test list  Operator only installed")
	operatorEmpty := NewOperator()
	err = operatorEmpty.GetList(true)
	assert.NoError(t, err)

}
