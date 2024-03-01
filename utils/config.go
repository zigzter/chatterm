package utils

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// IsAuthRequired checks whether the token is present in the Viper config.
// If it isn't, we need to authenticate.
func IsAuthRequired() bool {
	token := viper.GetString(TokenKey)
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
			viper.Set(UsernameKey, "")
			viper.Set(WatchedUsersKey, map[string]string{})
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

func WatchUser(username string) string {
	watchedUsers := viper.GetStringMap(WatchedUsersKey)
	var responseMsg string
	if watchedUsers[username] == true {
		delete(watchedUsers, username)
		responseMsg = fmt.Sprintf("Removed %s from watched users", username)
	} else {
		watchedUsers[username] = true
		responseMsg = fmt.Sprintf("Added %s to watched users", username)
	}
	SaveConfig(map[string]interface{}{
		WatchedUsersKey: watchedUsers,
	})
	// We need to refresh the stored config values
	SetFormatterConfigValues()
	return responseMsg
}
