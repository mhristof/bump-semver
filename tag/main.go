package tag

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Get(path string) []string {
	fmt.Println(fmt.Sprintf("path: %+v", path))

	r, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	tagrefs, err := r.Tags()

	fmt.Println(fmt.Sprintf("tagrefs: %+v", tagrefs))

	return []string{}
}
