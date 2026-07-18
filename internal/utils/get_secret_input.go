package utils

import (
	"golang.org/x/term"
	"os"
)

// GetSecretInput reads a password from stdin without echoing it to the terminal.
// This is critical for security when prompting for passwords.
func GetSecretInput() (string, error) {
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(password), nil
}
