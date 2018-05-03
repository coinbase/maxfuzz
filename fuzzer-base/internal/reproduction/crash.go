package reproduction

// Serialized and sent to/from Redis to pass around details for testcases to
// reproduce.

type Crash struct {
  Filename string
  Kind string
}
