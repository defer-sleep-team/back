package s3_service

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetSignedObjectURL генерирует подписанную ссылку на объект в S3
func (c Cloud) GetSignedObjectURL(bucketName string, objectKey string) (string, error) {
	req := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	// Генерируем подписанную ссылку
	presignClient := s3.NewPresignClient(c.S3)
	presignedRequest, err := presignClient.PresignGetObject(context.TODO(), req, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		log.Print("Error getting link: ", err)
		return "", err
	}
	return presignedRequest.URL, nil
}
