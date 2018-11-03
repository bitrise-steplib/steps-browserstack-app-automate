package appautomate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/lunny/log"
)

const (
	ipaEndpoint        = "https://api-cloud.browserstack.com/app-automate/upload"
	testRunnerEndpoint = "https://api-cloud.browserstack.com/app-automate/xcuitest/test-suite"
)

// XCUITests ...
type XCUITests struct {
	AccessKey  stepconf.Secret
	UserName   string
	httpClient *http.Client
}

type uploadIPAResponse struct {
	AppURL string `json:"app_url"`
}

type uploadTestRunnerResponse struct {
	TestURL string `json:"test_url"`
}

//
//// Public methods

// NewXCUITests returns a new XCUITests with the provided Acces Key and User Name
func NewXCUITests(accesKey stepconf.Secret, userName string) *XCUITests {
	return &XCUITests{
		AccessKey:  accesKey,
		UserName:   userName,
		httpClient: &http.Client{},
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("akosbirmacher1", "xD4vYRrzWyzZFVFCawGd")
	return req, err
}

// UploadIPA uploads your iOS app (.ipa file) to the BrowserStack servers using the REST API.
func (x XCUITests) UploadIPA(ipa string) (string, error) {
	file, err := os.Open(ipa)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(ipa))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, ipaEndpoint, body)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(x.UserName, string(x.AccessKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	var response uploadIPAResponse
	_, err = x.performRequest(req, &response)
	if err != nil {
		return "", err
	}

	return response.AppURL, nil
}

// UploadTestRunner uploads the test zip file to the BrowserStack servers using the REST API
func (x XCUITests) UploadTestRunner(runner string) (string, error) {
	file, err := os.Open(runner)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(runner))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, testRunnerEndpoint, body)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(x.UserName, string(x.AccessKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	var response uploadTestRunnerResponse
	_, err = x.performRequest(req, &response)
	if err != nil {
		return "", err
	}

	return response.TestURL, nil
}

//
//// Private methods

func (x *XCUITests) performRequest(req *http.Request, requestResponse interface{}) ([]byte, error) {
	response, err := x.httpClient.Do(req)
	if err != nil {
		// On error, any Response can be ignored
		return nil, fmt.Errorf("failed to perform request, error: %s", err)
	}

	// The client must close the response body when finished with it
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Warnf("Failed to close response body, error: %s", cerr)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, error: %s", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusMultipleChoices {
		return nil, fmt.Errorf("Response status: %d - Body: %s", response.StatusCode, string(body))
	}

	// Parse JSON body
	if requestResponse != nil {
		if err := json.Unmarshal([]byte(body), &requestResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response (%s), error: %s", body, err)
		}
	}
	return body, nil
}
