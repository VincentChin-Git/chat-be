package services

import (
	"chat-be/config"
	"chat-be/storage"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadSignature(filename string, validExts []string, errMsg string) (string, error) {

	if filename == "" {
		return "", errors.New("Filename Error")
	}

	// validate filetype
	extension := filepath.Ext(filename)
	if extension == "" {
		return "", errors.New("Filename Error")
	}

	extension = extension[1:]
	if extension == "" {
		return "", errors.New("Filename Error")
	}

	isExtValid := false

	for _, ext := range validExts {
		if ext == extension {
			isExtValid = true
		}
	}
	if !isExtValid {
		return "", errors.New(errMsg)
	}

	configGet := config.GetConfig()

	// Prepare the S3 request so a signature can be generated
	r, _ := storage.AwsInstance.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(configGet.AwsBucketName),
		Key:    aws.String(configGet.AwsBucketFolder + filename),
	})

	url, err := r.Presign(5 * time.Minute)
	if err != nil {
		fmt.Println("Failed to generate a pre-signed url: ", err)
		return "", errors.New("")
	}

	return url, nil
}

func UploadImgSignature(filename string) (string, error) {

	// config
	validExts := []string{"jpg", "JPG", "jpeg", "JPEG", "png", "PNG", "pneg", "PNEG"}
	errMsg := "Only file type of jpg, png are allowed"

	return uploadSignature(filename, validExts, errMsg)
}

func UploadVideoSignature(filename string) (string, error) {

	// config
	validExts := []string{"mp4", "MP4"}
	errMsg := "Only file type of mp4 are allowed"

	return uploadSignature(filename, validExts, errMsg)
}
