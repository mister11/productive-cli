package main

import (
	"os"

	"github.com/mister11/productive-cli/internal/app"
	"github.com/mister11/productive-cli/internal/log"
)

func main() {

	if err := app.NewProductiveCLI().Run(os.Args); err != nil {
		log.Error("CLI quit unexpectedly", err)
		os.Exit(-1)
	}
}
