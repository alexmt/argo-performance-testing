package main

import (
	"fmt"
	"os"

	"github.com/alexmt/argo-performance-testing/simulator/cmd/commands"
)

func main() {
	command := commands.NewCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
