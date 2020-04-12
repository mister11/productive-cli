package app

import (
	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/stdin/promptui"
	"github.com/urfave/cli/v2"
)

func CreateProductiveCliApp() *cli.App {
	stdin := promptui.NewPromptUiStdin()
	configManager := config.NewFileConfigManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	productiveClient := client.NewProductiveClient(configManager, dateTimeProvider)

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
							action.TrackFood(productiveClient, configManager, dateTimeProvider, trackFoodRequest)
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
							action.TrackProject(productiveClient, stdin, configManager, trackProjectRequest)
							return nil
						},
					},
				},
			},
			{
				Name:  "init",
				Usage: "Initializes user data",
				Action: func(c *cli.Context) error {
					action.Init(productiveClient, stdin, configManager)
					return nil
				},
			},
		},
	}
}
