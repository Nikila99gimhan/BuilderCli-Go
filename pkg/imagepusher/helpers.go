package imagepusher

import "strings"

// newStringReader is a tiny helper to pipe a string into a command's Stdin.
func newStringReader(s string) *strings.Reader {
	return strings.NewReader(s)
}
