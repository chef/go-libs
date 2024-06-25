package keyfetcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"text/template"

	"github.com/chef/go-libs/chef_licensing/api"
	"github.com/chef/go-libs/chef_licensing/config"
	"github.com/chef/go-libs/chef_licensing/spinner"
	"github.com/cqroot/prompt"
	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

func StartInteractions(startID string) (keys []string) {
	if startID == "" {
		startID = "start"
	}
	initializePromptInputs()
	// var performedInteractions []string
	currentID := startID
	previousID := ""
	interactions := getIntractions()

	for {
		action := interactions[currentID]
		if currentID == "" || currentID == "exit" {
			break
		}
		// performedInteractions = append(performedInteractions, currentID)
		previousID = currentID
		currentID = performInteraction(action)
	}
	if currentID != "exit" {
		log.Fatal("Something went wrong in the flow. The last interaction was " + previousID)
	}
	if lastUserInput != "" {
		keys = append(keys, lastUserInput)
	}

	// fmt.Println("Completed", performedInteractions)
	return
}

func UpdatePromptInputs(conf map[string]string) {
	v := reflect.ValueOf(&PromptInput).Elem()
	for key, value := range conf {
		field := v.FieldByName(key)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}

		field.SetString(value)
	}
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

func initializePromptInputs() {
	m := make(map[string]string)
	conf := config.GetConfig()
	m["ProductName"] = conf.ProductName
	m["ExecutableName"] = conf.ExecutableName
	if conf.ExecutableName == "chef" {
		m["UnitMeasure"] = "nodes"
	} else {
		m["UnitMeasure"] = "targets"
	}
	UpdatePromptInputs(m)
}

func getIntractions() map[string]ActionDetail {
	var intr Interaction
	err := yaml.Unmarshal(interactionsYAML, &intr)
	if err != nil {
		log.Fatal(err)
	}
	return intr.Actions
}

func performInteraction(action ActionDetail) (nextID string) {
	var methodName string
	if action.PromptType != "" {
		methodName = action.PromptType
	} else if action.Action != "" {
		methodName = action.Action
	}

	meth := reflect.ValueOf(action).MethodByName(methodName)
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

		err = tmpl.Execute(os.Stdout, PromptInput)
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
	isValid := ValidateKeyFormat(key)
	if isValid {
		return nil
	} else {
		return fmt.Errorf("%s: %w", key, ErrInvalidKeyFormat)
	}
}

func getLicense() api.LicenseClient {
	spinner, err := spinner.GetSpinner()
	if err != nil {
		log.Println("Unable to start the spinner")
	}
	_ = spinner.Start()
	spinner.Message("In progress")
	client, _ := api.GetClient().GetLicenseClient([]string{PromptInput.LicenseID})

	spinner.StopCharacter("✓")
	spinner.StopColors("green")

	spinner.Message("Done")

	_ = spinner.Stop()

	return *client
}
