package telemetry

import (
	"fmt"
)

type Telemetry struct {
	Payload_dir                  string
	Session_file                 string
	Installation_identifier_file string
	Enabled                      bool
	Dev_mode                     bool
	Host_os                      string
	Arch                         string
	Version                      string
}

func (t Telemetry) Setup() {

	// TODO validate required & correct keys
	// payload_dir #required
	// session_file # required
	// installation_identifier_file # required
	// enabled  # false, not required
	// dev_mode # false, not required
	fmt.Println("testing the result-----")
	fmt.Println(t)
	telemetry := t
	startUploadThread(telemetry)
}

func enabled(t Telemetry) bool {
	// def enabled?
	//   require_relative "telemetry/decision"
	//   config[:enabled] && !Telemetry::Decision.env_opt_out?
	// end
	return t.Enabled && envOptOut()

}
