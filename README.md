# Back-end-developer
Full Back-end developer with GO
```markdown
# Hospital Backend System

This is a backend system for hospital patient management developed with Go and Gin framework.

## Features
- Patient management
- Staff authentication
- Hospital information

## Installation

### Prerequisites
- Go 1.16 or higher
- PostgreSQL 13 or higher

### Setup
1. Clone the repository
   ```
   git clone https://github.com/Natthaphatpiw/Backend-with-GO-GIN.git
   ```

2. Install dependencies
   ```
   go mod download
   ```

3. Configure environment variables
   ```
   export DB_HOST=localhost
   export DB_USER=myuser
   export DB_PASSWORD=mypassword
   export DB_NAME=mydatabase
   export DB_PORT=5432
   ```

4. Run the application
   ```
   go run main.go
   ```

## API Documentation
API documentation is available at `/swagger/index.html` after starting the server.

## Running Tests
```
go test ./... -v
```
```
