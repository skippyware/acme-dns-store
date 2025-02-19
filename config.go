package main

import (
	"os"
)

type Config struct {
	ServerAddress string
	DnsZone string
	HostGateway string
	DnsStoragePath string
	ZoneStoragePath string
	ResolversStoragePath string
	AcmeDnsAddress string
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func LoadConfig() (*Config) {
	return &Config{
		ServerAddress: getEnv("HTTP_LISTEN", ":8080"),
		DnsZone: getEnv("DNS_ZONE", "test"),
		HostGateway: getEnv("HOST_GATEWAY", "127.0.0.1"),
		DnsStoragePath: getEnv("DNS_STORAGE_PATH", "/data/accounts/acme-dns-accounts.json"),
		ZoneStoragePath: getEnv("ZONE_STORAGE_PATH", "/data/zones/dynamic"),
	}
}
