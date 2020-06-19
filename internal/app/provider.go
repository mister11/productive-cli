package app

import (
	"github.com/urfave/cli/v2"
)

type commandsProvider struct {
	commandExecutor *commandsExecutor
}

func newProvider() *commandsProvider {
	return &commandsProvider{
		commandExecutor: newExecutor(),
	}
}

func (provider *commandsProvider) Init() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes user data",
		Action: func(c *cli.Context) error {
			return provider.commandExecutor.Init()
		},
	}
}

func (provider *commandsProvider) Track(subcommands ...*cli.Command) *cli.Command {
	return &cli.Command{
		Name:  "track",
		Usage: "Track food for any service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "d",
				Usage: "track particular day (format: YYYY-MM-DD)",
			},
		},
		Subcommands: subcommands,
	}
}

func (provider *commandsProvider) TrackFood() *cli.Command {
	return &cli.Command{
		Name:  "food",
		Usage: "Default 30 mintues for lunch",
		Action: func(c *cli.Context) error {
			return provider.commandExecutor.TrackFood(c)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "w",
				Usage: "track whole week",
			},
		},
	}
}

func (provider *commandsProvider) TrackProject() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Track project",
		Action: func(c *cli.Context) error {
			return provider.commandExecutor.TrackProject(c)
		},
	}
}
