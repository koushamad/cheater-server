# Set base image
FROM golang:1.16-alpine3.14

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN go mod download
RUN go mod vendor
RUN go build -o cheater

# Expose port 8080
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./cheater"]
