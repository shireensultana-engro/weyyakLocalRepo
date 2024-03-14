# Makefile for various commands

# Env Vars 
APP_NAME=user_management
APP_VERSION=1.1
NETWORK=roachnet

SERVICE_PORT=3000
DB_SERVER=localhost
DB_IMAGE=db
DB_PORT=26257
DB_USER=api_testing
DB_PASSWORD=strongpassword
DB_DATABASE=wk_user_management
DEFAULT_PAGE_SIZE=20

build:
	docker build -t ${APP_NAME}:${APP_VERSION} .

network:
	docker network create -d bridge ${NETWORK}

runapp:
	docker run --rm -d --net ${NETWORK} --name ${APP_NAME} \
		-e SERVICE_PORT=${SERVICE_PORT}  -e DB_SERVER=${DB_SERVER} \
		-e DB_PORT=${DB_PORT} -e DB_USER=${DB_USER} -e DB_DATABASE=${DB_DATABASE} \
		-e DB_PASSWORD=${DB_PASSWORD} -e DEFAULT_PAGE_SIZE=${DEFAULT_PAGE_SIZE} \
		-p ${SERVICE_PORT}:${SERVICE_PORT} ${APP_NAME}:${APP_VERSION}

rundb:
	docker run --rm -d --net ${NETWORK} --name $(DB_SERVER) \
		-e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-p ${DB_PORT}:${DB_PORT} $(DB_IMAGE)

clean:
	docker stop ${APP_NAME}
	docker stop ${DB_SERVER}
	docker network rm ${NETWORK}
