package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"dfss"
	"dfss/dfssd/api"
	"dfss/dfssd/gui"
	"dfss/dfssd/server"
	"github.com/visualfc/goqt/ui"
)

var (
	port int
)

func init() {

	flag.IntVar(&port, "p", 3000, "Network port used")

	flag.Usage = func() {
		fmt.Println("DFSS demonstrator v" + dfss.Version)
		fmt.Println("Debug tool to check remote transmissions")

		fmt.Println("\nUsage:")
		fmt.Println("  dfssd [flags] command")

		fmt.Println("\nThe commands are:")
		fmt.Println("  help     print this help")
		fmt.Println("  version  print dfss client version")
		fmt.Println("  nogui    start demonstrator server without GUI")
		fmt.Println("  gui      start demonstrator server with GUI")

		fmt.Println("\nFlags:")
		flag.PrintDefaults()

		fmt.Println()
	}
}

func main() {
	flag.Parse()
	command := flag.Arg(0)
	addrPort := "0.0.0.0:" + strconv.Itoa(port)

	switch command {
	case "help":
		flag.Usage()
	case "version":
		fmt.Println("v"+dfss.Version, runtime.GOOS, runtime.GOARCH)
	case "nogui":
		fmt.Println("Listening on " + addrPort)
		fn := func(v *api.Log) {
			fmt.Printf("[%d] %s: %s\n", v.Timestamp, v.Identifier, v.Log)
		}
		err := server.Listen(addrPort, fn)
		if err != nil {
			os.Exit(1)
		}
	default:
		ui.Run(func() {
			window := gui.NewWindow()
			go func() {
				err := server.Listen(addrPort, window.AddEvent)
				if err != nil {
					window.Log("!! " + err.Error())
				}
			}()
			window.Show()
		})
	}
}
