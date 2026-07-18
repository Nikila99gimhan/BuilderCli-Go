package cliapplication

// CliApplication is the top-level interface for the CLI application.
// Start runs the main interactive or non-interactive workflow and returns
// a non-nil error if the process fails, so callers can exit with code 1.
type CliApplication interface {
	Start() error
}
