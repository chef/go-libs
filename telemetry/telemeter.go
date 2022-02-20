package main

import (
	"fmt"
)

func Setup() {
	// TODO validate required & correct keys
	// :payload_dir #required
	// :session_file # required
	// :installation_identifier_file # required
	// :enabled  # false, not required
	// :dev_mode # false, not required

	cfg := map[string]string{"enabled": "true", "dev_mode": "false", "payload_dir": "/Users/ngupta/.chef-workstation/telemetry", "installation_identifier_file": "/Users/ngupta/.chef-workstation/installation_id", "session_file": "/Users/ngupta/.chef-workstation/telemetry/TELEMETRY_SESSION_ID"}
	// fmt.Println(cfg)

	fmt.Println("testing the result")
	for key, value := range cfg {
		fmt.Printf("%q is the key for the value %q\n", key, value)
	}
}
