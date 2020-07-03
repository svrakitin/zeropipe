package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DefaultServiceRoot = "_zeropipe._tcp"
const DefaultDomain = "local"

var rootCmd = &cobra.Command{
	Use:   "zeropipe",
	Short: "pipe over the network using zeroconf mDNS for receiver discovery",
}

func init() {
	viper.SetEnvPrefix("zeropipe")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringP("token", "t", "", "challenge token")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(recvCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
