package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/naerrow/proxmox-monitor-cli/config"
	"github.com/naerrow/proxmox-monitor-cli/internal/proxmox"
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
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)

		vmid := args[0]
		fmt.Printf("🚀 VM %s 시작 중...\n", vmid)

		if err := client.StartVM(vmid); err != nil {
			fmt.Println("❌ 오류:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ VM %s 시작 완료\n", vmid)
	},
}

var vmStopCmd = &cobra.Command{
	Use:   "stop [vmid]",
	Short: "VM 중지",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)

		vmid := args[0]
		fmt.Printf("🛑 VM %s 중지 중...\n", vmid)

		if err := client.StopVM(vmid); err != nil {
			fmt.Println("❌ 오류:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ VM %s 중지 완료\n", vmid)
	},
}

var vmDeleteCmd = &cobra.Command{
	Use:   "delete [vmid]",
	Short: "VM 삭제",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)

		vmid := args[0]

		fmt.Printf("⚠️  VM %s 를 정말 삭제할까요? (yes/no): ", vmid)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "yes" {
			fmt.Println("❌ 삭제 취소")
			return
		}

		fmt.Printf("🗑️  VM %s 삭제 중...\n", vmid)

		if err := client.DeleteVM(vmid); err != nil {
			fmt.Println("❌ 오류:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ VM %s 삭제 완료\n", vmid)
	},
}

func init() {
	vmCmd.AddCommand(vmStartCmd)
	vmCmd.AddCommand(vmStopCmd)
	vmCmd.AddCommand(vmDeleteCmd)
}
