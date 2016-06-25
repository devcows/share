package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteVersionCmd(t *testing.T) {
	runVersionCmd()
	assert.Nil(t, nil)
}
