package productive

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://api.productive.io/api/v2/"
const orgID = "1"

type Client struct {
	httpClient *http.Client
	baseUrl    *url.URL

	SessionService                *sessionService
	DealService                   *dealsService
	ServiceService                *serviceService
	OrganizationMembershipService *organizationMembershipService
	TimeEntryService              *timeEntryService
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	encodedURL, _ := url.Parse(baseURL)

	client := &Client{
		httpClient: httpClient,
		baseUrl:    encodedURL,
	}

	client.SessionService = newSessionService(client)
	client.DealService = newDealsService(client)
	client.ServiceService = newServiceService(client)
	client.OrganizationMembershipService = newOrganizationMembershipService(client)
	client.TimeEntryService = newTimeEntryService(client)
	return client
}

func (client *Client) NewRequest(
	method string,
	urlPath string,
	body interface{},
	headers map[string]string,
) (*http.Request, error) {
	url, err := client.baseUrl.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		json, err := ToJsonEmbedded(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func (client *Client) Do(request *http.Request) (io.Reader, error) {
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, Unauthorized
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	return bytes.NewReader(responseBody), nil
}

func getHeaders(token string) map[string]string {
	defaultHeaders := getDefaultHeaders()
	defaultHeaders["X-Auth-Token"] = token
	return defaultHeaders
}

func getDefaultHeaders() map[string]string {
	headers := map[string]string{}
	headers["Content-Type"] = "application/vnd.api+json"
	headers["X-Organization-Id"] = orgID
	return headers
}
