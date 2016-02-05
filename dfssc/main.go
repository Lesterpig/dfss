package main

import (
	"dfss"
	"flag"
	"fmt"
	osuser "os/user"
	"runtime"
	"time"

	"dfss/dfssc/user"
)

var (
	verbose  bool
	fca      string // Path to the CA
	fcert    string // Path to the certificate
	fkey     string // Path to the private key
	addrPort string // Address and port of the platform
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")
	flag.StringVar(&fca, "ca", "ca.pem", "Path to the root certificate")
	flag.StringVar(&fcert, "cert", "cert.pem", "Path to the user certificate")
	flag.StringVar(&fkey, "key", "priv_key.pem", "Path to the private key")
	flag.StringVar(&addrPort, "host", "127.0.0.1:9000", "Host of the DFSS platform (e.g. 127.0.0.1:9000)")

	flag.Usage = func() {
		fmt.Println("DFSS client command line v" + dfss.Version)
		fmt.Println("A tool to sign multiparty contracts")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssc [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  help      print this help")
		fmt.Println("  version   print dfss client version")
		fmt.Println("  register  register a new client")
		fmt.Println("  show <c>  print contract information from file c")

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
		fmt.Println("Registering a new user")
		// Initialize variables
		var country, mail, organization, unit, passphrase string
		var bits int

		u, err := osuser.Current()
		if err != nil {
			fmt.Println("An error occurred : ", err.Error())
			break
		}

		// Get all the necessary parameters
		readStringParam("Mail", "", &mail)
		readStringParam("Country", time.Now().Location().String(), &country)
		readStringParam("Organization", u.Name, &organization)
		readStringParam("Organizational unit", u.Name, &unit)
		readIntParam("Length of the key (2048 or 4096)", 2048, &bits)
		err = readPassword(&passphrase)
		if err != nil {
			fmt.Println("An error occurred : ", err.Error())
			break
		}

		recapUser(fca, fcert, fkey, addrPort, country, organization, unit, mail, bits)

		err = user.Register(fca, fcert, fkey, addrPort, passphrase, country, organization, unit, mail, bits)
		if err != nil {
			fmt.Println("An error occurred : ", err.Error())
		}

	case "show":
		showContract(flag.Arg(1))
	default:
		flag.Usage()
	}
}
