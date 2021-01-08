package cmd

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"github.com/spf13/cobra"
)

var recvCmd = &cobra.Command{
	Use:   "recv <id>",
	Short: "`recv` output over the network",
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

		_, err = io.Copy(os.Stdout, conn)
		return err
	},
}
