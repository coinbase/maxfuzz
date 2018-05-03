// +build unit

package helpers

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/spf13/afero"
  "github.com/sirupsen/logrus"
)

var logger = &logrus.Logger{
  Out: os.Stderr,
  Formatter: new(logrus.JSONFormatter),
  Hooks: make(logrus.LevelHooks),
  Level: logrus.PanicLevel,
}

func TestPreSyncChecks(t *testing.T) {
  testFs := afero.NewMemMapFs()
  preSyncChecks(testFs)

  // Ensure all directories are created
  result, err := afero.IsDir(testFs, "/root/sync/test/crashes")
  assert.True(t, result)
  assert.Nil(t, err)

  result, err = afero.IsDir(testFs, "/root/sync/test/hangs")
  assert.True(t, result)
  assert.Nil(t, err)

  result, err = afero.IsDir(testFs, "/root/sync/test/leaks")
  assert.True(t, result)
  assert.Nil(t, err)

  result, err = afero.IsDir(testFs, "/root/sync/test/timeouts")
  assert.True(t, result)
  assert.Nil(t, err)
}

func TestFilesystemSync(t *testing.T) {
  testFs := afero.NewMemMapFs()

  // Setup file to be copied
  content := []byte("Content")
  testFs.Mkdir("/root", 0755)
  afero.WriteFile(testFs, "/root/fileToMove", content, 0755)

  filesystemSync(testFs, "/root/fileToMove", "test/crashes/movedFile", logger)

  // Ensure contents of file is the same
  result, err := afero.ReadFile(testFs, "/root/sync/test/crashes/movedFile")
  assert.Nil(t, err)
  assert.Equal(t, result, content)
}

func TestFilesystemDownload(t *testing.T) {
  testFs := afero.NewMemMapFs()
  // Create the necessry dirs
  preSyncChecks(testFs)

  // Setup file to be copied
  content := []byte("Content")
  afero.WriteFile(testFs, "/root/sync/test/crashes/fileToMove", content, 0755)

  filesystemDownload(testFs, "test/crashes/fileToMove", "/root/movedFile", logger)

  // Ensure contents of file is the same
  result, err := afero.ReadFile(testFs, "/root/movedFile")
  assert.Nil(t, err)
  assert.Equal(t, result, content)
}

func TestFilesystemBackupExists(t *testing.T) {
  testFs := afero.NewMemMapFs()

  result := filesystemBackupExists(testFs, "test_backup.zip")
  assert.False(t, result)

  content := []byte("Content")
  afero.WriteFile(testFs, "/root/sync/test/test_backup.zip", content, 0755)
  result = filesystemBackupExists(testFs, "test/test_backup.zip")
  assert.True(t, result)
}
