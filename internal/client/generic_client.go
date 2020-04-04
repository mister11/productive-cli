package client

import (
	"io"
	"net/http"

	"github.com/mister11/productive-cli/internal/utils"
)

type GenericClient struct {
	baseURL string
	headers map[string]string
	client  *http.Client
}

func NewGenericClient(baseURL string, headers map[string]string) GenericClient {
	client := GenericClient{}
	client.baseURL = baseURL
	client.headers = headers
	client.client = &http.Client{}
	return client
}

func (client *GenericClient) Get(uri string) io.ReadCloser {
	req, err := http.NewRequest("GET", client.baseURL+uri, nil)
	if err != nil {
		utils.ReportError("Failed to create a request.", err)
	}
	for key, value := range client.headers {
		req.Header.Add(key, value)
	}
	resp, err := client.client.Do(req)
	if err != nil {
		utils.ReportError("Request "+req.RequestURI+" failed", err)
	}

	return resp.Body
}

func (client *GenericClient) Post(uri string, body io.Reader) io.ReadCloser {
	req, err := http.NewRequest("POST", client.baseURL+uri, body)
	if err != nil {
		utils.ReportError("Failed to create a request", err)
	}
	for key, value := range client.headers {
		req.Header.Add(key, value)
	}
	resp, err := client.client.Do(req)
	if err != nil {
		utils.ReportError("Request "+req.RequestURI+" failed", err)
	}

	return resp.Body
}
