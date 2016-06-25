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

var testSettings lib.SettingsShare

func TempFilename(prefix string, extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+extension)
}

func setup() {
	testSettings = lib.NewSettings()
	testSettings.Daemon.DatabaseFilePath = TempFilename("db_", ".db")

	err := lib.InitDB(testSettings)
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
	configFile := TempFilename("config_", ".toml")

	err := lib.InitSettings(configFile, &testSettings)
	assert.Nil(t, err)

	err = lib.InitDB(testSettings)
	assert.Nil(t, err)

	server := lib.Server{UUID: uuid.NewV4().String(), Path: "MyString", Port: 1234, ListIps: []string{"1", "2"}, CreatedAt: time.Now()}
	err = lib.StoreServer(server)
	assert.Nil(t, err)

	err = runServerCmd(configFile, &testSettings)
	assert.Nil(t, err)
}
