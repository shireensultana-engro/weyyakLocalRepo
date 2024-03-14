# Makefile for various commands

# Env Vars 
APP_NAME=refactor_weyyak_api
APP_VERSION=1.1
NETWORK=roachnet

SERVICE_PORT=3001
DB_SERVER=weyyak-ms-api-qa-postgresql.cklc9f36hm30.ap-south-1.rds.amazonaws.com
DB_IMAGE=cockroach:v21.1.1
DB_PORT=5432
DB_USER=weyyakadmin
DB_DATABASE=wk_content
DB_PASSWORD=3hb0kQwBr9mkwC
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
