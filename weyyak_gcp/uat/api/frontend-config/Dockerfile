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
    DB_SERVER=wyk2db-uat.engro.in \
    DB_SERVER_READER=wyk2db-uat.engro.in \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123UAT \
    DB_DATABASE=wyk_frontend_config \
    CDB_DATABASE=wyk_content \
    FDB_DATABASE=wk_frontend \
    USER_DB_DATABASE=wk_user_management \
    DEFAULT_PAGE_SIZE=50 \
    CMS=https://msapifo-uat.weyyak.z5.com/v1/ \
    UM=https://msapifo-uat.weyyak.z5.com/v1/ \
    IMAGES=https://content-uat.weyyak.com/ \
    GEO_LOCATION=https://weyyak-geoapi-uat.engro.in \
    VIDEO_API=https://api-weyyak.akamaized.net/get_info/ \
    AD_TAG_URL=https://storage.googleapis.com/z5xml-qa/mobile_apps_ads_ios.xml \
    BASE_URL=https://msapifo-uat.weyyak.z5.com/v1/ \
    S3_BUCKET=z5content-uat \
    S3_REGION=ap-south-1 \
    S3_ID=AKIAYOGUWMUMGAQLMW3U \
    S3_SECRET=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 \
    HOME_PAGE_FOR_TV=https://msapifo-uat.weyyak.z5.com/v1/ \
    REDIS_RDS=34.93.213.137:6379 \
    REDIS_CONTENT_KEY=GOAPIQA \
    REDIS_CACHE_URL=https://msapiuat-events.z5.com/v1/cache \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    RESIZE_IMAGE_URL=https://msapiuat-image.z5.com/crop? \
    JEAGER_SERVICE_NAME=Frontend-Config \
    BUCKET_NAME=wyk-content-uat \
    TYPE=service_account \
    PROJECT_ID=testing-400312 \
    PRIVATE_KEY_ID=dd03c9b33e8647b787124ad3a0775e5bcf95cf3b \
    PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC4nz+BJqcU1t0V\n9EogDpucUrxDOfWiS5z0NGTYmThH3UZF13D1Ape3Bv+aE+OWfSXkL+95+HHwSCHi\nBk8qCcluIPz7BxVYVgkP/9czg9cX4OHK7Aho7wofbznZqLr32wCnUc4Vhf0Wgd2I\npzld7U71oA2V6FuL5BMlQ/Kw/Sph/oMFLYy/PuBfL/YNnhU6jdaNWAibQM1JruW/\nHBYbcKONSrTs50WPSoW1ZZgXJFeqgAz9MSGuvBhFu6/LSUINCCrnyotkJY65p4yL\n+jpEdaqzbw5xkR/Sckk2t+DFZyRowamNXAr4iGsk/5KnKokv0fcaCLRriaHpSdXs\natckZmtZAgMBAAECggEABsZpSS8G8J/V6ylU0wpWMY1jtT/aQMNvlhKaJMyyoQiA\nvK1kSsl2kdPi61+ReYNMMayEqEAyxjOPcsDyhMNpLe5t+jRPVzeJC5pC5nQPH6Q0\nBQNWZ6tl/rNRNyiW++OAiaNZ4bZSDFJls88XLtg3jpH6DadCPMb44OQ2csJHnOLv\nK5dMzUenvh2V/zWVvNS2CaTEFqzWowbzIek/Cv2hWPizo6ZwgplTKQu1iFfbrBui\nYs8WvvpsxA8X216SyECCzo9eZv2lAwFtLedvf/4MfhNL4q6/0KTHDiTDWuWkx0f1\nzVeLbK8b4LfXplZomF9N4fcq/z8SVKpFE1DjuL9jEQKBgQDfiuCHw3K/hOAZCF+c\nNyZ5M2KIyooOv1mxZrVBifvZjlO6kji6L2L9BtIYryYRG0NLrXEpmVZz7IM3NuT1\nrCfuzKQ+l07dEtg6VngoYup1JDKrWTJD4v2cqwOJJAXdQNU03b+shTmXZvW5tJ2J\nbSibIEE6Ol9BFA3b69FV08m9PQKBgQDTbbC3bhIPa5MEo8j2kvM/S9KdJaomjST3\nywCl64sbQXzSFKH3hMB+h7zX2ztEU+o1Ls7P8lvFFQEg7R78kDZy6ulPosm1GFcE\nydOriruXDMHESBjD8YWGO+PevsJIp3j0SADFdvoJg/MHcdpJBI/JDZKUA5mNP52B\n+gD3MCKATQKBgHyh1lt7Oe/TqonqZDwZd4bdglNX8S8VunExHV+kCdmbA82ilqQf\npWYDNoHyrRuegp+f3NbfmhbZx7KyFWdvi3gVeoE3JQ4W4p1r9mQ1+hhDjUiBW4gD\n93gw1LDSd76K1hQ6ihIq2RgznE7kh2zGgnwyuIs5XkBPaQazbKwYf4LVAoGAUES3\nr5whTStpIHzSAhLeOKyfpDu1cndpjo3KjDN1l35wVg2xRBhpQGBmKIk54gH9y/0e\nVUJM5vDHgemkNvzFPzHCLBLAg3lfKBk7vEeqWnlkYxGAHXvnVoQMSfegKRczy4I+\nkLlPyicHME9gMRKDSDBX8su/EyoQsVTp4u6qWo0CgYEAoAH/l+OI6TbWmvA3oISu\n6lVDSgzp28FEJWsxZTleQHTXTa3qV2ePq7IS3Nd021Kojyn5gZ3U6R5nhg114LM+\nojqUSY3MQeXBPYfGaLTIDn/fy27Iq6R6UFVKpB3P61oQxwPSV3QcbbzzsvFcSFOZ\n5WXXSnGkpsn3iFOdCGKNsE8=\n-----END PRIVATE KEY-----\n" \
    CLIENT_EMAIL=weyyak-qa@testing-400312.iam.gserviceaccount.com \
    CLIENT_ID=106140453468102648729 \
    AUTH_URI=https://accounts.google.com/o/oauth2/auth \
    TOKEN_URI=https://oauth2.googleapis.com/token \
    AUTH_PROVIDER_X509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs \
    CLIENT_X509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/weyyak-content-images-uat%40weyyak-uat.iam.gserviceaccount.com \
    UNIVERSE_DOMAIN=googleapis.com \
    EN_BASE_URL=https://msapifo-uat.weyyak.z5.com/v1/en/ \
    AR_BASE_URL=https://msapifo-uat.weyyak.z5.com/v1/ar/


EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
