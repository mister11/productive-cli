package productive

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/url"
	"reflect"
	"time"
)

type DealResponse struct {
	ID   string `jsonapi:"primary,deals"`
	Name string `jsonapi:"attr,name"`
}

type dealsService struct {
	client *client
}

func newDealsService(client *client) *dealsService {
	return &dealsService{
		client: client,
	}
}

func (service *dealsService) SearchDeals(
	query string,
	startDate time.Time,
	endDate time.Time,
	headers map[string]string,
) ([]DealResponse, error) {
	startDateFormatted := formatDate(startDate)
	endDateFormatted := formatDate(endDate)
	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
		url.QueryEscape(query), startDateFormatted, endDateFormatted)

	req, err := service.client.NewRequest("GET", uri, nil, headers)
	if err != nil {
		return nil, err
	}
	dealsResponseBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	var dealsResponse []DealResponse
	dealResponseInterfaces, err := jsonapi.UnmarshalManyPayload(dealsResponseBody, reflect.TypeOf(new(DealResponse)))
	if err != nil {
		return nil, err
	}
	for _, dealResponseInterface := range dealResponseInterfaces {
		dealsResponse = append(dealsResponse, dealResponseInterface.(DealResponse))
	}
	return dealsResponse, nil
}
