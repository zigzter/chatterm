package utils

import (
	"log"

	"github.com/spf13/viper"
)

// Sets up the config, creating if necessary.
// Returns a boolean, true if authentication is required
func InitConfig() (requiresAuth bool) {
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
			return true
		} else {
			log.Println("Error reading config file:", err)
		}
	}
	return false
}

func SaveConfig(options map[string]interface{}) {
	for key, value := range options {
		viper.Set(key, value)
	}
	if err := viper.WriteConfig(); err != nil {
		log.Println("Error saving config:", err)
	}
}
