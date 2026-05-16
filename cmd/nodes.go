package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/naerrow/proxmox-monitor-cli/config"
	"github.com/naerrow/proxmox-monitor-cli/internal/proxmox"
	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "노드 목록 (CPU/메모리 포함)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)

		nodes, err := client.GetNodes()
		if err != nil {
			fmt.Println("❌ 오류:", err)
			os.Exit(1)
		}

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
	},
}
