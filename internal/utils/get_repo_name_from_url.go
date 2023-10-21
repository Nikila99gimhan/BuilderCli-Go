package utils

import "strings"

func GetRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
