package app

import (
	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/action/track"
	"github.com/urfave/cli/v2"
)


func trackFood(context *cli.Context) error {
	trackFoodRequest := action.TrackFoodRequest{
		IsWeekTracking: context.Bool("w"),
		Day:            context.String("d"),
	}
	track.NewTrackingManager().TrackFood(trackFoodRequest)
	return nil
}

func trackProject(context *cli.Context) error {
	trackProjectRequest := action.TrackProjectRequest{
		Day: context.String("d"),
	}
	track.NewTrackingManager().TrackProject(trackProjectRequest)
	return nil
}