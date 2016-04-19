package cmd

import (
	"dfss"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print dfss protocol version",
	Long:  "Print dfss protocol version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	},
}
