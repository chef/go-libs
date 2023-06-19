package licensing

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Key struct {
	Licenses          []Licenses `yaml:":licenses"`
	FileFormatVersion string     `yaml:":file_format_version"`
}

type Licenses struct {
	LicenseKey  string `yaml:":license_key"`
	LicenseType string `yaml:":license_type"`
	UpdateTime  string `yaml:":update_time"`
}

func CheckSoftwareEntitlement(softwareEntitlementID string, URL string) {
	var licenseKey []string
	home, _ := os.UserHomeDir()
	licenseFilePath := filepath.Join(home, ".chef/licenses.yaml")
	info, _ := os.Stat(licenseFilePath)
	if info != nil {
		fmt.Println("inside license fileeee")
		licenseKey = licenseFileFetch(licenseFilePath)
		client(licenseKey, softwareEntitlementID, URL)
		return
	}
	key, check := os.LookupEnv("CHEF_LICENSE_KEY")
	if check {
		fmt.Println("inside ENV var")
		licenseKey = append(licenseKey, key)
		client(licenseKey, softwareEntitlementID, URL)
		return
	}
	args := os.Args
	for k, v := range args {
		if v == "--chef-license-key" {
			fmt.Println("inside arggg1")
			if len(args) > k+1 {
				licenseKey = append(licenseKey, args[k+1])
				client(licenseKey, softwareEntitlementID, URL)
				return
			}
		} else if strings.HasPrefix(v, "--chef-license-key=") {
			fmt.Println("inside arggggg2")
			split := strings.Split(v, "=")
			licenseKey = append(licenseKey, split[len(split)-1])
			client(licenseKey, softwareEntitlementID, URL)
			return
		}
	}
	client(licenseKey, softwareEntitlementID, URL)
}

func licenseFileFetch(licenseFilePath string) []string {
	data, err := ioutil.ReadFile(licenseFilePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var li Key
	err = yaml.Unmarshal(data, &li)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	licenseKey := []string{}

	for i := 0; i < len(li.Licenses); i++ {
		licenseKey = append(licenseKey, li.Licenses[i].LicenseKey)
	}

	return licenseKey

}

func client(licenseKey []string, softwareEntitlementID string, URL string) {
	if len(licenseKey) == 0 {
		fmt.Println("You dont have license key, Please generate by running chef license command")
		os.Exit(0)
	} else {
		var opts = make(map[string]string)
		opts["licenseId"] = strings.Join(licenseKey, ",")
		opts["entitlementId"] = softwareEntitlementID
		invokeGetAPI(opts, URL)
	}
}
