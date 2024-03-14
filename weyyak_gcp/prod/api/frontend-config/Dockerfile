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
    DB_SERVER=34.18.48.27 \
    DB_SERVER_READER=34.18.48.27 \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=Yy017EUY5aSz \
    DB_DATABASE=wyk_frontend_config \
    CDB_DATABASE=wyk_content \
    FDB_DATABASE=wk_frontend \
    USER_DB_DATABASE=wk_user_management \
    DEFAULT_PAGE_SIZE=50 \
    CMS=https://msapifo-prod-me.weyyak.z5.com/v1/ \
    UM=https://msapifo-prod-me.weyyak.z5.com/v1/ \
    IMAGES=https://content.weyyak.com/ \
    GEO_LOCATION=https://geo.weyyak1.z5.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://msapifo-prod-me.weyyak.z5.com/ \
    S3_BUCKET=z5content-uat \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMGAQLMW3U \
    S3_SECRET=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 \
    REDIS_CONTENT_KEY=GOAPIPROD \
    REDIS_CACHE_URL=https://msapiprod-me-events.z5.com/v1/cache \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    RESIZE_IMAGE_URL=https://msapiuat-image.z5.com/crop? \
    JEAGER_SERVICE_NAME=Frontend-Config \
    BUCKET_NAME=weyyak-production \
    TYPE=service_account \
    PROJECT_ID=wyk-production \
    PRIVATE_KEY_ID=2d0304f0867a7a4f06d6c4369140923c82378ab2 \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDdw+3R+RyhqQRo\nSHDkAoP2J6Jg0qfy1OG1SQmB66NQMc8Jhp17G66iUatF+AqLNhcrXtXYHEMBmVPb\nsCpkrGXb6wJzUCkvHrj4IhF8xmrQdQOBbKEBF0pg9GjAT4UpQjUMxryXT1I4Dfau\ns/D9+rf/sTIgQID2n+uOQ1knC29LK63Er6+1JxgjbXbL/778ilJ8npl/nsPBITWo\nRQUgsTvOdcgsnLdgMdN0pb0XG5hTWiIQF/cYJDJ22YriYUwPXgTdiHW4qn65TKcT\nxg2YNzLJnw1Uqhb0RZVsraEkkW6RiVgWYDXwcWrCRede80HifvpSEWiCC3KnXiIt\nNFHJE447AgMBAAECggEABwBJgMiBi+T/G5+12Kzvp5TGvpHH9ZWc7pE4uJ5M0JpR\n8/YJALr1/2/enV3gT1bM0nSzAZia0PEbQaNFI1qB+LhpomRUeIVax5KjxLGq65vW\nGX7pclRe58Kvj+qyxIOvkxCvIYPCj7x5HjjWEd6ZcnwQng4LRD32PM6JgP8Oa2wN\nh+XLlge2VrdqidIW0R+koP7Avq59A5fVM/VPakSI21FV6PT3kpwjXBzGUKXt9POX\nWu8vLCj0cbvOcLgvBFKvysJyH/0eiQBg1XJJLKadSTGAdFZGl706e8wkOkJZK2Wj\n3/EJHRV/B8FqyN5WS6UusTCg5QsjxKoBHhvDtC2SsQKBgQD5dzBmtZQTDcunrToB\nIn/oeZWjlORsns3YggUIijJ1aNK7MKvA+jfrplJ8biRN3I+DiSPRnc6I/RJDemC5\nJP9ky96iIMfeuJL/eNVy91YBsam1ZzZFm1PFKRw0OEcga2jxa/2YRPFE2+1Q2cyf\n4Mw3WQVplq3Sv9pmOwKLwo3oOQKBgQDjkv5rydaIEBgwOE6iHA5ndVDBxbQRVu4t\nSpA9+F90oQEO7Xrbei4fs5UmAWdRCVhgmLRJsEQEdOwBbtqe+n3Kk8yb+/QhSX1q\n/6RutjAEM/KqOpIEUqgKyR7Y5A2AC0zpIU1TyZ2U+rTTHPRBAOWSy5kJ+NkhuQUP\nc2JsU9HiEwKBgCEQKPwT6NI1q95HWT65Qdaf9rM9kqDK02F0qhIdrt5czEE/DCSB\nhVPYMWqIdotTRjoavQKVNcB2OitzVspzGt5THujCC3t7XxA5BaE9IssKrwF58nl7\nQrkI39IT+2lSkxAcTfoWeRu1QljK5RHzi11ykQMTk2oxP1L5UzcOzBwRAoGAGM+I\nz2WU7waaLH+nCwN2Co9+u3F7fTx2ARgU+7ydY5C+FcuMTmtWpfwlMZyLkAktynI7\njaEa+UVqCYn1acmzdyd/8i2Y4xwpAUZXvf484+hp92clTjVYvrxIkarjUedpfi00\nSgM8G+btWerZMlEPtl5eE/k+au/J/nI888R7qGMCgYEAj+AkZEz/jcsOI5WKUW1A\nwOKKU5LFFT6eQlVNLq8NkY9hAyiHQiGRKT8hSi78LaEWqb7yqWw4V27Xcm5lYuNN\nKqC1yK6Aj4rFV8r1chNcagryJrlSZwKOaG1mS5sIPcqsCV3ZfPs4Rl17otqoro+D\n8dx7sV4TUup4uPiTkm4s96A=\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=msapifo-prod-service-account@wyk-production.iam.gserviceaccount.com \
    CLIENT_ID=117641798705593610361 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/msapifo-prod-service-account%40wyk-production.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com

EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
