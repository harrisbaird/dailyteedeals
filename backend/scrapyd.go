package backend

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/utils"
)

// ErrInvalidScrapydJob occurs when a scrapyd job doesb't exist.
var ErrInvalidScrapydJob = errors.New("Invalid Scrapyd job id")

// ScrapydSchedule schedules a scrapyd job with the given spider name.
// Returns a scrapyd job id and any error occured.
func ScrapydSchedule(spiderName string) (string, error) {
	resp, err := utils.HTTPPostForm(fmt.Sprintf("http://%s/schedule", config.App.ScrapydAddr), url.Values{"spider": []string{spiderName}})
	if err != nil {
		return "", err
	}

	return utils.HTTPReadResponseString(resp)
}

// ScrapydIsFinished checks if a job has finished.
// Returns whether the job has finished the ErrInvalidScrapydJob
// error if the job doesn't exist.
func ScrapydIsFinished(jobID string) (bool, error) {
	url := fmt.Sprintf("http://%s/status/%s", config.App.ScrapydAddr, jobID)
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
	feedURL := fmt.Sprintf("http://%s/download/%s", config.App.ScrapydAddr, jobID)
	return utils.HTTPGetToTempfile(feedURL)
}
