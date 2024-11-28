package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
)

func InitAWSClient() *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		//config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create the client with a local endpoint
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT"))
		o.UsePathStyle = true // Important for LocalStack
	})

	// Test s3 connection
	_, err = s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("unable to list buckets, %v", err)
	}

	return s3Client
}

func CreateBuckets(s3Client *s3.Client) {
	if os.Getenv("ENV") == "local" {
		// if Bucket exists, skip
		const b = "fg-analyzer-replay-uploads"
		_, err := s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
			Bucket: aws.String(b),
		})
		if err == nil {
			log.Println("Bucket already exists, skipping")
			return
		}

		_, err = s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String(b),
			CreateBucketConfiguration: &s3Types.CreateBucketConfiguration{
				LocationConstraint: s3Types.BucketLocationConstraint(os.Getenv("AWS_REGION")),
			},
		})
		if err != nil {
			log.Fatalf("unable to create bucket, %v", err)
		}
	}
}
