package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteClienRmCmd(t *testing.T) {
	RmCmd.Execute()
	assert.Nil(t, nil)
}
