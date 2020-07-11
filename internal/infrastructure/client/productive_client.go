package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mister11/productive-cli/internal/domain"
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/json"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/utils"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const baseURL = "https://api.productive.io/api/v2/"
const orgID = "1"

type HttpProductiveClient struct {
	client                HttpClient
	userConfigManager     UserSessionManager
	projectsConfigManager domain.TrackedProjectManager
}

func NewProductiveClient(
	userConfigManager UserSessionManager,
	projectsConfigManager domain.TrackedProjectManager,
) *HttpProductiveClient {
	client := &HttpProductiveClient{
		client:                NewHttpClient(baseURL),
		userConfigManager:     userConfigManager,
		projectsConfigManager: projectsConfigManager,
	}
	return client
}

func (client *HttpProductiveClient) Login(username string, password string) (*domain.LoginData, error) {
	headers, err := client.getHeaders()
	if err == nil {
		return nil, err
	}
	sessionRequest := &SessionRequest{
		ID:       "0",
		Email:    username,
		Password: password,
	}
	jsonBytes, err := json.ToJsonEmbedded(sessionRequest)
	if err != nil {
		return nil, err
	}
	log.Debug("Creating session for user %s", username)
	body := client.client.Post("sessions", bytes.NewReader(jsonBytes), headers)
	defer body.Close()
	var sessionResponse SessionResponse
	if err := json.FromJsonOne(body, sessionResponse); err != nil {
		return nil, err
	}
	log.Debug("Session created. Expiration date %s", sessionResponse.TokenExpirationDate)
	return &domain.LoginData{
		Token:               sessionResponse.Token,
		TokenExpirationDate: sessionResponse.TokenExpirationDate,
	}, nil
}

func (client *HttpProductiveClient) TrackFood(entries []tracking.FoodEntry) error {
	userConfig, err := client.userConfigManager.GetUserSession()
	if err != nil {
		return err
	}
	userID := userConfig.PersonID
	for _, entry := range entries {
		dayFormatted := utils.FormatDate(entry.Day)
		log.Info("Tracking food for " + dayFormatted)
		service, err := client.findFoodService(entry.Day)
		if err != nil {
			return err
		}
		timeEntry := NewTimeEntry("", "30", userID, service, dayFormatted)
		if err := client.createTimeEntry(timeEntry); err != nil {
			return err
		}
	}
	return nil
}

func (client *HttpProductiveClient) TrackProject(entry *tracking.ProjectEntry) error {
	dayFormatted := utils.FormatDate(entry.Day)
	durationFormatted, _ := utils.ParseTime(entry.Duration)
	notesFormatted := createNotes(entry.Notes)
	userConfig, err := client.userConfigManager.GetUserSession()
	if err != nil {
		return err
	}
	userID := userConfig.PersonID
	service := &Service{
		ID:   entry.Service.ID,
		Name: entry.Service.Name,
	}
	timeEntry := NewTimeEntry(notesFormatted, durationFormatted, userID, service, dayFormatted)
	return client.createTimeEntry(timeEntry)
}

func (client *HttpProductiveClient) GetOrganizationMemberships() ([]domain.OrganizationMembershipData, error) {
	headers, err := client.getHeaders()
	if err == nil {
		return nil, err
	}
	response := client.client.Get("organization_memberships", headers)
	defer response.Close()

	orgMembershipInterfaces, err := json.FromJsonMany(response, reflect.TypeOf(new(OrganizationMembership)))
	if err != nil {
		return nil, err
	}

	var orgMemberships []domain.OrganizationMembershipData
	for _, orgMembershipInterface := range orgMembershipInterfaces {
		orgMembership, ok := orgMembershipInterface.(*OrganizationMembership)
		if !ok {
			utils.ReportError("Failed to convert to OrganizationMembership", nil)
		}
		orgMemberships = append(orgMemberships, domain.OrganizationMembershipData{
			PersonID: orgMembership.User.ID,
		})
	}
	return orgMemberships, nil
}

func (client *HttpProductiveClient) SearchDeals(query string, day time.Time) ([]interface{}, error) {
	dayFormatted := utils.FormatDate(day)
	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
		url.QueryEscape(query), dayFormatted, dayFormatted)
	headers, err := client.getHeaders()
	if err == nil {
		return nil, err
	}
	response := client.client.Get(uri, headers)
	defer response.Close()

	dealInterfaces, err := json.FromJsonMany(response, reflect.TypeOf(new(Deal)))
	if err == nil {
		return nil, err
	}

	var deals []interface{}
	deals = append(deals, dealInterfaces...)

	return deals, nil
}

func (client *HttpProductiveClient) SearchServices(query string, dealID string, day time.Time) ([]interface{}, error) {
	dayFormatted := utils.FormatDate(day)
	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
		url.QueryEscape(query), dayFormatted, dayFormatted, dealID)

	headers, err := client.getHeaders()
	if err == nil {
		return nil, err
	}
	resp := client.client.Get(uri, headers)
	defer resp.Close()

	serviceInterfaces, err := json.FromJsonMany(resp, reflect.TypeOf(new(Service)))
	if err == nil {
		return nil, err
	}

	var services []interface{}
	services = append(services, serviceInterfaces...)

	return services, nil
}

