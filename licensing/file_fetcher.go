package licensing

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	FILE_VERSION = "4.0.0"
)

var LICENSE_TYPES []string = []string{"free", "trial", "commercial"}

type LicenseFileData struct {
	Licenses          []Licenses `yaml:":licenses"`
	FileFormatVersion string     `yaml:":file_format_version"`
	LicenseServerURL  string     `yaml:":license_server_url"`
}

type Licenses struct {
	LicenseKey  string `yaml:":license_key"`
	LicenseType string `yaml:":license_type"`
	UpdateTime  string `yaml:":update_time"`
}

func licenseFileFetch() []string {
	licenseKey := []string{}
	li := *ReadLicenseKeyFile()

	for i := 0; i < len(li.Licenses); i++ {
		licenseKey = append(licenseKey, li.Licenses[i].LicenseKey)
	}

	return licenseKey
}

func licenseFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".chef/licenses.yaml")
}

func ReadLicenseKeyFile() *LicenseFileData {
	li := &LicenseFileData{}
	filePath := licenseFilePath()
	info, _ := os.Stat(filePath)
	if info == nil {
		return li
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &li)
	if err != nil {
		log.Fatal(err)
	}
	return li
}

func fetchLicenseKeysBasedOnType(licenseType string) (out []string) {
	content := ReadLicenseKeyFile()
	for _, key := range content.Licenses {
		if key.LicenseType == licenseType {
			out = append(out, key.LicenseKey)
		}
	}
	return
}

func persistAndConcat(newKeys []string, licenseType string) {
	if !slices.Contains(LICENSE_TYPES, licenseType) {
		log.Fatal("License type " + licenseType + " is not a valid license type.")
	}

	license := Licenses{
		LicenseKey:  newKeys[0],
		LicenseType: ":" + licenseType,
		UpdateTime:  time.Now().Format("2006-01-02T15:04:05-07:00"),
	}

	fileContent := ReadLicenseKeyFile()

	var found bool
	for _, key := range fileContent.Licenses {
		if key.LicenseKey == license.LicenseKey {
			found = true
		}
	}

	if !found {
		fileContent.Licenses = append(fileContent.Licenses, license)
	}
	updateDefaultsOnLicenseFile(fileContent)
	saveLicenseFile(fileContent)
	appendLicenseKey(newKeys[0])
}

func updateDefaultsOnLicenseFile(content *LicenseFileData) {
	if content.FileFormatVersion == "" {
		content.FileFormatVersion = FILE_VERSION
	}

	if content.LicenseServerURL == "" {
		config := GetConfig()
		content.LicenseServerURL = config.LicenseServerURL
	}
}

func saveLicenseFile(content *LicenseFileData) {
	filepath := licenseFilePath()

	data, err := yaml.Marshal(&content)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}
