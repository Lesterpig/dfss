// Package dfssc is the dfss CLI client.
package main

import (
	"os"

	"dfss/dfssc/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
