package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Example of how to:
// - Clone a repository into memory
// - Get the HEAD reference
// - Using the HEAD reference, obtain the commit this reference is pointing to
// - Using the commit, obtain its history and print it
func main() {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	fmt.Println("========= git clone https://github.com/gorilych/test.git")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/gorilych/test.git",
	})
	if err != nil {
		fmt.Println(err)
	}

	// retrieve commit history with commits modifying README
	fName := "README.md"
	cIter, err := r.Log(&git.LogOptions{FileName: &fName})
	if err != nil {
		fmt.Println(err)
	}

	// find most recent commit
	var cmt *object.Commit

	cIter.ForEach(func(c *object.Commit) error {
		cmt = c
		return fmt.Errorf("first commit found")
	})
	fmt.Println("========= Most recent COMMMIT which modified README")
	fmt.Println(cmt.Hash, cmt.Committer.When, strings.TrimSpace(cmt.Message))

	for {
		time.Sleep(2 * time.Second)
		newcmts := []*object.Commit{}
		err = r.Fetch(&git.FetchOptions{})
		if err != nil {
			fmt.Println("========= No new commits: ", err)
		} else {
			fmt.Println("========= I see new commits!")
			cIter, err = r.Log(&git.LogOptions{FileName: &fName, All: true})
			if err != nil {
				fmt.Println("========= But have error ", err)
			}

			cIter.ForEach(func(c *object.Commit) error {
				if c.Hash == cmt.Hash {
					return fmt.Errorf("found all commits till previos one")
				}
				newcmts = append(newcmts, c)
				return nil
			})

			if len(newcmts) > 0 {
				// iterate over new commits in reverse order
				for i := len(newcmts) - 1; i >= 0; i-- {
					curcmt := newcmts[i]
					fmt.Println("========= Next COMMMIT which modified README")
					fmt.Println(curcmt.Hash, curcmt.Committer.When, strings.TrimSpace(curcmt.Message))
				}
				// update most recent commit
				cmt = newcmts[0]
			}
		}
	}

}
