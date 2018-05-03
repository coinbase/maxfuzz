package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "os"
  "strings"

  "maxfuzz/fuzzer-base/internal/helpers"
  "maxfuzz/fuzzer-base/internal/reproduction"

  "github.com/adjust/rmq"
  "github.com/sirupsen/logrus"
  "github.com/go-cmd/cmd"
)

type CrashConsumer struct {}

func (consumer *CrashConsumer) Consume(delivery rmq.Delivery) {
  var crash reproduction.Crash
  err := json.Unmarshal([]byte(delivery.Payload()), &crash)
  if err != nil {
    delivery.Reject()
    return
  }
  reproduce(crash.Filename, crash.Kind)

  delivery.Ack()
}

func reproduce(filename string, kind string) {
  log.WithFields(
    logrus.Fields{
      "message": fmt.Sprintf(
        "Reproducing %s, kind: %s",
        filename,
        kind,
      ),
    },
  ).Info()

  // Download testcase to reproduce
  downloadedFile := "/root/fuzz_in/input"
  helpers.GetBackup(filename, downloadedFile)

  fuzzerInstance := ""
  if (strings.Contains(filename, "slave")) {
    fuzzerInstance = "slave"
  } else {
    fuzzerInstance = "master"
  }

  sl := strings.Split(filename, "/")
  filename = sl[len(sl)-1]

  // Check if we're passing in as an arg, or sending as stdin
  commandA := os.Getenv("AFL_BINARY")
  commandB := ""
  if (strings.Contains(commandA, " @@")) {
    commandA = strings.Replace(commandA, " @@", "", 1)
    commandB = downloadedFile
  } else {
    commandB = commandA
    commandA = "/root/scripts/reproduce_stdin"
  }

  log.WithFields(
    logrus.Fields{
      "message": fmt.Sprintf(
        "Reproducing with command: %s %s",
        commandA,
        commandB,
      ),
    },
  ).Info()

  command := cmd.NewCmd(commandA, commandB)
  status := <-command.Start()
  var finalStdout bytes.Buffer
  var finalStderr bytes.Buffer

  finalStdout.WriteString("STDOUT:\n")
  for _, line := range status.Stdout {
		finalStdout.WriteString(fmt.Sprintf("%s\n", line))
	}

  finalStdout.WriteString("\nSTDERR:\n")
  for _, line := range status.Stderr {
    finalStderr.WriteString(fmt.Sprintf("%s\n", line))
  }

  finalFilename := fmt.Sprintf("/root/fuzz_out/%s/crashes/%s", fuzzerInstance, filename)
  f, err := os.Create("/root/tmp_out")
  helpers.Check("Could not create temporary out file: %v", err)

  numBytes, err := f.Write(finalStdout.Bytes())
  numBytes, err = f.Write(finalStderr.Bytes())
  helpers.Check("Could not write output to file: %v", err)
  log.WithFields(
    logrus.Fields{
      "message": fmt.Sprintf(
        "Wrote %d bytes to file",
        numBytes,
      ),
    },
  ).Info()
  f.Sync()
  f.Close()

  log.WithFields(
    logrus.Fields{
      "message": fmt.Sprintf(
        "Saving file: %s",
        finalFilename,
      ),
    },
  ).Info()
  err =  os.Rename("/root/tmp_out", finalFilename)
}
