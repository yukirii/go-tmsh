package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use: "tmsh",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of tmsh command.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tmsh v0.1.0")
	},
}

func init() {
	cobra.OnInitialize()

	RootCmdFlags := RootCmd.Flags()
	RootCmdFlags.StringP("user", "u", "", "TMSH SSH username [$TMSH_USER]")
	RootCmdFlags.StringP("password", "p", "", "TMSH SSH passsord [$TMSH_PASSWORD]")
	RootCmdFlags.StringP("host", "H", "", "TMSH SSH host [$TMSH_HOST]")
	RootCmdFlags.StringP("port", "P", "22", "TMSH SSH port [$TMSH_PORT]")

	viper.AutomaticEnv()
	viper.BindPFlag("TMSH_USER", RootCmdFlags.Lookup("user"))
	viper.BindPFlag("TMSH_PASSWORD", RootCmdFlags.Lookup("password"))
	viper.BindPFlag("TMSH_HOST", RootCmdFlags.Lookup("host"))
	viper.BindPFlag("TMSH_PORT", RootCmdFlags.Lookup("port"))

	RootCmd.AddCommand(versionCmd)
}
