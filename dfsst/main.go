package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"dfss"
	dapi "dfss/dfssd/api"
	"dfss/dfsst/server"
	"dfss/net"
)

var (
	verbose bool
	fca     string // Path to the CA
	fcert   string // Path to the certificate
	fkey    string // Path to the private key
	address string
	dbURI   string
	port    string
	demo    string
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")
	flag.StringVar(&fca, "ca", "ca.pem", "Path to the root certificate")
	flag.StringVar(&fcert, "cert", "cert.pem", "Path to the user certificate")
	flag.StringVar(&fkey, "key", "key.pem", "Path to the private key")
	flag.StringVar(&demo, "d", "", "Demonstrator address and port (empty string disables debug)")

	flag.StringVar(&port, "p", "9020", "Default port listening")
	flag.StringVar(&address, "a", "0.0.0.0", "Default address to bind for listening")

	flag.StringVar(&dbURI, "db", "mongodb://localhost/dfss", "Name of the environment variable containing the server url in standard MongoDB format")

	flag.Usage = func() {
		fmt.Println("DFSS TTP v" + dfss.Version)
		fmt.Println("Trusted third party resolver")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssp [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  start    start the TTP service")
		fmt.Println("           fill the `DFSS_TTP_PASSWORD` environment variable if the private key is enciphered")
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
	dapi.Configure(demo != "", demo, "ttp")

	switch command {
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	case "start":
		password := os.Getenv("DFSS_TTP_PASSWORD")
		srv := server.GetServer(fca, fcert, fkey, password, dbURI, verbose)

		addrPort := address + ":" + port
		fmt.Println("Listening on " + addrPort)
		dapi.DLog("TTP server started on " + addrPort)
		err := net.Listen(addrPort, srv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		flag.Usage()
	}

	dapi.DClose()
}
