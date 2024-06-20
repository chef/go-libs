package licensing

import (
	"fmt"
	"slices"
)

func hasUnrestrictedLicenseAdded(newKeys []string, licenseType string) bool {
	if isLicenseRestricted(licenseType) {
		// Existing license keys of same license type are fetched to compare if old license key or a new one is added.
		// However, if user is trying to add Free Tier License, and user has active trial license, we fetch the trial license key

		var existingLicenseKeysInFile []string
		if licenseType == "free" && doesUserHasActiveTrialLicense() {
			existingLicenseKeysInFile = fetchLicenseKeysBasedOnType(":trial")
		} else {
			existingLicenseKeysInFile = fetchLicenseKeysBasedOnType(":" + licenseType)
		}
		if existingLicenseKeysInFile[len(existingLicenseKeysInFile)-1] != newKeys[0] {
			promptLicenseAdditionRestricted(licenseType, existingLicenseKeysInFile)
			return false
		}

		return true
	} else {
		persistAndConcat(newKeys, licenseType)
	}

	return true
}

func isLicenseRestricted(licenseType string) (out bool) {
	allowed := AllowedLicencesForAddition()
	if !slices.Contains(allowed, licenseType) {
		out = true
	}

	return
}

func AllowedLicencesForAddition() []string {
	var license_types = []string{"free", "trial", "commercial"}
	currentTypes := currentLicenseTypes()
	// fmt.Println(license_types, currentTypes)

	// fmt.Printf("contains trial ? %t\n", slices.Contains(currentTypes, ":trial"))
	// fmt.Printf("contains free ? %t\n", slices.Contains(currentTypes, ":free"))

	if slices.Contains(currentTypes, ":trial") {
		// fmt.Println("inside trial removal")
		removeItem(&license_types, "trial")
	}
	if slices.Contains(currentTypes, ":free") || doesUserHasActiveTrialLicense() {
		// fmt.Println("Inside free removal")
		removeItem(&license_types, "free")
	}

	return license_types
}

func currentLicenseTypes() (out []string) {
	content := *ReadLicenseKeyFile()
	for _, license := range content.Licenses {
		out = append(out, license.LicenseType)
	}
	return
}

func doesUserHasActiveTrialLicense() (out bool) {
	content := *ReadLicenseKeyFile()
	for _, license := range content.Licenses {
		if license.LicenseType == "trial" && fetchLicenseClient([]string{license.LicenseKey}).Status == "Active" {
			out = true
		}
	}

	return
}

func removeItem(target *[]string, item string) {
	var out []string
	for _, str := range *target {
		if str != item {
			out = append(out, str)
		}
	}

	fmt.Println("Inside removeItem item: " + item)
	fmt.Println(out)

	*target = out
}
