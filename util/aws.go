package util

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// https://qiita.com/daijinload/items/dff973943d2a4967c78a

func DownloadFile(name string) error {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), "")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1"),
	}))

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String("kyopro-progress-watcher"),
		Key:    aws.String("/" + name),
	})
	if err != nil {
		return err
	}

	fmt.Printf("file downloaded, %d bytes\n", n)

	return nil
}

func UploadFile(name string) error {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	privateKey := os.Getenv("AWS_SECRET_KEY")
	region := "ap-northeast-1"
	fileName := name
	bucketName := "kyopro-progress-watcher"

	creds := credentials.NewStaticCredentials(accessKey, privateKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("/" + fileName),
		Body:   f,
	})
	if err != nil {
		return err
	}

	fmt.Println("Successfully upload file")

	return nil
}
