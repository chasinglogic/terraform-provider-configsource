package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	token   string
	baseURL string
	client  http.Client
}

type requestSpec struct {
	method string
	url    string
	body   interface{}
}

func New(token, baseURL string) *Client {
	return &Client{token: token, baseURL: baseURL, client: http.Client{}}
}

func (c *Client) Do(ctx context.Context, spec requestSpec, output interface{}) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s", c.baseURL, spec.url)

	req, err := http.NewRequestWithContext(ctx, spec.method, fullURL, nil)
	req.Header.Add("Accepts", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	httpResp, err := c.client.Do(req)
	if err != nil {
		return httpResp, nil
	}

	if output != nil {
		err = json.NewDecoder(httpResp.Body).Decode(&output)
	}

	return httpResp, err
}

func (c *Client) GetConfigValue(ctx context.Context, environmentName, key string) (ConfigValue, error) {
	var cv ConfigValue
	_, err := c.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("/api/v1/config-values/%s/%s", environmentName, key),
	}, &cv)
	return cv, err
}
