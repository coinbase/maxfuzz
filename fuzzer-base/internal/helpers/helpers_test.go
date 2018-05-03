// +build unit

package helpers

import (
  "testing"
  "os"
  "strings"

  "github.com/stretchr/testify/assert"
)

type Testcase struct {
  in string
  out string
  kind string
}

func TestGetenv(t *testing.T) {
  os.Setenv("TEST_ENV", "set")
  result := Getenv("TEST_ENV", "unset")
  assert.Equal(t, result, "set")

  os.Unsetenv("TEST_ENV")
  result = Getenv("TEST_ENV", "unset")
  assert.Equal(t, result, "unset")
}

var generatetests = []Testcase {
  {"/root/fuzz_out/master/crashes/testcase", "test/crashes/master_testcase", "crashes"},
  {"/root/fuzz_out/slave1/hangs/hang1", "test/hangs/slave1_hang1", "hangs"},
}

func TestGenerateTestcaseName(t *testing.T) {
  for _, tt := range generatetests {
    result, kind := GenerateTestcaseName(tt.in)
    //Strip the timestamp
    stripped := strings.Join(strings.Split(result, "_")[:2], "_")
    assert.Equal(t, stripped, tt.out)
    assert.Equal(t, kind, tt.kind)
  }
}

// Backup helpers not tested, as they simply pass through to a corresponding
// method in either filesystem_helpers.go or s3_helpers.go
