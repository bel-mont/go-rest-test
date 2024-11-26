# List all objects in the bucket
aws --endpoint-url=http://localhost:4566 s3 ls s3://fg-analyzer

# List with more details
aws --endpoint-url=http://localhost:4566 s3 ls s3://fg-analyzer --recursive

# List with human-readable sizes
aws --endpoint-url=http://localhost:4566 s3 ls s3://fg-analyzer --recursive --human-readable

# Get details about a specific file
aws --endpoint-url=http://localhost:4566 s3api head-object --bucket fg-analyzer --key ab.txt

# Download file (to a different name to avoid overwriting)
aws --endpoint-url=http://localhost:4566 s3 cp s3://fg-analyzer/ab.txt downloaded-ab.txt

# Delete a specific file
aws --endpoint-url=http://localhost:4566 s3 rm s3://fg-analyzer/ab.txt

# Delete all files in bucket
aws --endpoint-url=http://localhost:4566 s3 rm s3://fg-analyzer --recursive