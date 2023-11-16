# Use the official Golang image as the base image
FROM golang:1.18-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Install testing essentials
RUN apk update \
    && apk add --virtual build-dependencies \
    build-base \
    gcc \
    wget \
    git

# Download and cache the Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the application binary
RUN go build -o app main.go

# Copy the migrations directory into the container
COPY database/migrations /app/migrations

# Install migration package
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@master

# Expose port 8080 for the app to listen on
EXPOSE 8080

# The nc -z db 5432 command checks whether the db service is available
# by attempting to connect to db port. The sleep 1 command adds a delay
# between attempts to connect to the database. Once the database is available,
# the migration command is executed followed by the seeder: database/seeders/init_seeder.go
ENTRYPOINT ["sh", "-c", "while ! nc -z db ${DATABASE_PORT}; do sleep 1; done && migrate -path migrations -database postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@db:${DATABASE_PORT}/${DATABASE_DB}?sslmode=disable up && go run database/seeders/init_seeder.go && ./app"]



CMD ["./main"]

