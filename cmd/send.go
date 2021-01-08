package cmd

import (
	"context"
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
			idUUID := uuid.Must(uuid.NewRandom())
			id = idUUID.String()
			fmt.Println(id)
		} else {
			id = args[0]
		}

		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			return err
		}

		w := &multiWriteCloser{
			writers: make(map[string]io.WriteCloser),
		}
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

				addr := fmt.Sprintf("%s:%d", entry.AddrIPv4[0].String(), entry.Port)
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					log.Printf("dial: %s", err)
				} else {
					log.Printf("dial: %s", addr)
					w.Add(addr, conn)
				}
			}
		}()

		select {
		case <-t.C:
			_, err = io.Copy(w, os.Stdin)
		}

		return err
	},
}

func init() {
	sendCmd.Flags().DurationP("cooldown", "c", 1*time.Second, "receiver discovery cooldown")
	viper.BindPFlag("cooldown", sendCmd.Flags().Lookup("cooldown"))
}

type multiWriteCloser struct {
	writers map[string]io.WriteCloser
}

func (mw *multiWriteCloser) Write(p []byte) (n int, err error) {
	if len(mw.writers) == 0 {
		return 0, io.ErrShortWrite
	}
	for key, w := range mw.writers {
		if n, err = w.Write(p); err != nil {
			delete(mw.writers, key)
		}
	}
	return len(p), nil
}

func (mw *multiWriteCloser) Add(key string, w io.WriteCloser) {
	mw.writers[key] = w
}

func (mw *multiWriteCloser) Close() error {
	for _, w := range mw.writers {
		w.Close()
	}
	return nil
}
