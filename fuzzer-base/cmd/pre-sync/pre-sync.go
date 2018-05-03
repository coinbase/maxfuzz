// Syncs the fuzzer upon startup so it restarts from the last saved fuzzer state

package main
import (
  "fmt"
  "os"

  "maxfuzz/fuzzer-base/internal/helpers"

  "github.com/pierrre/archivefile/zip"
  log "github.com/sirupsen/logrus"
)

func main() {
  log.SetFormatter(&log.JSONFormatter{})
  fuzzer := helpers.GetFuzzerName()
  outFilePath := fmt.Sprintf("/root/%v_backup.zip", fuzzer)
  uploadName := fmt.Sprintf("%v_backup.zip", fuzzer)
  progress := func(path string) {
    //Don't be verbose when uncompressing
    //fmt.Println("Uncompressing %v", path)
  }
  if (helpers.BackupExists(uploadName)) {
    err := os.Remove("/root/fuzz_out")
    helpers.Check("Unable to remove /root/fuzz_out: %v", err)

    helpers.GetBackup(uploadName, outFilePath)
    err = zip.UnarchiveFile(outFilePath, "/root/", progress)
    helpers.Check("Unable to extract AFL state: %v", err)

    err = os.Remove(outFilePath)
    helpers.Check("Unable to remove afl state zip: %v", err)

    aflIoOptions := "/root/config/afl-io-options"
    err = os.Remove(aflIoOptions)
    helpers.Check("Unable to remove afl-io-options: %v", err)

    f, err := os.Create(aflIoOptions)
    helpers.Check("Unable to create afl-io-options: %v", err)

    defer f.Close()
    toWrite := []byte("-i- -o /root/fuzz_out")
    _, err = f.Write(toWrite)
    helpers.Check("Unable to write afl-io-options: %v", err)

    log.WithFields(log.Fields{"message": "Wrote new config to afl-io-options"}).Info()
    log.WithFields(log.Fields{"message": fmt.Sprintf("Downloaded backup (%v)", uploadName)}).Info()
  } else {
    log.WithFields(log.Fields{"message": fmt.Sprintf("No backup found (%v) starting fuzzing from scratch", uploadName)}).Info()
  }
}
