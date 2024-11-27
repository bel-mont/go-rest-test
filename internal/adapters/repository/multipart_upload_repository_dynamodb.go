package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/core/entities"
)

func NewMultipartUploadDynamoRepository(client *dynamodb.Client) BaseDynamoRepository[entities.MultipartUpload] {
	return NewBaseDynamoRepository[entities.MultipartUpload](client, "MultipartUploads")
}
