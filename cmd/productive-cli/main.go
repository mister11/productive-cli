package main

import (
	"fmt"
	"github.com/mister11/productive-cli/internal/interfaces/cli"
	"os"
)

func main() {

	if err := cli.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}
