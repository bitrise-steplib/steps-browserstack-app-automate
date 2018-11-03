package appautomate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

const (
	ipaEndpoint          = "https://api-cloud.browserstack.com/app-automate/upload"
	testRunnerEndpoint   = "https://api-cloud.browserstack.com/app-automate/xcuitest/test-suite"
	buildEndPoint        = "https://api-cloud.browserstack.com/app-automate/xcuitest/build"
	buildSummaryEndPoint = "https://api-cloud.browserstack.com/app-automate/xcuitest/builds/"
)

// XCUITests ...
type XCUITests struct {
	AccessKey  stepconf.Secret
	UserName   string
	httpClient *http.Client
}

// Test ...
type Test struct {
	Devices    []string `json:"devices"`
	App        string   `json:"app"`
	DeviceLogs bool     `json:"deviceLogs"`
	TestSuite  string   `json:"testSuite"`
}

type uploadIPAResponse struct {
	AppURL string `json:"app_url"`
}

type uploadTestRunnerResponse struct {
	TestURL string `json:"test_url"`
}

type executeTestResponse struct {
	Message string `json:"message"`
	BuildID string `json:"build_id"`
}

// BuildResult ...
type BuildResult struct {
	Build Build
	Error error
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

// ExecuteTest execute your test for the provided app
func (x XCUITests) ExecuteTest(appURL, testURL string, deviceLogs bool, devices []string) (string, string, error) {
	data := Test{
		App:        appURL,
		TestSuite:  testURL,
		DeviceLogs: deviceLogs,
		Devices:    devices,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	body := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, buildEndPoint, body)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(x.UserName, string(x.AccessKey))
	req.Header.Set("Content-Type", "application/json")

	var response executeTestResponse
	_, err = x.performRequest(req, &response)
	if err != nil {
		return "", "", err
	}

	return response.Message, response.BuildID, nil
}

// ListenForTestComplete ...
func (x XCUITests) ListenForTestComplete(buildID string, ch chan BuildResult) {
	var isDone bool
	u, err := url.Parse(buildID)
	if err != nil {
		ch <- BuildResult{Build: Build{}, Error: err}
	}

	base, err := url.Parse(buildSummaryEndPoint)
	if err != nil {
		ch <- BuildResult{Build: Build{}, Error: err}
	}

	req, err := http.NewRequest(http.MethodGet, base.ResolveReference(u).String(), nil)
	if err != nil {
		ch <- BuildResult{Build: Build{}, Error: err}
	}

	req.SetBasicAuth(x.UserName, string(x.AccessKey))

	for !isDone {
		var response Build
		_, err = x.performRequest(req, &response)
		if err != nil {
			ch <- BuildResult{Build: Build{}, Error: err}
		}

		if response.Status == "done" {
			log.Warnf("Still: %+v", response)
			ch <- BuildResult{Build: response, Error: nil}
			isDone = true
		}

		if !isDone {
			time.Sleep(30 * time.Second)
		}
	}
}

//
//// Private methods

func (x XCUITests) performRequest(req *http.Request, requestResponse interface{}) ([]byte, error) {
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
