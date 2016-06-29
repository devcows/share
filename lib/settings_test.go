package lib

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TempFilename(prefix string, extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+extension)
}

func TestCreateConfigFile(t *testing.T) {
	var testSettings SettingsShare

	configFile := TempFilename("config_", ".toml")
	err := CreateConfigFile(configFile, testSettings)
	assert.Nil(t, err, GetErrorMessage(err))

	_, err2 := os.Stat(configFile)
	assert.False(t, os.IsNotExist(err2), "The config file: %s doesn't exists!", configFile)
}

func TestInitSettings(t *testing.T) {
	var testSettings SettingsShare

	testSettings = NewSettings()
	testSettings.ShareDaemon.DatabaseFilePath = TempFilename("db_", ".db")

	configFile := TempFilename("config_", ".toml")
	err := CreateConfigFile(configFile, testSettings)
	assert.Nil(t, err, GetErrorMessage(err))

	// Create
	err = InitSettings(configFile, &testSettings)
	assert.Nil(t, err, GetErrorMessage(err))
	assert.NotNil(t, testSettings)

	// LOAD
	err = InitSettings(configFile, &testSettings)
	assert.Nil(t, err, GetErrorMessage(err))
	assert.NotNil(t, testSettings)
}

func TestNewSettings(t *testing.T) {
	testSettings := NewSettings()

	assert.NotNil(t, testSettings)
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
	testSettings := NewSettings()
	endPoint := ConfigServerEndPoint(testSettings)

	assert.NotNil(t, endPoint)
	assert.True(t, len(endPoint) > 0)
}

func TestUserHomeDir(t *testing.T) {
	listOs := []string{"linux", "windows", "darwin"}

	for i := 0; i < len(listOs); i++ {
		pathConfigUserDir := UserHomeDir(listOs[i])

		assert.NotNil(t, pathConfigUserDir)
		if listOs[i] == runtime.GOOS {
			assert.True(t, len(pathConfigUserDir) > 0)
		}
	}
}
