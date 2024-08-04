package utils

import (
	"fmt"
	"log"
	"os"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
)

// IsAuthRequired checks whether the token is present in the Viper config.
// If it isn't, we need to authenticate.
func IsAuthRequired() bool {
	token := viper.GetString(TokenKey)
	return token == ""
}

func createConfigDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func SetupPath() string {
	scope := gap.NewScope(gap.User, "chatterm")
	dirs, err := scope.ConfigDirs()
	if err != nil {
		log.Fatal(err)
	}
	var configPath string
	if len(dirs) > 0 {
		configPath = dirs[0]
	} else {
		configPath, _ = os.UserHomeDir()
	}
	if err := createConfigDir(configPath); err != nil {
		log.Fatal(err)
	}
	return configPath
}

// InitConfig sets up the config, creating if necessary.
func InitConfig() {
	configPath := SetupPath()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, creating a new one...")
			viper.Set(UsernameKey, "")
			viper.Set(WatchedUsersKey, map[string]string{})
			viper.Set(ShowBadgesKey, true)
			viper.Set(ShowTimestampsKey, true)
			viper.Set(HighlightSubsKey, true)
			viper.Set(HighlightRaidsKey, true)
			viper.Set(FirstTimeChatterColorKey, "#e64553")
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
