package keyfetcher

import (
	"flag"
	"os"

	"github.com/chef/go-libs/chef_licensing/api"
)

var licenseKeys []string

func GlobalFetchAndPersist() []string {
	// Load the existing licenseKeys from the license file
	for _, key := range licenseFileFetch() {
		appendLicenseKey(key)
	}

	newKeys := []string{fetchFromArg()}
	licenseType := validateAndFetchLicenseType(newKeys[0])
	if licenseType != "" && !hasUnrestrictedLicenseAdded(newKeys, licenseType) {
		appendLicenseKey(newKeys[0])
		return licenseKeys
	}

	newKeys = []string{fetchFromEnv()}
	licenseType = validateAndFetchLicenseType(newKeys[0])
	if licenseType != "" && !hasUnrestrictedLicenseAdded(newKeys, licenseType) {
		appendLicenseKey(newKeys[0])
		return licenseKeys
	}

	// Return keys if license keys are active and not expired or expiring
	// Return keys if there is any error in /client API call, and do not block the flow.
	// Client API possible errors will be handled in software entitlement check call (made after this)
	// client_api_call_error is set to true when there is an error in licenses_active? call
	isActive, startID := isLicenseActive(getLicenseKeys())
	if len(getLicenseKeys()) > 0 && isActive {
		return getLicenseKeys()
	}

	newKeys = fetchInteractively(startID)
	if len(newKeys) > 0 {
		licenseClient, _ := api.GetClient().GetLicenseClient(newKeys)
		persistAndConcat(newKeys, licenseClient.LicenseType)
		if (!licenseClient.IsExpired() && !licenseClient.IsExhausted()) || licenseClient.IsCommercial() {
			return licenseKeys
		}
	}

	return licenseKeys
}

func FetchLicenseType(licenseKeys []string) string {
	client, _ := api.GetClient().GetLicenseClient(licenseKeys)
	return client.LicenseType
}

func getLicenseKeys() []string {
	return licenseKeys
}

func appendLicenseKey(key string) {
	licenseKeys = append(licenseKeys, key)
}

func fetchFromArg() string {
	licenseKey := flag.String("chef-license-key", "", "Chef license key")

	flag.Parse()
	return *licenseKey
}

func fetchFromEnv() string {
	key, _ := os.LookupEnv("CHEF_LICENSE_KEY")

	return key
}

func fetchInteractively(startID string) []string {
	return StartInteractions(startID)
}

func validateAndFetchLicenseType(key string) string {
	var licenseType string
	if key == "" {
		return licenseType
	}

	isValid, _ := api.GetClient().ValidateLicenseAPI(key)
	if isValid {
		licenseType = FetchLicenseType([]string{key})
	}

	return licenseType
}
