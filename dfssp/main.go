package main

import (
	"dfss"
	"dfss/dfssp/autority"
	"flag"
	"fmt"
	"runtime"
)

var (
	verbose                      bool
	path, country, org, unit, cn string
	keySize, validity            int
	pid                          *autority.PlatformID
)

func init() {

	flag.BoolVar(&verbose, "v", false, "Print verbose messages")

	flag.StringVar(&path, "path", autority.GetHomeDir(), "Path for the platform's private key and root certificate")
	flag.StringVar(&country, "country", "France", "Country for the root certificate")
	flag.StringVar(&org, "org", "DFSS", "Organization for the root certificate")
	flag.StringVar(&unit, "unit", "INSA Rennes", "Organizational unit for the root certificate")
	flag.StringVar(&cn, "cn", "dfssp", "Common name for the root certificate")

	flag.IntVar(&keySize, "keySize", 512, "Encoding size for the private key")
	flag.IntVar(&validity, "validity", 21, "Root certificate's validity duration (days)")

	flag.Usage = func() {
		fmt.Println("DFSS platform v" + dfss.Version)
		fmt.Println("Users and contracts manager")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssp command [flags]")

		fmt.Println("\nThe commands are:")
		fmt.Println("  init     [cn, country, keySize, org, path, unit, validity]")
		fmt.Println("           create and save the platform's private key and root certificate")
		fmt.Println("  start    [path]")
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

	switch command {
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	case "init":
		err := autority.Initialize(keySize, validity, country, org, unit, cn, path)
		if err != nil {
			fmt.Println("An error occured during the initialization operation")
			fmt.Println(err)
			panic(err)
		}
	case "start":
		pid, err := autority.Start(path)
		if err != nil {
			fmt.Println("An error occured during the start operation")
			fmt.Println(err)
			panic(err)
		}
		// TODO: use pid
		_ = pid
	default:
		flag.Usage()
	}
}
