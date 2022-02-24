package telemetry

const OPT_OUT_FILE = "telemetry_opt_out"
const OPT_IN_FILE = "telemetry_opt_in"

func optOut() bool {
	// We check that the user has made a decision so that we can have a default setting for robots
	return userOptedOut() || envOptOut() || localOptOut() || made()
}

func made() bool {
	return true
}

func userOptedOut() bool {
	return true
}

func userOptedIn() bool {
	return true
}

func envOptOut() bool {
	return true
}

func localOptOut() bool {
	return true
}
