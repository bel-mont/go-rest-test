package entities

import "time"

type Replay struct {
	ID         string    `json:"id" dynamodbav:"id"`
	UserID     string    `json:"user_id" dynamodbav:"user_id"`
	UploadedAt time.Time `json:"uploaded_at" dynamodbav:"uploaded_at"`
	S3Bucket   string    `json:"s3_bucket" dynamodbav:"s3_bucket"`
	S3Path     string    `json:"s3_path" dynamodbav:"s3_path"`
	S3FileName string    `json:"s3_file_name" dynamodbav:"s3_file_name"`
	S3FileSize int64     `json:"s3_file_size" dynamodbav:"s3_file_size"`
}

func (u Replay) GetID() string {
	return u.ID
}

func (u Replay) SetID(id string) {
	u.ID = id
}
