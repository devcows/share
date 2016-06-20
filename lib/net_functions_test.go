package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFreePort(t *testing.T) {
	port := RandomFreePort()
	assert.Equal(t, true, port > 0, "Retrieve random free port")
}
