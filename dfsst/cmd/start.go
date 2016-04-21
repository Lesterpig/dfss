package cmd

import (
	"fmt"
	"os"
	"strconv"

	dapi "dfss/dfssd/api"
	"dfss/dfsst/server"
	"dfss/net"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the TTP service",
	Long: `Start the TTP of the DFSS project

Fill the DFSS_TTP_PASSWORD environment variable if the private key is enciphered`,
	Run: func(cmd *cobra.Command, args []string) {
		demo := viper.GetString("demo")
		dapi.Configure(demo != "", demo, "ttp")

		srv := server.GetServer()

		addrPort := viper.GetString("address") + ":" + strconv.Itoa(viper.GetInt("port"))
		fmt.Println("Listening on " + addrPort)
		dapi.DLog("TTP server started on " + addrPort)
		err := net.Listen(addrPort, srv)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		dapi.DClose()
	},
}
