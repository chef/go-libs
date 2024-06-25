package main

import (
	"fmt"

	cheflicensing "github.com/chef/go-libs/chef_licensing"
	licenseConfig "github.com/chef/go-libs/chef_licensing/config"
)

func main() {
	licenseConfig.SetConfig("Workstation", "x6f3bc76-a94f-4b6c-bc97-4b7ed2b045c0", "https://licensing-acceptance.chef.co/License", "chef")
	fmt.Println(cheflicensing.FetchAndPersist())
}
