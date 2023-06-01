# Use a base image with Go preinstalled
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the source code and go.mod/go.sum files into the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the Go application
RUN go build -o my-go-app

# # Set the command to run when the container starts
# CMD ["./my-go-app"]
