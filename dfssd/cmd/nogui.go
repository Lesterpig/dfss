package cmd

import (
	"dfss/dfssd/api"
	"dfss/dfssd/server"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var noguiCmd = &cobra.Command{
	Use:   "nogui",
	Short: "start demonstrator without GUI",
	Run: func(cmd *cobra.Command, args []string) {
		addrPort := "0.0.0.0:" + strconv.Itoa(viper.GetInt("port"))

		fmt.Println("Listening on " + addrPort)
		fn := func(v *api.Log) {
			fmt.Printf("[%d] %s: %s\n", v.Timestamp, v.Identifier, v.Log)
		}
		err := server.Listen(addrPort, fn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
