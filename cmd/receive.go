package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errChallengeFailed = errors.New("token challenge failed")

var recvCmd = &cobra.Command{
	Use:   "recv <id>",
	Short: "`recv` input from the network into standard output",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		defer ln.Close()

		instance, _ := uuid.NewRandom()
		service := fmt.Sprintf("%s.%s", id, DefaultServiceRoot)
		port := ln.Addr().(*net.TCPAddr).Port
		server, err := zeroconf.Register(instance.String(), service, DefaultDomain, port, nil, nil)
		if err != nil {
			return err
		}
		defer server.Shutdown()

		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		challengeToken := viper.GetString("token")
		if challengeToken != "" {
			tokenBytes := []byte(challengeToken)
			recvToken := make([]byte, len(tokenBytes))
			n, err := conn.Read(recvToken)
			if err != nil {
				return err
			}

			if n != len(tokenBytes) || bytes.Compare(tokenBytes, recvToken) != 0 {
				return errChallengeFailed
			}
		}

		_, err = io.Copy(os.Stdout, conn)

		return err
	},
}
