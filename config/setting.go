package config

import "time"

var TestSetting *Setting

// Setting
// global test settings
type Setting struct {

	// ScreenshotOnFail
	// used in find element strategy
	// takes screenshot and writes to artifacts
	// if unable to find webelement within timeout
	ScreenshotOnFail bool

	// TimeoutFind
	// used in find element strategy
	// controls timeout of performing driver.F("selector") find
	// 20 seconds timeout is approximation of 2 retries
	TimeoutFind time.Duration

	// TimeoutDelay
	// delay to retry find element request
	// 700 ms is an arbitrary value
	TimeoutDelay time.Duration

	// RefreshOnFindError
	// calls /session/{sessionId}/refresh
	// if find retry fails
	RefreshOnFindError bool

	// Artifact path
	// 
	// ArtifactRecordsPath
	// from app root a directory that stores
	// Google Chrome Recorder json files
	// for TestSteps generation
	// will check specified path for *.json records
	ArtifactRecordsPath string

	// ArtifactScreenshotsPath
	// from app root a directory where
	// ScreenshotOnFail, or d.Screenshot()
	// stores driver screnshots in *.jpg format
	ArtifactScreenshotsPath string
}

func DefaultSetting() *Setting {
	return &Setting{
		ScreenshotOnFail:   true,
		TimeoutFind:        30,
		TimeoutDelay:       700,
		RefreshOnFindError: true,
		ArtifactRecordsPath: "../artifacts/records",
	}
}

type SettingsFunc func(*Setting)

func WithTimeoutDelay(t time.Duration) SettingsFunc {
	return func(s *Setting) {
		s.TimeoutDelay = t
	}
}

func WithTimeoutFind(t time.Duration) SettingsFunc {
	return func(s *Setting) {
		s.TimeoutFind = t
	}
}
