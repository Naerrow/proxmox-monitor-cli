package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxmoxURL   string
	ProxmoxToken string
	ProxmoxNode  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		ProxmoxURL:   os.Getenv("PROXMOX_URL"),
		ProxmoxToken: os.Getenv("PROXMOX_TOKEN"),
		ProxmoxNode:  os.Getenv("PROXMOX_NODE"),
	}

	if cfg.ProxmoxURL == "" || cfg.ProxmoxToken == "" {
		log.Fatal("❌ PROXMOX_URL과 PROXMOX_TOKEN을 설정해주세요 (.env 파일 확인)")
	}

	return cfg
}
