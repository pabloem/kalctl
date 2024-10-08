package main

import (
	"fmt"
	"os"

	"github.com/pabloem/kalctl/commands"
)

func main() {
	err := commands.RunCommand(os.Args)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}
