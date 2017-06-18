package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/shiftky/go-tmsh"
	"github.com/spf13/cobra"
)

func showPoolInfo(pool *tmsh.Pool) {
	fmt.Println("\nPool Status")
	fmt.Println("  Name               :", pool.Name)
	fmt.Println("  Active Members     :", pool.ActiveMemberCount)
	fmt.Println("  Monitor Rule       :", pool.MonitorRule)
	fmt.Println("  Availability State :", pool.AvailabilityState)
	fmt.Println("  Enabled State      :", pool.EnabledState)
	fmt.Println("  Status Reason      :", pool.StatusReason)

	if pool.ActiveMemberCount > 0 {
		fmt.Println("\nPool Members")
		for _, member := range pool.PoolMembers {
			fmt.Println("  #", member.Name, "("+member.Addr+":"+strconv.Itoa(member.Port)+")")
			fmt.Println("    Monitor Rule       :", member.MonitorRule)
			fmt.Println("    Monitor Status     :", member.MonitorStatus)
			fmt.Println("    Enabled State      :", member.EnabledState)
			fmt.Println("    Availability State :", member.AvailabilityState)
			fmt.Println("    Status Reason      :", member.StatusReason, "\n")
		}
	}
}

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Configures a load balancing pool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var poolCreateCmd = &cobra.Command{
	Use:   "create [pool name]",
	Short: "Create a pool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.CreatePool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()

		pool, err := bigip.GetPool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showPoolInfo(pool)
	},
}

var poolShowCmd = &cobra.Command{
	Use:   "show [pool name]",
	Short: "Show a pool information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		pool, err := bigip.GetPool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showPoolInfo(pool)
	},
}

func init() {
	poolCmd.AddCommand(poolCreateCmd)
	poolCmd.AddCommand(poolShowCmd)
	RootCmd.AddCommand(poolCmd)
}
