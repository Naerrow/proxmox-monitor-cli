package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
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

			// 노드 출력
			fmt.Println("[ 노드 ]")
			nodes, err := client.GetNodes()
			if err != nil {
				fmt.Println("❌ 노드 조회 오류:", err)
			} else {
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"노드", "상태", "CPU 사용률", "메모리 사용량", "업타임"})
				for _, n := range nodes {
					cpuPercent := fmt.Sprintf("%.1f%%", n.CPU*100)
					memUsed := n.Mem / 1024 / 1024 / 1024
					memTotal := n.MaxMem / 1024 / 1024 / 1024
					mem := fmt.Sprintf("%dGB / %dGB", memUsed, memTotal)
					uptime := fmt.Sprintf("%dd %dh", n.Uptime/86400, (n.Uptime%86400)/3600)
					t.AppendRow(table.Row{n.Node, n.Status, cpuPercent, mem, uptime})
				}
				t.SetStyle(table.StyleLight)
				t.Render()
			}

			fmt.Println()

			// VM 출력
			fmt.Println("[ VM 목록 ]")
			vms, err := client.GetVMs()
			if err != nil {
				fmt.Println("❌ VM 조회 오류:", err)
			} else if len(vms) == 0 {
				fmt.Println("VM이 없습니다.")
			} else {
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"VMID", "이름", "상태", "CPU 사용률", "메모리 사용량"})
				for _, vm := range vms {
					cpuPercent := fmt.Sprintf("%.1f%%", vm.CPU*100)
					var mem string
					if vm.Status == "running" {
						memUsed := vm.Mem / 1024 / 1024
						memTotal := vm.MaxMem / 1024 / 1024
						mem = fmt.Sprintf("%dMB / %dMB", memUsed, memTotal)
					} else {
						mem = "-"
					}
					status := vm.Status
					if status == "running" {
						status = "🟢 running"
					} else {
						status = "🔴 stopped"
					}
					t.AppendRow(table.Row{vm.VMID, vm.Name, status, cpuPercent, mem})
				}
				t.SetStyle(table.StyleLight)
				t.Render()
			}

			select {
			case <-sig:
				fmt.Println("\n👋 종료합니다.")
				return
			case <-time.After(5 * time.Second):
			}
		}
	},
}
