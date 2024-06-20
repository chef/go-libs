package licensing

import (
	"fmt"

	"github.com/gookit/color"
)

// TODO: Fix the license overview
func printLicenseKeyOverview(keys []string) {
	client := fetchLicenseClient(keys)

	fmt.Println("------------------------------------------------------------")
	color.Bold.Println("License Details")
	color.Println("License key: ", client.License)
}
