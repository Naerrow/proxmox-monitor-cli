package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "전체 VM 목록과 상태",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: vms 명령어 구현 예정")
	},
}
