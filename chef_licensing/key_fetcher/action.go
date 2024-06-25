package keyfetcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chef/go-libs/chef_licensing/api"
	"github.com/chef/go-libs/chef_licensing/spinner"
	"github.com/cqroot/prompt"
	inputPrompt "github.com/cqroot/prompt/input"
	"github.com/gookit/color"
)

type PromptAttribute struct {
	TimeoutDuration int    `yaml:"timeout_duration"`
	TimeoutMessage  string `yaml:"timeout_message"`
}

type ActionDetail struct {
	Messages        []string          `yaml:"messages"`
	Options         []string          `yaml:"options,omitempty"`
	Action          string            `yaml:"action"`
	PromptType      string            `yaml:"prompt_type"`
	PromptAttribute PromptAttribute   `yaml:"prompt_attributes"`
	Paths           []string          `yaml:"paths"`
	ResponsePathMap map[string]string `yaml:"response_path_map"`
	Choice          string            `yaml:"choice"`
}

var lastUserInput string

func (ad ActionDetail) Say() string {
	renderMessages(ad.Messages)
	return ad.Paths[0]
}

func (ad ActionDetail) TimeoutSelect() string {
	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	done := make(chan struct{})
	var val string
	var err error
	go func() {
		val, err = prompt.New().Ask(ad.Messages[0]).
			Choose(ad.Options)
		checkPromptErr(err)
		close(done)
	}()

	select {
	case <-done:
		if err == nil {
			fmt.Println("Selected option: ", val)
			return ad.ResponsePathMap[val]
		}
	case <-timeoutContext.Done():
		fmt.Println(printInColor("red", "Prompt timed out. Use non-interactive flags or enter an answer within 60 seconds.", false, true))
		fmt.Println("Timeout!")
		os.Exit(1)
	}
	return ""
}

func (ad ActionDetail) Ask() string {
	val, err := prompt.New().Ask(ad.Messages[0]).
		Input("license-key", inputPrompt.WithValidateFunc(validateLicenseFormat))
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else if errors.Is(err, ErrInvalidKeyFormat) {
			fmt.Fprintln(os.Stderr, err)
		} else {
			panic(err)
		}
	}
	lastUserInput = val
	PromptInput.LicenseID = val

	return ad.Paths[0]
}

func (ad ActionDetail) Select() string {
	val1, err := prompt.New().Ask(ad.Messages[0]).
		Choose(ad.Options)
	checkPromptErr(err)

	return ad.ResponsePathMap[val1]
}

func (ad ActionDetail) SayAndSelect() string {
	renderMessages(ad.Messages)
	val1, err := prompt.New().Ask(ad.Choice).Choose(ad.Options)
	checkPromptErr(err)

	return ad.ResponsePathMap[val1]
}

func (ad ActionDetail) Warn() string {
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) Error() string {
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) Ok() string {
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) DoesLicenseHaveValidPattern() string {
	isValid := ValidateKeyFormat(lastUserInput)
	if isValid {
		return ad.ResponsePathMap["true"]
	} else {
		color.Warn.Println(ErrInvalidKeyFormat)
		return ad.ResponsePathMap["false"]
	}
}

func (ad ActionDetail) IsLicenseValidOnServer() string {
	spinner, err := spinner.GetSpinner()
	if err != nil {
		log.Println("Unable to start the spinner")
	}
	_ = spinner.Start()
	spinner.Message("In progress")
	isValid, message := api.GetClient().ValidateLicenseAPI(lastUserInput, true)

	if isValid {
		spinner.StopCharacter("✓")
		spinner.StopColors("green")
	} else {
		spinner.StopCharacter("✖")
		spinner.StopColors("red")
		PromptInput.FailureMessage = message.Error()
	}

	spinner.Message("Done")
	_ = spinner.Stop()

	return ad.ResponsePathMap[strconv.FormatBool(isValid)]
}

func (ad ActionDetail) FetchInvalidLicenseMessage() string {
	if PromptInput.FailureMessage == "" {
		_, message := api.GetClient().ValidateLicenseAPI(lastUserInput, true)
		PromptInput.FailureMessage = message.Error()
	}
	return ad.Paths[0]
}

