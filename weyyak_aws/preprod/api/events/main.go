package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

func main() {
	// Read ENV into config
	config, err := AppConfig(".env")
	if err != nil {
		log.Fatalf("Error reading config: %s.", err)
	}

	// Initialize Gin Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Connect to Redis Server for Caching and Pub/Sub
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: "",
		DB:       0,
	})

	// create elastic connection - we are using OpenSearch due to Elastic license changes
	// the configuration parameters are different and you need to check OpenSearch docs
	// OpenSearch is fork of Elastic 7.10.2, hence using elasticsearch Go library
	fmt.Println(config.ElasticUser)
	/*esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticUrl},
		Username:  config.ElasticUser,
		Password:  config.ElasticPassword,
		Header:    map[string][]string{},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Logger: &estransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true},
	})
	if err != nil {
		log.Fatalf("Error connecting to Elastic: %s.", err)
	}*/
	// connecting to elastic search with demo credentials
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://65.0.204.174:9200/",
		},
		Username: "elasticadmin",
		Password: "elasticadmin",
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error connecting to Elastic: %s.", err)
	}
	// Setup Logger
	log := NewLogger()

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("CONFIG", config)
		c.Set("REDIS_CLIENT", rdb)
		c.Set("CONTEXT", ctx)
		c.Set("LOG", log)
		c.Set("ELASTIC", esClient)
	})

	// Cache Services
	router.POST("/cache", CacheSetKey)
	router.GET("/cache/:key", CacheGetKey)
	router.DELETE("/cache/:key", CacheRemoveKey)

	// Log request -- to be removed
	router.POST("/log", Log)

	// Event logging
	router.POST("/event/:type", Event)
	// 	Activity logging
	router.POST("/watchingevent/:type", WatchingEvent)

	router.GET("/health", healthsvc)

	// log.Info().Msg("Starting server on :" + port)
	router.Run(":" + config.Port)
}
func healthsvc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": health()})
}
