package lib

import (
	"os"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var testSettings SettingsShare

func setup() {
	testSettings = NewSettings()
	testSettings.Daemon.DatabaseFilePath = TempFilename("db_", ".db")

	err := InitDB(testSettings)
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
	_, err2 := os.Stat(testSettings.Daemon.DatabaseFilePath)
	assert.False(t, os.IsNotExist(err2), "The database: %s doesn't exists!", testSettings.Daemon.DatabaseFilePath)
}

func TestStoreRemoveServer(t *testing.T) {
	initial_servers, err := ListServers()
	assert.Nil(t, err)

	server := Server{UUID: uuid.NewV4().String(), Path: "MyString", Port: 1234, ListIps: []string{"1", "2"}, CreatedAt: time.Now()}
	err2 := StoreServer(server)
	assert.Nil(t, err2)

	addedServers, err3 := ListServers()
	assert.Nil(t, err3)
	assert.Equal(t, len(initial_servers)+1, len(addedServers), "The servers doesn't incremented")

	err4 := RemoveServer(server.UUID)
	assert.Nil(t, err4)

	removedServers, err5 := ListServers()
	assert.Nil(t, err5)
	assert.Equal(t, len(initial_servers), len(removedServers), "The servers doesn't incremented")
}