func (ad ActionDetail) IsLicenseAllowed() string {
	client, error := api.GetClient().GetLicenseClient([]string{lastUserInput})
	if error != nil {
		log.Fatal(error)
	}
	licenseType := client.LicenseType
	PromptInput.LicenseType = licenseType
	if licenseType == "commercial" {
		PromptInput.IsCommercial = true
	}

	var isRestricted bool
	if IsLicenseRestricted(licenseType) {
		// Existing license keys needs to be fetcher to show details of existing license of license type which is restricted.
		// However, if user is trying to add Free Tier License, and user has active trial license, we fetch the trial license key
		var existingLicenseKeysInFile []string
		if licenseType == "free" && DoesUserHasActiveTrialLicense() {
			existingLicenseKeysInFile = FetchLicenseKeysBasedOnType(":trial")
		} else {
			existingLicenseKeysInFile = FetchLicenseKeysBasedOnType(":" + licenseType)
		}
		PromptInput.LicenseID = existingLicenseKeysInFile[len(existingLicenseKeysInFile)-1]
	} else {
		isRestricted = true
	}
	return ad.ResponsePathMap[strconv.FormatBool(isRestricted)]
}

func (ad ActionDetail) DetermineRestrictionType() string {
	var resType string
	if PromptInput.LicenseType == "free" && DoesUserHasActiveTrialLicense() {
		resType = "active_trial_restriction"
	} else {
		resType = PromptInput.LicenseType + "_restriction"
	}

	return ad.ResponsePathMap[resType]
}

func (ad ActionDetail) DisplayLicenseInfo() string {
	PrintLicenseKeyOverview([]string{lastUserInput})
	return ad.Paths[0]
}

func (ad ActionDetail) FetchLicenseTypeRestricted() string {
	var val string
	if IsLicenseRestricted("trial") && IsLicenseRestricted("free") {
		val = "trial_and_free"
	} else if IsLicenseRestricted("trial") {
		val = "trial"
	} else {
		val = "free"
	}
	return ad.ResponsePathMap[val]
}

func (ad ActionDetail) CheckLicenseExpirationStatus() string {
	licenseClient := getLicense()
	var status string
	if licenseClient.IsExpired() || licenseClient.HaveGrace() {
		status = "expired"
	} else if licenseClient.IsAboutToExpire() {
		expiresOn, err := time.Parse(time.RFC3339, licenseClient.ChangesOn)
		if err != nil {
			log.Fatal("Unknown expiration time received from the server: ", err)
		}

		expirationIn := int(time.Until(expiresOn).Hours() / 24)
		PromptInput.LicenseExpirationDate = expiresOn.Format(time.UnixDate)
		PromptInput.ExpirationInDays = strconv.Itoa(expirationIn)
		status = "about_to_expire"
	} else if licenseClient.IsExhausted() && (licenseClient.IsCommercial() || licenseClient.IsFree()) {
		status = "exhausted_license"
	} else {
		status = "active"
	}

	return ad.ResponsePathMap[status]
}

func (ad ActionDetail) FetchLicenseId() string {
	return ad.Paths[0]
}

func (ad ActionDetail) IsCommercialLicense() string {
	val := PromptInput.IsCommercial
	return ad.ResponsePathMap[strconv.FormatBool(val)]
}

func (ad ActionDetail) IsRunAllowedOnLicenseExhausted() string {
	val := PromptInput.IsCommercial

	return ad.ResponsePathMap[strconv.FormatBool(val)]
}

func (ad ActionDetail) FilterLicenseTypeOptions() string {
	var val string
	if IsLicenseRestricted("trial") && IsLicenseRestricted("free") || DoesUserHasActiveTrialLicense() {
		val = "ask_for_commercial_only"
	} else if IsLicenseRestricted("trial") {
		val = "ask_for_license_except_trial"
	} else if IsLicenseRestricted("free") {
		val = "ask_for_license_except_free"
	} else {
		val = "ask_for_all_license_type"
	}

	return ad.ResponsePathMap[val]
}

func (ad ActionDetail) SetLicenseInfo() string {
	lastUserInput = PromptInput.LicenseID
	return ad.Paths[0]
}
