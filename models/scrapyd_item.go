package models

import (
	"strconv"
	"time"
)

type ScrapydItem struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	URL          string         `json:"url"`
	ArtistName   string         `json:"artist_name"`
	ArtistUrls   []string       `json:"artist_urls"`
	Prices       map[string]int `json:"prices"`
	ImageURL     string         `json:"image_url"`
	Tags         []string       `json:"tags"`
	FabricColors []string       `json:"fabric_colors"`
	Active       bool           `json:"active"`
	Deal         bool           `json:"deal"`
	LastChance   bool           `json:"last_chance"`
	Valid        bool           `json:"valid"`
	ExpiresAt    time.Time      `json:"expires_at"`
}

func (item *ScrapydItem) StringPrices() map[string]string {
	out := make(map[string]string)
	for k, v := range item.Prices {
		out[k] = strconv.Itoa(v)
	}
	return out
}
