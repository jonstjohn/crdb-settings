FROM golang:1.22 AS build-stage

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to /app
COPY go.mod go.sum ./

# Install go modules
RUN go mod download

# Copy source code
#COPY main.go cmd/ pkg/ ./
COPY . /app

# Compile command
RUN CGO_ENABLED=0 GOOS=linux go build -o /crdb-settings

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /crdb-settings /crdb-settings

ENTRYPOINT ["/crdb-settings"]
CMD ["settings", "update"]