package helpers

import (
  "fmt"
  "os"
  "strings"
  "time"

  "maxfuzz/fuzzer-base/internal/reproduction"

  "github.com/spf13/afero"
  "github.com/howeyc/fsnotify"
  "github.com/sirupsen/logrus"
)

var log = BasicLogger()
var fs = afero.NewOsFs()
var fuzzer = Getenv("FUZZER_NAME", "test")
var revision = Getenv("GIT_SHA", "no_git")
var reproductionQueue = reproduction.NewProducer(
  Getenv("REDIS_QUEUE_URL", ""),
  fuzzer,
)

func Getenv(key, def string) string {
  temp := os.Getenv(key)
  if len(temp) == 0 {
    return def
  }
  return temp
}

func GetFuzzerName() string {
  return fuzzer
}

func S3Enabled() (bool) {
  return Getenv("MAXFUZZ_ENABLE_S3", "0") == "1"
}

func BasicLogger() *logrus.Logger {
  logger := &logrus.Logger{
    Out: os.Stderr,
    Formatter: new(logrus.JSONFormatter),
    Hooks: make(logrus.LevelHooks),
    Level: logrus.DebugLevel,
  }

  return logger
}

func GenerateTestcaseName(filename string) (string, string) {
  // Define filenames as fuzzInstance_name_timestamp_gitSha
  num := int32(time.Now().Unix())
  sl := strings.Split(filename, "/")
  filename = sl[len(sl)-1]
  fuzzInstance := sl[len(sl)-3]
  testcaseType := sl[len(sl)-2]
  if (Getenv("MAXFUZZ_ENV", "fuzzer") == "reproducer") {
    // If we're reproducing filenames are already structured
    return fmt.Sprintf(
      "%s/%s/%s.output",
      fuzzer,
      testcaseType,
      filename,
    ), testcaseType
  }
  return fmt.Sprintf(
    "%v/%v/%v_%v_%v_%v",
    fuzzer,
    testcaseType,
    fuzzInstance,
    filename,
    num,
    revision,
  ), testcaseType
}

// Error helpers

func exitErrorf(msg string, args ...interface{}) {
  log.WithFields(
    logrus.Fields{"message": fmt.Sprintf(msg+"\n", args...)},
  ).Fatal()
  os.Exit(1)
}

func Check(msg string, err error) {
  if err != nil {
    exitErrorf(msg, err)
  }
}

// Backup Helpers

func BackupExists(filename string) bool {
  if (S3Enabled()) {
    return s3BackupExists(filename)
  } else {
    return filesystemBackupExists(fs, filename)
  }
}

func GetBackup(location string, destination string) {
  // Retrieve a file from backup dir or S3
  if (S3Enabled()) {
    s3Download(fs, location, destination, log)
  } else {
    filesystemDownload(fs, location, destination, log)
  }
}

func RegularBackup(fileName string) {
  // Compresses the entire fuzzer state at a regular interval and saves it
  // Either to S3 or local fileystem
  if (S3Enabled()) {
    s3RegularBackup(fs, fileName, log)
  } else {
    filesystemRegularBackup(fs, fileName, log)
  }
}

// File helpers

func WatchFile(w *fsnotify.Watcher) {
  // Watches a crash directory, and uploads when it finds a new addition
  for {
    select {
    case ev := <-w.Event:
      if ev.IsCreate() {
        if (!strings.Contains(ev.Name, "README.txt")) {
          uploadName, crashType := GenerateTestcaseName(ev.Name)
          if (S3Enabled()) {
            s3Upload(fs, ev.Name, uploadName, log)
          } else {
            filesystemSync(fs, ev.Name, uploadName, log)
          }
          if (os.Getenv("NO_REPRODUCTION") != "1" &&
          os.Getenv("MAXFUZZ_ENV") == "fuzzer"){
            reproduction.Produce(
              reproductionQueue,
              uploadName,
              crashType,
            )
          }
        }
      }
    case err := <-w.Error:
      log.Error("error:", err)
    }
  }
}

func Exists(filename string) bool {
  exists, err := afero.Exists(fs, filename)
  Check("File existence check fail: %v", err)
  return exists
}
