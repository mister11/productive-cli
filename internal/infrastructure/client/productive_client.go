package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/jsonapi"
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
	sessionRequest := &SessionRequest{
		ID:       "0",
		Email:    username,
		Password: password,
	}
	jsonBytes, err := json.ToJsonEmbedded(sessionRequest)
	if err != nil {
		return nil, err
	}
	log.Info("Creating session for user %s", username)
	body := client.client.Post("sessions", bytes.NewReader(jsonBytes), getDefaultHeaders())
	defer body.Close()
	var sessionResponse SessionResponse
	if err := jsonapi.UnmarshalPayload(body, &sessionResponse); err != nil {
		return nil, err
	}
	log.Info("Session created. Expiration date %s", sessionResponse.TokenExpirationDate)
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
		log.Info("Tracking food for %s", dayFormatted)
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
	if err != nil {
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
	if err != nil {
		return nil, err
	}
	response := client.client.Get(uri, headers)
	defer response.Close()

	dealInterfaces, err := json.FromJsonMany(response, reflect.TypeOf(new(Deal)))
	if err != nil {
		return nil, err
	}

	var deals []interface{}
	for _, dealInterface := range dealInterfaces {
		apiDeal := dealInterface.(*Deal)
		deal := &domain.Deal{
			ID:   apiDeal.ID,
			Name: apiDeal.Name,
		}
		deals = append(deals, deal)
	}

	return deals, nil
}

func (client *HttpProductiveClient) SearchServices(query string, dealID string, day time.Time) ([]interface{}, error) {
	dayFormatted := utils.FormatDate(day)
	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
		url.QueryEscape(query), dayFormatted, dayFormatted, dealID)

	headers, err := client.getHeaders()
	if err != nil {
		return nil, err
	}
	resp := client.client.Get(uri, headers)
	defer resp.Close()

	serviceInterfaces, err := json.FromJsonMany(resp, reflect.TypeOf(new(Service)))
	if err != nil {
		return nil, err
	}

	var services []interface{}
	for _, serviceInterface := range serviceInterfaces {
		apiService := serviceInterface.(*Service)
		service := &domain.Service{
			ID:   apiService.ID,
			Name: apiService.Name,
		}
		services = append(services, service)
	}

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
	if err != nil {
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

func (client *HttpProductiveClient) getHeaders() (map[string]string, error) {
	defaultHeaders := getDefaultHeaders()
	sessionData, err := client.userConfigManager.GetUserSession()
	if err != nil {
		return nil, errors.New("unable to read session token")
	}
	defaultHeaders["X-Auth-Token"] = sessionData.Token
	return defaultHeaders, nil
}

func getDefaultHeaders() map[string]string {
	headers := map[string]string{}
	headers["Content-Type"] = "application/vnd.api+json"
	headers["X-Organization-Id"] = orgID
	return headers
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
