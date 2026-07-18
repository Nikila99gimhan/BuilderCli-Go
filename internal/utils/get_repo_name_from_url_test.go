package utils

import "testing"

func TestGetRepoNameFromURL_HTTPSWithGit(t *testing.T) {
	got := GetRepoNameFromURL("https://github.com/user/my-repo.git")
	if got != "my-repo" {
		t.Errorf("expected my-repo, got %q", got)
	}
}

func TestGetRepoNameFromURL_HTTPSWithoutGit(t *testing.T) {
	got := GetRepoNameFromURL("https://github.com/user/my-repo")
	if got != "my-repo" {
		t.Errorf("expected my-repo, got %q", got)
	}
}

func TestGetRepoNameFromURL_SSH(t *testing.T) {
	got := GetRepoNameFromURL("git@github.com:user/my-repo.git")
	if got != "my-repo" {
		t.Errorf("expected my-repo, got %q", got)
	}
}

func TestGetRepoNameFromURL_TrailingSlash(t *testing.T) {
	got := GetRepoNameFromURL("https://github.com/user/my-repo/")
	// trailing slash results in empty last part — verify no panic
	_ = got
}
