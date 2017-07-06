package backend_test

import (
	"os"
	"testing"

	. "github.com/harrisbaird/dailyteedeals/backend"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

const (
	spiderName  = "myspidername"
	jobID       = "myjobid"
	validPath   = "/status/" + jobID
	invalidPath = "/404"
)

func TestScrapydSchedule(t *testing.T) {
	defer gock.Off()
	gock.New(config.ScrapydURL()).Post("/schedule").BodyString("spider=" + spiderName).Reply(200).BodyString(jobID)

	newJobID, err := ScrapydSchedule(spiderName)

	st.Expect(t, err, nil)
	st.Expect(t, jobID, newJobID)
	st.Expect(t, gock.IsDone(), true)
}

func TestIsFinished(t *testing.T) {
	defer gock.Off()
	testCases := []struct {
		name         string
		url          string
		response     string
		wantFinished bool
		wantErr      error
	}{
		{"finished", validPath, "finished", true, nil},
		{"pending", validPath, "pending", false, nil},
		{"running", validPath, "running", false, nil},
		{"invalid", validPath, "invalid", false, ErrInvalidScrapydJob},
	}

	for _, tt := range testCases {
		t.Run(tt.response, func(t *testing.T) {
			gock.New(config.ScrapydURL()).Get(validPath).Reply(200).BodyString(tt.response)
			finished, err := ScrapydIsFinished(jobID)
			st.Expect(t, finished, tt.wantFinished)
			st.Expect(t, err, tt.wantErr)
			st.Expect(t, gock.IsDone(), true)
			gock.Clean()
		})
	}
}

func TestDownloadFeed(t *testing.T) {
	defer gock.Off()

	fixturePath := utils.ProjectRootPath("backend", "testdata", "feed", "new_product.json")
	fixture, err := os.Open(fixturePath)
	st.Expect(t, err, nil)

	gock.New(config.ScrapydURL()).Get("/download/" + jobID).Reply(200).File(fixturePath)

	feed, err := ScrapydDownloadFeed(jobID)

	st.Expect(t, err, nil)
	st.Expect(t, filesIdentical(fixture, feed), true)
	st.Expect(t, gock.IsDone(), true)

	if feed != nil {
		feed.Close()
		os.Remove(feed.Name())
	}
}
