param (
    [string]$Dir = "u"  # Default is to run "up" migrations
)

# Define the path to the .env file and the migrations folder
$envFilePath = ".\.env"
$migrationsFolder = "scripts\migrate"

# Load .env file into a hashtable
$envVars = @{}
Get-Content $envFilePath | ForEach-Object {
    if ($_ -match "^(.*)=(.*)$") {
        $envVars[$matches[1]] = $matches[2]
    }
}

# Extract each variable and use it to build the database connection string
$dbUser = $envVars["DB_USER"]
$dbPassword = $envVars["DB_PASSWORD"]
$dbName = $envVars["DB_NAME"]
$dbHost = $envVars["DB_HOST"]
$dbPort = $envVars["DB_PORT"]

# Build the connection string with proper formatting for Goose
$connectionString = "postgres://{0}:{1}@{2}:{3}/{4}?sslmode=disable" -f $dbUser, $dbPassword, $dbHost, $dbPort, $dbName

# Define Goose executable path, assuming it’s installed and in PATH
$goosePath = "goose"

# Run Goose command based on the $Dir parameter
if ($Dir -eq "u") {
    $gooseCommand = "& $goosePath -dir $migrationsFolder postgres $connectionString up"
} elseif ($Dir -eq "d") {
    $gooseCommand = "& $goosePath -dir $migrationsFolder postgres $connectionString down"
} else {
    Write-Output "Invalid Direction parameter. Use 'up' or 'down'."
    exit 1
}

# Execute the Goose command
Invoke-Expression $gooseCommand
