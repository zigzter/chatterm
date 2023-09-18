package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/zigzter/chatterm/utils"
)

var (
	Channel  string
	Username string
	Oauth    string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to a Twitch chat",
	Long:  `Connects to a Twitch chat`,
	Run: func(cmd *cobra.Command, args []string) {
		go utils.EstablishWSConnection(Channel, Username, Oauth)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}

func init() {
	connectCmd.Flags().StringVarP(&Channel, "channel", "c", "", "The Twitch channel to join")
	connectCmd.Flags().StringVarP(&Username, "username", "u", "", "Your username on Twitch")
	connectCmd.Flags().StringVarP(&Oauth, "oauth", "o", "", "The Oath string, in format oauth:xyz123")
	connectCmd.MarkFlagRequired("channel")
	connectCmd.MarkFlagRequired("username")
	connectCmd.MarkFlagRequired("oauth")
	connectCmd.Println(Channel, Username, Oauth)

	rootCmd.AddCommand(connectCmd)
}
