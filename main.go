package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-browserstack-app-automate/appautomate"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/briandowns/spinner"
)

// TestType ...
type TestType string

// const ...
const (
	XCUITests TestType = "XCUITests"
)

// Config ...
type Config struct {
	AccesKey            stepconf.Secret `env:"access_key,required"`
	UserName            string          `env:"user_name,required"`
	TestType            string          `env:"test_type,opt[XCUITests]"`
	XCUITestsIPA        string          `env:"xcuitests_ipa"`
	XCUITestsRunner     string          `env:"xcuitests_runner"`
	XCUITestsDeviceLogs bool            `env:"xcuitests_device_logs"`
	XCUITestsDevices    string          `env:"xcuitests_devices"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	switch cfg.TestType {
	case string(XCUITests):
		xcuitests := appautomate.NewXCUITests(cfg.AccesKey, cfg.UserName)

		// Upload app
		log.Infof("Uploading IPA")
		appURL, err := xcuitests.UploadIPA(cfg.XCUITestsIPA)
		if err != nil {
			failf("Failed to upload IPA (%s), error: %s", cfg.XCUITestsIPA, err)
		}

		log.Printf("Upload success")
		log.Donef("Uploaded app URL => %s", appURL)
		fmt.Println()

		// Upload test runner
		log.Infof("Uploading UITest runner")
		testURL, err := xcuitests.UploadTestRunner(cfg.XCUITestsRunner)
		if err != nil {
			failf("Failed to upload test runner (%s), error: %s", cfg.XCUITestsRunner, err)
		}

		log.Printf("Upload success")
		log.Donef("Uploaded test URL => %s", testURL)
		fmt.Println()

		// Execute test
		log.Infof("Executing test")
		devices := strings.Split(cfg.XCUITestsDevices, "|")
		message, buildID, err := xcuitests.ExecuteTest(appURL, testURL, cfg.XCUITestsDeviceLogs, devices)
		if err != nil {
			failf("Failed to execute test, error: %s", err)
		}

		log.Printf("Respond: %s", message)
		log.Donef("Build ID => %s", buildID)
		fmt.Println()

		if message != "Success" {
			failf("Failed to execute test")
		}

		log.Infof("Test running")
		{
			ch := make(chan appautomate.BuildResult)

			s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
			s.Start()

			go xcuitests.ListenForTestComplete(buildID, ch)
			buildResult := <-ch
			s.Stop()

			log.Printf("buildResult: %+v", buildResult)
		}

	}

}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
