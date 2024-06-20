package licensing

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/cqroot/prompt"
	inputPrompt "github.com/cqroot/prompt/input"
	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

//go:embed interactions.yml
var interactionsYAML []byte

type Interaction struct {
	FileFormatVersion string                  `yaml:":file_format_version"`
	Actions           map[string]ActionDetail `yaml:"interactions"`
}

type TemplateConfig struct {
	ProductName           string
	UnitMeasure           string
	ChefExecutableName    string
	FailureMessage        string
	IsCommercial          bool
	LicenseType           string
	LicenseID             string
	LicenseExpirationDate string
	ExpirationInDays      string
}

var input TemplateConfig
var lastUserInput string

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

func (ad ActionDetail) Say() string {
	fmt.Println("Say function called")
	renderMessages(ad.Messages)
	return ad.Paths[0]
}

// TODO: Implement the timeout
func (ad ActionDetail) TimeoutSelect() string {
	fmt.Println("TimeoutSelect function called")
	val, err := prompt.New().Ask(ad.Messages[0]).
		Choose(ad.Options)
	checkPromptErr(err)
	fmt.Println("Selected option: ", val)
	return ad.ResponsePathMap[val]
}

func (ad ActionDetail) Ask() string {
	fmt.Println("Ask function called")
	val, err := prompt.New().Ask(ad.Messages[0]).
		Input("license-key", inputPrompt.WithValidateFunc(validateLicenseFormat))
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
		} else if errors.Is(err, ErrInvalidKeyFormat) {
			fmt.Fprintln(os.Stderr, err)
		} else {
			panic(err)
		}
	}
	lastUserInput = val
	input.LicenseID = val

	return ad.Paths[0]
}

func (ad ActionDetail) Select() string {
	fmt.Println("Select function called")
	val1, err := prompt.New().Ask(ad.Messages[0]).
		Choose(ad.Options)
	checkPromptErr(err)

	return ad.ResponsePathMap[val1]
}

func (ad ActionDetail) SayAndSelect() string {
	fmt.Println("SayAndSelect called....")
	renderMessages(ad.Messages)
	val1, err := prompt.New().Ask(ad.Choice).Choose(ad.Options)
	checkPromptErr(err)

	return ad.ResponsePathMap[val1]
}

func (ad ActionDetail) Warn() string {
	fmt.Println("Warn function called")
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) Error() string {
	fmt.Println("Error function called")
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) Ok() string {
	fmt.Println("Ok function called")
	renderMessages(ad.Messages)

	return ad.Paths[0]
}

func (ad ActionDetail) DoesLicenseHaveValidPattern() string {
	isValid := validateKeyFormat(lastUserInput)
	if isValid {
		return ad.ResponsePathMap["true"]
	} else {
		// fmt.Errorf("%s: %w", lastUserInput, ErrInvalidKeyFormat)
		color.Warn.Println(ErrInvalidKeyFormat)
		return ad.ResponsePathMap["false"]
	}
}

func (ad ActionDetail) IsLicenseValidOnServer() string {
	spinner, err := GetSpinner()
	if err != nil {
		log.Println("Unable to start the spinner")
	}
	_ = spinner.Start()
	spinner.Message("In progress")
	isValid, message := validateLicenseWithServer(lastUserInput, true)

	if isValid {
		spinner.StopCharacter("✓")
		spinner.StopColors("green")
		// return ad.ResponsePathMap["true"]
	} else {
		spinner.StopCharacter("✖")
		spinner.StopColors("red")
		input.FailureMessage = message
		// return ad.ResponsePathMap["false"]
	}
	// time.Sleep(2 * time.Second)

	spinner.Message("Done")

	_ = spinner.Stop()

	return ad.ResponsePathMap[strconv.FormatBool(isValid)]
}

func (ad ActionDetail) FetchInvalidLicenseMessage() string {
	if input.FailureMessage == "" {
		_, message := validateLicenseWithServer(lastUserInput, true)
		input.FailureMessage = message
	}
	return ad.Paths[0]
}

