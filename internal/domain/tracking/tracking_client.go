package tracking

import (
	"github.com/mister11/productive-cli/internal/domain"
	"time"
)

type TrackingClient interface {
	TrackFood(entries []FoodEntry) error
	TrackProject(entry *ProjectEntry) error
	Login(username string, password string) (*domain.LoginData, error)
	GetOrganizationMemberships() ([]domain.OrganizationMembershipData, error)
	SearchDeals(query string, day time.Time) ([]interface{}, error)
	SearchServices(query string, dealID string, day time.Time) ([]interface{}, error)
}
