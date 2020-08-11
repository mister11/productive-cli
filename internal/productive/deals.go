package productive

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/url"
	"reflect"
	"time"
)

type Deal struct {
	ID   string `jsonapi:"primary,deals"`
	Name string `jsonapi:"attr,name"`
}

type dealsService struct {
	client *Client
}

func newDealsService(client *Client) *dealsService {
	return &dealsService{
		client: client,
	}
}

func (service *dealsService) SearchDeals(
	query string,
	startDate time.Time,
	endDate time.Time,
	token string,
) ([]Deal, error) {
	startDateFormatted := formatDate(startDate)
	endDateFormatted := formatDate(endDate)
	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
		url.QueryEscape(query), startDateFormatted, endDateFormatted)

	req, err := service.client.NewRequest("GET", uri, nil, getHeaders(token))
	if err != nil {
		return nil, err
	}
	dealsResponseBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	var dealsResponse []Deal
	dealResponseInterfaces, err := jsonapi.UnmarshalManyPayload(dealsResponseBody, reflect.TypeOf(new(Deal)))
	if err != nil {
		return nil, err
	}
	for _, dealResponseInterface := range dealResponseInterfaces {
		dealsResponse = append(dealsResponse, dealResponseInterface.(Deal))
	}
	return dealsResponse, nil
}
