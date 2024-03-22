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
COPY --from=builder /build/templates/ar /
COPY --from=builder /build/templates/en /
ENV SERVICE_PORT=3000 \
    DB_SERVER=wyk2db-uat.engro.in \
    DB_SERVER_READER=wyk2db-uat.engro.in \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123UAT \
    DB_DATABASE=wk_user_management \
    CONTENT_DB_DATABASE=wyk_content \
    FRONTEND_DB_DATABASE=wyk_frontend_config \
    TEMPLATE_URL=/templates/ \
    EMAILIMAGEBASEURL=https://s3.ap-south-1.amazonaws.com/mailtemp/ \
    EMAILHEADIMAGEFILENAME=logo.png \
    EMAILCONTENTIMAGEFILENAME=devices.png \
    REDIRECTION_URL=https://wyk2qa.weyyak.com/ \
    ADMIN_MAIL=marathon007@mailnesia.com \
    DEFAULT_PAGE_SIZE=20 \
    AWS_REGION=ap-south-1 \
    ACCESS_SECRET=AKIAYOGUWMUMEEQD6CPW \
    REFRESH_SECRET=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 \ 
    PASSWORDCHANGEURL=https://api-uat.z5backoffice.com/password? \
    ReCAPTCHA_SECRET_web=6LdNV9siAAAAAAIJ9j7sBr6oNgstd7Tx_qLHwQ9x \
    ReCAPTCHA_SECRET_ios=6Ld3UtsiAAAAAKBDjrIBbsZfD3ujFQyT01XIpntd \
    ReCAPTCHA_SECRET_android=6LedaJwdAAAAAE35WnRWzqcQxuQauQX5HCFoCT5y \
    BASE_URL=//wyk2.weyyak.com \
    EGYPTBASE_URL=https://uat-weyyak1.z5.com \
    SUBSCRIPTION_URL=https://zpapi.wyk.z5.com/orders/ \
    USER_DELETE_URL=https://zpapi.wyk.z5.com/payment/registration/delete?id= \
    SES_REGION=ap-south-1 \
    SES_ID=AKIAYOGUWMUMK2O4DT6B \
    SES_SECRET=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 \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    JEAGER_SERVICE_NAME=User \
    TWITTER_CONSUMER_KEY=9eZDfmSeYROq2unPSqEXbIKrH \
    TWITTER_CONSUMER_SECRET_KEY=WqPZJ3uloFI876gKosGb4zujv2cW8TD4X5jPeWgU7pAd9Mxbmq \
    BOPASSWORDCHANGEURL=https://wyk2bo.weyyak.com/reset-password? 

EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
