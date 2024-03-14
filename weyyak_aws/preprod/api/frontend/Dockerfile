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
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /dist/main /
ENV SERVICE_PORT=3003 \
    DB_SERVER=weyyak-preprod.cyj14gqdm97b.me-south-1.rds.amazonaws.com \
    DB_SERVER_READER=weyyak-preprod.cyj14gqdm97b.me-south-1.rds.amazonaws.com \
    DB_PORT=5432 \
    DB_USER=weyyak_aurora \
    DB_PASSWORD=M5Ltay9sDY93khvmcpNE \
    DB_DATABASE=wk_frontend \
    CONTENT_DB_DATABASE=wyk_content \
    FRONTEND_CONFIG_DB_DATABASE=wyk_frontend_config \
    USER_DB_DATABASE=wk_user_management \
    DEFAULT_PAGE_SIZE=20 \
    CONTENT_COMMON_API_URL= \
    CMS=https://apiqa.wyk.z5.com/v1/ \
    UM=https://apiqa.wyk.z5.com/v1/ \
    IMAGES=https://contents-uat.weyyak.z5.com/ \
    GEO_LOCATION=https://geo.weyyak1.z5.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa/ \
    USER_ACTIVITY_URL=https://msapiqa-events.z5.com/event/activity \
    CONTENT_TYPE_URL=https://msapifo-preprod-me.weyyak.z5.com/v1/en/contents/contentType?contentType= \
    CONTENT_TYPE_URL_PAGINATION=&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0 \
    USER_LOG_URL=https://msapiqa-events.z5.com/event/log \
    S3_BUCKET=z5content-uat \
    S3_URLFORCONFIG=https://z5content-uat.s3.ap-south-1.amazonaws.com/configpreprod.json \
    CONFIG_FILE_NAME=configpreprod.json \
    CONFIG_KEY=configpreprod \
    REDIS_CACHE_URL=https://msapipreprod-me-events.z5.com/cache \
    REDIS_CONTENT_KEY=GOAPIPREPROD \
    DOTNET_URL=https://uat-api.weyyak.z5.com/v1/ar/oauth2/tokendata?access_token= 
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
