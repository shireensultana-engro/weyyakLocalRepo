package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"frontend_config/config"
	"frontend_config/contentType"
	"frontend_config/country"
	"frontend_config/docs"
	"frontend_config/language"
	l "frontend_config/logger"
	"frontend_config/menu"
	"frontend_config/page"
	"frontend_config/playlist"
	"frontend_config/slider"
	"frontend_config/user"

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
	router.Use(CORSMiddleware())

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

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	withoutPassDsn := "postgres://" + os.Getenv("DB_USER") + ":******@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	log.Info().Msg(withoutPassDsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()
	cdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("CDB_DATABASE")
	cdb, err := gorm.Open("postgres", cdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer cdb.Close()
	udsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("USER_DB_DATABASE")
	udb, err := gorm.Open("postgres", udsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer udb.Close()
	fdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("FDB_DATABASE")
	fdb, err := gorm.Open("postgres", fdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer fdb.Close()
	sdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/wyk_sync"
	// log.Info().Msg(dsn)
	sdb, err := gorm.Open("postgres", sdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer sdb.Close()
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.243.128.6:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	db.DB().SetMaxIdleConns(10)
	cdb.DB().SetMaxIdleConns(10)
	udb.DB().SetMaxIdleConns(10)
	fdb.DB().SetMaxIdleConns(10)
	sdb.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	db.SingularTable(true)
	cdb.LogMode(true)
	cdb.SingularTable(true)
	udb.LogMode(true)
	udb.SingularTable(true)
	fdb.LogMode(true)
	fdb.SingularTable(true)
	sdb.LogMode(true)
	sdb.SingularTable(true)

	// Swagger info
	docs.SwaggerInfo.Title = "Weyyak Frontend Config APIs"
	docs.SwaggerInfo.Description = "List of APIs for Frontend Config"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("CDB", cdb)
		c.Set("FDB", fdb)
		c.Set("LOG", log)
		c.Set("UDB", udb)
		c.Set("FDB", fdb)
		c.Set("SDB", sdb)
		//c.Set("REDIS", "redis")
		c.Set("CONTEXT", ctx)
		c.Set("REDIS_CLIENT", rdb)
	})

	countrySvc := &country.HandlerService{}
	countrySvc.Bootstrap(router)

	userSvc := &user.HandlerService{}
	userSvc.Bootstrap(router)

	languageSvc := &language.HandlerService{}
	languageSvc.Bootstrap(router)

	configSvc := &config.HandlerService{}
	configSvc.Bootstrap(router)

	sliderSvc := &slider.HandlerService{}
	sliderSvc.Bootstrap(router)

	pageSvc := &page.HandlerService{}
	pageSvc.Bootstrap(router)

	contentType := &contentType.HandlerService{}
	contentType.Bootstrap(router)

	menuSvc := &menu.HandlerService{}
	menuSvc.Bootstrap(router)

	playlistSvc := &playlist.HandlerService{}
	playlistSvc.Bootstrap(router)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	l.JSON(c, http.StatusOK, gin.H{"status": health()})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3002")
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
