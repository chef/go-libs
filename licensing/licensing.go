package licensing

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var licenseKeys []string

func getLicenseKeys() []string {
	return licenseKeys
}

func appendLicenseKey(key string) {
	licenseKeys = append(licenseKeys, key)
}

func CheckSoftwareEntitlement() {
	var licenseKey []string
	licenseKey = licenseFileFetch()
	client(licenseKey)

	key, check := os.LookupEnv("CHEF_LICENSE_KEY")
	if check {
		licenseKey = append(licenseKey, key)
		client(licenseKey)
		return
	}
	args := os.Args
	for k, v := range args {
		if v == "--chef-license-key" {
			if len(args) > k+1 {
				licenseKey = append(licenseKey, args[k+1])
				client(licenseKey)
				return
			}
		} else if strings.HasPrefix(v, "--chef-license-key=") {
			split := strings.Split(v, "=")
			licenseKey = append(licenseKey, split[len(split)-1])
			client(licenseKey)
			return
		}
	}
	client(licenseKey)
}

func FetchAndPersist() []string {
	// Load the existing licenseKeys from the license file
	for _, key := range licenseFileFetch() {
		appendLicenseKey(key)
	}

	newKeys := []string{fetchFromArg()}
	licenseType := validateAndFetchLicenseType(newKeys[0])
	fmt.Println("License type is: " + licenseType)
	if licenseType != "" && !hasUnrestrictedLicenseAdded(newKeys, licenseType) {
		// licenseKeys = append(licenseKeys, newKeys[0])
		appendLicenseKey(newKeys[0])
		return licenseKeys
	}

	newKeys = []string{fetchFromEnv()}
	fmt.Println("key from env", newKeys)
	licenseType = validateAndFetchLicenseType(newKeys[0])
	if licenseType != "" && !hasUnrestrictedLicenseAdded(newKeys, licenseType) {
		// licenseKeys = append(licenseKeys, newKeys[0])
		appendLicenseKey(newKeys[0])
		return licenseKeys
	}

	// # Return keys if license keys are active and not expired or expiring
	// # Return keys if there is any error in /client API call, and do not block the flow.
	// # Client API possible errors will be handled in software entitlement check call (made after this)
	// # client_api_call_error is set to true when there is an error in licenses_active? call
	// return @license_keys if (!@license_keys.empty? && licenses_active? && ChefLicensing::Context.license.license_type.downcase == "commercial") || client_api_call_error
	if len(getLicenseKeys()) > 0 && isLicenseActive(getLicenseKeys()) {
		return getLicenseKeys()
	}

	return licenseKeys
}

func fetchFromArg() string {
	log.Println("License Key fetcher examining CLI arg checks")
	licenseKey := flag.String("chef-license-key", "", "Chef license key")

	flag.Parse()
	return *licenseKey
}

func fetchFromEnv() string {
	log.Println("License Key fetcher examining ENV checks")
	key, _ := os.LookupEnv("CHEF_LICENSE_KEY")
	// if check {
	// 	validateAndFetchLicenseType(key)
	// }

	return key
}

func fetchInteractively() string {

	return ""
}

func validateLicenseWithServer(key string, suppress bool) (bool, string) {
	var opts = make(map[string]string)
	opts["licenseId"] = key
	response := invokeGetAPI("validate", opts, suppress).(*ValidateResponse)
	return response.Data, response.Message
}

func validateAndFetchLicenseType(key string) string {
	var licenseType string
	if key == "" {
		return licenseType
	}

	isValid, _ := validateLicenseWithServer(key, false)
	if isValid {
		licenseType = fetchLicenseType([]string{key})
	}

	return licenseType
}

func fetchLicenseClient(keys []string) *Client {
	config := GetConfig()
	var opts = make(map[string]string)

	opts["licenseId"] = strings.Join(keys, ",")
	opts["entitlementId"] = config.EntitlementID
	response := invokeGetAPI("client", opts).(*ClientResponse)

	return &response.Data.Client
}

func fetchLicenseType(licenseKey []string) string {
	client := fetchLicenseClient(licenseKey)
	return client.License
}

func client(licenseKey []string) {
	config := GetConfig()
	if len(licenseKey) == 0 {
		log.Fatal("You dont have license key, Please generate by running chef license command")
	} else {
		var opts = make(map[string]string)
		opts["licenseId"] = strings.Join(licenseKey, ",")
		opts["entitlementId"] = config.EntitlementID
		invokeGetAPI("client", opts)
	}
}
