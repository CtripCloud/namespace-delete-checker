# Build the manager binary
FROM golang:1.12.7 as builder

# Copy in the go src
WORKDIR /namespace-delete-checker
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o delete-checker main.go

# Copy the controller-manager into a thin image
FROM ubuntu:18.04
WORKDIR /app
COPY --from=builder /namespace-delete-checker/delete-checker /app


