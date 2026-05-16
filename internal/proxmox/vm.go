package proxmox

import (
	"fmt"
	"net/http"
)

type VM struct {
	VMID   int     `json:"vmid"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Mem    int64   `json:"mem"`
	MaxMem int64   `json:"maxmem"`
	Node   string  `json:"node"`
}

func (c *Client) GetVMs() ([]VM, error) {
	var vms []VM
	if err := c.get(fmt.Sprintf("/nodes/%s/qemu", c.Node), &vms); err != nil {
		return nil, fmt.Errorf("VM 목록 조회 실패: %w", err)
	}
	return vms, nil
}

func (c *Client) StartVM(vmid string) error {
	return c.post(fmt.Sprintf("/nodes/%s/qemu/%s/status/start", c.Node, vmid))
}

func (c *Client) StopVM(vmid string) error {
	return c.post(fmt.Sprintf("/nodes/%s/qemu/%s/status/stop", c.Node, vmid))
}

func (c *Client) DeleteVM(vmid string) error {
	url := fmt.Sprintf("%s/api2/json/nodes/%s/qemu/%s", c.BaseURL, c.Node, vmid)

	req, err := http.NewRequest("DELETE", url, nil)
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
		return fmt.Errorf("VM 삭제 실패: %s", resp.Status)
	}
	return nil
}

func (c *Client) post(path string) error {
	url := fmt.Sprintf("%s/api2/json%s", c.BaseURL, path)

	req, err := http.NewRequest("POST", url, nil)
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
	return nil
}
