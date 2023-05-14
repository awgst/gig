package pkg

import (
	"github.com/spf13/viper"
)

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
