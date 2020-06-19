package app

import (
	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/action/track"
	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/stdin/promptui"
	"github.com/urfave/cli/v2"
)

type commandsExecutor struct {
	loginManager *action.LoginManger
	trackingManager *track.TrackingManager
}

func newExecutor() *commandsExecutor {
	stdin := promptui.NewPromptUiStdin()
	configManager := config.NewFileConfigManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	productiveClient := client.NewProductiveClient(configManager)
	return &commandsExecutor{
		loginManager: action.NewLoginManager(productiveClient, stdin, configManager),
		trackingManager: track.NewTrackingManager(productiveClient, stdin, configManager, dateTimeProvider),
	}
}

func (executor *commandsExecutor) Init() error {
	executor.loginManager.Init()
	return nil
}

func (executor *commandsExecutor) TrackFood(context *cli.Context) error {
	trackFoodRequest := action.TrackFoodRequest{
		IsWeekTracking: context.Bool("w"),
		Day:            context.String("d"),
	}
	executor.trackingManager.TrackFood(trackFoodRequest)
	return nil
}

func (executor *commandsExecutor) TrackProject(context *cli.Context) error {
	trackProjectRequest := action.TrackProjectRequest{
		Day: context.String("d"),
	}
	executor.trackingManager.TrackProject(trackProjectRequest)
	return nil
}