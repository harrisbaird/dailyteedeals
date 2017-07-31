package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

var (
	ErrNotSuccessful = errors.New("HTTP Get request was not successful")
	httpClient       = pester.New()
)

func init() {
	httpClient.MaxRetries = 5
	httpClient.Backoff = pester.ExponentialBackoff
}

func SetHTTPTestMode() {
	httpClient.MaxRetries = 0
	httpClient.Backoff = pester.DefaultBackoff
}

// TODO: Create common tempdir?

func HTTPGetToTempfile(url string) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "download")
	if err != nil {
		return nil, err
	}

	resp, err := httpGet(url)
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, err
	}

	defer resp.Body.Close()

	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, err
	}

	tmpfile.Seek(0, 0)
	return tmpfile, nil
}

func HTTPGetBytes(url string) ([]byte, error) {
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return HTTPReadResponse(resp)
}

func HTTPGetString(url string) (string, error) {
	resp, err := httpGet(url)
	if err != nil {
		return "", err
	}
	return HTTPReadResponseString(resp)
}

func HTTPReadResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HTTPReadResponseString(resp *http.Response) (string, error) {
	b, err := HTTPReadResponse(resp)
	return string(b), err
}

func HTTPPostForm(url string, data url.Values) (*http.Response, error) {
	return httpClient.PostForm(url, data)
}

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func httpGet(url string) (*http.Response, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
