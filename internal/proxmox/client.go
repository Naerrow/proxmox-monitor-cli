package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	Token      string
	Node       string
	httpClient *http.Client
}

type apiResponse struct {
	Data json.RawMessage `json:"data"`
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

func (c *Client) get(path string, result interface{}) error {
	url := fmt.Sprintf("%s/api2/json%s", c.BaseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "PVEAPIToken="+c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API 오류: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("응답 읽기 실패: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("JSON 파싱 실패: %w", err)
	}

	return json.Unmarshal(apiResp.Data, result)
}
