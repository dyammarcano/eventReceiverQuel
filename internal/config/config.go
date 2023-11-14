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
	Config struct {
		Hostname string `json:"-" mapstructure:"-" yaml:"-"`
		Azure    Azure  `json:"azure" mapstructure:"azure" yaml:"azure"`
	}

	Azure struct {
		EventHub EventHub `json:"eventHub" mapstructure:"eventHub" yaml:"eventHub"`
		Storage  Storage  `json:"storage" mapstructure:"storage" yaml:"storage"`
	}

	EventHub struct {
		NameSpace NameSpace `json:"nameSpace" mapstructure:"nameSpace" yaml:"nameSpace"`
		Topics    Topics    `json:"topics" mapstructure:"topics" yaml:"topics"`
	}

	NameSpace struct {
		GroupID          string `json:"groupID" mapstructure:"groupID" yaml:"groupID"`
		ConnectionString string `json:"connectionString" mapstructure:"connectionString" yaml:"connectionString"`
		AccountName      string `json:"accountName" mapstructure:"accountName" yaml:"accountName"`
		SharedAccessKey  string `json:"sharedAccessKey" mapstructure:"sharedAccessKey" yaml:"sharedAccessKey"`
		ConsumerGroup    string `json:"consumerGroup" mapstructure:"consumerGroup" yaml:"consumerGroup"`
	}

	Topics struct {
		ReceivedFile          string `json:"receivedFile" mapstructure:"receivedFile" yaml:"receivedFile"`
		Error                 string `json:"error" mapstructure:"error" yaml:"error"`
		ErrorValidacaoNegocio string `json:"errorValidacaoNegocio" mapstructure:"errorValidacaoNegocio" yaml:"errorValidacaoNegocio"`
	}

	Storage struct {
		AccountKey  string `json:"accountKey" mapstructure:"accountKey" yaml:"accountKey"`
		AccountName string `json:"accountName" mapstructure:"accountName" yaml:"accountName"`
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
	{
		if c.Azure.Storage.AccountName == "" {
			return errors.New("azure storage account [name] is empty")
		}

		if c.Azure.Storage.AccountKey == "" {
			return errors.New("azure storage account [key] is empty")
		}
	}

	{
		if c.Azure.EventHub.NameSpace.AccountName == "" {
			return errors.New("azure storage account [name] is empty")
		}

		if c.Azure.EventHub.NameSpace.SharedAccessKey == "" {
			return errors.New("azure storage account [key] is empty")
		}

		if c.Azure.EventHub.NameSpace.ConnectionString == "" {
			return errors.New("azure storage account [connection string] is empty")
		}
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
