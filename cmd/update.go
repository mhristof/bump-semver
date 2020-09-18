package cmd

import (
	"fmt"
	"runtime"

	"github.com/mhristof/go-update"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the binary with a new version",
		Run: func(cmd *cobra.Command, args []string) {
			url := fmt.Sprintf("https://github.com/mhristof/semver/releases/latest/download/semver.%s", runtime.GOOS)
			//
			updates, updateFunc, err := update.Check(url)
			if err != nil {
				panic(err)
			}

			if updates {
				fmt.Println("New version downloaded!")
				updateFunc()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}
