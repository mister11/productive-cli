package app

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/mister11/productive-cli/internal/action"
	"gitlab.com/mister11/productive-cli/internal/client"
)

func CreateProductiveCliApp() *cli.App {
	productiveClient := client.NewProductiveClient()

	return &cli.App{
		Name:                 "Productive CLI",
		Usage:                "Manage Productive from your terminal!",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "track",
				Usage: "Track time for any service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "d",
						Usage: "track particular day (format: YYYY-MM-DD)",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "food",
						Usage: "Default 30 mintues for lunch",
						Action: func(c *cli.Context) error {
							trackFoodRequest := action.TrackFoodRequest{
								IsWeekTracking: c.Bool("w"),
								Day:            c.String("d"),
							}
							action.TrackFood(productiveClient, trackFoodRequest)
							return nil
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "w",
								Usage: "track whole week",
							},
						},
					},
					{
						Name:  "project",
						Usage: "Track project",
						Action: func(c *cli.Context) error {
							trackProjectRequest := action.TrackProjectRequest{
								Day: c.String("d"),
							}
							action.TrackProject(productiveClient, trackProjectRequest)
							return nil
						},
					},
				},
			},
			{
				Name:  "init",
				Usage: "Initializes user data",
				Action: func(c *cli.Context) error {
					action.Init(productiveClient)
					return nil
				},
			},
		},
	}
}
