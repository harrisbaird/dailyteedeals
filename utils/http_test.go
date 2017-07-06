package utils_test

import (
	"os"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"io/ioutil"

	. "github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestHTTPGetString(t *testing.T) {
	defer gock.Off()

	testCases := []struct {
		name         string
		url          string
		wantResponse string
		wantErr      bool
	}{
		{"valid", "http://scrapyd:6900/status/jobID", "finished", false},
		{"404", "http://google.com/trigger404", "404 error", false},
		{"error", "invalid", "", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gock.New("http://scrapyd:6900").Get("/status/jobID").Reply(200).BodyString("finished")
			gock.New("http://google.com").Get("/trigger404").Reply(404).BodyString("404 error")

			resp, err := HTTPGetString(tt.url)
			st.Expect(t, err != nil, tt.wantErr)
			st.Expect(t, resp, tt.wantResponse)

			gock.Clean()
		})
	}
}

func TestHTTPGetToTempfile(t *testing.T) {
	defer gock.Off()

	testCases := []struct {
		name         string
		url          string
		wantResponse string
		wantErr      bool
	}{
		{"valid", "http://scrapyd:6900/status/jobID", "finished", false},
		{"404", "http://google.com/trigger404", "404 error", false},
		{"error", "invalid", "", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gock.New("http://scrapyd:6900").Get("/status/jobID").Reply(200).BodyString("finished")
			gock.New("http://google.com").Get("/trigger404").Reply(404).BodyString("404 error")

			tmp, err := HTTPGetToTempfile(tt.url)
			st.Expect(t, err != nil, tt.wantErr)

			if tmp != nil {
				b, err := ioutil.ReadFile(tmp.Name())
				st.Expect(t, err != nil, false)
				st.Expect(t, string(b), tt.wantResponse)

				tmp.Close()
				os.Remove(tmp.Name())
			}

			gock.Clean()
		})
	}
}
