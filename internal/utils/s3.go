package utils

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadToS3(file *multipart.FileHeader) (string, error) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	client := s3.NewFromConfig(cfg)

	f, _ := file.Open()
	defer f.Close()

	key := "avatars/" + uuid.New().String() + "_" + file.Filename

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("my-bucket"),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", err
	}

	url := "https://my-bucket.s3.amazonaws.com/" + key
	return url, nil
}