func (ad ActionDetail) IsLicenseAllowed() string {
	licenseType := fetchLicenseType([]string{lastUserInput})
	input.LicenseType = licenseType
	if licenseType == "commercial" {
		input.IsCommercial = true
	}

	var isRestricted bool
	if isLicenseRestricted(licenseType) {
		// Existing license keys needs to be fetcher to show details of existing license of license type which is restricted.
		// However, if user is trying to add Free Tier License, and user has active trial license, we fetch the trial license key
		var existingLicenseKeysInFile []string
		if licenseType == "free" && doesUserHasActiveTrialLicense() {
			existingLicenseKeysInFile = fetchLicenseKeysBasedOnType(":trial")
		} else {
			existingLicenseKeysInFile = fetchLicenseKeysBasedOnType(":" + licenseType)
		}
		input.LicenseID = existingLicenseKeysInFile[len(existingLicenseKeysInFile)-1]
	} else {
		isRestricted = true
	}
	return ad.ResponsePathMap[strconv.FormatBool(isRestricted)]
}

func (ad ActionDetail) DetermineRestrictionType() string {
	var resType string
	if input.LicenseType == "free" && doesUserHasActiveTrialLicense() {
		resType = "active_trial_restriction"
	} else {
		resType = input.LicenseType + "_restriction"
	}

	return ad.ResponsePathMap[resType]
}

func (ad ActionDetail) DisplayLicenseInfo() string {
	printLicenseKeyOverview([]string{lastUserInput})
	return ad.Paths[0]
}

func (ad ActionDetail) FetchLicenseTypeRestricted() string {
	var val string
	if isLicenseRestricted("trial") && isLicenseRestricted("free") {
		val = "trial_and_free"
	} else if isLicenseRestricted("trial") {
		val = "trial"
	} else {
		val = "free"
	}
	return ad.ResponsePathMap[val]
}

func (ad ActionDetail) CheckLicenseExpirationStatus() string {
	license := getLicense()
	var status string
	if isExpired(license) || haveGrace(license) {
		status = "expired"
	} else if isAboutToExpire(license) {
		expiresOn, err := time.Parse(time.RFC3339, license.ChangesOn)
		if err != nil {
			log.Fatal("Unknown expiration time received from the server: ", err)
		}

		expirationIn := int(time.Until(expiresOn).Hours() / 24)
		input.LicenseExpirationDate = expiresOn.Format(time.UnixDate)
		input.ExpirationInDays = strconv.Itoa(expirationIn)
		status = "about_to_expire"
	} else if isExhausted(license) && (license.License == "commercial" || license.License == "free") {
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
	val := input.IsCommercial
	return ad.ResponsePathMap[strconv.FormatBool(val)]
}

func (ad ActionDetail) IsRunAllowedOnLicenseExhausted() string {
	val := input.IsCommercial

	return ad.ResponsePathMap[strconv.FormatBool(val)]
}

func (ad ActionDetail) FilterLicenseTypeOptions() string {
	var val string
	if isLicenseRestricted("trial") && isLicenseRestricted("free") || doesUserHasActiveTrialLicense() {
		val = "ask_for_commercial_only"
	} else if isLicenseRestricted("trial") {
		val = "ask_for_license_except_trial"
	} else if isLicenseRestricted("free") {
		val = "ask_for_license_except_free"
	} else {
		val = "ask_for_all_license_type"
	}

	return ad.ResponsePathMap[val]
}

func checkPromptErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

func updateInputs(conf map[string]string) {
	input.ProductName = conf["ProductName"]
	input.ChefExecutableName = conf["ChefExecutableName"]
	if conf["ChefExecutableName"] == "chef" {
		input.UnitMeasure = "nodes"
	} else {
		input.UnitMeasure = "targets"
	}
}

func GetIntractions() map[string]ActionDetail {
	m := make(map[string]string)
	m["ProductName"] = "Workstation"
	m["ChefExecutableName"] = "chef"
	updateInputs(m)

	var intr Interaction
	err := yaml.Unmarshal(interactionsYAML, &intr)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(interactionsYAML))
	return intr.Actions
}

