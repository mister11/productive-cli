package app

import (
	"github.com/urfave/cli/v2"
)

type ProductiveCLI struct {
	app *cli.App
}

func NewProductiveCLI() *ProductiveCLI {
	return &ProductiveCLI{
		app: createProductiveCliApp(),
	}
}

func (cli *ProductiveCLI) Run(args[] string) error {
	return cli.app.Run(args)
}

func createProductiveCliApp() *cli.App {
	return &cli.App{
		Name:                 "Productive CLI",
		Usage:                "Manage Productive from your terminal!",
		EnableBashCompletion: true,
		BashComplete: cli.DefaultAppComplete,
		Commands: []*cli.Command{
			trackCommand(
				trackFoodSubCommand(),
				trackProjectSubCommand(),
			),
		},
	}
}
