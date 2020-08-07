package cli

import (
	"github.com/mister11/productive-cli/internal/service"
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/urfave/cli/v2"
)


func trackFood(context *cli.Context) error {
	trackFoodRequest := tracking.TrackFoodRequest{
		IsWeekTracking: context.Bool("w"),
		Day:            context.String("d"),
	}
	return service.NewTrackingService().TrackFood(trackFoodRequest)
}

func trackProject(context *cli.Context) error {
	trackProjectRequest := tracking.TrackProjectRequest{
		Day: context.String("d"),
	}
	return service.NewTrackingService().TrackProject(trackProjectRequest)
}