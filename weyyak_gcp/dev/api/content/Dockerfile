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
ENV SERVICE_PORT=3001 \
    DB_SERVER=10.87.0.5 \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123DEV \
    DB_DATABASE=wyk_content \
    FRONTEND_DB_DATABASE=wk_frontend \
    DEFAULT_PAGE_SIZE=20 \
    CONTENT_COMMON_API_URL=https://qa-ms-api-contentcommon.wyk.z5.com/ \
    CMS=https://apiqa.wyk.z5.com/v1/ \
    UM=https://apiqa.wyk.z5.com/v1/ \
    IMAGES=https://weyyak-content-dev.engro.in/ \
    GEO_LOCATION=https://geo.weyyak1.z5.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://apiqa.wyk.z5.com/v1/ \
    IMAGERY_URL=https://weyyak-content-dev.engro.in/ \
    AWS_REGION=ap-south-1 \
    ACCESS_SECRET=AKIAYOGUWMUMEEQD6CPW \
    S3_BUCKET=z5content-uat \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMGAQLMW3U \
    S3_SECRET=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 \
    REDIS_CACHE_URL=https://backofficeapi-events-dev.engro.in/v1/cache \
    REDIS_CONTENT_KEY=GOAPIQA \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    JEAGER_SERVICE_NAME=Content \
    BUCKET_NAME=weyyak-dev \
    TYPE=service_account \
    PROJECT_ID=engro-project-392708 \
    PRIVATE_KEY_ID=ee585d8aec442766ecd64561ab33c23f2ed79e71 \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCvoYm3Mij1OE6B\nQFPONNGIU6Lt1Vo9eUqxnmot2+g5S3rH4FzgtyUYuDWWczQmw/Eo201NIR99hxhu\nRJfhw9wwynC7ioGTG17BAhbjQV97hJK1Ucj8B45Ggs6VOJ2e6sWAE+JDG6CdEHUV\nXxWUFRwr0FUyP1aMrE4GGPXAXnWHEJ3FLCfpY20DJ2SlvprNlARZqhQ3XGxjbqpf\ns4Lb5f6Kce7+Q6VIG4SLQ0X/bNSqwppDreqxh3sqFfIgBGVsxgUIaraDUMB7C188\n+h90xwXoFxKZ5b4rzk7VDdQSU644MZrpt5a+/WFB/G30lzzO8MsQrtCuMAidN4Vi\nlVejk7SZAgMBAAECggEAALC+P3i8YbtlDpvDwJKjVLl8cVomesmiZiiY4wGDbA88\nhjLe9pittgdKWaMuDC12whn7l68mxins6gkNvkrUb6IqF4ijugHwaUqSb7YC87O5\nlmoiA8CdWfBbcsl4UCa83EQnvmsKQnfNOAZpya1cER2GajZGxzKJs3svYyBT7rIb\nujDyDWTSe5xMKWuCg9o6Z+wN0izY/MYwWMnx4h5hv4yKJO5HMCsyibz4+QTvQMnU\nOjajAPcTEH9PSeA78rLNEdcdl6AK60lAGMU4vfgo3k4Eah1W9PW8Rt1XmbFIPQUZ\nlxT9O8awaUcKNaBqvkYd/fl4yk8GyILBEK9JYkbSQQKBgQDn0NvmnyvGZS2pJtoA\nZyJ2A4TLfbYuy8PXDe7uVYxzEk6kMAKXazevxa9hC+CzZa00Vg1/FLxUG7mLN639\n72s2+hv4gnrjO2CqXrziziPeHaGZbPoAn0LWUe1pFc6ykdru1zRQ1sjH+whvOTU6\n++4Z0CL37uNdSrAUHBuE/rLk2QKBgQDB9CRKZIz7hTotDq8CISPqWikEbfYSPu3n\nVvF/eEyCuBG158tujWJSShnVTw8mo0dQxST8lyBVETf464tefTs8tk5P9gTlfhne\nNb2I+bVTzBAlstEWjzXm9gc83pO1ayBIRhseZH5xUHEWyaptc0IGPvelykm4k3Rc\n11Rd8k51wQKBgFvLEXiF6NCbufP0JXjpXfT7ObdWZTCskTT0GHuZto7OXgOX2z7Y\nPunETskP6f8/sNE7jgaXUGfBvrwI5UXba0oOycqi6ERYe9OaxhLNa9iFRTs6sthS\nOCYjtO3aq8l8BgRzhQzJsCG4HsPkFFbx8XuZNm+U0cVcUNL1ulP0gyQpAoGAEqjL\nSQDtc1E47JbKDXlkvCoaYhu2HcFCIAyERVAshDLzOT3om2UsvV/1pkjp2zzAIqwk\nx8w+RxDWArOef+OqbEfZtuWjSPTdxjwDVQc5A6vRLd9LMFFszSCER73eFgkNXQap\nmz3t0Db9+kRCii8ZfgJWzvNsLj65ZiktK7ChDcECgYAkljodDDYejMYKKwD1FDMH\nPuWFLoT23V1J9c9IcxNQscT6pNSa+Wa8jFAkhIgmrWu6bTmRt8hHlaBJF5vH5Ins\nRgEkKGUF/jjzZJ7CKjw8S855VBJMctjNVCvmMJJvpgfJkC5iBmHvguL2J03SR8Wu\nbEpoMQHNEiWXfG8PHDPzUw==\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=xmsbucket@engro-project-392708.iam.gserviceaccount.com \
    CLIENT_ID=115836177213859121422 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/xmsbucket%40engro-project-392708.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com
    
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
