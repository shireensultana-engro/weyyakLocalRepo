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
    DB_SERVER=wyk2db-uat.engro.in \
    DB_SERVER_READER=wyk2db-uat.engro.in \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123UAT \
    DB_DATABASE=wyk_content \
    FRONTEND_DB_DATABASE=wk_frontend \
    DEFAULT_PAGE_SIZE=20 \
    CONTENT_COMMON_API_URL=https://qa-ms-api-contentcommon.wyk.z5.com/ \
    CMS=https:https://msapifo-uat.weyyak.z5.com/v1/ \
    UM=https://msapifo-uat.weyyak.z5.com/v1/ \
    IMAGES=https://content-uat.weyyak.com/ \
    GEO_LOCATION=https://weyyak-geoapi-uat.engro.in \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=//storage.googleapis.com/z5xml-qa/mobile_apps_ads_ios.xml \
    BASE_URL=//msapifo-uat.weyyak.z5.com/v1/ \
    IMAGERY_URL=https://content-uat.weyyak.com/\
    AWS_REGION=ap-south-1 \
    ACCESS_SECRET=AKIAYOGUWMUMEEQD6CPW \
    S3_BUCKET=z5content-uat \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMGAQLMW3U \
    S3_SECRET=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 \
    RESIZE_IMAGE_URL=https://msapiuat-image.z5.com/crop? \
    REDIS_CACHE_URL=//msapiuat-events.z5.com/v1/cache \
    REDIS_CONTENT_KEY=GOAPIUAT \
    REDIS_RDS=34.93.213.137:6379 \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    JEAGER_SERVICE_NAME=Content \
    BUCKET_NAME=wyk-content-uat \
    TYPE=service_account \
    PROJECT_ID=testing-400312 \
    PRIVATE_KEY_ID=4ba9f2e2438c5fa83c162336657431ad1f92f6d2 \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCrrOzb27NS3vj2\n/mspxBm/22giPXnCND0DYKLRcQuzsYv0elnbqmcGSJfCqs9C59PGJzRS6RV+ie/N\nhgflwUx3I3zMtd1fgVsjZFYYwmIDsaShkaK0B3eYWGXqYw8swG67qctk3GubNYNl\n04p6qIsUZgTpmH8Jjvn7abwntDvfxYbK9QVC2rjXzzA0pABfrgo3CVZJiXO8t4Ly\na8P5jkNY8Uze6Us+L0XxYL8T08vD09Cqde2kEEMBM4R6okdtgj3Vp8G287Dj0OUe\nMBGoj6YgsH32ZOXxxeQp003tMETaJbwQx6HTKr4N0CDwQnQkMbVUN2N8wOd4e0+2\nG+U36cGVAgMBAAECggEAFn+JDZ8TNyRleD5gs46G2VqFoRxxXSlqEuE9NTlyu8/k\nHtv8nrRhirSaFDbnsUWfE/QwqpTv7i9hhTZayUS1zVSR7GSrvZ0UNo/Vq1T+HKx7\n03i52+IGov54DL7X+ZjBFPLsPCxEJd5eI/Vpy9KpYg5PTSsLqv2udmulmYZzOktP\nYeV/qAaV/h/uQa+yTkxz9q0lixganx+ZSiC/3iTLwQLTI+Em8ayjVcIGQ/A9j6X1\nVCOxHBvy3bcIgZe+ZImwoWvko8ryaHWrdCKz7zVgXPZ9aT6B+VW0qqJGsHS0F5m7\nK0EC8fkdMlRufEiw6DChWUmspg7FYNW3fL7boAXemQKBgQDoOzEwr9khlO32ZXSs\nqIKRGNoL5pZPekVHPc/LI6713Vwg4g52xmtT0ZwgjkUB9QF4CimYVGLHytZG6P0G\nSBAdf4JMeeuBkJtmkXnYdJAlbwRNTiHWz409yAJ9hIyPafLZFKYxMLAqj95wnBxc\nMGq9accLaLIUtG8WGfSUrs5fwwKBgQC9PxTrAMl+ewm2O+a86du+BGsiy8fUvQZX\nJ9xayx9ARjEJXv1cgD4z59mQDn6gzBLrDcH+KY2ZSZUmvPof5LkXUlXXplJWh1Qj\nYvpMx2IOdu2OFFfydtyvq/JbXaEMrvUGU3+pvCF1e7Wxf+jlCTZM4yKwg1Ba9FyT\nCUaPlJFbxwKBgC0wv4y622TWh0voSEEE9Ytoq52fPGaw42ROme3svrInZjMb6jag\nu+fupRQMu077L1L9n0R+P06joPjhg8NCKKik1GUvYG2xBxx5eJ1vaVFvfgXRC3Ky\npsh78Egej/+kXVZy1zhBQja2ElIVfstNvKepOst0jxrKVceWO2rnbU9jAoGAdDH6\nNvxpuyXyZZjL6GwyRq5R1bCHRqC09uh7jKewzXcLfrR7HcOD7bzKQYAU0cfbScVN\nui9rSJX8ZSec794woxgjqt/tKEG5MG0CQAgftb/hxd3Jzg6bG6WYje6kBrSZr0Ov\nW9kuNgM6IPznU1FfrL+9OeG2gdIN0R3d3CSdR1sCgYBxpi1DXXBCXVeU1wIb9XBA\nwewiFSAabF/UtiF7CHkGSMN1lMe/R1AFKM8Irrqbbm0jl00BZ5fgVYV/wVaYZtDw\nPZQmGeO3yGi6FanLnBaxE/bKjk+RkaORM8QoaYGghX59TNoFzHNE1rF0w1lMdrlN\nnsFelOtLls3xNrtNNMxHeg==\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=weyyak-qa@testing-400312.iam.gserviceaccount.com \
    CLIENT_ID=106140453468102648729 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=//www.googleapis.com/robot/v1/metadata/x509/weyyak-content-images-uat%40weyyak-uat.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com \
    EN_BASE_URL=https://msapifo-uat.weyyak.z5.com/v1/en/ \
    AR_BASE_URL=https://msapifo-uat.weyyak.z5.com/v1/ar/
    
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
