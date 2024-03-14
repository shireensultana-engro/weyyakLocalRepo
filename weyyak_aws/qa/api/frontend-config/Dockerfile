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
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /dist/main /
ENV SERVICE_PORT=3004 \
    DB_SERVER=msapiqa-rds.z5.com \
    DB_PORT=5432 \
    DB_USER=weyyak_aurora \
    DB_PASSWORD=M5Ltay9sDY93khvmcpNE \
    DB_DATABASE=wyk_frontend_config \
    CDB_DATABASE=wyk_content \
    FDB_DATABASE=wk_frontend \
    USER_DB_DATABASE=wk_user_management \
    DEFAULT_PAGE_SIZE=50 \
    CMS=https://apiqa.wyk.z5.com/v1/ \
    UM=https://apiqa.wyk.z5.com/v1/ \
    IMAGES=https://contents-uat.weyyak.z5.com/ \
    GEO_LOCATION=https://geo.weyyak1.z5.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa/ \
    S3_BUCKET=z5content-uat \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMGAQLMW3U \
    S3_SECRET=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 \
    REDIS_CONTENT_KEY=GOAPIQA \
    REDIS_CACHE_URL=https://msapiqa-events.z5.com/cache \
    DOTNET_URL=https://uat-api.weyyak.z5.com/v1/ar/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    JEAGER_SERVICE_NAME=Frontend-Config
# BASE_URL=http://localhost:3004
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
