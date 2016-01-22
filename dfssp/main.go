package main

import (
	"dfss"
	dapi "dfss/dfssd/api"
	"dfss/dfssp/authority"
	"dfss/dfssp/server"
	"dfss/net"
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	verbose, demo                                      bool
	path, country, org, unit, cn, port, address, dbURI string
	keySize, validity                                  int
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")
	flag.BoolVar(&demo, "d", false, "Enable demonstrator")

	flag.StringVar(&port, "p", "9000", "Default port listening")
	flag.StringVar(&address, "a", "0.0.0.0", "Default address to bind for listening")

	flag.StringVar(&path, "path", authority.GetHomeDir(), "Path for the platform's private key and root certificate")
	flag.StringVar(&country, "country", "France", "Country for the root certificate")
	flag.StringVar(&org, "org", "DFSS", "Organization for the root certificate")
	flag.StringVar(&unit, "unit", "INSA Rennes", "Organizational unit for the root certificate")
	flag.StringVar(&cn, "cn", "dfssp", "Common name for the root certificate")

	flag.IntVar(&keySize, "keySize", 512, "Encoding size for the private key")
	flag.IntVar(&validity, "validity", 21, "Root certificate's validity duration (days)")

	flag.StringVar(&dbURI, "db", "mongodb://localhost/dfss", "Name of the environment variable containing the server url in standard MongoDB format")

	flag.Usage = func() {
		fmt.Println("DFSS platform v" + dfss.Version)
		fmt.Println("Users and contracts manager")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssp [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  init     [cn, country, keySize, org, path, unit, validity]")
		fmt.Println("           create and save the platform's private key and root certificate")
		fmt.Println("  start    [path, db, a, p]")
		fmt.Println("           start the platform after loading its private key and root certificate")
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
	dapi.Switch(demo)

	switch command {
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	case "init":
		err := authority.Initialize(keySize, validity, country, org, unit, cn, path)
		if err != nil {
			fmt.Println("An error occured during the initialization operation:", err)
			os.Exit(1)
		}
		dapi.DLog("Private key generated !")
	case "start":
		srv := server.GetServer(path, dbURI, verbose)
		fmt.Println("Listening on " + address + ":" + port)
		dapi.DLog("Platform server started on " + address + ":" + port)
		err := net.Listen(address+":"+port, srv)
		if err != nil {
			fmt.Println(err)
		}
	default:
		flag.Usage()
	}

	dapi.DClose()
}
