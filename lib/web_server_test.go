package lib

import (
	"io/ioutil"
	"net/http"
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

func createTmpHandler(t *testing.T) http.Handler {
	testFile, err := createTmpFile()
	assert.Nil(t, err)

	handler := CreateHandler(testFile)
	assert.NotNil(t, handler)

	return handler
}

func TestCreateHandler(t *testing.T) {
	createTmpHandler(t)
}

func TestServerDaemon(t *testing.T) {
	handler := createTmpHandler(t)
	port, err := RandomFreePort()
	assert.Nil(t, err)

	go ServerDaemon(port, handler)
	assert.Nil(t, nil)
}

func TestStartServer(t *testing.T) {
	testFile, err := createTmpFile()
	assert.Nil(t, err)

	port, err := RandomFreePort()
	assert.Nil(t, err)

	server := Server{UUID: uuid.NewV4().String(), Port: port, Path: testFile, CreatedAt: time.Now()}

	go StartServer(&server)
	assert.Nil(t, nil)
}
