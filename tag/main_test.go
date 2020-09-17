package tag

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var cases = []struct {
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
		eval(fmt.Sprintf("git -C %s init", folder))
		eval(fmt.Sprintf("touch %s/test", folder))
		eval(fmt.Sprintf("git -C %s add .", folder))
		eval(fmt.Sprintf("git -C %s commit -m init", folder))

		for _, tag := range test.tags {
			eval(fmt.Sprintf("git -C %s tag %s", folder, tag))
		}

		tags := Get(folder)

		assert.Equal(t, tags, test.exp, test.name)

	}
}
