# Wait for LocalStack to be ready
Write-Host "Waiting for LocalStack to be ready..."
Start-Sleep -Seconds 5

# Create S3 bucket
Write-Host "Creating S3 bucket..."
aws --endpoint-url=http://localhost:4566 s3 mb s3://fg-analyzer

# List buckets
Write-Host "Listing buckets..."
aws --endpoint-url=http://localhost:4566 s3 ls

# Check CORS rules
aws --endpoint-url=http://localhost:4566 s3api get-bucket-cors --bucket fg-analyzer-replay-uploads

# Set CORS
aws --endpoint-url=http://localhost:4566 s3api put-bucket-cors --bucket fg-analyzer-replay-uploads --cors-configuration '{"CORSRules": [{"AllowedHeaders": ["*"], "AllowedMethods": ["GET", "PUT", "POST"], "AllowedOrigins": ["*"]}]}'

# Set CORS (Powershell double quotes)
aws --endpoint-url=http://localhost:4566 s3api put-bucket-cors --bucket fg-analyzer-replay-uploads --cors-configuration "{\"CORSRules\": [{\"AllowedHeaders\": [\"*\"], \"AllowedMethods\": [\"GET\", \"PUT\", \"POST\"], \"AllowedOrigins\": [\"*\"]}]}"

# Set CORS with file
aws --endpoint-url=http://localhost:4566 s3api put-bucket-cors --bucket fg-analyzer-replay-uploads --cors-configuration file://localstack/s3-cors.json