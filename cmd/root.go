package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mhristof/semver/tag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:     "semver",
	Short:   "Create semver releases",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		list, abs := tags()

		major, err := cmd.Flags().GetBool("major")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("Cannot get major flag")
		}

		minor, err := cmd.Flags().GetBool("minor")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("Cannot get minor flag")
		}

		patch, err := cmd.Flags().GetBool("patch")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Cannot get patch flag")
		}

		var lastTag string
		if len(list) == 0 {
			lastTag = "v0.0.0"
		} else {
			lastTag = list[len(list)-1]
		}

		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Panic("cannot retrieve auto value")
		}

		if auto {
			major, minor, patch = tag.FindNext(lastTag)
		}

		// default to minor release
		minor = minor || !major && !minor && !patch

		next := tag.Increment(lastTag, major, minor, patch)

		gitCmd := fmt.Sprintf("git -C %s tag %s", abs, next)
		if silent, _ := cmd.Flags().GetBool("silent"); !silent {
			fmt.Println(gitCmd)
		}

		if dryrun, _ := cmd.Flags().GetBool("dryrun"); dryrun {
			return
		}

		tag.Eval(gitCmd)
	},
}

func tags() ([]string, string) {
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

	return tag.Get(abs), abs
}

// Verbose Increase verbosity.
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
	rootCmd.Flags().BoolP("minor", "m", false, "Perform a minor release. Default if auto is disabled.")
	rootCmd.Flags().BoolP("patch", "p", false, "Perform a patch release")
	rootCmd.Flags().BoolP("silent", "s", false, "Disable all output")
	rootCmd.Flags().BoolP("auto", "a", true, "Autodetect next version based on commit messages")

	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run mode")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
