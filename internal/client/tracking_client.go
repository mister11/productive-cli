package client

import (
	"github.com/mister11/productive-cli/internal/client/model"
	"time"
)

type TrackingClient interface {
	GetOrganizationMembership() []model.OrganizationMembership
	CreateTimeEntry(timeEntry *model.TimeEntry)
	SearchDeals(query string, dat time.Time) []interface{}
	SearchService(query string, dealID string, day time.Time) []interface{}
}
