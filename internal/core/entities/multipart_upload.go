package entities

import "time"

type MultipartUpload struct {
	ID             string         `json:"id" dynamodbav:"id"`
	UserID         string         `json:"user_id" dynamodbav:"user_id"`
	Status         string         `json:"status" dynamodbav:"status"` // 'in_progress', 'completed', 'aborted'
	FileName       string         `json:"file_name" dynamodbav:"file_name"`
	FileSize       int64          `json:"file_size" dynamodbav:"file_size"`
	S3Bucket       string         `json:"s3_bucket" dynamodbav:"s3_bucket"`
	S3Key          string         `json:"s3_key" dynamodbav:"s3_key"`
	TotalParts     int            `json:"total_parts" dynamodbav:"total_parts"`
	CompletedParts map[int]string `json:"completed_parts" dynamodbav:"completed_parts"` // part number -> ETag
	StartedAt      time.Time      `json:"started_at" dynamodbav:"started_at"`
	UpdatedAt      time.Time      `json:"updated_at" dynamodbav:"updated_at"`
	ExpiresAt      time.Time      `json:"expires_at" dynamodbav:"expires_at"`
}

func (u MultipartUpload) GetID() string {
	return u.ID
}

func (u MultipartUpload) SetID(id string) {
	u.ID = id
}
