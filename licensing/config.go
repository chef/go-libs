package licensing

import (
	"time"

	"github.com/theckman/yacspin"
)

type LicenseConfig struct {
	ProductName      string
	EntitlementID    string
	LicenseServerURL string
}

var cfg *LicenseConfig

func SetConfig(c *LicenseConfig) {
	cfg = c
}

func GetConfig() *LicenseConfig {
	return cfg
}

func GetSpinner() (*yacspin.Spinner, error) {
	SpinnerConfig := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          "License validation",
		SuffixAutoColon: true,
		// StopCharacter:   "✓",
		// StopColors:      []string{"fgGreen"},
	}

	return yacspin.New(SpinnerConfig)
}
