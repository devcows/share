package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteClienPsCmd(t *testing.T) {
	PsCmd.Execute()
	assert.Nil(t, nil)
}
