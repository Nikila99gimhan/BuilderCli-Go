package cliapplication

import "fmt"

const Green = "\033[32m"
const Bold = "\033[1m"
const Reset = "\033[0m"

type Display struct{}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) PrintLogo() {
	logo := `       	____        _ __    __         
	/ __ )__  __(_) /___/ /__  _____
 / __  / / / / / / __  / _ \/ ___/
/ /_/ / /_/ / / / /_/ /  __/ /    
/_____/\__,_/_/_/\__,_/\___/_/     

																	
	` // Your ASCII art goes here.
	fmt.Println(Green + Bold + logo + Reset)
}

func (d *Display) PrintGreenBold(message string) {
	fmt.Println(Green + Bold + message + Reset)
}

func (d *Display) Print(message string) {
	fmt.Println(message)
}
