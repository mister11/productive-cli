package main

import (
	"os"

	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/interfaces"
)

func main() {

	if err := interfaces.CLI().Run(os.Args); err != nil {
		log.Error("CLI quit unexpectedly", err)
		os.Exit(-1)
	}
}
