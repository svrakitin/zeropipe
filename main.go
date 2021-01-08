package main

import (
	"os"

	"github.com/svrakitin/zeropipe/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
