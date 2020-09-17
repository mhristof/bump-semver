package tag

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Get(path string) []string {
	fmt.Println(fmt.Sprintf("path: %+v", path))

	r, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	tags, err := r.TagObjects()
	if err != nil {
		panic(err)
	}

	err = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t)
		return nil
	})

	return []string{}
}
