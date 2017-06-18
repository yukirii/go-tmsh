package cmd

import (
	"fmt"
	"os"

	"github.com/shiftky/go-tmsh"
	"github.com/spf13/cobra"
)

func showNodeInfo(node *tmsh.Node) {
	fmt.Println("\nNode Status")
	fmt.Println("  Name           :", node.Name)
	fmt.Println("  Addr           :", node.Addr)
	fmt.Println("  Enabled State  :", node.EnabledState)
	fmt.Println("  Monitor Rule   :", node.MonitorRule)
	fmt.Println("  Monitor Status :", node.MonitorStatus, "\n")
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Configures a node",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var nodeCreateCmd = &cobra.Command{
	Use:   "create [node name] [IP address]",
	Short: "Create a node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.CreateNode(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()

		node, err := bigip.GetNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showNodeInfo(node)
	},
}

var nodeDeleteCmd = &cobra.Command{
	Use:   "delete [node name]",
	Short: "Delete a node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.DeleteNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()
	},
}

var nodeShowCmd = &cobra.Command{
	Use:   "show [node name]",
	Short: "Show a node information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		node, err := bigip.GetNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showNodeInfo(node)
	},
}

var nodeEnableCmd = &cobra.Command{
	Use:   "enable [node name]",
	Short: "Enable a node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.EnableNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()

		node, err := bigip.GetNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showNodeInfo(node)
	},
}

var nodeDisableCmd = &cobra.Command{
	Use:   "disable [node name]",
	Short: "Disable a node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.DisableNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()

		node, err := bigip.GetNode(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showNodeInfo(node)
	},
}

func init() {
	nodeCmd.AddCommand(nodeCreateCmd)
	nodeCmd.AddCommand(nodeDeleteCmd)
	nodeCmd.AddCommand(nodeShowCmd)
	nodeCmd.AddCommand(nodeEnableCmd)
	nodeCmd.AddCommand(nodeDisableCmd)
	RootCmd.AddCommand(nodeCmd)
}
