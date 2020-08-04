package cli

import (
	"github.com/urfave/cli/v2"
)

func trackCommand(subcommands ...*cli.Command) *cli.Command {
	return &cli.Command{
		Name:  "track",
		Usage: "Track time for any service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "d",
				Usage: "track particular day (format: YYYY-MM-DD)",
			},
		},
		Subcommands: subcommands,
	}
}

func trackFoodSubCommand() *cli.Command {
	return &cli.Command{
		Name:  "food",
		Usage: "Default 30 minutes for lunch",
		Action: func(c *cli.Context) error {
			return trackFood(c)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "w",
				Usage: "track whole week",
			},
		},
	}
}

func trackProjectSubCommand() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Track project",
		Action: func(c *cli.Context) error {
			return trackProject(c)
		},
	}
}
