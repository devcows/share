package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/devcows/share/lib"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testSettings lib.SettingsShare
	configFile   string
)

func TempFilename(prefix string, extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+extension)
}

func setup() {
	configFile = TempFilename("config_", ".toml")
	testSettings = lib.NewSettings()
	testSettings.Daemon.DatabaseFilePath = TempFilename("db_", ".db")

	err := lib.CreateConfigFile(configFile, testSettings)
	if err != nil {
		panic(err)
	}
}

func shutdown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestExecuteServerCmd(t *testing.T) {
	err := lib.InitSettings(configFile, &testSettings)
	assert.Nil(t, err, GetErrorMessage(err))

	err = lib.InitDB(testSettings)
	assert.Nil(t, err, GetErrorMessage(err))

	server := lib.Server{UUID: uuid.NewV4().String(), Path: "MyString", Port: 1234, ListIps: []string{"1", "2"}, CreatedAt: time.Now()}
	err = lib.StoreServer(server)
	assert.Nil(t, err, GetErrorMessage(err))

	err = runServerCmd(configFile, &testSettings)
	assert.Nil(t, err, GetErrorMessage(err))
}
