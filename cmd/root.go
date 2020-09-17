package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mhristof/semver/log"
	"github.com/mhristof/semver/tag"
	"github.com/spf13/cobra"
)

var (
	version string
)

var rootCmd = &cobra.Command{
	Use:     "semver",
	Short:   "Create semver releases",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		pwd, err := os.Getwd()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("Cannot get pwd")

		}

		abs, err := filepath.Abs(pwd)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"pwd": pwd,
			}).Panic("Cannot get abs path")

		}

		list := tag.Get(abs)
		major, err := cmd.Flags().GetBool("major")
		minor, err := cmd.Flags().GetBool("minor")
		patch, err := cmd.Flags().GetBool("patch")

		// minor is the default increment and we need to turn if off if one
		// of the other levels are set.
		minor = minor && !(major || patch)

		next := tag.Increment(list[len(list)-1], major, minor, patch)

		gitCmd := fmt.Sprintf("git -C %s tag %s", abs, next)
		if dry, err := cmd.Flags().GetBool("dryrun"); err == nil && dry {
			fmt.Println(gitCmd)
		} else {
			tag.Eval(gitCmd)
		}
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
	rootCmd.Flags().BoolP("major", "M", false, "Perform a major release")
	rootCmd.Flags().BoolP("minor", "m", true, "Perform a minor release")
	rootCmd.Flags().BoolP("patch", "p", false, "Perform a patch release")

	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run mode")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
