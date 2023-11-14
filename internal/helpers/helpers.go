package helpers

import "github.com/spf13/viper"

// GetString viper get string and check if it is empty
func GetString(key string) (string, error) {
	val := viper.GetString(key)
	if val == "" {
		return "", nil
	}
	return val, nil
}
