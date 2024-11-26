package repository

import (
	"context"
)

// Entity interface defines the basic requirements for an entity
type Entity interface {
	GetID() string
	SetID(id string)
}

type Repository[T Entity] interface {
	Create(ctx context.Context, item T) (T, error)
	Get(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, item T) error
	Delete(ctx context.Context, id string) error
	QueryByIndex(ctx context.Context, indexName, keyName, keyValue string) ([]T, error)
	Scan(ctx context.Context) ([]T, error)
}