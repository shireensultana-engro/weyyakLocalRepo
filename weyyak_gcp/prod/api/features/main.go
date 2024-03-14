package main

import (
	"feature/competition"
	"feature/docs"
	"feature/marathon"
	"fmt"
	"net/http"
	"os"

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
	cdbsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("CDB_DATABASE")
	// log.Info().Msg(dsn)
	cdb, err := gorm.Open("postgres", cdbsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer cdb.Close()
	// db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	fdb.DB().SetMaxIdleConns(10)
	fcdb.DB().SetMaxIdleConns(10)
	udb.DB().SetMaxIdleConns(10)
	cdb.DB().SetMaxIdleConns(10)
	cdb.SingularTable(true)
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
		c.Set("CDB", cdb)
		c.Set("UDB", udb)
		c.Set("LOG", log)
		c.Set("REDIS", "redis")
	})

	// Boostrap services
	competitionSvc := &competition.HandlerService{}
	competitionSvc.Bootstrap(router)

	marathonSvc := &marathon.HandlerService{}
	marathonSvc.Bootstrap(router)
	// --- Development Only ---
	setupQuotes(db)

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

// Utility function to populate some records into DB
// Should not be used in production.
// Production should use sql scripts to create DB with default tables and data
func setupQuotes(db *gorm.DB) {
	// check if table exists
	// if table exists, return
	if !db.HasTable(&competition.CompetitionUsers{}) {
		db.AutoMigrate(&competition.CompetitionUsers{})
	}
	if !db.HasTable(&competition.AgeGroup{}) {
		db.AutoMigrate(&competition.AgeGroup{})
	}
}

// /*func GenerateToken(c *gin.Context) {
// 	// Token Generation and Store into DB
// 	ucdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + "/wk_user_management"
// 	pgxConn, _ := pgx.Connect(context.TODO(), ucdsn)
// 	manager := manage.NewDefaultManager()
// 	manager.MapAccessGenerate(generates.NewAccessGenerate())
// 	udb := c.MustGet("UDB").(*gorm.DB)
// 	// use PostgreSQL token store with pgx.Connection adapter
// 	adapter := pgx4adapter.NewConn(pgxConn)
// 	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
// 	defer tokenStore.Close()
// 	defer pgxConn.Close(context.Background())

// 	clientStore, _ := pg.NewClientStore(adapter)

// 	manager.MapTokenStorage(tokenStore)
// 	manager.MapClientStorage(clientStore)

// 	srv := server.NewDefaultServer(manager)

// 	if c.Request.FormValue("grant_type") == "password" {
// 		/* Grant Type Password setting value to get the userid for given username and password inputs*/
// 		srv.PasswordAuthorizationHandler = func(username, password string) (string, error) {
// 			type UserDetails struct {
// 				UserId string `json:"userId"`
// 			}
// 			var user UserDetails
// 			udb.Raw("SELECT id as user_id FROM public.user WHERE lower(USER_NAME) = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR lower(USER_NAME) = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?", strings.ToLower(c.Request.FormValue("username")), c.Request.FormValue("username"), c.Request.FormValue("username"), ("+" + c.Request.FormValue("username")), ("+" + c.Request.FormValue("username")), ("+" + username), strings.TrimLeft(c.Request.FormValue("username"), "0"), strings.TrimLeft(c.Request.FormValue("username"), "0"), strings.TrimLeft(c.Request.FormValue("username"), "0")).Scan(&user)
// 			fmt.Println("authroization password handler...", username, password)
// 			return user.UserId, nil
// 		}
// 	}
// 	srv.SetAllowGetAccessRequest(true)
// 	srv.SetClientInfoHandler(server.ClientFormHandler)
// 	// if err := db.Table("user").Where("id=?", c.MustGet("userid")).Update("last_activity_at", time.Now()).Error; err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
// 	// 	return
// 	// }
// 	// srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
// 	// 	// log.Println("Internal Error:", err.Error())
// 	// 	return
// 	// })

// 	// srv.SetResponseErrorHandler(func(re *errors.Response) {
// 	// 	// log.Println("Response Error:", re.Error.Error())
// 	// })
// 	err := srv.HandleTokenRequest(c.Writer, c.Request)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is not valid", "Status": http.StatusBadRequest})
// 	}
// }
// func GenerateTokenWithUserId(c *gin.Context, userid string) {
// 	// Token Generation and Store into DB
// 	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + "/wk_user_management"
// 	pgxConn, _ := pgx.Connect(context.TODO(), dsn)
// 	manager := manage.NewDefaultManager()
// 	manager.MapAccessGenerate(generates.NewAccessGenerate())
// 	// use PostgreSQL token store with pgx.Connection adapter
// 	adapter := pgx4adapter.NewConn(pgxConn)
// 	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
// 	defer tokenStore.Close()
// 	defer pgxConn.Close(context.Background())

// 	clientStore, _ := pg.NewClientStore(adapter)
// 	manager.MapTokenStorage(tokenStore)
// 	manager.MapClientStorage(clientStore)

// 	srv := server.NewDefaultServer(manager)

// 	if c.Request.FormValue("grant_type") == "password" {
// 		/* Grant Type Password setting value to get the userid for given username and password inputs*/
// 		srv.PasswordAuthorizationHandler = func(username, password string) (string, error) {
// 			/* here changed to userd id getting from login api query */
// 			return userid, nil
// 		}
// 	}
// 	srv.SetAllowGetAccessRequest(true)
// 	srv.SetClientInfoHandler(server.ClientFormHandler)
// 	err := srv.HandleTokenRequest(c.Writer, c.Request)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is not valid", "Status": http.StatusBadRequest})
// 	}
// }
