package main

import (
	"dfss"
	"flag"
	"fmt"
	"runtime"

	dapi "dfss/dfssd/api"
)

var (
	verbose   bool
	demo      string
	fca       string // Path to the CA
	fcert     string // Path to the certificate
	fkey      string // Path to the private key
	addrPort  string // Address and port of the platform
	localPort int    // Port to open for P2P communication
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")
	flag.StringVar(&fca, "ca", "ca.pem", "Path to the root certificate")
	flag.StringVar(&fcert, "cert", "cert.pem", "Path to the user certificate")
	flag.StringVar(&fkey, "key", "key.pem", "Path to the private key")
	flag.StringVar(&addrPort, "host", "localhost:9000", "Host of the DFSS platform")
	flag.StringVar(&demo, "d", "", "Demonstrator address and port (empty string disables debug)")
	flag.IntVar(&localPort, "port", 9005, "Port to use for P2P communication between clients")

	flag.Usage = func() {
		fmt.Println("DFSS client command line v" + dfss.Version)
		fmt.Println("A tool to sign multiparty contracts")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssc [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  help       print this help")
		fmt.Println("  version    print dfss client version")
		fmt.Println("  register   register a new client")
		fmt.Println("  auth       authenticate a new client")
		fmt.Println("  new        create a new contract")
		fmt.Println("  fetch      get a contract hosted on the platform")
		fmt.Println("  show <c>   print contract information from file c")
		fmt.Println("  export <c> export certificate and private key of the user to file c")
		fmt.Println("  import <c> import private key and certificate from file c")
		fmt.Println("  sign <c>   sign contract from file c")

		fmt.Println("\nFlags:")
		flag.PrintDefaults()

		fmt.Println()
	}
}

type command struct {
	nbArgs int
	fn     func([]string)
}

var commands = map[string]command{
	"version": command{0, func([]string) {
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	}},
	"register": command{0, registerUser},
	"auth":     command{0, authUser},
	"new":      command{0, newContract},
	"fetch":    command{0, fetchContract},
	"show":     command{1, showContract},
	"export":   command{1, exportConf},
	"import":   command{1, importConf},
	"sign":     command{1, signContract},
}

func main() {
	flag.Parse()
	arg := flag.Arg(0)
	dapi.Configure(demo != "", demo, "client")

	c, ok := commands[arg]

	if !ok || flag.NArg()-1 < c.nbArgs {
		flag.Usage()
		return
	}

	args := flag.Args()
	args = append(args, "")
	c.fn(args[1:])
}
