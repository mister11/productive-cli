package cli

import (
	"github.com/urfave/cli/v2"
)

func Run(args[] string) error {
	return createProductiveCliApp().Run(args)
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
