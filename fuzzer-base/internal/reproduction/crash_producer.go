package reproduction

// Defines methods for opening up a connection to a Redis Queue, and sending
// crash definitions (crash.go) so a crash reproducer can pick them up.

import (
  "fmt"
  "encoding/json"
  "os"

  "github.com/sirupsen/logrus"
  "github.com/adjust/rmq"
)

var log = &logrus.Logger{
  Out: os.Stderr,
  Formatter: new(logrus.JSONFormatter),
  Hooks: make(logrus.LevelHooks),
  Level: logrus.DebugLevel,
}

func NewProducer(redisUrl, fuzzerName string) rmq.Queue  {
  var connection rmq.Connection
  if (os.Getenv("MAXFUZZ_ENV") == "test" || os.Getenv("NO_REPRODUCTION") == "1") {
    connection = rmq.NewTestConnection()
  } else {
    connection = rmq.OpenConnection("crash stream", "tcp", redisUrl, 1)
  }
  crashQueue := connection.OpenQueue(fuzzerName)
  return crashQueue
}

func Produce(q rmq.Queue, filename string, kind string) {
  crash := Crash{filename, kind}
  crashBytes, err := json.Marshal(crash)
  if (err != nil) {
    log.WithFields(
      logrus.Fields{"message": fmt.Sprintf("Could not marshal Crash struct into json: %v", err)},
    ).Fatal()
    os.Exit(1)
  }
  q.PublishBytes(crashBytes)
}
