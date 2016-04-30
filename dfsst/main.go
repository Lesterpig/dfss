// Package dfsst is the dfss trusted third-party (ttp).
package main

import (
	"os"

	"dfss/dfsst/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
