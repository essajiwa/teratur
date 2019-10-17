package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	configPath = filepath.Join(filepath.Join(os.Getenv("GOPATH"), "src/github.com/essajiwa/teratur"), "files/etc/teratur/teratur.development.yaml")
)

func TestMain(t *testing.T) {
	testInit(t)
	testGet(t)
	testError(t)
}

func testInit(t *testing.T) {
	// Test with no error
	err := Init(WithConfigFile(configPath))
	assert.NoError(t, err)
}

func testGet(t *testing.T) {
	assert.NotNil(t, Get())
}

func testError(t *testing.T) {
	// Set env as non dev and should return error
	os.Setenv("ENV", "production")
	err := Init()
	assert.Error(t, err)
}
