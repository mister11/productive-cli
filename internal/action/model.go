package action

type TrackFoodRequest struct {
	IsWeekTracking bool
	Day            string
}

type TrackProjectRequest struct {
	Day string
}

func (trackFoodRequest TrackFoodRequest) IsValid() bool {
	return !(trackFoodRequest.IsWeekTracking && trackFoodRequest.Day != "")
}
