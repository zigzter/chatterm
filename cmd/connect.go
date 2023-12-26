package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/zigzter/chatterm/utils"
)

var Channel string

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to a Twitch chat",
	Long:  `Connects to a Twitch chat`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitConfig()
		// username := viper.GetString("username")
		// oauth := viper.GetString("oauth")
		// go utils.EstablishWSConnection(Channel, username, oauth)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}

func init() {
	connectCmd.Flags().StringVarP(&Channel, "channel", "c", "", "The Twitch channel to join")
	connectCmd.MarkFlagRequired("channel")
	connectCmd.Println(Channel)

	rootCmd.AddCommand(connectCmd)
}
