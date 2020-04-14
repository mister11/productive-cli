package client

import (
	"io"
	"net/http"

	"github.com/mister11/productive-cli/internal/utils"
)

type HttpClient struct {
	baseURL string
	client  *http.Client
}

func NewHttpClient(baseURL string) HttpClient {
	client := HttpClient{}
	client.baseURL = baseURL
	client.client = &http.Client{}
	return client
}

func (client *HttpClient) Get(uri string, headers map[string]string) io.ReadCloser {
	req, err := http.NewRequest("GET", client.baseURL+uri, nil)
	if err != nil {
		utils.ReportError("Failed to create a request.", err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := client.client.Do(req)
	if err != nil {
		utils.ReportError("Request "+req.RequestURI+" failed", err)
	}

	return resp.Body
}

func (client *HttpClient) Post(uri string, body io.Reader, headers map[string]string) io.ReadCloser {
	req, err := http.NewRequest("POST", client.baseURL+uri, body)
	if err != nil {
		utils.ReportError("Failed to create a request", err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := client.client.Do(req)
	if err != nil {
		utils.ReportError("Request "+req.RequestURI+" failed", err)
	}

	return resp.Body
}
