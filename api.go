package main

import (
	"context"
	"log"
	"net/http"
	"fmt"
	"sync"
	"encoding/json"
	"errors"
	"io"
	// "time"
	"github.com/nrdcg/goacmedns"
	"github.com/nrdcg/goacmedns/storage"
)

type Api struct {
	mutex sync.RWMutex
	storage goacmedns.Storage
	config *Config
}

func NewApi(config *Config) (*Api, error) {
	ctx := context.Background()
	st := storage.NewFile(config.DnsStoragePath, 0o600)

	api := &Api{
		storage: st,
		config: config,
	}

	accounts, _ := api.storage.FetchAll(ctx)
	err := GenerateZone(api.config, accounts)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func (a *Api) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "READY\n")
}

func (a *Api) FetchAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	accounts, err := a.storage.FetchAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	}
}

func (a *Api) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	domain := r.PathValue("domain")

	account, err := a.storage.Fetch(ctx, domain)
	if err != nil {
		if errors.Is(err, storage.ErrDomainNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(account)
	}
}

func (a *Api) Put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	a.mutex.Lock()
	defer a.mutex.Unlock()
	domain := r.PathValue("domain")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	var account goacmedns.Account

	err = json.Unmarshal(body, &account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	a.storage.Put(ctx, domain, account)
	err = a.storage.Save(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	accounts, _ := a.storage.FetchAll(ctx)
	err = GenerateZone(a.config, accounts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// time.Sleep(10 * time.Second)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
