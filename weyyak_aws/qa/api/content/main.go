package main

import (
	"context"
	"fmt"

	// "goroach/handler"
	"content/actor"
	"content/agegroup"
	"content/anchor"
	"content/content"
	"content/dashboard"
	"content/digitalRights"
	"content/docs"
	"content/faq"
	"content/genre"
	"content/mrssfeed"
	"content/multitiercontent"
	"content/productname"
	"content/tags"
	"content/viewactivity"
	"net/http"
	"os"

	"content/channel"
	"content/director"
	"content/episode"
	"content/fragments"
	l "content/logger"
	"content/musiccomposer"
	"content/seasonorepisode"
	"content/singer"
	"content/sitemaps"
	"content/songwriter"
	"content/subscriptionplan"
	"content/writer"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	db     *gorm.DB
	router *gin.Engine
	log    zerolog.Logger
)

func main() {
	// Initialize Dependencies
	// Service Port, Database, Logger, Cache, Message Queue etc.
	router := gin.Default()

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(os.Getenv("JEAGER_URL"))))
	if err != nil {
		log.Error().Err(err)
	}

	// Create a resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(os.Getenv("JEAGER_SERVICE_NAME")),
	)

	// Create a tracer provider with the Jaeger exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Register the tracer provider as the global provider
	otel.SetTracerProvider(tp)

	// Create a Gin router
	// Middleware to start a new span for each request
	router.Use(l.JaegerMiddleware)

	router.Use(CORSMiddleware())
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	// Database
	// dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
	// 	os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
	// 	os.Getenv("DB_DATABASE"), os.Getenv("DB_PASSWORD"),
	// )
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	// log.Info().Msg(dsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()
	fdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("FRONTEND_DB_DATABASE")
	// log.Info().Msg(dsn)
	fdb, err := gorm.Open("postgres", fdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer fdb.Close()
	fcdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/wyk_frontend_config"
	// log.Info().Msg(dsn)
	fcdb, err := gorm.Open("postgres", fcdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer fcdb.Close()
	ucdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/wk_user_management"
	// log.Info().Msg(dsn)
	udb, err := gorm.Open("postgres", ucdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer udb.Close()
	sdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/wyk_sync"
	// log.Info().Msg(dsn)
	sdb, err := gorm.Open("postgres", sdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer sdb.Close()
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "ms-api-qa.jserxp.ng.0001.aps1.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	db.DB().SetMaxIdleConns(10)
	fdb.DB().SetMaxIdleConns(10)
	fcdb.DB().SetMaxIdleConns(10)
	udb.DB().SetMaxIdleConns(10)
	sdb.DB().SetMaxIdleConns(10)
	// db.LogMode(true)
	db.SingularTable(true)
	fdb.SingularTable(true)
	fcdb.SingularTable(true)
	udb.SingularTable(true)
	sdb.SingularTable(true)
	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Weyyak APIs"
	docs.SwaggerInfo.Description = "List of APIs for weyyak project"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("FDB", fdb)
		c.Set("FCDB", fcdb)
		c.Set("UDB", udb)
		c.Set("SDB", sdb)
		c.Set("LOG", log)
		c.Set("CONTEXT", ctx)
		c.Set("REDIS_CLIENT", rdb)
	})

	// Boostrap services
	contentSvc := &content.ContentService{}
	contentSvc.Bootstrap(router)

	actorSvc := &actor.HandlerService{}
	actorSvc.Bootstrap(router)

	tagsSvc := &tags.HandlerService{}
	tagsSvc.Bootstrap(router)

	productSvc := &productname.HandlerService{}
	productSvc.Bootstrap(router)

	agegroupSvc := &agegroup.HandlerService{}
	agegroupSvc.Bootstrap(router)

	// artistSvc := &artist.HandlerService{}
	// artistSvc.Bootstrap(router)
	mrssfeedSvc := &mrssfeed.HandlerService{}
	mrssfeedSvc.Bootstrap(router)

	multitierSvc := &multitiercontent.HandlerService{}
	multitierSvc.Bootstrap(router)

	digitalRightsSvc := &digitalRights.HandlerService{}
	digitalRightsSvc.Bootstrap(router)

	sitemapsSvc := &sitemaps.HandlerService{}
	sitemapsSvc.Bootstrap(router)

	genreSvc := &genre.HandlerService{}
	genreSvc.Bootstrap(router)

	directorSvc := &director.HandlerService{}
	directorSvc.Bootstrap(router)

	musiccomposerSvc := &musiccomposer.HandlerService{}
	musiccomposerSvc.Bootstrap(router)

	singerSvc := &singer.HandlerService{}
	singerSvc.Bootstrap(router)

	songwriterSvc := &songwriter.HandlerService{}
	songwriterSvc.Bootstrap(router)

	writerSvc := &writer.HandlerService{}
	writerSvc.Bootstrap(router)

	seasonepisodecontentvarianceSvc := &seasonorepisode.HandlerService{}
	seasonepisodecontentvarianceSvc.Bootstrap(router)

	viewactivitySvc := &viewactivity.HandlerService{}
	viewactivitySvc.Bootstrap(router)

	subscriptionplanSvc := &subscriptionplan.HandlerService{}
	subscriptionplanSvc.Bootstrap(router)

	episodeSvc := &episode.HandlerService{}
	episodeSvc.Bootstrap(router)

	fragSvc := &fragments.HandlerService{}
	fragSvc.Bootstrap(router)

	dashboard := &dashboard.HandlerService{}
	dashboard.Bootstrap(router)

	channelSvc := &channel.HandlerService{}
	channelSvc.Bootstrap(router)

	anchorSvc := &anchor.HandlerService{}
	anchorSvc.Bootstrap(router)

	faq := &faq.HandlerService{}
	faq.Bootstrap(router)
	// --- Development Only ---
	setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	l.JSON(c, http.StatusOK, gin.H{"status": health()})
}

// Utility function to populate some records into DB
// Should not be used in production.
// Production should use sql scripts to create DB with default tables and data
func setupQuotes(db *gorm.DB) {
	// check if table exists
	// if table exists, return
	db.AutoMigrate(&anchor.Anchor{})
	if !db.HasTable(&content.Content{}) {
		db.AutoMigrate(&content.Content{})

		// quotes := []quote.QuoteModel{
		// 	{Author: "Gandhi", Quote: "The best way to find yourself is to lose yourself in the service of others."},
		// 	{Author: "Duke Ellington", Quote: "A problem is a chance for you to do your best."},
		// 	{Author: "Steve Prefontaine", Quote: "To give anything less than your best, is to sacrifice the gift."},
		// 	{Author: "Peter Drucker", Quote: "The best way to predict the future is to create it."},
		// }

		// for i := range quotes {
		// 	db.Create(&quotes[i])
		// }
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		//fmt.Println(c.Request.Method)

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
