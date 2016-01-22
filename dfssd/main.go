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
		fmt.Println("DFSS demonstrator v" + dfss.Version)
		fmt.Println("Debug tool to check remote transmissions")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssd [flags] command")

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
