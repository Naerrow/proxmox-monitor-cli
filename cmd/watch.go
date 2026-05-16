package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/naerrow/proxmox-monitor-cli/config"
	"github.com/naerrow/proxmox-monitor-cli/internal/proxmox"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "5초마다 자동 갱신",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		fmt.Println("👀 watching... (종료: Ctrl+C)")

		for {
			fmt.Print("\033[H\033[2J")
			fmt.Printf("🕐 %s  (5초마다 갱신, 종료: Ctrl+C)\n\n", time.Now().Format("2006-01-02 15:04:05"))

			fmt.Println("[ 노드 ]")
			printNodes(client)

			fmt.Println()

			fmt.Println("[ VM 목록 ]")
			printVMs(client)

			select {
			case <-sig:
				fmt.Println("\n👋 종료합니다.")
				return
			case <-time.After(5 * time.Second):
			}
		}
	},
}
