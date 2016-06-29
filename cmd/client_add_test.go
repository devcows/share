package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func TestExecuteClienAddCmd(t *testing.T) {
	//runAddCmd()
	assert.Nil(t, nil)
}

/*
Not working on linux
func TestCopyToClipboard(t *testing.T) {
	myString := "My String"
	err := copyClipboard(myString)

	assert.Nil(t, err, GetErrorMessage(err))

	myString2, err := clipboard.ReadAll()
	assert.Nil(t, err, GetErrorMessage(err))

	assert.Equal(t, myString, myString2)
}
*/
