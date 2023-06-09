# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /transactions-email-processor

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /transactions-email-processor /transactions-email-processor

# Copy project assets
COPY email/dist/email_template.html email/

USER nonroot:nonroot

ENV SENDER_MAIL="sender@gmail.com"
ENV PASSWORD="secret-password"
ENV RECIPIENT_MAIL="recipient@example.com"

ENTRYPOINT ["/transactions-email-processor"]
