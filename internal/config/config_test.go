//go:build windows

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInstance(t *testing.T) {
	cfg := NewConfig()

	err := cfg.LoadConfigFile("C:\\arqprod_local\\cfg\\config.yaml")
	assert.Nil(t, err)

	err = cfg.Validate() //this case is used to verify oracle, azure and kafka values
	assert.Nil(t, err)
}
