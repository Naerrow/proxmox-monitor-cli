package proxmox

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL    string
	Token      string
	Node       string
	httpClient *http.Client
}

func NewClient(url, token, node string) *Client {
	return &Client{
		BaseURL: url,
		Token:   token,
		Node:    node,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api2/json%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)
	return req, nil
}
