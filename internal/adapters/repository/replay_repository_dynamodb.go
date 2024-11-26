package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go-rest-test/internal/core/entities"
)

func NewReplayRepository(client *dynamodb.Client) BaseDynamoRepository[entities.Replay] {
	return NewBaseDynamoRepository[entities.Replay](client, "Replays")
}
