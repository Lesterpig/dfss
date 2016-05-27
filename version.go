// Package dfss is the root of the dfss architecture.
package dfss

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version represents the current version of the DFSS software suite
const Version = "0.3.1"

// VersionCmd is the cobra command common to all dfss modules
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version of dfss protocol",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DFSS v"+Version, runtime.GOOS, runtime.GOARCH)
	},
}
