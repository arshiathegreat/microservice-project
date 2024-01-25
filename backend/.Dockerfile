# Use the official Go image as the base image
FROM golang:1.20

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules files to the working directory
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application files to the working directory
COPY . ./

# Build the Go app
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Start the Go app when the container starts
CMD ["./main"]
