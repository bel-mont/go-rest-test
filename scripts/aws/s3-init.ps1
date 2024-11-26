# Wait for LocalStack to be ready
Write-Host "Waiting for LocalStack to be ready..."
Start-Sleep -Seconds 5

# Create S3 bucket
Write-Host "Creating S3 bucket..."
aws --endpoint-url=http://localhost:4566 s3 mb s3://fg-analyzer

# List buckets
Write-Host "Listing buckets..."
aws --endpoint-url=http://localhost:4566 s3 ls