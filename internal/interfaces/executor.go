package interfaces

import (
	"github.com/mister11/productive-cli/internal/application"
	tracking2 "github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/urfave/cli/v2"
)


func trackFood(context *cli.Context) error {
	trackFoodRequest := tracking2.TrackFoodRequest{
		IsWeekTracking: context.Bool("w"),
		Day:            context.String("d"),
	}
	return application.NewTrackingService().TrackFood(trackFoodRequest)
}

func trackProject(context *cli.Context) error {
	trackProjectRequest := tracking2.TrackProjectRequest{
		Day: context.String("d"),
	}
	return application.NewTrackingService().TrackProject(trackProjectRequest)
}