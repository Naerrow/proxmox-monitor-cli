package proxmox

import "fmt"

type Node struct {
	Node   string  `json:"node"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	MaxCPU int     `json:"maxcpu"`
	Mem    int64   `json:"mem"`
	MaxMem int64   `json:"maxmem"`
	Uptime int64   `json:"uptime"`
}

func (c *Client) GetNodes() ([]Node, error) {
	var nodes []Node
	if err := c.get("/nodes", &nodes); err != nil {
		return nil, fmt.Errorf("노드 목록 조회 실패: %w", err)
	}
	return nodes, nil
}
