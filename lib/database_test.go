package lib

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var test_settings SettingsShare

func TempFilename() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), "foo"+hex.EncodeToString(randBytes)+".sqlite3")
}

func setup() {
	test_settings = NewSettings()
	test_settings.Daemon.DatabaseFilePath = TempFilename()

	err := InitDB(test_settings)
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

func TestOpenDatabase(t *testing.T) {
	db, err := OpenDatabase()
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestInitDatabase(t *testing.T) {
	_, err2 := os.Stat(test_settings.Daemon.DatabaseFilePath)
	assert.Equal(t, false, os.IsNotExist(err2), "The database: %s doesn't exists!", test_settings.Daemon.DatabaseFilePath)
}

func TestStoreRemoveServer(t *testing.T) {
	initial_servers, err := ListServers()
	assert.Nil(t, err)

	server := Server{Path: "MyString", Port: 1234, ListIps: []string{"1", "2"}}
	id, err2 := StoreServer(server)
	assert.Nil(t, err2)
	assert.NotNil(t, id)

	added_servers, err3 := ListServers()
	assert.Nil(t, err3)
	assert.Equal(t, len(initial_servers)+1, len(added_servers), "The server doesn't incremented")

	err4 := RemoveServer(id)
	assert.Nil(t, err4)

	removed_servers, err5 := ListServers()
	assert.Nil(t, err5)
	assert.Equal(t, len(initial_servers), len(removed_servers), "The server doesn't incremented")
}
