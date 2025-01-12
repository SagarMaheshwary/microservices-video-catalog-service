package aws

import (
	"time"

	awslib "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
)

func NewSession() (*session.Session, error) {
	c := config.Conf.AWS

	s, err := session.NewSession(&awslib.Config{
		Region:      awslib.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	})

	if err != nil {
		logger.Error("Unable to create aws session: %v", err)

		return nil, err
	}

	return s, nil
}

func CreateGetObjectPresignedUploadUrl(key string) (string, error) {
	c := config.Conf.AWS
	s, err := NewSession()

	if err != nil {
		return "", err
	}

	svc := s3.New(s)

	r, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: awslib.String(c.S3Bucket),
		Key:    awslib.String(key),
	})

	url, err := r.Presign(
		time.Duration(time.Duration(c.S3PresignedURLExpiry) * time.Minute),
	)

	if err != nil {
		logger.Error("Unable to create presigned upload url: %v", err)

		return "", err
	}

	return url, err
}
