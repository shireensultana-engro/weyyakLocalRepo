
# weyyak-ms-api

### Need to point CONTENT_COMMON_API_URL to your local ip address in the Dockerfile.dev  when running any service that depends on content_common_service.

# ports microserices running on 

User - 3000
Content - 3001
Masterdata - 3002
Frontend - 3003
Frontend Config - 3004
Events - 3005
partner - 3006


## build a service
	CD to the directory corresponding to $SERVICE_NAME

	docker build -f Dockerfile.dev -t $DOCKER_IMAGE_NAME .

	example: 
	
		cd content_common_service
		docker build -f Dockerfile.dev -t wk-content-common .

## run a service
	CD to the directory corresponding to $SERVICE_NAME
	
	docker run -p $PORT:$PORT $DOCKER_IMAGE_NAME
	
	example:
	
		docker run -p 3002:3002 wk-content-common

