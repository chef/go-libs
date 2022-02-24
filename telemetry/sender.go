package telemetry

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func startUploadThread(t Telemetry) {
	// Find the files before we spawn the thread - otherwise
	// we may accidentally pick up the current run's session file if it
	// finishes before the thread scans for new files
	sessionFiles := findSessionFiles(t)
	// sender := Sender.new(session_files, config)
	// Thread.new{sender.run}
	fmt.Println("-----------sessionFiles-------")
	fmt.Println(sessionFiles)
	run(t, sessionFiles)

}

func findSessionFiles(t Telemetry) []string {
	fmt.Println("Looking for telemetry data to submit")
	sessionSearch := path.Join(t.Payload_dir, "telemetry-payload-*.yml")
	sessionFiles, _ := filepath.Glob(sessionSearch)
	fmt.Println("-----------sessionSearch-------", sessionSearch)
	fmt.Println("-----------sessionFiles-------", sessionFiles)

	fmt.Println("Found #{session_files.length} sessions to submit")
	return sessionFiles

}

func run(t Telemetry, sessionFiles []string) {
	fmt.Println("t is ------", t)
	fmt.Println("sessionfile is ------", sessionFiles)
	if enabled(t) {
		fmt.Println("Telemetry enabled, beginning upload of previous session(s)")
		if t.Dev_mode {
			os.Setenv("CHEF_TELEMETRY_ENDPOINT", "https://telemetry-acceptance.chef.io")
		}
		for i := 0; i < len(sessionFiles); i++ {
			fmt.Println("Array lenth ------", sessionFiles[i])
			processSession(sessionFiles[i])
		}

	} else {
		// If telemetry is not enabled, just clean up and return. Even though
		// the telemetry gem will not send if disabled, log output saying that we're submitting
		// it when it has been disabled can be alarming.
		fmt.Println("Telemetry disabled, clearing any existing session captures without sending them.")
		for i := 0; i < len(sessionFiles); i++ {
			fmt.Println("Array lenth ------", sessionFiles[i])
			err := os.RemoveAll(sessionFiles[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err := os.RemoveAll(t.Session_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Terminating, nothing more to do.")
}

func processSession(sessionFilePath string) {

	// 	def process_session(path)
	// 	require 'byebug'
	// 	byebug
	// 	Telemeter::Log.info("Processing telemetry entries from #{path}")
	// 	content = load_and_clear_session(path)
	// 	submit_session(content)
	//   end
	fmt.Println("Processing telemetry entries from")
	// content := loadAndClearSession(sessionFilePath)
	// submitSession(content)

}

func loadAndClearSession(sessionFilePath string) {

	// def load_and_clear_session(path)
	//     require 'byebug'
	//     byebug
	//     content = File.read(path)
	//     We'll remove it now instead of after we parse or submit it -
	//     if we fail to deliver, we don't want to be stuck resubmitting it if the problem
	//     was due to payload. This is a trade-off - if we get a transient error, the
	//     payload will be lost.
	//     TODO: Improve error handling so we can intelligently decide whether to
	//            retry a failed load or failed submit.
	//     FileUtils.rm_rf(path)
	//     YAML.load(content)
	//   end
	content, err := ioutil.ReadFile("file.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:")
	fmt.Println(string(content))

	err = os.RemoveAll(sessionFilePath)
	if err != nil {
		log.Fatal(err)
	}
	// YAML.load(content)

}
