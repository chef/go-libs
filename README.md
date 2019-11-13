# go-libs
[![Build status](https://badge.buildkite.com/a5dfa44b20a6ec189a93bcbda031db452f1d964fa6836f7065.svg)](https://buildkite.com/chef/chef-go-libs-master-verify)
[![Code coverage](https://img.shields.io/badge/coverage-0.0%25-brightgreen)](https://buildkite.com/chef/chef-go-libs-master-verify)

A collection of Go libraries used across the Chef ecosystem

## `featflag`
An implementation to manage feature flags in Go.

### Simple usage:
```go
package main

import "github.com/chef/go-libs/featflag"

func main() {
	fmt.Printf("A global feature flag '%s' that is currently: %s",
		featflag.ChefFeatAnalyze.String(),
		featflag.ChefFeatAnalyze.Enabled(),
	)
}
```
Look at [featflag/example](featflag/example) for a working example.

## `config`
An abstraction of the Chef Workstation configuration file. (`config.toml`)

### Simple usage:
```go
package main

import "github.com/chef/go-libs/config"

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println("unable to read the config", err)
	}

	fmt.Println("the log level of my config is: ", cfg.Log.Level)
}
```

## `credentials`
An abstraction of the Chef credentials file. (`credentials`)

### Simple usage:
```go
package main

import "github.com/chef/go-libs/credentials"

func main() {
	profile := "prod"
	creds, err := credentials.New(profile)
	if err != nil {
		fmt.Println("unable to read the credentials", err)
	}

	fmt.Printf("The client_name from the '%s' profile is: %s", creds.ClientName)
}
```
