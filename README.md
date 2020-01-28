# go-libs
[![Build status](https://badge.buildkite.com/19949c499939e46053e5b4c573d7e6bba9a0b78a870a07501b.svg)](https://buildkite.com/chef/chef-go-libs-master-verify)
[![Code coverage](https://img.shields.io/badge/coverage-97.4%25-brightgreen)](https://buildkite.com/chef/chef-go-libs-master-code-coverage)

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

## `distgen`
A simple generator for creating easily distributable Go packages.

### Simple usage:
```go
package main
//go:generate go run github.com/chef/go-libs/distgen
```
Look at the [distgen README](distgen/README.md) for more examples.
