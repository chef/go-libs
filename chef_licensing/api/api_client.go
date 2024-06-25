package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	config "github.com/chef/go-libs/chef_licensing/config"
)

const (
	CLIENT_VERSION = "v1"
)

type APIClient struct {
	URL        string
	httpclient *http.Client
	Headers    map[string]string
}

var apiClient *APIClient

func (c *APIClient) BaseURL() string {
	return fmt.Sprintf("%s/%s/", c.URL, CLIENT_VERSION)
}

func NewClient() *APIClient {
	conf := config.GetConfig()

	apiClient = &APIClient{
		URL:        conf.LicenseServerURL,
		httpclient: &http.Client{},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return apiClient
}

func GetClient() *APIClient {
	if apiClient == nil {
		apiClient = NewClient()
		return apiClient
	} else {
		return apiClient
	}
}

func (c *APIClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

func (c *APIClient) doGETRequest(endpoint string, queryParams map[string]string) (*http.Response, error) {
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		q := urlObj.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		urlObj.RawQuery = q.Encode()
	}
	return c.doRequest("GET", urlObj.String(), nil)
}

func (c *APIClient) doPOSTRequest(endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	var err error

	if body != nil {
		reqBody, err = c.encodeJSON(body)
		if err != nil {
			return nil, err
		}

	}
	return c.doRequest("POST", endpoint, reqBody)
}

func (c *APIClient) doRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL() + endpoint

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	return c.httpclient.Do(req)
}

func (c *APIClient) decodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *APIClient) encodeJSON(v interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
