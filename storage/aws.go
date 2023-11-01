package storage

import (
	"chat-be/config"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var AwsInstance *s3.S3

func SetupAws() {

	configGet := config.GetConfig()

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(configGet.AwsAccessKey, configGet.AwsSecretKey, ""),
		Region:           aws.String("ap-southeast-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})

	if err != nil {
		fmt.Println(err.Error(), "setup aws session error")
	}

	AwsInstance = s3.New(sess)

}
