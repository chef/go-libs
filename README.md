# go-libs
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
