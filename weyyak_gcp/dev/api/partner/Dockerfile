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
ENV SERVICE_PORT=3006 \
    DB_SERVER=10.87.0.5 \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123DEV \
    DB_DATABASE=wyk_content \
    DEFAULT_PAGE_SIZE=50 \
    API_VERSION=/v3 \
    BASE_URL=https://apiqafo.engro.in \
    IMAGE_URL=https://weyyak-content-dev.engro.in/ \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    LOGIN_API=https://apiqafo.engro.in/v1/oauth2/token \
    CONTENT_URL=https://apistg.weyyak.z5.com/v1/get_info/ \
    IMAGE_URL_GO=https://weyyak-content-dev.engro.in/v1/ \
    DETAILS_BACKGROUND=/details-background\
    POSTER_IMAGE=/poster-image\
    MOBILE_DETAILS_BACKGROUND=/mobile-details-background\
    OVERLAY_POSTER_IMAGE=/overlay-poster-image\
    DUBBLING_SCRIPT=/dubbling-script\
    WATCH_NOW=cce8db39-3d54-442c-8f21-f6c1aa11d396\
    SUBTITLING_SCRIPT=/subtitling-script \
    REDIS_CACHE_URL=https://backofficeapi-events-dev.engro.in/v1/cache \
    REDIS_CONTENT_KEY=GOAPIQA \
    IMAGERY_URL=https://weyyak-content-dev.engro.in/

EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
