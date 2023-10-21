package main

import (
	"cliapp/pkg/cliapplication"
)

func main() {
	var app cliapplication.CliApplication = cliapplication.NewCliApplicationImpl()
	app.Start()

}
