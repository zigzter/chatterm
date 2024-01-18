package utils

import (
	"log"

	"github.com/spf13/viper"
)

// IsAuthRequired checks whether the token is present in the Viper config.
// If it isn't, we need to authenticate.
func IsAuthRequired() bool {
	token := viper.GetString("token")
	return token == ""
}

// InitConfig sets up the config, creating if necessary.
func InitConfig() {
	viper.SetConfigName("chatterm")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, creating a new one...")
			viper.Set("username", "")
			if err := viper.SafeWriteConfig(); err != nil {
				log.Println("Error creating config file:", err)
			}
		} else {
			log.Println("Error reading config file:", err)
		}
	}
}

func SaveConfig(options map[string]interface{}) {
	for key, value := range options {
		viper.Set(key, value)
	}
	if err := viper.WriteConfig(); err != nil {
		log.Println("Error saving config:", err)
	}
}
