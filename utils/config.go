package utils

import (
	"log"
	"strings"

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

// StoreUserState processes the GLOBALUSERSTATE IRC message,
// storing any relevant data to Viper
func StoreUserState(input string) {
	parts := strings.SplitN(input, ">", 2)
	metadata := parts[1]
	keyValPairs := strings.Split(metadata, ";")
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]
			switch key {
			case "user-id":
				viper.Set("userId", value)
			case "color":
				viper.Set("color", value)
			}
		}
	}
	if err := viper.WriteConfig(); err != nil {
		log.Println("Error saving config:", err)
	}
}

// StoreRoomState processes the ROOMSTATE IRC message,
// storing any relevant data to Viper
func StoreRoomState(input string) {
	parts := strings.SplitN(input, "\n", 2)
	metadata := parts[1]
	keyValPairs := strings.Split(metadata, ";")
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]
			switch key {
			case "room-id":
				viper.Set("channelid", value)
			}
		}
	}
	if err := viper.WriteConfig(); err != nil {
		log.Println("Error saving config:", err)
	}
}
