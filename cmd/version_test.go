package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteVersionCmd(t *testing.T) {
	err := VersionCmd.Execute()
	assert.Nil(t, err)
}
