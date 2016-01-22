package main

import (
	"dfss"
	"flag"
	"fmt"
	"runtime"
)

var (
	verbose bool
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")

	flag.Usage = func() {
		fmt.Println("DFSS client command line v" + dfss.Version)
		fmt.Println("A tool to sign multiparty contracts")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssc command [flags]")

		fmt.Println("\nThe commands are:")
		fmt.Println("  help     print this help")
		fmt.Println("  version  print dfss client version")

		fmt.Println("\nFlags:")
		flag.PrintDefaults()

		fmt.Println()
	}
}

func main() {
	flag.Parse()
	command := flag.Arg(0)

	switch command {
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	default:
		flag.Usage()
	}
}