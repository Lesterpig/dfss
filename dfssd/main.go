package main

import (
	"dfss"
	"flag"
	"fmt"
	"os"
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
		fmt.Println("  start    start demonstrator server")

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
	case "start":
		err := listen("localhost:3000")
		if err != nil {
			os.Exit(1)
		}
	default:
		flag.Usage()
	}
}