func (client *HttpProductiveClient) findFoodService(day time.Time) (*Service, error) {
	project, err := client.findProjectInfo("Operations general", "Food", day)
	if err != nil {
		return nil, err
	}
	return project.service, nil
}

func (client *HttpProductiveClient) findProjectInfo(
	dealName string,
	serviceName string,
	day time.Time,
) (*Project, error) {
	deals, err := client.SearchDeals(dealName, day)
	if err != nil {
		return nil, err
	}
	if len(deals) != 1 {
		return nil, errors.New("Found none of multiple deals with name " + dealName)
	}
	deal := deals[0].(*Deal)
	services, err := client.SearchServices(serviceName, deal.ID, day)
	if err != nil {
		return nil, err
	}
	if len(services) != 1 {
		return nil, errors.New("Found none of multiple services with name " + serviceName)
	}
	service := services[0].(*Service)
	return &Project{deal: deal, service: service}, nil
}

func (client *HttpProductiveClient) createTimeEntry(timeEntry *TimeEntry) error {
	headers, err := client.getHeaders()
	if err == nil {
		return err
	}
	jsonBytes, err := json.ToJsonEmbedded(timeEntry)
	if err != nil {
		return err
	}
	body := client.client.Post("time_entries", bytes.NewReader(jsonBytes), headers)
	defer body.Close()
	return nil
}

//func (client *HttpProductiveClient) CreateFoodTimeEntry(day time.Time, userID string) {
//	dayFormatted := utils.FormatDate(day)
//	log.Info("Tracking food for " + dayFormatted)
//	service := client.findFoodService(day)
//	timeEntry := model.NewTimeEntry("", "30", userID, service, dayFormatted)
//	client.createTimeEntry(timeEntry)
//}
//
//func (client *HttpProductiveClient) CreateProjectTimeEntry(
//	service *model.Service,
//	day time.Time,
//	duration string,
//	notes string,
//	userID string,
//) {
//	dayFormatted := utils.FormatDate(day)
//	durationFormatted, _ := utils.ParseTime(duration)
//	timeEntry := model.NewTimeEntry(notes, durationFormatted, userID, service, dayFormatted)
//	client.createTimeEntry(timeEntry)
//}
//

//func (client *HttpProductiveClient) GetOrganizationMembership() []model.OrganizationMembership {
//	response := client.client.Get("organization_memberships", client.getHeaders())
//	defer response.Close()
//
//	orgMembershipInterfaces := json.FromJsonMany(response, reflect.TypeOf(new(model.OrganizationMembership)))
//
//	var orgMemberships []model.OrganizationMembership
//	for _, orgMembershipInterface := range orgMembershipInterfaces {
//		orgMembership, ok := orgMembershipInterface.(*model.OrganizationMembership)
//		if !ok {
//			utils.ReportError("Failed to convert to OrganizationMembership", nil)
//		}f
//		orgMemberships = append(orgMemberships, *orgMembership)
//	}
//	return orgMemberships
//}
//
//func (client *HttpProductiveClient) SearchDeals(query string, day time.Time) []interface{} {
//	dayFormatted := utils.FormatDate(day)
//	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
//		url.QueryEscape(query), dayFormatted, dayFormatted)
//
//	response := client.client.Get(uri, client.getHeaders())
//	defer response.Close()
//
//	dealInterfaces := json.FromJsonMany(response, reflect.TypeOf(new(model.Deal)))
//
//	var deals []interface{}
//	deals = append(deals, dealInterfaces...)
//
//	return deals
//}
//
//func (client *HttpProductiveClient) SearchServices(query string, dealID string, day time.Time) []interface{} {
//	dayFormatted := utils.FormatDate(day)
//	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
//		url.QueryEscape(query), dayFormatted, dayFormatted, dealID)
//
//	resp := client.client.Get(uri, client.getHeaders())
//	defer resp.Close()
//
//	serviceInterfaces := json.FromJsonMany(resp, reflect.TypeOf(new(model.Service)))
//
//	var services []interface{}
//	services = append(services, serviceInterfaces...)
//
//	return services
//}
//
//func (client *HttpProductiveClient) FindProjectInfo(dealQuery string, serviceQuery string, day time.Time) (*model.Deal, *model.Service) {
//	deal := client.SearchDeals(dealQuery, day)[0].(*model.Deal)
//	service := client.SearchServices(serviceQuery, deal.ID, day)[0].(*model.Service)
//	return deal, service
//}

func (client *HttpProductiveClient) getHeaders() (map[string]string, error) {
	sessionData, err := client.userConfigManager.GetUserSession()
	if err != nil {
		return nil, errors.New("unable to read session token")
	}
	headers := map[string]string{}
	headers["Content-Type"] = "application/vnd.api+json"
	headers["X-Auth-Token"] = sessionData.Token
	headers["X-Organization-Id"] = orgID
	return headers, nil
}

func createNotes(notes []string) string {
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
