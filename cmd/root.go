package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mhristof/semver/log"
	"github.com/mhristof/semver/tag"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Create semver releases",
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		abs, err := filepath.Abs(pwd)
		if err != nil {
			panic(err)
		}

		tag.Get(abs)
	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
