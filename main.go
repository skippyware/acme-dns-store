package main

import (
	"os"
	"net/http"
	"path/filepath"
	"log"
)

func main() {
	config := LoadConfig()

	err := os.MkdirAll(config.ZoneStoragePath, 0755)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = os.MkdirAll(filepath.Dir(config.DnsStoragePath), 0755)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	api, err := NewApi(config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// http.HandleFunc("/hello", hello)
    mux.HandleFunc("GET /health", api.HealthCheck)
	mux.HandleFunc("GET /", api.FetchAll)
	mux.HandleFunc("GET /{domain}", api.Fetch)
	mux.HandleFunc("POST /{domain}", api.Put)

	http.ListenAndServe(config.ServerAddress, mux)
}
