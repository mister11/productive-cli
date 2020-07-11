package main

import (
	"fmt"
	"os"

	"github.com/mister11/productive-cli/internal/interfaces"
)

func main() {

	if err := interfaces.CLI().Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}
