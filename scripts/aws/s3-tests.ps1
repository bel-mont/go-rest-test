# List all objects in the bucket
aws --endpoint-url=http://localhost:4566 s3 ls s3://my-test-bucket

# List with more details
aws --endpoint-url=http://localhost:4566 s3 ls s3://my-test-bucket --recursive

# List with human-readable sizes
aws --endpoint-url=http://localhost:4566 s3 ls s3://my-test-bucket --recursive --human-readable

# Get details about a specific file
aws --endpoint-url=http://localhost:4566 s3api head-object --bucket my-test-bucket --key ab.txt

# Download file (to a different name to avoid overwriting)
aws --endpoint-url=http://localhost:4566 s3 cp s3://my-test-bucket/ab.txt downloaded-ab.txt

# Delete a specific file
aws --endpoint-url=http://localhost:4566 s3 rm s3://my-test-bucket/ab.txt

# Delete all files in bucket
aws --endpoint-url=http://localhost:4566 s3 rm s3://my-test-bucket --recursive