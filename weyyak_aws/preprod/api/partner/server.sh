export SERVICE_PORT=3006 
export DB_SERVER=msapiqa-rds.z5.com 
export DB_PORT=5432 
export DB_USER=weyyak_aurora 
export DB_PASSWORD=M5Ltay9sDY93khvmcpNE 
export DB_DATABASE=wyk_content 
export DEFAULT_PAGE_SIZE=50 
export API_VERSION=/v3 
export BASE_URL=https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa 
export IMAGE_URL=https://contents-uat.weyyak.z5.com/ 
export VIDEO_API=https://api-weyyak.akamaized.net/get_info/ 
export LOGIN_API=https://msapifo.weyyak.z5.com/v1/oauth2/token 
export CONTENT_URL=https://apistg.weyyak.z5.com/v1/get_info/ 
export IMAGE_URL_GO=https://z5content-uat.s3.amazonaws.com/v1/ 
export DETAILS_BACKGROUND=/details-background
export POSTER_IMAGE=/poster-image
export MOBILE_DETAILS_BACKGROUND=/mobile-details-background
export OVERLAY_POSTER_IMAGE=/overlay-poster-image
export DUBBLING_SCRIPT=/dubbling-script
export WATCH_NOW=cce8db39-3d54-442c-8f21-f6c1aa11d396
export SUBTITLING_SCRIPT=/subtitling-script
export REDIS_CONTENT_KEY=GOAPIQA
export REDIS_CACHE_URL=http://localhost:3005/cache  
export IMAGERY_URL=https://contents-uat.weyyak.z5.com/

go run *.go