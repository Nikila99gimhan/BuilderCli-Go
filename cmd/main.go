package main

import (
	"cliapp/pkg/cliapplication"
	"fmt"
	"os"
)

func main() {
	var app cliapplication.CliApplication = cliapplication.NewCliApplicationImpl()
	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}
}
