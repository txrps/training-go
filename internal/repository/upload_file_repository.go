package repository

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadFileS3(s3Client *s3.Client, bucket, objectKey string, file io.Reader, size int64, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	return err
}

func GetPresignedURLs3(s3Client *s3.Client, bucket, objectKey string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	presignClient := s3.NewPresignClient(s3Client)
	result, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func DeleteFileS3(s3Client *s3.Client, bucket, objectKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	return err
}
