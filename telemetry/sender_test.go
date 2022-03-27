package telemetry

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var Expected = []string{"/Users/ngupta/.chef-workstation/telemetry/telemetry-payload-1.yml"}


	var telTests =  struct {
	PayloadDir                 string
	SessionFile                string
	InstallationIdentifierFile string
	Enabled                    bool
	DevMode                    bool
	HostOs                     string
	Arch                       string
	WorkstationVersion         string
}{
	PayloadDir: "/Users/ngupta/.chef-workstation/telemetry",
	SessionFile: "/Users/ngupta/.chef-workstation/telemetry/TELEMETRY_SESSION_ID",
	InstallationIdentifierFile: "/Users/ngupta/.chef-workstation/installation_id",
	Enabled: true,
	DevMode: false,
	HostOs: "darwin",
	Arch: "amd64",
	WorkstationVersion: "22.2.802",
}

func TestFindSessionFiles(t *testing.T) {
	output := findSessionFiles(telTests)
	//fmt.Println(output)
	if output[0] == Expected[0] {
		fmt.Println("PASS")
	}
}

func TestRun(t *testing.T) {
	err := run(telTests, Expected)
	if err!= nil {
		fmt.Println(err)
	}
}

func TestLoadAndClearSession(t *testing.T) {
	filename, _ := filepath.Abs(Expected[0])
	yfile, err := ioutil.ReadFile(filename)
	var config TelemetryPayload

	if err != nil {

		log.Fatal(err)
	}
	err2 := yaml.Unmarshal(yfile, &config)

	if err2 != nil {

		log.Fatal(err2)
	}
	var result TelemetryPayload
	result = loadAndClearSession(Expected)
	if &result != &config {
		fmt.Println(err)
	}

}

