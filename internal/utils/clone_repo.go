package utils

import (
	"github.com/go-git/go-git/v5"
	"os"
)

// Clone clones the given repoURL using the native go-git library, removing
// the dependency on a host-installed `git` binary.
func Clone(repoURL string) error {
	repoName := GetRepoNameFromURL(repoURL)
	if _, err := os.Stat(repoName); !os.IsNotExist(err) {
		return nil // Already exists locally, skip clone.
	}

	_, err := git.PlainClone(repoName, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	return err
}
