export SERVICE_PORT=3000 
export DB_SERVER=msapiqa-rds.z5.com 
export DB_SERVER_READER=msapiqa-rds-ro.z5.com 
export DB_PORT=5432 
export DB_USER=weyyak_aurora 
export DB_PASSWORD=M5Ltay9sDY93khvmcpNE 
export DB_DATABASE=wk_user_management 
export CONTENT_DB_DATABASE=wyk_content 
export FRONTEND_DB_DATABASE=wyk_frontend_config 
export TEMPLATE_URL=/templates/ 
export EMAILIMAGEBASEURL=https://s3.ap-south-1.amazonaws.com/mailtemp/ 
export EMAILHEADIMAGEFILENAME=logo.png 
export EMAILCONTENTIMAGEFILENAME=devices.png 
export REDIRECTION_URL=https://wyk2qa.weyyak.com/ 
export ADMIN_MAIL=marathon007@mailnesia.com 
export DEFAULT_PAGE_SIZE=20 
export AWS_REGION=ap-south-1 
export ACCESS_SECRET=AKIAYOGUWMUMEEQD6CPW 
export REFRESH_SECRET=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  
export PASSWORDCHANGEURL=https://wyk2boqa1.weyyak.com/password? 
export ReCAPTCHA_SECRET_web=6LdNV9siAAAAAAIJ9j7sBr6oNgstd7Tx_qLHwQ9x 
export ReCAPTCHA_SECRET_ios=6Ld3UtsiAAAAAKBDjrIBbsZfD3ujFQyT01XIpntd 
export ReCAPTCHA_SECRET_android=6LeaMtgkAAAAAH7T5B_u-EO-75uknIXPWY-6aiZn 
export BASE_URL=https://wyk2qa.weyyak.com 
export EGYPTBASE_URL=https://qa-weyyak1.z5.com 
export SUBSCRIPTION_URL=https://zpapi.wyk.z5.com/orders/ 
export USER_DELETE_URL=https://zpapi.wyk.z5.com/payment/registration/delete?id= 
export SES_REGION=ap-south-1 
export SES_ID=AKIAYOGUWMUMK2O4DT6B 
export SES_SECRET=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 
export DOTNET_URL=https://uat-api.weyyak.z5.com/v1/ar/oauth2/tokendata?access_token= 
export TWITTER_CONSUMER_KEY=9eZDfmSeYROq2unPSqEXbIKrH 
export TWITTER_CONSUMER_SECRET_KEY=WqPZJ3uloFI876gKosGb4zujv2cW8TD4X5jPeWgU7pAd9Mxbmq 
export SESSION_SECRET="70Ie7PiMuS8JUIl1n-CcEP07Les5Y7Nk-eBc8x0jaHLz8ilfow" 
export JEAGER_URL=https://jaeger-tracer.weyyak.com/api/traces
export JEAGER_SERVICE_NAME=User

# export PATH=$(go env GOPATH)/bin:$PATH
# swag init -g main.go
go run *.go


# go get -u github.com/swaggo/swag/cmd/swag