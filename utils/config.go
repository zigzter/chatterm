package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("chatterm")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, creating a new one...")
			viper.Set("username", "")
			viper.Set("oauth", "")
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println("Error creating config file:", err)
			}
		} else {
			fmt.Println("Error reading config file:", err)
		}
	}
}
