package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DefaultServiceRoot = "_zeropipe._tcp"
const DefaultDomain = "local."

var rootCmd = &cobra.Command{
	Use:   "zeropipe",
	Short: "pipe over the network using zeroconf mDNS for receiver discovery",
}

func init() {
	viper.SetEnvPrefix("zeropipe")
	viper.AutomaticEnv()

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(recvCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
