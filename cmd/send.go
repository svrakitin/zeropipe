package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cooldown time.Duration

var sendCmd = &cobra.Command{
	Use:   "send [id]",
	Short: "`send` input over the network",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id string
		if len(args) == 0 {
			idUUID, _ := uuid.NewRandom()
			id = idUUID.String()
			fmt.Println(id)
		} else {
			id = args[0]
		}

		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			return err
		}

		w := &multiWriteCloser{}
		defer w.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		service := fmt.Sprintf("%s.%s", id, DefaultServiceRoot)
		entries := make(chan *zeroconf.ServiceEntry)
		err = resolver.Browse(ctx, service, DefaultDomain, entries)
		if err != nil {
			return err
		}

		cooldown := viper.GetDuration("cooldown")
		t := time.NewTimer(cooldown)

		go func() {
			for entry := range entries {
				if !t.Stop() {
					<-t.C
				}
				t.Reset(cooldown)

				addr := fmt.Sprintf("%v:%v", entry.AddrIPv4[0].String(), entry.Port)
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					log.Print(err)
				} else {
					w.Add(conn)
				}
			}
		}()

		select {
		case <-t.C:
			challengeToken := viper.GetString("token")
			if challengeToken != "" {
				w.Write([]byte(challengeToken))
			}
			_, err = io.Copy(w, os.Stdin)
		}

		if err == errNoReceivers {
			return nil
		}

		return err
	},
}

func init() {
	sendCmd.Flags().DurationP("cooldown", "c", 1*time.Second, "receiver discovery cooldown")
	viper.BindPFlag("cooldown", sendCmd.Flags().Lookup("cooldown"))
}

var errNoReceivers = errors.New("no receivers")

type multiWriteCloser struct {
	writers []io.WriteCloser
}

func (mw *multiWriteCloser) Write(p []byte) (n int, err error) {
	if len(mw.writers) == 0 {
		return 0, errNoReceivers
	}
	for i, w := range mw.writers {
		if n, err = w.Write(p); err != nil || n != len(p) {
			mw.writers = append(mw.writers[:i], mw.writers[i+1:]...)
		}
	}
	return len(p), nil
}

func (mw *multiWriteCloser) Add(w io.WriteCloser) {
	mw.writers = append(mw.writers, w)
}

func (mw *multiWriteCloser) Close() error {
	for _, w := range mw.writers {
		w.Close()
	}
	return nil
}
