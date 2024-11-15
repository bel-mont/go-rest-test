# Navigate to the directory containing the Dockerfile
$dockerfilePath = "./docker/Dockerfile.postgres"
$contextPath = "./docker"
$imageName = "server-postgres"
$containerName = "server-postgres-container"

# Build the Docker image
docker build -t $imageName -f $dockerfilePath $contextPath

# Check if a container with the same name is already running and stop it
if (docker ps -a --format "{{.Names}}" | Select-String -Pattern $containerName) {
    Write-Host "Stopping and removing existing container..."
    docker stop $containerName
    docker rm $containerName
}

# Run the Docker container
docker run -d -p 5432:5432 --name $containerName $imageName

# Output the status
Write-Host "Docker container '$containerName' is now running."
