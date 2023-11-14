package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	errEmptyPath      = "config file path cannot be empty"
	errFileNotFound   = "config file not found in: %s"
	errUnsupportedExt = "config file extension not supported: %s"
	errConfigRead     = "error reading config file: %v"
	errDecodeStruct   = "unable to decode into struct, %v"
)

type (
	//azure:
	//	eventHub:
	//	eventHub: eventHub
	//	accountName: accountName
	//	accountKey: accountKey
	//	topic: topic
	//	consumerGroup: consumerGroup

	Config struct {
		Hostname string `json:"-" mapstructure:"-" yaml:"-"`
		EventHub struct {
			AccountName   string `json:"accountName" mapstructure:"accountName" yaml:"accountName"`
			AccountKey    string `json:"accountKey" mapstructure:"accountKey" yaml:"accountKey"`
			Topic         string `json:"topic" mapstructure:"topic" yaml:"topic"`
			ConsumerGroup string `json:"consumerGroup" mapstructure:"consumerGroup" yaml:"consumerGroup"`
		} `json:"eventHub" mapstructure:"eventHub" yaml:"eventHub"`
	}
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfigFile(filePath string) error {
	if filePath == "" {
		return errors.New(errEmptyPath)
	}

	fp := filepath.Clean(filePath)

	if !fileExists(fp) {
		return fmt.Errorf(errFileNotFound, fp)
	}

	if !c.isSupportedFileExtension(fp) {
		return fmt.Errorf(errUnsupportedExt, fp)
	}

	viper.SetConfigFile(fp)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf(errConfigRead, err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf(errDecodeStruct, err)
	}

	// set default values if not defined
	if err := c.defaultGlobalValue(); err != nil {
		return err
	}

	return nil
}

func (c *Config) isSupportedFileExtension(filePath string) bool {
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")

	for _, supportedExt := range viper.SupportedExts {
		if ext == supportedExt {
			return true
		}
	}

	return false
}

func (c *Config) Validate() error {
	if c.EventHub.AccountName == "" {
		return errors.New("event hub account [name] is empty")
	}

	if c.EventHub.AccountKey == "" {
		return errors.New("event hub account [key] is empty")
	}

	if c.EventHub.Topic == "" {
		return errors.New("event hub [topic] is empty")
	}

	return nil
}

func (c *Config) defaultGlobalValue() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if hostname[len(hostname)-1] == '.' {
		hostname = hostname[:len(hostname)-1]
	}

	c.Hostname = hostname

	return nil
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}
