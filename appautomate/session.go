package appautomate

// Session details of your test sessions. Each session is the execution of your test suite on individual devices.
type Session struct {
	BuildID                 string                  `json:"build_id"`
	SessionID               string                  `json:"session_id"`
	Device                  string                  `json:"device"`
	StartTime               string                  `json:"start_time"`
	Duration                float64                 `json:"duration"`
	SessionAppDetails       SessionAppDetails       `json:"app_details"`
	SessionTestSuiteDetails SessionTestSuiteDetails `json:"test_suite_details"`
	TestCount               int                     `json:"test_count"`
	TestDetails             map[string]interface{}  `json:"test_details"`
	SessionTestStatus       SessionTestStatus       `json:"test_status"`
}

// SessionAppDetails ...
type SessionAppDetails struct {
	URL      string      `json:"url"`
	BundleID string      `json:"bundle_id"`
	CustomID interface{} `json:"custom_id"`
	Version  string      `json:"version"`
	Name     string      `json:"name"`
}

// SessionTestSuiteDetails ...
type SessionTestSuiteDetails struct {
	URL      string      `json:"url"`
	BundleID string      `json:"bundle_id"`
	CustomID interface{} `json:"custom_id"`
	Name     string      `json:"name"`
}

// // TestDetails ...   TODO
// type TestDetails struct {
// 	BullsEyeUITestsBullsEyeUITests2 struct {
// 		TestScreenshot struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testScreenshot"`
// 		TestGameStyleSwitch struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testGameStyleSwitch"`
// 		TestCase struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testCase"`
// 	} `json:"BullsEyeUITests/BullsEyeUITests_2"`
// 	BullsEyeUITestsBullsEyeUITests struct {
// 		TestScreenshot struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testScreenshot"`
// 		TestGameStyleSwitch struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testGameStyleSwitch"`
// 		TestCase struct {
// 			StartTime          string  `json:"start_time"`
// 			Status             string  `json:"status"`
// 			TestID             string  `json:"test_id"`
// 			Duration           float64 `json:"duration"`
// 			InstrumentationLog string  `json:"instrumentation_log"`
// 			DeviceLog          string  `json:"device_log"`
// 			Video              string  `json:"video"`
// 			NetworkLog         string  `json:"network_log"`
// 		} `json:"testCase"`
// 	} `json:"BullsEyeUITests/BullsEyeUITests"`
// }

// SessionTestStatus ...
type SessionTestStatus struct {
	SUCCESS  int `json:"SUCCESS"`
	FAILED   int `json:"FAILED"`
	IGNORED  int `json:"IGNORED"`
	TIMEDOUT int `json:"TIMEDOUT"`
	QUEUED   int `json:"QUEUED"`
}
