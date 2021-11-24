package tag

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cases := []struct {
		name string
		tags []string
		exp  []string
	}{
		{
			name: "simple tags",
			tags: []string{
				"v1.0.0",
				"v1.2.0",
			},
			exp: []string{
				"v1.0.0",
				"v1.2.0",
			},
		},
		{
			name: "skip non semver tags",
			tags: []string{
				"v1.0.0",
				"not-semver",
				"v1.2.0",
			},
			exp: []string{
				"v1.0.0",
				"v1.2.0",
			},
		},
	}

	for _, test := range cases {
		folder, err := ioutil.TempDir("", "sampledir")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(folder)

		fmt.Println(folder)
		Eval(fmt.Sprintf("git -C %s init", folder))
		Eval(fmt.Sprintf("touch %s/test", folder))
		Eval(fmt.Sprintf("git -C %s add .", folder))
		Eval(fmt.Sprintf("git -C %s commit -m init", folder))

		for _, tag := range test.tags {
			Eval(fmt.Sprintf("git -C %s tag %s", folder, tag))
		}

		tags := Get(folder)

		assert.Equal(t, tags, test.exp, test.name)

	}
}

func TestIncrement(t *testing.T) {
	cases := []struct {
		name    string
		version string
		major   bool
		minor   bool
		patch   bool
		exp     string
	}{
		{
			name:    "increase major",
			version: "v1.1.1",
			major:   true,
			exp:     "v2.0.0",
		},
		{
			name:    "increase minor",
			version: "v1.1.1",
			minor:   true,
			exp:     "v1.2.0",
		},
		{
			name:    "increase patch",
			version: "v1.1.1",
			patch:   true,
			exp:     "v1.1.2",
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.exp, Increment(test.version, test.major, test.minor, test.patch), test.name)
	}
}

func tmpDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "semver")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestFindNext(t *testing.T) {
	cases := []struct {
		name  string
		cmds  []string
		repo  string
		patch bool
		minor bool
		major bool
	}{
		{
			name: "next bug",
			cmds: strings.Split(heredoc.Doc(`
				git init
				git commit --allow-empty -m initial.import
				git tag v0.1.0
				git commit --allow-empty -m bug:-test`), "\n"),
			repo:  tmpDir(t),
			patch: true,
		},
		{
			name: "next minor",
			cmds: strings.Split(heredoc.Doc(`
				git init
				git commit --allow-empty -m initial.import
				git tag v0.1.0
				git commit --allow-empty -m feature:-test
				git commit --allow-empty -m bug:-test`), "\n"),
			repo:  tmpDir(t),
			minor: true,
		},
	}

	for _, test := range cases {
		err := os.Chdir(test.repo)
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range test.cmds {
			_, err := Eval(c)
			if err != nil {
				t.Fatal(err)
			}
		}

		major, minor, patch := FindNext("v0.1.0")

		assert.Equal(t, test.patch, patch, test.name+" patch")
		assert.Equal(t, test.minor, minor, test.name+" minor")
		assert.Equal(t, test.major, major, test.name+" major")
	}
}
