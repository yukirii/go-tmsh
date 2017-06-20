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

var createCmd = &cobra.Command{
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

var deleteCmd = &cobra.Command{
	Use:   "delete [pool name]",
	Short: "Delete a pool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.DeletePool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bigip.Save()
	},
}

var showCmd = &cobra.Command{
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

var addMonitorCmd = &cobra.Command{
	Use:   "add-monitor [pool name] [monitor name]",
	Short: "Add monitor to pool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		err := bigip.AddMonitorToPool(args[0], args[1])
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

var addMemberCmd = &cobra.Command{
	Use:   "add-member [pool name] [node name] [port num] [monitor name]",
	Short: "Add member node to pool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = bigip.AddPoolMember(args[0], args[1], args[3], port)
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

var deleteMemberCmd = &cobra.Command{
	Use:   "delete-member [pool name] [node name] [port num]",
	Short: "Delete member node from pool",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = bigip.DeletePoolMember(args[0], args[1], port)
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

var enableMemberCmd = &cobra.Command{
	Use:   "enable-member [pool name] [node name] [port num]",
	Short: "Enable member node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = bigip.EnablePoolMember(args[0], args[1], port)
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

var disableMemberCmd = &cobra.Command{
	Use:   "disable-member [pool name] [node name] [port num]",
	Short: "Disable member node",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Usage()
			os.Exit(2)
		}

		bigip := NewSession()
		defer bigip.Close()

		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = bigip.DisablePoolMember(args[0], args[1], port)
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

func init() {
	poolCmd.AddCommand(createCmd)
	poolCmd.AddCommand(deleteCmd)
	poolCmd.AddCommand(showCmd)

	poolCmd.AddCommand(addMonitorCmd)

	poolCmd.AddCommand(addMemberCmd)
	poolCmd.AddCommand(deleteMemberCmd)
	poolCmd.AddCommand(enableMemberCmd)
	poolCmd.AddCommand(disableMemberCmd)

	RootCmd.AddCommand(poolCmd)
}
