package main

import (
	"dfss"
	dapi "dfss/dfssd/api"
	"flag"
	"fmt"
	"runtime"
)

var (
	verbose, demo        bool
	port, address, dbURI string
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")
	flag.BoolVar(&demo, "d", false, "Enable demonstrator")

	flag.StringVar(&port, "p", "9010", "Default port listening")
	flag.StringVar(&address, "a", "0.0.0.0", "Default address to bind for listening")

	flag.StringVar(&dbURI, "db", "mongodb://localhost/dfss", "Name of the environment variable containing the server url in standard MongoDB format")

	flag.Usage = func() {
		fmt.Println("DFSS TTP v" + dfss.Version)
		fmt.Println("Trusted third party resolver")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssp [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  start    [db, a, p]")
		fmt.Println("           start the TTP service")
		fmt.Println("  help     print this help")
		fmt.Println("  version  print dfss protocol version")

		fmt.Println("\nFlags:")
		flag.PrintDefaults()

		fmt.Println()
	}
}

func main() {
	flag.Parse()
	command := flag.Arg(0)
	dapi.Switch(demo)

	switch command {
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	case "start":
		// srv := server.GetServer(dbURI, verbose)
		// fmt.Println("Listening on " + address + ":" + port)
		// dapi.DLog("TTP server started on " + address + ":" + port)
		// err := net.Listen(address+":"+port, srv)
		// if err != nil {
		//	  fmt.Println(err)
		// }
	default:
		flag.Usage()
	}

	dapi.DClose()
}
