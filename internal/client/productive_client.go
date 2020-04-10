package client

import (
	"bytes"
	"fmt"
	"github.com/mister11/productive-cli/internal/client/model"
	"net/url"
	"reflect"
	"time"

	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/json"
	"github.com/mister11/productive-cli/internal/utils"
)

const baseURL = "https://api.productive.io/api/v2/"
const orgID = "1"

type ProductiveClient struct {
	client GenericClient
	configManager config.ConfigManager
}

func NewProductiveClient(configManager config.ConfigManager) *ProductiveClient {
	client := &ProductiveClient{}
	client.client = NewGenericClient(baseURL)
	client.configManager = configManager
	return client
}

func (client *ProductiveClient) CreateTimeEntry(timeEntry *model.TimeEntry) {
	jsonBytes := json.ToJsonEmbedded(timeEntry)
	body := client.client.Post("time_entries", bytes.NewReader(jsonBytes), client.getHeaders())
	defer body.Close()
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
	dayFormatted := datetime.Format(day)

	uri := fmt.Sprintf("deals?filter[query]=%s&filter[date][lt_eq]=%s&filter[end_date][gt_eq]=%s",
		url.QueryEscape(query), dayFormatted, dayFormatted)

	response := client.client.Get(uri, client.getHeaders())
	defer response.Close()

	dealInterfaces := json.FromJsonMany(response, reflect.TypeOf(new(model.Deal)))

	var deals []interface{}
	deals = append(deals, dealInterfaces...)

	return deals
}

func (client *ProductiveClient) SearchService(query string, dealID string, day time.Time) []interface{} {
	dayFormatted := datetime.Format(day)

	uri := fmt.Sprintf(`services?filter[name]=%s&filter[after]=%s&filter[before]=%s&filter[deal_id]=%s`,
		url.QueryEscape(query), dayFormatted, dayFormatted, dealID)

	resp := client.client.Get(uri, client.getHeaders())
	defer resp.Close()

	serviceInterfaces := json.FromJsonMany(resp, reflect.TypeOf(new(model.Service)))

	var services []interface{}
	services = append(services, serviceInterfaces...)

	return services
}

func (client *ProductiveClient) getHeaders() map[string]string {
	headers := map[string]string{}
	headers["Content-Type"] = "application/vnd.api+json"
	headers["X-Auth-Token"] = client.configManager.GetToken()
	headers["X-Organization-Id"] = orgID
	return headers
}
