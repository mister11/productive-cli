package cli

import (
	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service"
	"github.com/urfave/cli/v2"
)


func trackFood(context *cli.Context) error {
	trackFoodRequest := service.TrackFoodRequest{
		IsWeekTracking: context.Bool("w"),
		Day:            context.String("d"),
	}
	client := productive.NewClient(nil)
	return service.NewFoodTrackingService(client).TrackFood(trackFoodRequest)
}

func trackProject(context *cli.Context) error {
	trackProjectRequest := service.TrackProjectRequest{
		Day: context.String("d"),
	}
	client := productive.NewClient(nil)
	return service.NewProjectTrackingService(client).TrackProject(trackProjectRequest)
}