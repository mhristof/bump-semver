package tag

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	isemver "github.com/Masterminds/semver"
	"github.com/mhristof/semver/log"
	"golang.org/x/mod/semver"
)

func Eval(command string) []string {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"command": command,
			"stderr":  string(stderr.Bytes()),
		}).Panic("Error executing command")
	}

	return strings.Split(strings.TrimSuffix(string(stdout.Bytes()), "\n"), "\n")
}

type BySemVer []string

func (a BySemVer) Len() int           { return len(a) }
func (a BySemVer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySemVer) Less(i, j int) bool { return semver.Compare(a[i], a[j]) < 0 }

func Get(path string) []string {

	var ret []string
	for _, tag := range Eval(fmt.Sprintf("git -C %s show-ref --tag", path)) {
		parts := strings.Split(tag, " ")
		vTag := filepath.Base(parts[1])

		if !semver.IsValid(vTag) {
			log.WithFields(log.Fields{
				"vTag": vTag,
			}).Debug("Not a semver, skipping")

			continue
		}
		ret = append(ret, vTag)

	}

	sort.Sort(BySemVer(ret))
	return ret
}

func Increment(version string, major, minor, patch bool) string {
	v, err := isemver.NewVersion(version)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Panic("Cannot convert to semver")

	}

	var newV isemver.Version
	switch {
	case major:
		newV = v.IncMajor()
	case minor:
		newV = v.IncMinor()
	case patch:
		newV = v.IncPatch()
	default:
		log.WithFields(log.Fields{
			"version": version,
			"major":   major,
			"minor":   minor,
			"patch":   patch,
		}).Panic("Not sure what to do.")
	}

	return "v" + newV.String()
}
