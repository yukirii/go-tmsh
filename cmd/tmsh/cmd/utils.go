package cmd

import (
	"fmt"
	"os"

	"github.com/shiftky/go-tmsh"
	"github.com/spf13/viper"
)

func NewSession() *tmsh.BigIP {
	bigip, err := tmsh.NewSession(
		viper.GetString("tmsh_host"),
		viper.GetString("tmsh_port"),
		viper.GetString("tmsh_user"),
		viper.GetString("tmsh_password"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return bigip
}
