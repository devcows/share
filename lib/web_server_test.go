package lib

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func createTmpFile() (string, error) {
	testFile := TempFilename("test_file", ".txt")

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile(testFile, d1, 0644)
	return testFile, err
}

func createTmpFileHandler(t *testing.T) http.Handler {
	testFile, err := createTmpFile()
	assert.Nil(t, err, GetErrorMessage(err))

	handler := CreateHandler(testFile)
	assert.NotNil(t, handler)

	return handler
}

func createTmpFolderHandler(t *testing.T) http.Handler {
	testFile, err := createTmpFile()
	assert.Nil(t, err, GetErrorMessage(err))

	handler := CreateHandler(filepath.Dir(testFile))
	assert.NotNil(t, handler)

	return handler
}

func TestCreateFileHandler(t *testing.T) {
	createTmpFileHandler(t)
}

func TestCreateFolderHandler(t *testing.T) {
	createTmpFolderHandler(t)
}

func TestServerDaemon(t *testing.T) {
	handler := createTmpFileHandler(t)
	port, err := RandomFreePort()
	assert.Nil(t, err, GetErrorMessage(err))

	go ServerDaemon(port, handler)
	assert.Nil(t, nil)
}

func TestStartServer(t *testing.T) {
	testFile, err := createTmpFile()
	assert.Nil(t, err, GetErrorMessage(err))

	port, err := RandomFreePort()
	assert.Nil(t, err, GetErrorMessage(err))

	server := Server{UUID: uuid.NewV4().String(), Port: port, Path: testFile, CreatedAt: time.Now()}

	StartServer(&server)
	assert.Nil(t, nil)
}
