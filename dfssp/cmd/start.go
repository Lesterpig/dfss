package cmd

import (
	"fmt"
	"os"

	dapi "dfss/dfssd/api"
	"dfss/dfssp/server"
	"dfss/net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start [path, db, address, port, cert-validity]",
	Short: "start the platform after loading its private key and root certificate",
	Run: func(cmd *cobra.Command, args []string) {
		address := viper.GetString("address")
		port := viper.GetString("port")

		srv := server.GetServer()

		fmt.Println("Listening on " + address + ":" + port)
		dapi.DLog("Platform server started on " + address + ":" + port)
		err := net.Listen(address+":"+port, srv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}
