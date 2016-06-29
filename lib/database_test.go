package lib

import (
	"os"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var testSettings SettingsShare

func GetErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func setup() {
	testSettings = NewSettings()
	testSettings.ShareDaemon.DatabaseFilePath = TempFilename("db_", ".db")

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
	assert.Nil(t, err, GetErrorMessage(err))
	assert.NotNil(t, db)
}

func TestInitDatabase(t *testing.T) {
	_, err2 := os.Stat(testSettings.ShareDaemon.DatabaseFilePath)
	assert.False(t, os.IsNotExist(err2), "The database: %s doesn't exists!", testSettings.ShareDaemon.DatabaseFilePath)
}

func TestStoreRemoveServer(t *testing.T) {
	initial_servers, err := ListServers()
	assert.Nil(t, err, GetErrorMessage(err))

	server := Server{UUID: uuid.NewV4().String(), Path: "MyString", ListIps: []string{"1", "2"}, CreatedAt: time.Now()}
	err = StoreServer(server)
	assert.Nil(t, err, GetErrorMessage(err))

	addedServers, err := ListServers()
	assert.Nil(t, err, GetErrorMessage(err))
	assert.Equal(t, len(initial_servers)+1, len(addedServers), "The servers doesn't incremented")

	err = RemoveServer(server.UUID)
	assert.Nil(t, err, GetErrorMessage(err))

	removedServers, err := ListServers()
	assert.Nil(t, err, GetErrorMessage(err))
	assert.Equal(t, len(initial_servers), len(removedServers), "The servers doesn't incremented")
}
