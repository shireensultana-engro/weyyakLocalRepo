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
COPY --from=builder /build/templates /
ENV SERVICE_PORT=3007 \
    DB_SERVER=10.87.0.5 \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123DEV \
    DB_DATABASE=wyk_features \
    DEFAULT_PAGE_SIZE=50 \
    API_VERSION=/v3 \
    CDB_DATABASE=wyk_content\
    SES_REGION=ap-south-1 \
    SES_ID=AKIAYOGUWMUMK2O4DT6B \
    SES_SECRET=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 \
    BASE_URL=https://wyk2qa.weyyak.com

EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
