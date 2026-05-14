package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "VM 제어 (start/stop/delete)",
}

var vmStartCmd = &cobra.Command{
	Use:   "start [vmid]",
	Short: "VM 시작",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: VM %s 시작 구현 예정\n", args[0])
	},
}

var vmStopCmd = &cobra.Command{
	Use:   "stop [vmid]",
	Short: "VM 중지",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: VM %s 중지 구현 예정\n", args[0])
	},
}

var vmDeleteCmd = &cobra.Command{
	Use:   "delete [vmid]",
	Short: "VM 삭제",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: VM %s 삭제 구현 예정\n", args[0])
	},
}

func init() {
	vmCmd.AddCommand(vmStartCmd)
	vmCmd.AddCommand(vmStopCmd)
	vmCmd.AddCommand(vmDeleteCmd)
}
