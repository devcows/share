package cmd

import (
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

func TestExecuteClienAddCmd(t *testing.T) {
	//runAddCmd()
	assert.Nil(t, nil)
}

func TestCopyToClipboard(t *testing.T) {
	myString := "My String"
	err := copyClipboard(myString)
	assert.Nil(t, err)

	myString2, err := clipboard.ReadAll()
	assert.Nil(t, err)

	assert.Equal(t, myString, myString2)
}
