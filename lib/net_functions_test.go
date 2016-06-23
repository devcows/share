package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFreePort(t *testing.T) {
	port, err := RandomFreePort()
	assert.Nil(t, err)
	assert.Equal(t, true, port > 0, "Retrieve random free port")
}

func TestOpenUpnpPort(t *testing.T) {
	opened := OpenUpnpPort(5000)
	assert.NotNil(t, opened, "Return true or false open upnp port")
}

func TestGetLocalIps(t *testing.T) {
	listIps := GetLocalIps(5000)
	assert.NotNil(t, listIps, "GetLocalIps return array list")
	assert.True(t, len(listIps) > 0, "GetLocalIps empty array")
}

func TestGetPublicIps(t *testing.T) {
	listIps := GetPublicIps(5000)
	assert.NotNil(t, listIps, "GetPublicIps return array list")
	assert.True(t, len(listIps) > 0, "GetPublicIps empty array")
}
