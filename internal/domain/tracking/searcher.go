package tracking

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/input"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/client/model"
	"time"
)

type Searcher interface {
	SearchProjects(day time.Time) (*model.Deal, error)
	SearchServices(day time.Time, project *model.Deal) (*model.Service, error)
	SearchService(serviceName string, dealID string, day time.Time) (*model.Service, error)
}

type productiveSearcher struct {
	prompt         input.Prompt
	trackingClient client.TrackingClient
}

func NewProjectSearcher(
	prompt input.Prompt,
	trackingClient client.TrackingClient,
) Searcher {
	return &productiveSearcher{
		prompt:         prompt,
		trackingClient: trackingClient,
	}
}

func (searcher *productiveSearcher) SearchProjects(day time.Time) (*model.Deal, error) {
	projectQuery, err := searcher.prompt.Input("Enter project name")
	if err != nil {
		return nil, err
	}
	deals := searcher.trackingClient.SearchDeals(projectQuery, day)
	return searcher.prompt.SelectOne("Select project", deals).(*model.Deal), nil
}

func (searcher *productiveSearcher) SearchServices(day time.Time, project *model.Deal) (*model.Service, error) {
	serviceQuery, err := searcher.prompt.Input("Enter service name")
	if err != nil {
		return nil, err
	}
	services := searcher.trackingClient.SearchServices(serviceQuery, project.ID, day)
	return searcher.prompt.SelectOne("Select service", services).(*model.Service), nil
}

func (searcher *productiveSearcher) SearchService(serviceName string, dealID string, day time.Time) (*model.Service, error) {
	services := searcher.trackingClient.SearchServices(serviceName, dealID, day)
	if len(services) != 1 {
		return nil, errors.New("multiple service return when 1 expected")
	}
	return services[0].(*model.Service), nil
}
