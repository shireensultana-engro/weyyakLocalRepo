package main

import (
	"fmt"
	"net/http"
	"os"

	"masterdata/ageRatings"
	"masterdata/contentType"
	"masterdata/country"
	"masterdata/digitalRights"
	"masterdata/docs"
	"masterdata/genre"
	"masterdata/language"
	"masterdata/ott"

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

func main() {
	// Initialize Dependencies
	// Service Port, Database, Logger, Cache, Message Queue etc.
	router := gin.Default()
	router.Use(CORSMiddleware())
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	// Database
	// dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
	// 	os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
	// 	os.Getenv("DB_DATABASE"), os.Getenv("DB_PASSWORD"),
	// )
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	withoutPassDsn := "postgres://" + os.Getenv("DB_USER") + ":******@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	log.Info().Msg(withoutPassDsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()

	db.LogMode(true)
	db.SingularTable(true)

	// Swagger info
	docs.SwaggerInfo.Title = "Weyyak APIs(Masterdata Module)"
	docs.SwaggerInfo.Description = "List of APIs for master data module"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("LOG", log)
		c.Set("REDIS", "redis")
	})

	// Boostrap services
	contentType := &contentType.HandlerService{}
	contentType.Bootstrap(router)
	// router.GET("/onetier", contentType.GetAllConentOneTier)
	// router.GET("/multitier", contentType.GetAllConentMultiTier)

	countrySvc := &country.HandlerService{}
	countrySvc.Bootstrap(router)

	digitalRightsSvc := &digitalRights.HandlerService{}
	digitalRightsSvc.Bootstrap(router)

	languageSvc := &language.HandlerService{}
	languageSvc.Bootstrap(router)

	ottSvc := &ott.HandlerService{}
	ottSvc.Bootstrap(router)

	ageRatingSvc := &ageRatings.HandlerService{}
	ageRatingSvc.Bootstrap(router)

	genreSvc := &genre.HandlerService{}
	genreSvc.Bootstrap(router)
	// --- Development Only ---
	// setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}

func healthsvc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": health()})
}

// Utility function to populate some records into DB
// Should not be used in production.
// Production should use sql scripts to create DB with default tables and data
func setupQuotes(db *gorm.DB) {
	// check if table exists
	// if table exists, return
	if !db.HasTable(&country.Country{}) {
		db.AutoMigrate(&country.Country{})

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
