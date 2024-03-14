FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine

COPY --from=builder /dist/main /
ENV SERVICE_PORT=3002 \
    DB_SERVER=prod2-db.weyyak.com \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=Yy017EUY5aSz \
    DB_DATABASE=wk_read_only \
    DEFAULT_PAGE_SIZE=20 \
    API_VERSION=/v3
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
