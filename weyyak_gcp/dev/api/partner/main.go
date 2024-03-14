package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	// "goroach/handler"
	"masterdata/content"
	"masterdata/docs"
	"masterdata/subscription"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	db     *gorm.DB
	router *gin.Engine
	log    zerolog.Logger
)

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
func main() {
	// Initialize Dependencies
	// Service Port, Database, Logger, Cache, Message Queue etc.
	router := gin.Default()
	router.Use(logMiddleware())
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
	defer fdb.Close()
	ucdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/wk_user_management"
	// log.Info().Msg(dsn)
	udb, err := gorm.Open("postgres", ucdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer udb.Close()
	// db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	fdb.DB().SetMaxIdleConns(10)
	fcdb.DB().SetMaxIdleConns(10)
	udb.DB().SetMaxIdleConns(10)
	fdb.SingularTable(true)
	fcdb.SingularTable(true)
	udb.SingularTable(true)
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
		c.Set("LOG", log)
		c.Set("REDIS", "redis")
	})

	// Boostrap services
	episodeSvc := &content.HandlerService{}
	episodeSvc.Bootstrap(router)

	plansSvc := &subscription.HandlerService{}
	plansSvc.Bootstrap(router)

	// db.AutoMigrate(&subscription.PlanDetails{})
	// --- Development Only ---
	//setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": health()})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3006")
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

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the HTTP method and request URL
		start := time.Now()
		traceID := uuid.New().String()
		c.Request.Header.Set("X-Trace-ID", traceID)

		c.Next()

		end := time.Now()
		duration := end.Sub(start)
		method := c.Request.Method
		httpstatus := c.Writer.Status()
		var logType string

		switch httpstatus {
		case 200:
			logType = "info"
		case 400:
			logType = "alert"
		case 404:
			logType = "fatal"
		case 500:
			logType = "warning"
		}

		Loki(c, logType, method, traceID, duration)
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func Loki(c *gin.Context, logType string, msg string, TraceId string, duration time.Duration) {
	// config, _ := AppConfig(".env")
	v := strconv.FormatInt(time.Now().UnixNano(), 10)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	msg = msg + " " + "traceId=" + TraceId + " " + c.Request.URL.Path

	Lokijson := map[string]interface{}{
		"streams": []interface{}{
			map[string]interface{}{
				"stream": map[string]interface{}{
					"status":     c.Writer.Status(),
					"message":    msg,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"referer":    c.Request.Referer(),
					"duration":   duration.String(),
					"hostname":   hostname,
					"client_ip":  GetOutboundIP(),
					"user_agent": c.Request.UserAgent(),
					"level":      logType,
					"job":        "Partner",
				},
				"values": []interface{}{
					[]interface{}{
						string(v),
						msg,
					},
				},
			},
		},
	}

	jsonValue, _ := json.Marshal(Lokijson)
	url := "http://3.110.118.98:3100/loki/api/v1/push"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	fmt.Println(resp, err)
}
