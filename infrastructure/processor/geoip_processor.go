package processor

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"1litw/sqlc"
)

const (
	ipAPIBatchURL = "http://ip-api.com/batch?fields=60121&lang=en"
	batchSize     = 100
	ticker        = 4 * time.Second
)

type GeoIPProcessor struct {
	clickRepo ClickRepository
}

type ClickRepository interface {
	GetUnprocessedClicks(ctx context.Context, limit int64) ([]sqlc.GetUnprocessedClicksRow, error)
	UpdateClickGeoInfo(ctx context.Context, arg sqlc.UpdateClickGeoInfoParams) error
}

func NewGeoIPProcessor(clickRepo ClickRepository) *GeoIPProcessor {
	return &GeoIPProcessor{
		clickRepo: clickRepo,
	}
}

func (p *GeoIPProcessor) Start() {
	log.Println("Starting GeoIPProcessor...")
	ticker := time.NewTicker(ticker)

	go func() {
		for {
			<-ticker.C
			p.processBatch()
		}
	}()
}

type ipAPIRequest struct {
	Query string `json:"query"`
}

type ipAPIResponse struct {
	Query      *string  `json:"query"`
	Status     *string  `json:"status"`
	Country    *string  `json:"country"`
	RegionName *string  `json:"regionName"`
	City       *string  `json:"city"`
	Lat        *float64 `json:"lat"`
	Lon        *float64 `json:"lon"`
	ISP        *string  `json:"isp"`
	AS         *string  `json:"as"`
}

func (p *GeoIPProcessor) processBatch() {
	clicks, err := p.clickRepo.GetUnprocessedClicks(context.Background(), batchSize)
	if err != nil {
		log.Printf("Error getting unprocessed clicks: %v", err)
		return
	}

	if len(clicks) == 0 {
		return
	}

	ipAddresses := make([]string, len(clicks))
	for i, click := range clicks {
		ipAddresses[i] = click.IPAddress.String
	}

	log.Printf("Processing %d clicks with IP addresses: %v", len(clicks), ipAddresses)

	jsonData, err := json.Marshal(ipAddresses)
	if err != nil {
		log.Printf("Error marshalling IP addresses: %v", err)
		return
	}

	resp, err := http.Post(ipAPIBatchURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error querying ip-api.com: %v", err)
		return
	}
	defer resp.Body.Close()

	var geoInfos []ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoInfos); err != nil {
		log.Printf("Error decoding ip-api.com response: %v", err)
		return
	}

	for _, info := range geoInfos {
		params := sqlc.UpdateClickGeoInfoParams{
			IsSuccess:  *info.Status == "success",
			IPAddress:  buildNullString(info.Query),
			Country:    buildNullString(info.Country),
			RegionName: buildNullString(info.RegionName),
			City:       buildNullString(info.City),
			Lat:        buildNullFloat64(info.Lat),
			Lon:        buildNullFloat64(info.Lon),
			Isp:        buildNullString(info.ISP),
			AsInfo:     buildNullString(info.AS),
		}

		if err := p.clickRepo.UpdateClickGeoInfo(context.Background(), params); err != nil {
			log.Printf("Error updating click geo info for IP %s: %v", *info.Query, err)
		}
	}
}

func buildNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func buildNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}
