package cmd

import (
	"fmt"
	"os"

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
	Short: "Print the version number of tmsh-cli command.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tmsh v0.2.0")
	},
}

var execCmd = &cobra.Command{
	Use:   "exec [tmsh command]",
	Short: "Execute any command of TMSH",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		ret, err := bigip.ExecuteCommand(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(ret)
	},
}

func init() {
	cobra.OnInitialize()

	RootCmd.PersistentFlags().StringP("user", "u", "", "TMSH SSH username [$TMSH_USER]")
	RootCmd.PersistentFlags().StringP("password", "p", "", "TMSH SSH passsord [$TMSH_PASSWORD]")
	RootCmd.PersistentFlags().StringP("host", "H", "", "TMSH SSH host [$TMSH_HOST]")
	RootCmd.PersistentFlags().StringP("port", "P", "22", "TMSH SSH port [$TMSH_PORT]")

	flags := RootCmd.PersistentFlags()

	viper.AutomaticEnv()
	viper.BindPFlag("TMSH_USER", flags.Lookup("user"))
	viper.BindPFlag("TMSH_PASSWORD", flags.Lookup("password"))
	viper.BindPFlag("TMSH_HOST", flags.Lookup("host"))
	viper.BindPFlag("TMSH_PORT", flags.Lookup("port"))

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(execCmd)
}
