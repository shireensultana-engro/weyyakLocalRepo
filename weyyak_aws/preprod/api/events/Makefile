# Makefile for various commands

# Env Vars 
APP_NAME=events
APP_VERSION=1.0
NETWORK=eventnet
SERVICE_PORT=8880

REDIS_IMAGE=redis:6.2.4-alpine
REDIS_PORT=6379
REDIS_SERVER=redis-server

all: network build runredis runapp 

delnw:
	docker network rm ${NETWORK}

nw:
	docker network create -d bridge ${NETWORK}

build:
	docker build -t ${APP_NAME}:${APP_VERSION} .

runredis: 
	docker run --rm -d --net ${NETWORK} --name $(REDIS_SERVER) \
		-p ${REDIS_PORT}:${REDIS_PORT} -v redis.conf:/usr/local/etc/redis/redis.conf $(REDIS_IMAGE)

runapp:
	docker run --rm -d --net ${NETWORK} --name ${APP_NAME} \
		-e SERVICE_PORT=${SERVICE_PORT} -e REDIS_SERVER=${REDIS_SERVER}:${REDIS_PORT} \
		-e USE_CACHE=true -e LOG_MODE=C \
		-p ${SERVICE_PORT}:${SERVICE_PORT} ${APP_NAME}:${APP_VERSION}

clean:
	if docker network rm ${NETWORK}; then echo ''; fi
	if docker stop ${APP_NAME}; then echo ''; fi
	if docker stop ${REDIS_SERVER}; then echo ''; fi

esup:
	docker run --rm -d --name es01 --net eventnet -p 9200:9200 -p 9300:9300 -e discovery.type=single-node docker.elastic.co/elasticsearch/elasticsearch:7.13.3
	docker run --rm -d --name kib01 --net eventnet -p 5601:5601 -e ELASTICSEARCH_HOSTS=http://es01:9200 docker.elastic.co/kibana/kibana:7.13.3

esdown:
	if docker stop es01; then echo ''; fi
	if docker stop kib01; then echo ''; fi