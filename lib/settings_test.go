package lib

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TempFilename(prefix string, extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+extension)
}

func TestCreateConfigFile(t *testing.T) {
	var test_settings SettingsShare

	configFile := TempFilename("config_", ".toml")
	err := CreateConfigFile(configFile, test_settings)
	assert.Nil(t, err)

	_, err2 := os.Stat(configFile)
	assert.False(t, os.IsNotExist(err2), "The config file: %s doesn't exists!", configFile)
}

func TestInitSettings(t *testing.T) {
	var test_settings SettingsShare

	configFile := TempFilename("config_", ".toml")
	// Create
	err := InitSettings(configFile, &test_settings)
	assert.Nil(t, err)
	assert.NotNil(t, test_settings)

	// LOAD
	err2 := InitSettings(configFile, &test_settings)
	assert.Nil(t, err2)
	assert.NotNil(t, test_settings)
}

func TestNewSettings(t *testing.T) {
	test_settings := NewSettings()

	assert.NotNil(t, test_settings)
}

func TestConfigFolder(t *testing.T) {
	pathConfigFolder := ConfigFolder()

	assert.NotNil(t, pathConfigFolder)
	assert.True(t, len(pathConfigFolder) > 0)
}

func TestConfigFile(t *testing.T) {
	pathConfigFile := ConfigFile()

	assert.NotNil(t, pathConfigFile)
	assert.True(t, len(pathConfigFile) > 0)
}

func TestConfigFileSQLITE(t *testing.T) {
	pathConfigDb := ConfigFileSQLITE()

	assert.NotNil(t, pathConfigDb)
	assert.True(t, len(pathConfigDb) > 0)
}

func TestConfigServerEndPoint(t *testing.T) {
	test_settings := NewSettings()
	endPoint := ConfigServerEndPoint(test_settings)

	assert.NotNil(t, endPoint)
	assert.True(t, len(endPoint) > 0)
}

func TestUserHomeDir(t *testing.T) {
	pathConfigUserDir := UserHomeDir()

	assert.NotNil(t, pathConfigUserDir)
	assert.True(t, len(pathConfigUserDir) > 0)
}
