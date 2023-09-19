package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zigzter/chatterm/utils"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set the config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfig()

		var username string
		var oauth string

		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		viper.Set("username", username)

		fmt.Print("Paste oauth token: ")
		fmt.Scanln(&oauth)
		viper.Set("oauth", oauth)

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error saving config:", err)
		} else {
			fmt.Println("Config saved successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
