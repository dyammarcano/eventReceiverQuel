package cmd

import (
	"github.com/dyammarcano/template-go/internal/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallExternalAPI(t *testing.T) {
	cmd := mock.CobraMockCommand()

	viper.Set("url", "https://httpbin.org/get")

	err := CallExternalAPI(cmd, []string{})
	assert.Nilf(t, err, "Error should be nil")
}
