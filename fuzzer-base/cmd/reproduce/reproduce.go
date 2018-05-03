package main

import (
  "os"
  "time"

  "maxfuzz/fuzzer-base/internal/helpers"

  "github.com/howeyc/fsnotify"
  "github.com/adjust/rmq"
)

var log = helpers.BasicLogger()

func main() {
  slaveWatcher, err := fsnotify.NewWatcher()
  helpers.Check("Unable to create slave watcher: %v", err)
  masterWatcher, err := fsnotify.NewWatcher()
  helpers.Check("Unable to create master watcher: %v", err)

  go helpers.WatchFile(masterWatcher);
  go helpers.WatchFile(slaveWatcher);

  os.MkdirAll("/root/fuzz_out/master/crashes", 0755)
  os.MkdirAll("/root/fuzz_out/slave/crashes", 0755)

  // Add crash and hang directories to file watchers
  err = masterWatcher.Watch("/root/fuzz_out/master/crashes")
  helpers.Check("Error watching folder: %v", err)

  err = slaveWatcher.Watch("/root/fuzz_out/slave/crashes")
  helpers.Check("Error watching folder: %v", err)

  connection := rmq.OpenConnection(
    "crash stream", "tcp",
    helpers.Getenv("REDIS_QUEUE_URL", ""), 1)

  crashQueue := connection.OpenQueue(helpers.GetFuzzerName())

  crashConsumer := &(CrashConsumer{})
  crashQueue.StartConsuming(1, time.Second)
  crashQueue.AddConsumer("crash consumer", crashConsumer)

  select {}
}
