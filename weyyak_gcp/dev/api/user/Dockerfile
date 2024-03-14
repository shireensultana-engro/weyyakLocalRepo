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
    DB_SERVER=10.87.0.5 \
    DB_SERVER_READER=10.87.0.5 \
    DB_PORT=5432 \
    DB_USER=appuser \
    DB_PASSWORD=weyyakbo345123DEV \
    DB_DATABASE=wk_user_management \
    CONTENT_DB_DATABASE=wyk_content \
    FRONTEND_DB_DATABASE=wyk_frontend_config \
    TEMPLATE_URL=/templates/ \
    EMAILIMAGEBASEURL=https://s3.ap-south-1.amazonaws.com/mailtemp/ \
    EMAILHEADIMAGEFILENAME=logo.png \
    EMAILCONTENTIMAGEFILENAME=devices.png \
    REDIRECTION_URL=https://weyyak-dev.engro.in/ \
    ADMIN_MAIL=marathon007@mailnesia.com \
    DEFAULT_PAGE_SIZE=20 \
    AWS_REGION=ap-south-1 \
    ACCESS_SECRET=AKIAYOGUWMUMEEQD6CPW \
    REFRESH_SECRET=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 \ 
    PASSWORDCHANGEURL=https://weyyak-dev.engro.in/password? \
    ReCAPTCHA_SECRET_web=6LdNV9siAAAAAAIJ9j7sBr6oNgstd7Tx_qLHwQ9x \
    ReCAPTCHA_SECRET_ios=6Ld3UtsiAAAAAKBDjrIBbsZfD3ujFQyT01XIpntd \
    ReCAPTCHA_SECRET_android=6LeaMtgkAAAAAH7T5B_u-EO-75uknIXPWY-6aiZn \
    BASE_URL=https://weyyak-dev.engro.in \
    EGYPTBASE_URL=https://qa-weyyak1.z5.com \
    SUBSCRIPTION_URL=https://zpmsapi-dev.engro.in/orders/ \
    USER_DELETE_URL=https://zpmsapi-dev.engro.in/payment/registration/delete?id= \
    SES_REGION=ap-south-1 \
    SES_ID=AKIAYOGUWMUMK2O4DT6B \
    SES_SECRET=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 \
    DOTNET_URL=https://api-backoffice-production.weyyak.com/oauth2/tokendata?access_token= \
    JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces \
    JEAGER_SERVICE_NAME=User \
    TWITTER_CONSUMER_KEY=9eZDfmSeYROq2unPSqEXbIKrH \
    TWITTER_CONSUMER_SECRET_KEY=WqPZJ3uloFI876gKosGb4zujv2cW8TD4X5jPeWgU7pAd9Mxbmq \
    AWS_SNS_REGION=us-east-1          

EXPOSE ${SERVICE_PORT}

# Command to run
ENTRYPOINT ["/main"]
