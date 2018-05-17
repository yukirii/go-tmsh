package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shiftky/go-tmsh"
	"github.com/spf13/viper"
)

func NewSession() *tmsh.BigIP {
	key := []byte{}
	keyPath := viper.GetString("tmsh_identity_file")
	if keyPath != "" {
		var err error
		key, err = ioutil.ReadFile(keyPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	bigip, err := tmsh.GenSession(
		viper.GetString("tmsh_host"),
		viper.GetString("tmsh_port"),
		viper.GetString("tmsh_user"),
		viper.GetString("tmsh_password"),
		key,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return bigip
}
