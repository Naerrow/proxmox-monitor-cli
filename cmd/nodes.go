package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "노드 목록 (CPU/메모리 포함)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: nodes 명령어 구현 예정")
	},
}
