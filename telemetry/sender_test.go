package telemetry

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var expected = []string{"/Users/ngupta/.chef-workstation/telemetry/telemetry-payload-1.yml"}


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
	if output[0] == expected[0] {
		fmt.Println("PASS")
	}
}

func TestRun(t *testing.T) {
	err := run(telTests, expected)
	if err!= nil {
		fmt.Println(err)
	}
}

func TestLoadAndClearSession(t *testing.T) {
	filename, _ := filepath.Abs(expected[0])
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
	result = loadAndClearSession(expected)
	if &result != &config {
		fmt.Println(err)
	}

}

