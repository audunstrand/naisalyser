package main

import (
	"os"

	"github.com/navikt/naisalyser/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
