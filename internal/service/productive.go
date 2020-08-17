package service

import (
	"errors"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/productive"
)

type ProductiveService struct {
	client *productive.Client
}

func NewProductiveService(client *productive.Client) *ProductiveService {
	return &ProductiveService{
		client: client,
	}
}

func (s *ProductiveService) CreateFoodTimeEntry(entry FoodEntry, sessionData *UserSessionData) error {
	foodService, err := s.FindFoodService(entry.Day, sessionData)
	if err != nil {
		return err
	}
	log.Debug("Creating food time entry for day %v", entry.Day)
	return s.client.TimeEntryService.CreateTimeEntry(
		"", "30", sessionData.PersonID, foodService, entry.Day, sessionData.Token,
	)
}

func (s *ProductiveService) CreateProjectTimeEntry(entry ProjectEntry, sessionData *UserSessionData) error {
	log.Debug("Creating project time entry for day %v", entry.Day)
	return s.client.TimeEntryService.CreateTimeEntry(
		formatNotes(entry.Notes), entry.Duration,
		sessionData.PersonID, entry.Service, entry.Day, sessionData.Token,
	)
}

func (s *ProductiveService) FindFoodService(day time.Time, sessionData *UserSessionData) (*productive.Service, error) {
	deals, err := s.client.DealService.SearchDeals("Operations general", day, &day, sessionData.Token)
	if err != nil {
		return nil, err
	}
	if len(deals) > 1 {
		return nil, errors.New("multiple 'Operations general' deals found")
	}
	if len(deals) == 0 {
		return nil, errors.New("no 'Operations general' deals found")
	}
	deal := deals[0]
	services, err := s.client.ServiceService.SearchServices("Food", deal.ID, day, day, sessionData.Token)
	if err != nil {
		return nil, err
	}
	if len(services) > 1 {
		return nil, errors.New("multiple 'Operations general/Food' services found")
	}
	if len(services) == 0 {
		return nil, errors.New("no 'Operations general/Food' services found")
	}
	return &services[0], nil
}

func (s *ProductiveService) FindSavedProject(project *TrackedProject, day time.Time, sessionData *UserSessionData) *TrackedProject {
	services, err := s.client.ServiceService.SearchServices(project.ServiceName, project.DealID, day, day, sessionData.Token)
	if err != nil {
		return nil
	}
	// sometimes, Productive will return multiple matches even if there's a exact match
	// so, we try and find exact match
	for _, service := range services {
		if service.Name == project.ServiceName {
			return project
		}
	}
	return nil
}

func (s *ProductiveService) FindDeals(dealQuery string, day time.Time, sessionData *UserSessionData) ([]productive.Deal, error) {
	deals, err := s.client.DealService.SearchDeals(dealQuery, day, &day, sessionData.Token)
	// end_date in Productive can be null so we cover this here
	// it can be some other error, but we assume that one for simplicity
	if len(deals) == 0 {
		deals, err = s.client.DealService.SearchDeals(dealQuery, day, nil, sessionData.Token)
		if err != nil {
			return nil, err
		}
	}
	return deals, err
}

func (s *ProductiveService) FindServices(serviceQuery string, deal *productive.Deal, day time.Time, sessionData *UserSessionData) ([]productive.Service, error) {
	return s.client.ServiceService.SearchServices(serviceQuery, deal.ID, day, day, sessionData.Token)
}

func (s *ProductiveService) Login(username string, password string) (*productive.SessionResponse, error) {
	return s.client.SessionService.Login(username, password)
}

func (s *ProductiveService) GetOrganizationMemberships(sessionData *UserSessionData) ([]productive.OrganizationMembership, error) {
	return s.client.OrganizationMembershipService.FetchAll(sessionData.Token)
}

func formatNotes(notes []string) string {
	if len(notes) == 0 {
		return ""
	}
	var notesHTML strings.Builder
	notesHTML.WriteString("<ul>")
	for _, note := range notes {
		notesHTML.WriteString("<li>" + note + "</li>")
	}
	notesHTML.WriteString("</ul>")
	return notesHTML.String()
}
