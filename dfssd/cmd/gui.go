package cmd

import (
	"strconv"

	"dfss/dfssd/gui"
	"dfss/dfssd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/visualfc/goqt/ui"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "start the demonstrator with a gui",
	Run: func(cmd *cobra.Command, args []string) {
		addrPort := "0.0.0.0:" + strconv.Itoa(viper.GetInt("port"))
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
	},
}
