package client

import (
	"time"

	"github.com/mister11/productive-cli/internal/infrastructure/client/model"
)

type TrackingClient interface {
	CreateFoodTimeEntry(day time.Time, userID string)
	CreateProjectTimeEntry(service *model.Service, day time.Time, duration string, notes string, userID string)
	GetOrganizationMembership() []model.OrganizationMembership
	SearchDeals(query string, day time.Time) []interface{}
	SearchServices(query string, dealID string, day time.Time) []interface{}
	FindProjectInfo(dealQuery string, serviceQuery string, day time.Time) (*model.Deal, *model.Service)
}
