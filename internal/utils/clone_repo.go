package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func Clone(repoURL string) error {
	repoName := GetRepoNameFromURL(repoURL)
	if _, err := os.Stat(repoName); !os.IsNotExist(err) {
		fmt.Println("Repository already exists locally. Skipping clone.")
		return nil
	}

	cmd := exec.Command("git", "clone", repoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
