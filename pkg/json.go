// Package pkg implements list function and variable that can be used by other packages
package pkg

import (
	"github.com/spf13/viper"
)

// ReadJsonString reads a string from a JSON file
// Accepts key as string
// Returns string, error
func ReadJsonString(key string) (string, error) {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("gig")

	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}

	return viper.GetString(key), nil
}

// ReadJsonBool reads a bool from a JSON file
// Accepts key as string
// Returns bool, error
func ReadJsonBool(key string) (bool, error) {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("gig")

	err := viper.ReadInConfig()
	if err != nil {
		return false, err
	}

	return viper.GetBool(key), nil
}
