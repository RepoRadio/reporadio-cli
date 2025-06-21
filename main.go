package main

import (
	"fmt"
	"os"

	"github.com/reporadio/reporadio-cli/internal"
)

func main() {
	if err := internal.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
