package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(viper.GetString("aws.region")),
			Credentials: credentials.NewStaticCredentials(
				viper.GetString("aws.credentials.accessKeyId"),     // id
				viper.GetString("aws.credentials.secretAccessKey"), // secret
				"",
			), // token can be left blank for now
		},
	)
	if err != nil {
		panic(err)
	}
	return sess
}

func uploadFileAWS(header *multipart.FileHeader) string {

	fileS3, err := header.Open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(viper.GetString("aws.region")),
		Key:                  aws.String(header.Filename),
		Body:                 fileS3,
		ACL:                  aws.String("public-read"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return aws.StringValue(&result.Location)
}

func downloadFileAWS(fileName string) {
	//First, create a temporary file
	f, err := os.Create(fileName)
	if err != nil {
		// Do your error handling here
		return
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("aws.region")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		// Do your error handling here
		return
	}
}
