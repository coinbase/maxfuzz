// Manages crash file uploading & syncing, and logging, for all AFL fuzzers

package main
import (
  "fmt"
  "os"
  "time"
  "io/ioutil"
  "bufio"
  "strconv"
  "strings"

  "maxfuzz/fuzzer-base/internal/helpers"

  "github.com/howeyc/fsnotify"
  log "github.com/sirupsen/logrus"
)

func main() {

  // Setup file watchers & uploaders
  slaveWatcher, err := fsnotify.NewWatcher()
  helpers.Check("Unable to create slave watcher: %v", err)
  masterWatcher, err := fsnotify.NewWatcher()
  helpers.Check("Unable to create master watcher: %v", err)

  go helpers.WatchFile(masterWatcher);
  go helpers.WatchFile(slaveWatcher);

  // Wait for fuzzers to initialize
  for (
    !helpers.Exists("/root/fuzz_out/master/crashes") ||
    !helpers.Exists("/root/fuzz_out/slave/crashes")) {
    time.Sleep(time.Second * 10)
  }

  // Add crash and hang directories to file watchers
  err = masterWatcher.Watch("/root/fuzz_out/master/crashes")
  helpers.Check("Error watching folder: %v", err)
  err = masterWatcher.Watch("/root/fuzz_out/master/hangs")
  helpers.Check("Error watching folder: %v", err)

  err = slaveWatcher.Watch("/root/fuzz_out/slave/crashes")
  helpers.Check("Error watching folder: %v", err)
  err = slaveWatcher.Watch("/root/fuzz_out/slave/hangs")
  helpers.Check("Error watching folder: %v", err)

  // Ensure we backup the fuzz_out dir regularly
  go helpers.RegularBackup("/root/fuzz_out")

  // Ensure the fuzzers are outputting stats before we start logging
  for (
    !helpers.Exists("/root/fuzz_out/master/fuzzer_stats") ||
    !helpers.Exists("/root/fuzz_out/slave/fuzzer_stats")) {
    time.Sleep(time.Second * 5)
  }

  // Setup afl stats logging
  log.SetFormatter(&log.JSONFormatter{})
  for {
    liveFuzzers, err := ioutil.ReadDir("/root/fuzz_out")
    helpers.Check("Failed to read /root/fuzz_out %v", err)

    for _, fuzzerInstance := range liveFuzzers {
      summary_map := make(map[string]string)
      file, err := os.Open(
        fmt.Sprintf("/root/fuzz_out/%v/fuzzer_stats", fuzzerInstance.Name()),
      )
      helpers.Check("Opening fuzzer stats failed %v", err)
      defer file.Close()

      // This adds lines like "key : val" to summary_map[key] = val
      scanner := bufio.NewScanner(file)
      for scanner.Scan() {
        spl := strings.Split(scanner.Text(), ":")
        k := strings.TrimSpace(spl[0])
        v := strings.TrimSpace(spl[1])
        summary_map[k] = v
      }

      cycles, err := strconv.Atoi(summary_map["cycles_done"])
      execs, err := strconv.Atoi(summary_map["execs_done"])
      execsPerSecond, err := strconv.ParseFloat(
        summary_map["execs_per_sec"],
        64,
      )
      crashes, err := strconv.Atoi(summary_map["unique_crashes"])
      hangs, err := strconv.Atoi(summary_map["unique_hangs"])
      pending_faves, err := strconv.Atoi(summary_map["pending_favs"])
      pending_total, err := strconv.Atoi(summary_map["pending_total"])

      log.WithFields(log.Fields{
        "fuzzer": summary_map["afl_banner"],
        "cycles_done": cycles,
        "execs_done": execs,
        "execs_per_second": execsPerSecond,
        "crashes": crashes,
        "hangs": hangs,
        "pending_faves": pending_faves,
        "pending_total": pending_total,
      }).Info()
    }
    time.Sleep(30*time.Second)
  }
}
