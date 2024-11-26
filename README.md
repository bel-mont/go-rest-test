# Localhost
## Requirements
1. Golang
2. For windows, install "Makefile" through something like Chocolatey.
3. Python

## Running 
1. Start localstack with `docker-compose`
2. Start your golang program `go run ./cmd/server/main.go`

# TODO
1. Switch to production Tailwind
2. 


# Troubleshooting
If the imports are not working in Goland, make sure to enable Go Modules in the IDE settings.

## Database migrations and seeding
Run the db-migrate.ps1 powershell to update your DB.
Translate to shell if using Linux/macOs.

[Goose](https://github.com/pressly/goose) is used to run the migrations, install them in your machine.

## References
Check out https://github.com/go-gitea/gitea/tree/main to improve project organization.