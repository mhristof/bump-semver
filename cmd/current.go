package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	currentCmd = &cobra.Command{
		Use:   "current",
		Short: "Show the current tag",
		Run: func(cmd *cobra.Command, args []string) {
			list, _ := tags()
			fmt.Println(list[len(list)-1])
		},
	}
)

func init() {
	rootCmd.AddCommand(currentCmd)
}
