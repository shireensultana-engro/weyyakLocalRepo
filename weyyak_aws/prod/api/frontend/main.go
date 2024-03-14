package main

import (
	"fmt"

	// "goroach/handler"
	"net/http"
	"os"

	"frontend_service/config"
	"frontend_service/content"
	"frontend_service/country"
	"frontend_service/docs"
	"frontend_service/menu"

	"github.com/gin-gonic/gin"
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
	router.Use(CORSMiddleware())
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	// Database
	//In local sslmode=disabled required
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()
	cdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("CONTENT_DB_DATABASE")
	cdb, err := gorm.Open("postgres", cdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer cdb.Close()
	fcdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("FRONTEND_CONFIG_DB_DATABASE")
	fcdb, err := gorm.Open("postgres", fcdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer fcdb.Close()
	udsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("USER_DB_DATABASE")
	udb, err := gorm.Open("postgres", udsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer udb.Close()
	db.DB().SetMaxIdleConns(10)
	cdb.DB().SetMaxIdleConns(10)
	fcdb.DB().SetMaxIdleConns(10)
	udb.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	db.SingularTable(true)
	cdb.LogMode(true)
	cdb.SingularTable(true)
	fcdb.LogMode(true)
	fcdb.SingularTable(true)
	udb.LogMode(true)
	udb.SingularTable(true)
	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Weyyak Frontend APIs"
	docs.SwaggerInfo.Description = "List of Frontend APIs for weyyak project"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("CDB", cdb)
		c.Set("FCDB", fcdb)
		c.Set("UDB", udb)
		c.Set("LOG", log)
		c.Set("REDIS", "redis")
	})
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Println("hostname:", name)
	// Boostrap services
	menuSvc := &menu.HandlerService{}
	menuSvc.Bootstrap(router)

	contentSvc := &content.HandlerService{}
	contentSvc.Bootstrap(router)

	configSvc := &config.HandlerService{}
	configSvc.Bootstrap(router)

	countrySvc := &country.HandlerService{}
	countrySvc.Bootstrap(router)

	// --- Development Only ---
	//setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	fcdb := c.MustGet("FCDB").(*gorm.DB)
	udb := c.MustGet("UDB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	frontendDB:=db.DB().Stats()
	frontendconfigDB:=fcdb.DB().Stats()
	userDB:=udb.DB().Stats()
	contentDB:=cdb.DB().Stats()
	c.JSON(http.StatusOK, gin.H{"status": health(),"db_frontend":frontendDB,"db_frontendconfig":frontendconfigDB,"db_user":userDB,"db_content":contentDB})
}

// Utility function to populate some records into DB
// Should not be used in production.
// Production should use sql scripts to create DB with default tables and data
// func setupQuotes(db *gorm.DB) {
// 	// check if table exists
// 	// if table exists, return
// 	if !db.HasTable(&page.Page{}) {
// 		db.AutoMigrate(&page.Page{})

// 		// quotes := []quote.QuoteModel{
// 		// 	{Author: "Gandhi", Quote: "The best way to find yourself is to lose yourself in the service of others."},
// 		// 	{Author: "Duke Ellington", Quote: "A problem is a chance for you to do your best."},
// 		// 	{Author: "Steve Prefontaine", Quote: "To give anything less than your best, is to sacrifice the gift."},
// 		// 	{Author: "Peter Drucker", Quote: "The best way to predict the future is to create it."},
// 		// }

// 		// for i := range quotes {
// 		// 	db.Create(&quotes[i])
// 		// }
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3003")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, user-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
