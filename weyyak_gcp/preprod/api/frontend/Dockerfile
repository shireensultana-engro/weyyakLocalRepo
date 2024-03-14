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
	DB_SERVER=prod2-db.weyyak.com \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=Yy017EUY5aSz \
	DB_DATABASE=wk_frontend \
	CONTENT_DB_DATABASE=wyk_content \
	FRONTEND_CONFIG_DB_DATABASE=wyk_frontend_config \
	USER_DB_DATABASE=wk_user_management \
	DEFAULT_PAGE_SIZE=20 \
	CONTENT_COMMON_API_URL= \
	CMS=https://msapifo-prod-me.weyyak.z5.com/v1/ \
    UM=https://msapifo-prod-me.weyyak.z5.com/v1/ \
    IMAGES=https://content.weyyak.com/ \
    GEO_LOCATION=https://geo.weyyak1.z5.com \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://s3.ap-south-1.amazonaws.com/z5xml/mobile_apps_ads_ios.xml \
    BASE_URL=https://msapifo-prod-me.weyyak.z5.com/ \
    USER_ACTIVITY_URL=https://msapiprod-me-events.z5.com/event/activity \
	USER_LOG_URL=https://msapiprod-me-events.z5.com/event/log \
    CONTENT_TYPE_URL=https://msapifo-prod-me.weyyak.z5.com/v1/en/contents/contentType?contentType= \
    CONTENT_TYPE_URL_PAGINATION=&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0 \
    S3_BUCKET=z5content \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMLWREWLMX \
    S3_SECRET=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 \
    S3_URLFORCONFIG=https://z5content.s3.ap-south-1.amazonaws.com/configprod.json \
    REDIS_CACHE_URL=https://msapiprod-me-events.z5.com/v1/cache \
    REDIS_CONTENT_KEY=GOAPIPROD \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    GCP_URLFORCONFIG=https://storage.googleapis.com/weyyak-prod/configprodmigration.json \
    BUCKET_NAME=weyyak-prod \
    TYPE=service_account \
    PROJECT_ID=weyyak-preprod \
    PRIVATE_KEY_ID=ca04f32ec04c7f308872a2230187e8731b7f87ea \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1nukPNB8e8iXM\nMhHr1iuvSGwienHPd4avprk0yUIXAlBZzwvbcK9i8V5yfpZzz6RUQcwshPOs5k9r\n3hBMy7zTGiZeyh1tPHdSumn3c4o7vL90RivGKff0VFvbPk4GdcuUFrEOJEH5gMS3\nYBTJtKKhOxKK3hqG/e0WQVjPJidfZwIKFthq9+z+d/4GMpldJAI3CPRpw9l+xzoC\n+vOueZ0aUCaSMvLKqgVzsKp5+YxGAEZbdxeYPMZGffJlZVedBwFNnyELBL8uKmVi\nmnCABwMjCTRwL3bPSgJ9mHLm2FiIK3heJ6Tg5HFjjIHIrxcdbVG57lKoXOt2wKed\n23l24T9tAgMBAAECggEAAvM8+unWbG6qjzmvLPtn1kzLpXEoEEEd8ssxMqJIqCOM\nLHCGOubJnZXZ4evNMbH3BcjHirUcWTvluUW2Rh4GiA/KIdEKIdoXL1bzORTMvG7d\nhoOI/69agNtAgwIp6ZTO+K24QODQnBrNPtccJ7cXaamqoFI4XgsHc7Q2jsfNC2bp\nAaIi4ZZHLhhQf94KFmfqOsVMhX7nmaBjaVZrpIfSM+5g0ESKbYbaLgdg/yVwpbdQ\nrLjolOvZw5r8e8ZVstdfU/GwihHuNsbgTbU511IeUYd+YmxoCZ1fkJq0Xf0uJ6Cy\nz1byOXfFOps8RurZhR1hkUknfeBaTBGVrujlHcoV6wKBgQDgBPKgOlzEWh/I5Epz\nviJZKa3TTJS2kinIGDYjbiJNp5NucQZfkJu0xBn0vztwmFnIdYnIOV7kiWuWulVM\nzjC3KoSBiC+GVGjABukuU/dlcWpbSttRuKg90gJ/gOtF4FuYLWhZkJGM2iCndanv\nkFmylCMoq6aiPnC73VGX2mfvpwKBgQDPjHTAzU1RMWiymF9yO25RZ6jzcyB6hPXf\n2NG2YJ0luM41pZMx5DlRFi8ky6YGK2gFwNsyBBXhRh5AdciGD96FwbqtLEv6eKCn\nC1BxZceYAdA/P6Aa6h/4Wv7J67THKEYbwzGgPYE5Js+jmJXG3tKhtPqnlWa7Zx6H\nLS5uNM7aywKBgF4z1m9oe3AaUflhfqlzT/BcpXsQXgz0I9u/yqxVeNlc2ZN8tehj\n4AZA3IVeETnE5yRzwM/QyEWkP/jvPEWDA1tS5sutoAaF4lK11UKlDoi7C7V+IgIY\ne68ba++AH++PbBTvK01WjM5FP6wLv70832tH/gzxOa5KQY/OfqwzrLdLAoGBAMoa\nAoLQJ+7dRw9KEv9AYf9BCqLtw32qxWYRUrzePYhC+gIBVmEp1KpiCMwyxluRnvyj\nPI7qrYes6L5qMzZgc5YZ/LauwNmI5x9ihBW4P3CEq407XqN2wmTr7tke/e1FCWf1\nXfiki5XkdiLe7VI3HjI68i2H7P6lvnNxCppkL92bAoGAERQ76Setyiz6VOiy+AJ6\nRVcQg2PDpoP0woWj7AOstCwbP91AT1h0Wq/aXRm1lk10Yvq0zm8RNzAUkPrfqLwx\n1pC5SQPrut0h+RZaQPRzUJERxfzMzej/WStGh51E9gRyFdSKfL/iOQ6ZT8PysEwy\ntcDR1sRoS24TmmtlgwyP91o=\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=weyyak-preprod@weyyak-preprod.iam.gserviceaccount.com \
    CLIENT_ID=104965307123923002514 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/weyyak-preprod%40weyyak-preprod.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com
EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
