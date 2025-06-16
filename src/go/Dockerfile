# Start from a minimal base image containing the Go runtime
FROM golang:1.24-bookworm AS build

ARG GITHUB_REPOSITORY

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init --dir ./cmd/kube-ondemand-sidecar-injector/,./internal --output ./internal/docs/

# Run tests
RUN go test ./... -coverpkg=${GITHUB_REPOSITORY}/internal/controllers/injector,${GITHUB_REPOSITORY}/internal/kube,${GITHUB_REPOSITORY}/internal/logging -coverprofile=coverage.out

RUN go tool cover -html=coverage.out -o coverage.html

# Build the Go application
RUN go build -o main cmd/kube-ondemand-sidecar-injector/main.go

# Start from a new, lightweight base image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Set the current working directory inside the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /app/main .
COPY --from=build /app/coverage.html .

# Expose the port that the Go application listens on
EXPOSE 8080

RUN addgroup ondemandsidecarinjector \
       && adduser ondemandsidecarinjector -G ondemandsidecarinjector -D 

USER ondemandsidecarinjector

# Command to run the Go application
CMD ["/app/main"]
