package main

import (
	"fmt"

	licensing "github.com/chef/go-libs/licensing"
)

func main() {
	c := &licensing.LicenseConfig{
		ProductName:      "Workstation",
		EntitlementID:    "x6f3bc76-a94f-4b6c-bc97-4b7ed2b045c0",
		LicenseServerURL: "https://licensing-acceptance.chef.co/License",
	}
	licensing.SetConfig(c)
	// licensing.CheckSoftwareEntitlement()
	fmt.Println(licensing.FetchAndPersist())
	// fmt.Println(*licensing.ReadLicenseKeyFile())
	// fmt.Println(licensing.AllowedLicencesForAddition())
	// cc := licensing.ReadLicenseKeyFile()
	// fmt.Println(cc)
	// licensing.IsLicenseActive([]string{"tmns-b215381f-44df-4bbc-99f1-dda03b7dc637-2983"})
}
