package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pmon",
	Short: "Proxmox 모니터링 CLI 도구",
	Long:  "터미널에서 Proxmox 서버 상태 확인 및 VM 제어",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(vmsCmd)
	rootCmd.AddCommand(vmCmd)
	rootCmd.AddCommand(watchCmd)
}
