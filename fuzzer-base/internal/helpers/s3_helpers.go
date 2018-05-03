package helpers

import (
  "fmt"
  "time"

  "github.com/spf13/afero"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/pierrre/archivefile/zip"
  "github.com/sirupsen/logrus"
)

var bucket = aws.String(Getenv("CRASH_BUCKET", "fuzz-static-development"))
var environment = Getenv("MAXFUZZ_ENV", "test")
var region = aws.String(Getenv("AWS_REGION", "us-east-1"))
var awsSession, err = session.NewSession(
  &aws.Config{Region: region},
)
var uploader = NewS3Uploader()
var downloader = NewS3Downloader()
var service = s3.New(awsSession)

func NewS3Downloader() *s3manager.Downloader {
  downloader := s3manager.NewDownloader(awsSession)
  return downloader
}

func NewS3Uploader() *s3manager.Uploader {
  uploader := s3manager.NewUploader(awsSession)
  Check("Unable to create session: %v", err)
  return uploader
}

func s3RegularBackup(fs afero.Fs, backupDir string, log *logrus.Logger) {
  time.Sleep(10*time.Minute)
  outFilePath := fmt.Sprintf("/root/%v_backup.zip", fuzzer)
  progress := func(path string) {
    // Add prints here to shop zip progress
  }
  err = zip.ArchiveFile(backupDir, outFilePath, progress)
  Check("Unable to zip backup directory: %v", err)
  s3Upload(fs, outFilePath, fmt.Sprintf("%v_backup.zip", fuzzer), log)
  err = fs.Remove(outFilePath)
  Check("Unable to remove file: %v", err)
}

func s3Upload(fs afero.Fs, location string, key string, log *logrus.Logger) {
  file, err := fs.Open(location)
  Check("Unable to open file %q, %v", err)

  defer file.Close();

  if (environment != "test") {
    _, err = uploader.Upload(&s3manager.UploadInput{
      Bucket: bucket,
      Key: aws.String(key),
      Body: file,
    })
  }

  Check(fmt.Sprintf("Unable to upload %s to %s", location, *bucket), err)
  log.WithFields(
    logrus.Fields{"message": fmt.Sprintf("Synced file: %s", location)},
  ).Info()
}

func s3Download(fs afero.Fs, key string, location string, log *logrus.Logger) {
  file, err := fs.Create(location)
  Check("Unable to open file %q, %v", err)

  defer file.Close();

  _, err = downloader.Download(file,
    &s3.GetObjectInput{
      Bucket: bucket,
      Key: aws.String(key),
    })
  Check(fmt.Sprintf("Unable to download item %s", key), err)

  log.WithFields(
    logrus.Fields{"message": fmt.Sprintf("Downloaded file: %s", location)},
  ).Info()
}

func s3BackupExists(filename string) bool {
  resp, err := service.ListObjects(&s3.ListObjectsInput{Bucket: bucket, Prefix: &filename})
  Check(fmt.Sprintf("Unable to list items in bucket %s", *bucket), err)

  for _, item := range resp.Contents {
    fmt.Println("Name:         ", *item.Key)
    fmt.Println("Last modified:", *item.LastModified)
  }

  return len(resp.Contents) > 0
}
