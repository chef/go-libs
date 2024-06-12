package licensing

import (
	"log"
	"time"
)

// func prompt_license_addition_restricted(license_type, existing_license_keys_in_file)
func promptLicenseAdditionRestricted(licenseType string, existingLicenseKeysInFile []string) {
	log.Println("License Key fetcher - prompting license addition restriction")

}

func isLicenseActive(keys []string) (out bool) {
	log.Println("License Key fetcher - checking if licenses are active")

	spinner, err := GetSpinner()
	if err != nil {
		log.Println("Unable to start the spinner")
	}
	_ = spinner.Start()
	spinner.Message("In progress")
	license := *fetchLicenseClient(keys)

	// Intentional lag of 2 seconds when license is expiring or expired
	if isExpiringOrExpired(license) {
		time.Sleep(2 * time.Second)
	}

	if isExpired(license) || haveGrace(license) {
		// if ChefLicensing::Context.local_licensing_service?
		//   config[:start_interaction] = :prompt_license_expired_local_mode
		// else
		//   config[:start_interaction] = :prompt_license_expired
		// end
		// prompt_fetcher.config = config
		// false
		out = false
	} else if isAboutToExpire(license) {
		// config[:start_interaction] = :prompt_license_about_to_expire
		// prompt_fetcher.config = config
		out = false
	} else if isExhausted(license) && (license.License == "commercial" || license.License == "free") {
		// config[:start_interaction] = :prompt_license_exhausted
		// prompt_fetcher.config = config
		// false
		out = false
	} else {
		// If license is not expired or expiring, return true. But if the license is not commercial, warn the user.
		if license.License != "commercial" {
			// config[:start_interaction] = :warn_non_commercial_license unless license.license_type.downcase == "commercial"
		}
		out = true
	}
	if out {
		spinner.StopCharacter("✓")
		spinner.StopColors("green")
	} else {
		spinner.StopCharacter("X")
		spinner.StopColors("red")
	}

	time.Sleep(2 * time.Second)

	spinner.Message("Done")

	_ = spinner.Stop()

	return out
}
