package client

import (
	"bytes"
	"fmt"
	"github.com/mister11/productive-cli/internal/domain/config"
	"net/url"
	"reflect"
	"time"

	"github.com/mister11/productive-cli/internal/infrastructure/client/model"
	"github.com/mister11/productive-cli/internal/infrastructure/json"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/utils"
)

const baseURL = "https://api.productive.io/api/v2/"
const orgID = "1"

type ProductiveClient struct {
	client        HttpClient
	configManager config.Manager
}

func NewProductiveClient(configManager config.Manager) TrackingClient {
	client := &ProductiveClient{}
	client.client = NewHttpClient(baseURL)
	client.configManager = configManager
	return client
}

func (client *ProductiveClient) CreateFoodTimeEntry(day time.Time, userID string) {
	dayFormatted := utils.FormatDate(day)
	log.Info("Tracking food for " + dayFormatted)
	service := client.findFoodService(day)
	timeEntry := model.NewTimeEntry("", "30", userID, service, dayFormatted)
	client.createTimeEntry(timeEntry)
}

func (client *ProductiveClient) CreateProjectTimeEntry(
	service *model.Service,
	day time.Time,
	duration string,
	notes string,
	userID string,
) {
	dayFormatted := utils.FormatDate(day)
	durationFormatted, _ := utils.ParseTime(duration)
	timeEntry := model.NewTimeEntry(notes, durationFormatted, userID, service, dayFormatted)
	client.createTimeEntry(timeEntry)
}

func (client *ProductiveClient) createTimeEntry(timeEntry *model.TimeEntry) {
	jsonBytes := json.ToJsonEmbedded(timeEntry)
	body := client.client.Post("time_entries", bytes.NewReader(jsonBytes), client.getHeaders())
	defer body.Close()
}

func (client *ProductiveClient) findFoodService(day time.Time) *model.Service {
	_, service := client.findProjectInfo("Operations general", "Food", day)
	return service
}

func (client *ProductiveClient) findProjectInfo(
	dealName string,
	serviceName string,
	day time.Time,
) (*model.Deal, *model.Service) {
	deal := client.SearchDeals(dealName, day)[0].(*model.Deal)
	service := client.SearchServices(serviceName, deal.ID, day)[0].(*model.Service)
	return deal, service
}

func (client *ProductiveClient) GetOrganizationMembership() []model.OrganizationMembership {
	response := client.client.Get("organization_memberships", client.getHeaders())
	defer response.Close()

	orgMembershipInterfaces := json.FromJsonMany(response, reflect.TypeOf(new(model.OrganizationMembership)))

	var orgMemberships []model.OrganizationMembership
	for _, orgMembershipInterface := range orgMembershipInterfaces {
		orgMembership, ok := orgMembershipInterface.(*model.OrganizationMembership)
		if !ok {
			utils.ReportError("Failed to convert to OrganizationMembership", nil)
		}
		orgMemberships = append(orgMemberships, *orgMembership)
	}
	return orgMemberships
}

func (client *ProductiveClient) SearchDeals(query string, day time.Time) []interface{} {
	dayFormatted := utils.FormatDate(day)
	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
		url.QueryEscape(query), dayFormatted, dayFormatted)

	response := client.client.Get(uri, client.getHeaders())
	defer response.Close()

	dealInterfaces := json.FromJsonMany(response, reflect.TypeOf(new(model.Deal)))

	var deals []interface{}
	deals = append(deals, dealInterfaces...)

	return deals
}

func (client *ProductiveClient) SearchServices(query string, dealID string, day time.Time) []interface{} {
	dayFormatted := utils.FormatDate(day)
	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
		url.QueryEscape(query), dayFormatted, dayFormatted, dealID)

	resp := client.client.Get(uri, client.getHeaders())
	defer resp.Close()

	serviceInterfaces := json.FromJsonMany(resp, reflect.TypeOf(new(model.Service)))

	var services []interface{}
	services = append(services, serviceInterfaces...)

	return services
}

func (client *ProductiveClient) FindProjectInfo(dealQuery string, serviceQuery string, day time.Time) (*model.Deal, *model.Service) {
	deal := client.SearchDeals(dealQuery, day)[0].(*model.Deal)
	service := client.SearchServices(serviceQuery, deal.ID, day)[0].(*model.Service)
	return deal, service
}

func (client *ProductiveClient) getHeaders() map[string]string {
	headers := map[string]string{}
	headers["Content-Type"] = "application/vnd.api+json"
	headers["X-Auth-Token"] = client.configManager.GetToken()
	headers["X-Organization-Id"] = orgID
	return headers
}
