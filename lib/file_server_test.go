package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCreateHandler(t *testing.T) {
	handler := createHandler()
	assert.NotNil(t, handler)
}
