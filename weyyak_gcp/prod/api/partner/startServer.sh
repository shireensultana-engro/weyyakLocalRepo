export SERVICE_PORT=3006 
export DB_SERVER=34.18.48.27
export DB_SERVER_READER=34.18.48.27
export DB_PORT=5432 
export DB_USER=appuser 
export DB_PASSWORD=Yy017EUY5aSz 
export DB_DATABASE=wyk_content 
export DEFAULT_PAGE_SIZE=50 
export API_VERSION=/v3 
export BASE_URL=https://msapifo-prod-me.weyyak.com 
export IMAGE_URL=https://content.weyyak.com/ 
export VIDEO_API=https://api-weyyak.akamaized.net/get_info/
export LOGIN_API=https://msapifo-prod-me.weyyak.z5.com/v1/oauth2/token 
export CONTENT_URL=https://swagger.weyyak.com/v1/get_info/ 
export IMAGE_URL_GO=https://content.weyyak.com/ 
export DETAILS_BACKGROUND=/details-background
export POSTER_IMAGE=/poster-image
export MOBILE_DETAILS_BACKGROUND=/mobile-details-background
export OVERLAY_POSTER_IMAGE=/overlay-poster-image
export DUBBLING_SCRIPT=/dubbling-script
export WATCH_NOW=3e53035d-b315-4fb9-8cb7-6a419538b6b8
export SUBTITLING_SCRIPT=/subtitling-script
export REDIS_CONTENT_KEY=GOAPIPROD
export REDIS_CACHE_URL=http://localhost:3005/cache  
export IMAGERY_URL=https://content.weyyak.com/

# export PATH=$(go env GOPATH)/bin:$PATH
# swag init -g main.go
go run *.go