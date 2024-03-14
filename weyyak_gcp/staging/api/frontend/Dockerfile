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
	DB_SERVER=stg-psqldb.weyyak.com \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=WMFPfp6UGIQYIAY \
	DB_DATABASE=wk_frontend \
	CONTENT_DB_DATABASE=wyk_content \
	FRONTEND_CONFIG_DB_DATABASE=wyk_frontend_config \
	USER_DB_DATABASE=wk_user_management \
	DEFAULT_PAGE_SIZE=20 \
	CONTENT_COMMON_API_URL= \
	CMS=https://msapifo-stg-me.weyyak.z5.com/v1/ \
    UM=https://msapifo-stg-me.weyyak.z5.com/v1/ \
    IMAGES=https://contents-stg.weyyak.z5.com/ \
    GEO_LOCATION=https://geo-stg.weyyak.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://msapifo-stg-me.weyyak.z5.com/ \
    USER_ACTIVITY_URL=https://msapistage-me-events.z5.com/event/activity \
	USER_LOG_URL=https://msapistage-me-events.z5.com/event/log \
    CONTENT_TYPE_URL=https://msapifo-stg-me.weyyak.z5.com/v1/en/contents/contentType?contentType= \
    CONTENT_TYPE_URL_PAGINATION=&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0 \
    S3_BUCKET=z5content \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMLWREWLMX \
    S3_SECRET=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 \
    S3_URLFORCONFIG=https://z5content.s3.ap-south-1.amazonaws.com/configprod.json \
    REDIS_CACHE_URL=https://msapistage-me-events.z5.com/v1/cache \
    REDIS_CONTENT_KEY=GOAPISTAGE \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    GCP_URLFORCONFIG=https://storage.googleapis.com/weyyak-stage/configstagemigration.json \
    BUCKET_NAME=weyyak-stage \
    TYPE=service_account \
    PROJECT_ID=wyk-stage \
    PRIVATE_KEY_ID=fd9feef48aacb563b85d9a1e4b957a325c1e9187 \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCuZx2VsTN66BXo\nPTQZ6iqFz5DkLWQpSaJ+gnO244Qy7uDV973FF77HIGK71DD0oFCnLClDEEBe0qcY\nprXHUNH7fYawBs0mQ6LrlI9vJkcUEHuKEA0lFIXuL93g/wmymQYsedVF5cKOVhwS\n8cWpIkhY+SzwTIol4vZWT9VBmplliZgsM3XNDtHc7FlqEgV4YkaAuLNEUbPPe59Z\nH+r6CwU/1tCodsWY29TspZgM6A+2BQtqnMBaG0ZmEUFFFVoj5Zny4XbwCM++RD+j\nBd2b5ahz3HyXMXGT/cgVYIJSu0ACUerTzoPqiTsBuVKT6v4Bhgum74uInkn/0CsM\nFcUkieJjAgMBAAECggEAQZRltAm3fXJy9mq5Qjti9QEpjxuga8QUFbUxFOb1qX9Y\nzSyEz+2PX+pGJjHynLIB0fLRqzWbtuxISAU8GnBCHWrEsEFGeCp7w1EQJHAEzz8p\nbS9++d//QqL+MGJGf8F56FSBxXBRmds4YElIlcJd6RL9eniVhAyVq7wGSU58ewjE\nzVcmNoSldsTDxFOarKasLPxQP5vWr0PTQK05TDp0KM9iLjnvZ+rz3NVPGRAwYgZt\nFme2OMJV9KlyIhPJAlEDBIhKIbbeHxQSQw9jq99apB9xBv/S3B0IqdESwKW3l0fk\n69QRQwytlq97wb3RroVYvDjbvmTDWfA5KxUv/5pEeQKBgQDa4w0dLnaWEvOSD63D\nf67rV5zEpCsAXQEc8dgmZiQWkzMeZzfeVgE6C5Fq+jmRWGdtD7UiLZ42gL5acvHV\nxQqziqOWDJ2sj8BCEkJADW//J0J3a46z0RDB8PBNmKXe/JY6Gfec0bCUbit78hKc\nLkYrtIbk2KmIIKg0khlAYJa+KwKBgQDL+TKxIT9o9RssaG1mNTU4T9tiDd1Tlu+E\nFheVMVjkLftOA29je+be1zKqfCzi+Pjki1YZfIdMIFBafhQsQzOfIjzycGatbCAo\nL8dnm7GfGn+x8k9yCauQw15+VGjy9crwIf8H+WrotBbFLDA9ln0wJyn5g7aDpl5b\n8dNO278IqQKBgC9iFVWq2Dr2kUmeIUR/4HLLvYEH9NMYTwRgCiWN2OQWK2u6dmqn\n/Y3+MRWSY2HmXEVMsC5c57HuvsY/ucBH7mAdc/oocv2HBNyIrVTGU+dNLWNColqk\nuFbzfGGf2NBjp3WeZ27IKcyH6Nvk9Ehg8AnIA6eIi3KPdzoiIltATYGPAoGAeqRK\nW4B/LO97ODdN2A+UAVjastldQ35UY1oZMeZq6T/Qg3bDSdDc3w4Lc8CYCJzkoh45\nVbevbhFb3x1Bqd96MWuXs7nrssh5iVA5GE3f/ackMm1cYVf5JWCi91hvhe5N6Ba4\nziWkp2ho7IMGOnG6sfJKZbFRU1T8FVKQ2DB+a4kCgYEApv8Gslu4qaBeQiy2P+f5\nc569kMiw3HRSfTNhDON5Q+wQLUjzBykUC1z9XLBvYUBrrUjampvUYZFuOir8TBg8\nixXdojAjLKLz3OYFEEi7VAMnqkG/zYtKtZis+OiJH3U45aimxTitODF2GU2OmMU9\nlrG5eNebeOOn7E8YXCzmeWA=\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=weyyak-stage@wyk-stage.iam.gserviceaccount.com \
    CLIENT_ID=102280627070321861141 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/weyyak-stage%40wyk-stage.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
