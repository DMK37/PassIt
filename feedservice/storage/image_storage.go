package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ImageStorage interface {
	UploadImage(image *multipart.FileHeader, userId string) (string, error)
}

type S3ImageStorage struct {
	S3Client *s3.Client
}

func NewS3ImageStorage() (*S3ImageStorage, error) {

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	return &S3ImageStorage{S3Client: s3Client}, nil
}

func (s *S3ImageStorage) UploadImage(image *multipart.FileHeader, userId string) (string, error) {

	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("passit-bucket"),
		Key:    aws.String(userId + "/" + image.Filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", "passit-bucket", "us-east-1", userId+"/"+image.Filename)

	return imageURL, nil
}
