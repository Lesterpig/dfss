package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dfss/dfssc/common"
	"dfss/dfssp/authority"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ttpCmd = &cobra.Command{
	Use:   "ttp",
	Short: "create and save the TTP's private key and certificate",
	Long: `This command creates a new private key and a related certificate for a Trusted Third Party (TTP).
You must run the init command before this one in order to compute root certificate.

TTP credentials are saved in the ttp folder. You may want to move this folder in another secure place.
In order to provide one TTP per signature, a list must be stored for the platform.
This list is generated and updated with the ttp command, along with the future TTP address and port.
If no list is provided, no TTP will be proposed to signers, and they won't be able to use the resolution protocol.

You can customize the list location and the TTP address with "-t" and "-a" flags.
For example, to setup a TTP that will run on ttp.example.com:3000
	dfssp ttp -t ttp.example.com:3000 -a storefile.data

You can setup as many TTPs as you want, but beware certificate and private key are erased between each call of this command.`,
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("cn", cmd.Flags().Lookup("cn"))
		_ = viper.BindPFlag("validity", cmd.Flags().Lookup("validity"))
		_ = viper.BindPFlag("country", cmd.Flags().Lookup("country"))
		_ = viper.BindPFlag("organization", cmd.Flags().Lookup("org"))
		_ = viper.BindPFlag("unit", cmd.Flags().Lookup("unit"))
		_ = viper.BindPFlag("key_size", cmd.Flags().Lookup("key"))
		_ = viper.BindPFlag("ttps", cmd.Flags().Lookup("ttps"))
		_ = viper.BindPFlag("ttp_addr", cmd.Flags().Lookup("addr"))

		path := viper.GetString("path")

		pid, err := authority.Start(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Bad root CA or key; please use the `init` command before the `ttp` one.\n", err)
			os.Exit(1)
		}
		ttpPath := filepath.Join(path, "ttp")
		v := common.SubViper("key_size", "validity", "country", "organization", "unit", "cn")
		v.Set("path", ttpPath)
		hash, err := authority.Initialize(v, pid.RootCA, pid.Pkey)
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during TTP credentials generation:", err)
			os.Exit(1)
		}

		// Add this ttp to the ttp holder
		holder, err := authority.NewTTPHolder(viper.GetString("ttps"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during TTP list load:", err)
			os.Exit(1)
		}

		holder.Add(viper.GetString("ttp_addr"), hash)
		err = holder.Save(viper.GetString("ttps"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occured during TTP list save:", err)
			os.Exit(1)
		}
	},
}
