package appautomate

// Build ...
type Build struct {
	BuildID           string            `json:"build_id"`
	Framework         string            `json:"framework"`
	Status            string            `json:"status"`
	InputCapabilities InputCapabilities `json:"input_capabilities"`
	StartTime         string            `json:"start_time"`
	AppDetails        AppDetails        `json:"app_details"`
	TestSuiteDetails  TestSuiteDetails  `json:"test_suite_details"`
	Duration          string            `json:"duration"`
	Devices           map[string]Device `json:"devices"`
}

// InputCapabilities ...
type InputCapabilities struct {
	Devices         []string    `json:"devices"`
	DeviceLogs      bool        `json:"deviceLogs"`
	App             string      `json:"app"`
	TestSuite       string      `json:"testSuite"`
	SetEnvVariables interface{} `json:"setEnvVariables"`
}

// TestSuiteDetails ...
type TestSuiteDetails struct {
	URL      string `json:"url"`
	BundleID string `json:"bundle_id"`
	Name     string `json:"name"`
}

// AppDetails ...
type AppDetails struct {
	URL      string `json:"url"`
	BundleID string `json:"bundle_id"`
	Version  string `json:"version"`
	Name     string `json:"name"`
}

// Device ...
type Device struct {
	SessionID      string     `json:"session_id"`
	Status         string     `json:"status"`
	SessionDetails string     `json:"session_details"`
	TestStatus     TestStatus `json:"test_status"`
	Name           string
}

// TestStatus ...
type TestStatus struct {
	SUCCESS  int `json:"SUCCESS"`
	FAILED   int `json:"FAILED"`
	IGNORED  int `json:"IGNORED"`
	TIMEDOUT int `json:"TIMEDOUT"`
	QUEUED   int `json:"QUEUED"`
}
