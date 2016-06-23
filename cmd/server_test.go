package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteServerCmd(t *testing.T) {
	go ServerCmd.Execute()
	assert.Nil(t, nil)
}
