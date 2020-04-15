package client

import (
	"github.com/mister11/productive-cli/internal/client/model"
)

type TrackingClient interface {
	GetOrganizationMembership() []model.OrganizationMembership
	CreateTimeEntry(timeEntry *model.TimeEntry)
	SearchDeals(query string, dayFormatted string) []interface{}
	SearchServices(query string, dealID string, dayFormatted string) []interface{}
}
