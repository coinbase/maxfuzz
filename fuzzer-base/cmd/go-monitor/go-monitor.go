// Manages crash file syncing & uploading for go fuzzers

package main
import (
  "time"

  "maxfuzz/fuzzer-base/internal/helpers"

  "github.com/howeyc/fsnotify"
)

func main() {
  // Setup file watchers & uploaders
  crashWatcher, err := fsnotify.NewWatcher()
  helpers.Check("Unable to create crash watcher: %v", err)

  go helpers.WatchFile(crashWatcher);

  // Wait for fuzzers to initialize
  for (!helpers.Exists("/root/fuzz_out/crashers")) {
    time.Sleep(time.Second * 10)
  }

  // Add crash and hang directories to file watchers
  err = crashWatcher.Watch("/root/fuzz_out/crashers")
  helpers.Check("Error watching folder: %v", err)

  // Ensure we backup the fuzz_out dir regularly
  go helpers.RegularBackup("/root/fuzz_out")

  // For now, let go-fuzz log itself. It's not JSON, but it'll do.
  // TODO: Write JSON logging for go-fuzz
  select {}
}
