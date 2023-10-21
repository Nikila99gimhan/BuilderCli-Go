package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
