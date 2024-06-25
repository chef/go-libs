package cheflicensing

import keyfetcher "github.com/chef/go-libs/chef_licensing/key_fetcher"

func FetchAndPersist() []string {
	return keyfetcher.GlobalFetchAndPersist()
}
