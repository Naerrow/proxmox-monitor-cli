package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/naerrow/proxmox-monitor-cli/config"
	"github.com/naerrow/proxmox-monitor-cli/internal/proxmox"
	"github.com/spf13/cobra"
)

var vmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "전체 VM 목록과 상태",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		client := proxmox.NewClient(cfg.ProxmoxURL, cfg.ProxmoxToken, cfg.ProxmoxNode)
		printVMs(client)
	},
}

func printVMs(client *proxmox.Client) {
	vms, err := client.GetVMs()
	if err != nil {
		fmt.Println("❌ VM 조회 오류:", err)
		os.Exit(1)
	}

	if len(vms) == 0 {
		fmt.Println("VM이 없습니다.")
		return
	}

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
