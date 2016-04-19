package cmd

import (
	"dfss"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version of dfss protocol",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DFSS v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	},
}
