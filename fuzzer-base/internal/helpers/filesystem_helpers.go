package helpers

import (
  "fmt"
  "time"

  "github.com/spf13/afero"
  "github.com/pierrre/archivefile/zip"
  "github.com/sirupsen/logrus"
)

var localBackupDir = "/root/sync"

func preSyncChecks(fs afero.Fs) {
  exists, err := afero.Exists(fs, fmt.Sprintf("/root/sync/%v", fuzzer))
  Check("File existence check fail: %v", err)
  if(!exists) {
    fuzzerDir := fmt.Sprintf("/root/sync/%v", fuzzer)
    err := fs.MkdirAll(fmt.Sprintf("%v/crashes", fuzzerDir), 0755)
    Check("Could not create fuzzer sync directory: %v", err)
    err = fs.Mkdir(fmt.Sprintf("%v/hangs", fuzzerDir), 0755)
    Check("Could not create fuzzer sync directory: %v", err)
    err = fs.Mkdir(fmt.Sprintf("%v/leaks", fuzzerDir), 0755)
    Check("Could not create fuzzer sync directory: %v", err)
    err = fs.Mkdir(fmt.Sprintf("%v/timeouts", fuzzerDir), 0755)
    Check("Could not create fuzzer sync directory: %v", err)
  }
}

func filesystemSync(fs afero.Fs, location string, destination string, log *logrus.Logger) {
  preSyncChecks(fs)

  exists, err := afero.Exists(fs, location)
  Check("File existence check fail: %v", err)
  if(exists) {
    data, err := afero.ReadFile(fs, location)
    Check("Error reading file from fuzzer output dir: %v", err)
    err = afero.WriteFile(fs, fmt.Sprintf("/root/sync/%v", destination), data, 0755)
    Check("Error writing file to sync dir: %v", err)
    log.WithFields(
      logrus.Fields{"message": fmt.Sprintf("Synced file: %s", location)},
    ).Info()
  }
}

func filesystemDownload(fs afero.Fs, location string, destination string, log *logrus.Logger) {
  preSyncChecks(fs)
  actualLocation := fmt.Sprintf("%s/%s", localBackupDir, location)

  exists, err := afero.Exists(fs, actualLocation)
  Check("File existence check fail: %v", err)
  if (exists) {
    data, err := afero.ReadFile(fs, actualLocation)
    Check("Error reading file from sync dir: %v", err)
    err = afero.WriteFile(fs, destination, data, 0755)
    Check("Error saving file to local dir: %v", err)
    log.WithFields(
      logrus.Fields{"message": fmt.Sprintf("Downloaded file: %s", location)},
    ).Info()
  }
}

func filesystemBackupExists(fs afero.Fs, filename string) bool {
  result, err := afero.Exists(
    fs,
    fmt.Sprintf("%s/%s", localBackupDir, filename),
  )
  Check("File existence check fail: %v", err)
  return result
}

func filesystemRegularBackup(fs afero.Fs, localBackupDir string, log *logrus.Logger) {
  time.Sleep(10*time.Minute)
  outFilePath := fmt.Sprintf("/root/%v_backup.zip", fuzzer)
  progress := func(path string) {
    // Add prints here to show zip progress
  }
  err = zip.ArchiveFile(localBackupDir, outFilePath, progress)
  Check("Unable to zip backup directory: %v", err)
  filesystemSync(fs, outFilePath, fmt.Sprintf("%v_backup.zip", fuzzer), log)
}
