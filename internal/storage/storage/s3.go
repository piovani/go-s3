package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/piovani/go-class/internal/config"
)

type S3Client struct {
	client     *s3.Client
	bucketName string
	endpoint   string
}

func NewS3Client(ctx context.Context, cfg *config.Config) (*S3Client, error) {
	staticProvider := credentials.NewStaticCredentialsProvider(
		cfg.AWSAccessKey,
		cfg.AWSSecretKey,
		"", // Session token
	)

	awsCfg, err := awsCfg.LoadDefaultConfig(ctx,
		awsCfg.WithRegion(cfg.AWSRegion),
		awsCfg.WithCredentialsProvider(staticProvider),
	)

	if err != nil {
		return nil, err
	}

	awsS3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWSEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWSEndpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Client{
		client:     awsS3Client,
		bucketName: cfg.BucketName,
		endpoint:   cfg.AWSEndpoint,
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, fileName string, fileData []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, fileName)
	if s.endpoint != "" {
		url = fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, fileName)
	}

	return url, nil
}

func (s *S3Client) List(ctx context.Context) ([]string, error) {
	files := []string{}

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			files = append(files, aws.ToString(obj.Key))
		}
	}

	return files, nil
}

func (s *S3Client) Download(ctx context.Context, fileName string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (s *S3Client) Delete(ctx context.Context, fileName string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	return err
}
