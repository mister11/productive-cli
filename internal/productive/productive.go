package productive

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://api.productive.io/api/v2/"

type client struct {
	httpClient *http.Client
	baseUrl    *url.URL

	SessionService                *sessionService
	DealService                   *dealsService
	ServiceService                *serviceService
	OrganizationMembershipService *organizationMembershipService
}

func NewClient(httpClient *http.Client) *client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	encodedURL, _ := url.Parse(baseURL)

	client := &client{
		httpClient: httpClient,
		baseUrl:    encodedURL,
	}

	client.SessionService = newSessionService(client)
	client.DealService = newDealsService(client)
	client.ServiceService = newServiceService(client)
	client.OrganizationMembershipService = newOrganizationMembershipService(client)
	return client
}

func (client *client) NewRequest(
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

func (client *client) Do(request *http.Request) (io.Reader, error) {
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
