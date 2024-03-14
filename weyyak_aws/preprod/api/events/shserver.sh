export GIN_MODE=release
export SERVICE_PORT=3005 
export REDIS_SERVER=localhost:6379
export USE_CACHE=true 
export LOG_MODE=FILE
export ELASTICSEARCH_URL=http://10.33.82.63:9200/
export ELASTIC_USER=elastic
export ELASTIC_PASSWORD=changeme
go run *.go