# export SERVICE_PORT=3006 
# export DB_SERVER=34.100.228.116
# export DB_SERVER_READER=34.100.228.116
# export DB_PORT=5432 
# export DB_USER=appuser 
# export DB_PASSWORD=weyyakbo345123UAT 
# export DB_DATABASE=wyk_content 
# export DEFAULT_PAGE_SIZE=50 
# export API_VERSION=/v3 
# export BASE_URL=https://msapifo.weyyak.z5.com 
# export IMAGE_URL=https://content-uat.weyyak.com/ 
# export VIDEO_API=https://api-weyyak.akamaized.net/get_info/ 
# export LOGIN_API=https://msapifo.weyyak.z5.com/v1/oauth2/token 
# export CONTENT_URL=https://apistg.weyyak.z5.com/v1/get_info/ 
# export IMAGE_URL_GO=https://content-uat.weyyak.com/ 
# export DETAILS_BACKGROUND=/details-background
# export POSTER_IMAGE=/poster-image
# export MOBILE_DETAILS_BACKGROUND=/mobile-details-background
# export OVERLAY_POSTER_IMAGE=/overlay-poster-image
# export DUBBLING_SCRIPT=/dubbling-script
# export WATCH_NOW=cce8db39-3d54-442c-8f21-f6c1aa11d396
# export SUBTITLING_SCRIPT=/subtitling-script
# export REDIS_CONTENT_KEY=GOAPIQA
# export REDIS_CACHE_URL=http://localhost:3005/cache  
# export IMAGERY_URL=https://content-uat.weyyak.com/

# go run *.go



export SERVICE_PORT=3007
export DB_SERVER=34.100.228.116
export DB_PORT=5432 
export DB_USER=appuser 
export DB_PASSWORD=weyyakbo345123UAT 
export DB_DATABASE=wyk_content \
export FRONTEND_DB_DATABASE=wk_frontend 
export DEFAULT_PAGE_SIZE=50 
export API_VERSION=/v3 
export BASE_URL=https://msapifo-uat.weyyak.z5.com 
export IMAGE_URL=https://weyyak-content-uat.engro.in/ 
export VIDEO_API=https://api-weyyak.akamaized.net/get_info/ 
export LOGIN_API=https://msapifo-uat.weyyak.z5.com/v1/oauth2/token  
export CONTENT_URL=https://apistg.weyyak.z5.com/v1/get_info/ 
export IMAGE_URL_GO=https://weyyak-content-uat.engro.in/ 
export DETAILS_BACKGROUND=/details-background
export POSTER_IMAGE=/poster-image
export MOBILE_DETAILS_BACKGROUND=/mobile-details-background
export OVERLAY_POSTER_IMAGE=/overlay-poster-image
export DUBBLING_SCRIPT=/dubbling-script
export WATCH_NOW=cce8db39-3d54-442c-8f21-f6c1aa11d396
# export WATCH_NOW=57fb3648-4bfa-e911-8259-02e5c03f648e
export SUBTITLING_SCRIPT=/subtitling-script
export REDIS_CONTENT_KEY=GOAPIUAT
export REDIS_CACHE_URL=https://msapiuat-events.z5.com/v1/cache  
export IMAGERY_URL=https://weyyak-content-uat.engro.in/ 
export S3_BUCKET=z5content-uat
# export PATH=$(go env GOPATH)/bin:$PATH
# swag init -g main.go
go run *.go