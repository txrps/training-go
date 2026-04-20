package storage

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewStorageClient() *s3.Client {
	cfg := aws.Config{
		Region: "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("DELL_OBJECT_ACCESS_KEY"),
			os.Getenv("DELL_OBJECT_SECRET_KEY"),
			"",
		),
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(os.Getenv("DELL_OBJECT_ENDPOINT"))
	})
}
