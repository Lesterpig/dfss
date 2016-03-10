package main

import (
	"dfss"
	"flag"
	"fmt"
	"runtime"
)

var (
	verbose   bool
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
		fmt.Println("  show <c>   print contract information from file c")
		fmt.Println("  export <c> export certificate and private key of the user to file c")
		fmt.Println("  import <c> import private key and certificate from file c")
		fmt.Println("  sign <c>   sign contract from file c")

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
	case "register":
		registerUser()
	case "auth":
		authUser()
	case "new":
		newContract()
	case "show":
		showContract(flag.Arg(1))
	case "export":
		exportConf(flag.Arg(1))
	case "import":
		importConf(flag.Arg(1))
	case "sign":
		signContract(flag.Arg(1))
	default:
		flag.Usage()
	}
}
