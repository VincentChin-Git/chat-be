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

type uploadRes struct {
	UploadUrl string `json:"uploadUrl,omitempty"`
	ViewUrl   string `json:"viewUrl,omitempty"`
}

func uploadSignature(filename string, validExts []string, errMsg string) (uploadRes, error) {

	if filename == "" {
		return uploadRes{}, errors.New("Filename Error")
	}

	// validate filetype
	extension := filepath.Ext(filename)
	if extension == "" {
		return uploadRes{}, errors.New("Filename Error")
	}

	extension = extension[1:]
	if extension == "" {
		return uploadRes{}, errors.New("Filename Error")
	}

	isExtValid := false

	for _, ext := range validExts {
		if ext == extension {
			isExtValid = true
		}
	}
	if !isExtValid {
		return uploadRes{}, errors.New(errMsg)
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
		return uploadRes{}, errors.New("")
	}
	fullPath := configGet.AwsBucketPrefix + filename

	return uploadRes{UploadUrl: url, ViewUrl: fullPath}, nil
}

func UploadImgSignature(filename string) (uploadRes, error) {

	// config
	validExts := []string{"jpg", "JPG", "jpeg", "JPEG", "png", "PNG", "pneg", "PNEG"}
	errMsg := "Only file type of jpg, png are allowed"

	return uploadSignature(filename, validExts, errMsg)
}

func UploadVideoSignature(filename string) (uploadRes, error) {

	// config
	validExts := []string{"mp4", "MP4"}
	errMsg := "Only file type of mp4 are allowed"

	return uploadSignature(filename, validExts, errMsg)
}
