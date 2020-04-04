package main

import (
	"os"

	"github.com/mister11/productive-cli/internal/app"
)

func main() {

	cliApp := app.CreateProductiveCliApp()

	err := cliApp.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
