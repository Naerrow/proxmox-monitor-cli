package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "5초마다 자동 갱신",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: watch 명령어 구현 예정")
	},
}
