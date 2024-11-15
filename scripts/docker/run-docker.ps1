# Set up variables for image and container names
$imageName = "server-postgres"
$containerName = "server-postgres-container"

# Check if the container exists (running or stopped)
$containerExists = docker ps -a --filter "name=$containerName" --format "{{.Names}}" | Select-String -Pattern $containerName

if ($containerExists) {
    # If the container exists, check if it's running
    $isRunning = docker ps --filter "name=$containerName" --format "{{.Names}}" | Select-String -Pattern $containerName
    if ($isRunning) {
        Write-Host "Docker container '$containerName' is already running."
    } else {
        Write-Host "Starting existing Docker container '$containerName'..."
        docker start $containerName
    }
} else {
    # If the container does not exist, create and run it
    Write-Host "Creating and starting Docker container '$containerName'..."
    docker run -d -p 5432:5432 --name $containerName $imageName
}

# Output the status
Write-Host "Docker container '$containerName' is now running."
