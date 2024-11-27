package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/core/entities"
)

func NewPlayerDynamoRepository(client *dynamodb.Client) BaseDynamoRepository[entities.Player] {
	return NewBaseDynamoRepository[entities.Player](client, "Players")
}
