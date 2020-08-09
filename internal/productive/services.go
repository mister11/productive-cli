package productive

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/url"
	"reflect"
	"time"
)

type ServiceResponse struct {
	ID   string `jsonapi:"primary,services"`
	Name string `jsonapi:"attr,name"`
}

type serviceService struct {
	client *client
}

func newServiceService(client *client) *serviceService {
	return &serviceService{
		client: client,
	}
}

func (service *serviceService) SearchServices(
	query string,
	dealID string,
	startDate time.Time,
	endDate time.Time,
	headers map[string]string,
) ([]ServiceResponse, error) {
	startDateFormatted := formatDate(startDate)
	endDateFormatted := formatDate(endDate)
	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
		url.QueryEscape(query), startDateFormatted, endDateFormatted, dealID)
	req, err := service.client.NewRequest("GET", uri, nil, headers)
	if err != nil {
		return nil, err
	}
	servicesResponseBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	serviceResponseInterfaces, err := jsonapi.UnmarshalManyPayload(
		servicesResponseBody,
		reflect.TypeOf(new(ServiceResponse)),
	)
	if err != nil {
		return nil, err
	}
	var servicesResponse []ServiceResponse
	for _, serviceResponseInterface := range serviceResponseInterfaces {
		servicesResponse = append(servicesResponse, serviceResponseInterface.(ServiceResponse))
	}
	return servicesResponse, nil
}
