# Use the official Golang image as the base image
FROM golang:1.17-alpine AS build

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Install the dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the application
RUN go build -o app

# Use a smaller base image for the final image
FROM alpine:3.14

# Set the working directory to /app
WORKDIR /app

# Copy the application binary from the build stage to the final image
COPY --from=build /app/app .

# Expose port 8080 for the application
EXPOSE 8080

# Start the application
CMD ["./app"]
