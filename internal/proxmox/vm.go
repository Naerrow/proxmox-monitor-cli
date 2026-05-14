package proxmox

type VM struct {
	VMID   int     `json:"vmid"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Mem    int64   `json:"mem"`
	MaxMem int64   `json:"maxmem"`
	Node   string  `json:"node"`
}

// TODO: GetVMs(), StartVM(), StopVM(), DeleteVM() 구현 예정
