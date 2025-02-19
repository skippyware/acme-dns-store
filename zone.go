package main

import (
	"os"
	"text/template"
	"github.com/nrdcg/goacmedns"
	"math/rand"
	"bytes"
	"path/filepath"
	"time"
	"fmt"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateZone(config *Config, accounts map[string]goacmedns.Account) error {
	zoneTemplate := `; .{{ .DnsZone }} zone
{{ .DnsZone }}.                           IN       SOA         ns-local.{{ .DnsZone }}.  admin.ns-local.{{ .DnsZone }}. {{ .Serial }} 7200 3600 1209600 3600
{{ .DnsZone }}.                           IN       NS          ns-local.{{ .DnsZone }}.
{{ .DnsZone }}.                           IN       A           {{ .Gateway }}
*.{{ .DnsZone }}.                         IN       A           {{ .Gateway }}
{{- range $domain, $account := .Accounts }}
_acme-challenge.{{ $domain }}.  IN       CNAME       {{ $account.FullDomain }}.
{{- end }}
`

	tmpl, err := template.New("zone").Parse(zoneTemplate)
	if err != nil {
		return err
	}

	now := time.Now()
	now.Format("20060102")

	data := map[string]any{
		"Serial": fmt.Sprintf("%s%s", now.Format("20060102"), randString(2, "0123456789")),
		"DnsZone": config.DnsZone,
		"Gateway": config.HostGateway,
		"Accounts": accounts,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		filepath.Join(config.ZoneStoragePath, "db." + config.DnsZone),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(buf.Bytes())

	return nil
}

func randString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
