package tag

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mhristof/semver/log"
	"golang.org/x/mod/semver"
)

func eval(command string) []string {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return strings.Split(strings.TrimSuffix(string(stdout.Bytes()), "\n"), "\n")
}

type BySemVer []string

func (a BySemVer) Len() int           { return len(a) }
func (a BySemVer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySemVer) Less(i, j int) bool { return semver.Compare(a[i], a[j]) < 0 }

func Get(path string) []string {
	fmt.Println(fmt.Sprintf("path: %+v", path))

	var ret []string
	for _, tag := range eval(fmt.Sprintf("git -C %s show-ref --tag", path)) {
		parts := strings.Split(tag, " ")
		vTag := filepath.Base(parts[1])

		if !semver.IsValid(vTag) {
			log.WithFields(log.Fields{
				"vTag": vTag,
			}).Debug("Not a semver, skipping")

			continue
		}
		fmt.Println(fmt.Sprintf("vTag: %+v", vTag))

		ret = append(ret, vTag)

	}

	sort.Sort(BySemVer(ret))
	return ret
}