func StartInteractions(startID string) {
	if startID == "" {
		startID = "start"
	}

	var performedInteractions []string
	currentID := startID
	previousID := ""
	interactions := GetIntractions()
	// fmt.Println(interactions, previousID)

	// count := 1
	for {
		action := interactions[currentID]
		// fmt.Println("\nPerforming action::::::::::::::::::::: ", currentID)
		if currentID == "" || currentID == "exit" {
			break
		}
		performedInteractions = append(performedInteractions, currentID)
		previousID = currentID
		currentID = performInteraction(action)
		// if count == 23 {
		// 	fmt.Println("Stopping because of the counter, ", count)
		// 	break
		// }
		// count += 1
	}
	// fmt.Println(performedInteractions)
	if currentID != "exit" {
		log.Fatal("Something went wrong in the flow. The last interaction was " + previousID)
	}
}

func performInteraction(action ActionDetail) (nextID string) {
	var methodName string
	if action.PromptType != "" {
		methodName = action.PromptType
	} else if action.Action != "" {
		methodName = action.Action
	}

	meth := reflect.ValueOf(action).MethodByName(methodName)
	// fmt.Println("--------------------------------------------", action)
	returnVals := meth.Call(nil)

	if len(returnVals) > 0 {
		if returnValue, ok := returnVals[0].Interface().(string); ok {
			nextID = returnValue
		}
	} else {
		log.Fatal("Something went wrong with the interactions")
	}

	return
}

func renderMessages(messages []string) {
	if len(messages) == 0 {
		return
	}

	for _, message := range messages {
		tmpl, err := template.New("actionMessage").Funcs(template.FuncMap{
			"printHyperlink": printHyperlink,
			"printInColor":   printInColor,
			"printBoldText":  printBoldText,
		}).Parse(message)
		if err != nil {
			log.Fatalf("error parsing template: %v", err)
		}
		fmt.Println("")

		err = tmpl.Execute(os.Stdout, input)
		if err != nil {
			log.Fatalf("error executing template: %v", err)
		}
	}
}

func printHyperlink(url string) string {
	return color.Style{color.FgGreen, color.OpUnderscore}.Sprintf(url)
}

func printInColor(selColor, text string, options ...bool) string {
	output := color.Style{}
	var underline bool
	var bold bool

	if len(options) == 1 {
		underline = options[0]
	}
	if len(options) > 1 {
		bold = options[1]
	}

	switch selColor {
	case "red":
		output = append(output, color.FgRed)
	case "green":
		output = append(output, color.FgGreen)
	case "blue":
		output = append(output, color.FgBlue)
	case "yellow":
		output = append(output, color.FgYellow)
	}

	if underline {
		output = append(output, color.OpUnderscore)
	}
	if bold {
		output = append(output, color.OpBold)
	}

	return output.Sprintf(text)
}

func printBoldText(text1, text2 string) string {
	return color.Bold.Sprintf(text1 + " " + text2)
}

func validateLicenseFormat(key string) error {
	isValid := validateKeyFormat(key)
	if isValid {
		return nil
	} else {
		return fmt.Errorf("%s: %w", key, ErrInvalidKeyFormat)
	}
}

func getLicense() Client {
	spinner, err := GetSpinner()
	if err != nil {
		log.Println("Unable to start the spinner")
	}
	_ = spinner.Start()
	spinner.Message("In progress")
	client := fetchLicenseClient([]string{input.LicenseID})

	spinner.StopCharacter("✓")
	spinner.StopColors("green")

	spinner.Message("Done")

	_ = spinner.Stop()

	return *client
}
