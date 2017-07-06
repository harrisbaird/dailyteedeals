package backend

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/utils"
)

// ErrInvalidScrapydJob occurs when a scrapyd job doesb't exist.
var ErrInvalidScrapydJob = errors.New("Invalid Scrapyd job id")

// ScrapydSchedule schedules a scrapyd job with the given spider name.
// Returns a scrapyd job id and any error occured.
func ScrapydSchedule(spiderName string) (string, error) {
	resp, err := utils.HTTPPostForm(config.ScrapydURL()+"/schedule", url.Values{"spider": []string{spiderName}})
	if err != nil {
		return "", err
	}

	return utils.HTTPReadResponseString(resp)
}

// ScrapydIsFinished checks if a job has finished.
// Returns whether the job has finished the ErrInvalidScrapydJob
// error if the job doesn't exist.
func ScrapydIsFinished(jobID string) (bool, error) {
	url := fmt.Sprintf("%s/status/%s", config.ScrapydURL(), jobID)
	body, err := utils.HTTPGetString(url)
	if err != nil {
		return false, err
	}

	if body == "invalid" {
		return false, ErrInvalidScrapydJob
	}

	return body == "finished", nil
}

func ScrapydDownloadFeed(jobID string) (*os.File, error) {
	feedURL := fmt.Sprintf("%s/download/%s", config.ScrapydURL(), jobID)
	return utils.HTTPGetToTempfile(feedURL)
}

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
