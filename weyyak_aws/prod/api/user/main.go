package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"user/admin"
	"user/docs"
	"user/geoblock"
	"user/register"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"

	_ "github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"user/generates"
	"user/manage"

	// "github.com/go-oauth2/oauth2/v4/manage"
	"user/common"

	"user/server"

	"github.com/dghubble/oauth1"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/thanhpk/randstr"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var (
	db     *gorm.DB
	router *gin.Engine
	log    zerolog.Logger
)

type FinalResponses struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}
type DeviceLimit struct {
	Value string `json:"value"`
}

var hmacSampleSecret []byte

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
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer db.Close()
	cdsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("CONTENT_DB_DATABASE")
	// log.Info().Msg(dsn)
	cdb, err := gorm.Open("postgres", cdsn)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer cdb.Close()
	lgdb := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("FRONTEND_DB_DATABASE")
	// log.Info().Msg(dsn)
	fdb, err := gorm.Open("postgres", lgdb)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	defer fdb.Close()
	db.DB().SetMaxIdleConns(10)
	cdb.DB().SetMaxIdleConns(10)
	fdb.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	cdb.LogMode(true)
	fdb.LogMode(true)
	db.SingularTable(true)
	cdb.SingularTable(true)
	fdb.SingularTable(true)

	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Weyyak User Management APIs"
	docs.SwaggerInfo.Description = "List of APIs for weyyak project user management api's"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// router.GET("/authorize", func(c *gin.Context) {
	// 	// fmt.Println("kkkkkkkkkk")
	// 	err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
	// 	if err != nil {
	// 		// fmt.Println("jjjjjjjj")
	// 		// fmt.Println(err)
	// 		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	// 	}
	// })
	// Setup Middleware for Database and Log
	router.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("CDB", cdb)
		c.Set("FDB", fdb)
		c.Set("LOG", log)
		c.Set("REDIS", "redis")
	})

	router.POST("/oauth2/token", Login)

	// Boostrap services
	registerSvc := &register.HandlerService{}
	registerSvc.Bootstrap(router)

	adminSvc := &admin.HandlerService{}
	adminSvc.Bootstrap(router)

	geoblockSvc := &geoblock.HandlerService{}
	geoblockSvc.Bootstrap(router)

	// --- Development Only ---
	// setupQuotes(db)

	// Start the service
	router.GET("/health", healthsvc)
	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)

}

func GenerateToken(c *gin.Context) {
	// Token Generation and Store into DB
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	pgxConn, _ := pgx.Connect(context.TODO(), dsn)
	manager := manage.NewDefaultManager()
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	db := c.MustGet("DB").(*gorm.DB)
	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()
	defer pgxConn.Close(context.Background())

	clientStore, _ := pg.NewClientStore(adapter)

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)

	if c.Request.FormValue("grant_type") == "password" {
		/* Grant Type Password setting value to get the userid for given username and password inputs*/
		srv.PasswordAuthorizationHandler = func(username, password string) (string, error) {
			type UserDetails struct {
				UserId string `json:"userId"`
			}
			var user UserDetails
			db.Raw("SELECT id as user_id FROM public.user WHERE lower(USER_NAME) = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR lower(USER_NAME) = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?", strings.ToLower(c.Request.FormValue("username")), c.Request.FormValue("username"), c.Request.FormValue("username"), ("+" + c.Request.FormValue("username")), ("+" + c.Request.FormValue("username")), ("+" + username), strings.TrimLeft(c.Request.FormValue("username"), "0"), strings.TrimLeft(c.Request.FormValue("username"), "0"), strings.TrimLeft(c.Request.FormValue("username"), "0")).Scan(&user)
			fmt.Println("authroization password handler...", username, password)
			return user.UserId, nil
		}
	}
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	// if err := db.Table("user").Where("id=?", c.MustGet("userid")).Update("last_activity_at", time.Now()).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
	// 	return
	// }
	// srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
	// 	// log.Println("Internal Error:", err.Error())
	// 	return
	// })

	// srv.SetResponseErrorHandler(func(re *errors.Response) {
	// 	// log.Println("Response Error:", re.Error.Error())
	// })
	err := srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is not valid", "Status": http.StatusBadRequest})
	}
}
func GenerateTokenWithUserId(c *gin.Context, userid string) {
	// Token Generation and Store into DB
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_DATABASE")
	pgxConn, _ := pgx.Connect(context.TODO(), dsn)
	manager := manage.NewDefaultManager()
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()
	defer pgxConn.Close(context.Background())

	clientStore, _ := pg.NewClientStore(adapter)

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)

	if c.Request.FormValue("grant_type") == "password" {
		/* Grant Type Password setting value to get the userid for given username and password inputs*/
		srv.PasswordAuthorizationHandler = func(username, password string) (string, error) {
			/* here changed to userd id getting from login api query */
			return userid, nil
		}
	}
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	err := srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is not valid", "Status": http.StatusBadRequest})
	}
}

// func trimFirstRune(s string) string {
// 	_, i := utf8.DecodeRuneInString(s)
// 	return s[i:]
// }

// Login -  Login user
// POST /oauth2/token
// @Summary User login with generate token
// @Description User login with generate token
// @Tags Login
// @Accept  multipart/form-data
// @Produce  json
// @Param   grant_type formData string true  "Enter Grant Type"
// @Param   username formData string false  "Enter Username"
// @Param   password formData string false  "Enter Password"
// @Param   DeviceId formData string false  "Enter DeviceId"
// @Param   DeviceName formData string false  "Enter DeviceName"
// @Param   DevicePlatform formData string false  "Enter DevicePlatform"
// @Param   deviceId formData string false  "Enter DeviceId twitter"
// @Param   deviceName formData string false  "Enter DeviceName twitter"
// @Param   devicePlatform formData string false  "Enter DevicePlatform twitter"
// @Param   refresh_token formData string false  "Enter Refresh Token"
// @Param   facebook_token formData string false  "Enter Facebook Token"
// @Param   device_code formData string false  "Enter Device Code"
// @Param   apple_token formData string false  "Enter Apple Token"
// @Param   token formData string false  "Enter Token"
// @Param   tokensecret formData string false  "Enter Token secret"
// @Success 200
// @Router /oauth2/token [post]
func Login(c *gin.Context) {
	//ToDo: By host name differentiating
	db := c.MustGet("DB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	type UserResponse struct {
		UserId string
	}
	var count int
	var device common.Device
	var userDevice common.UserDevice
	var platform common.Platform
	if c.Request.FormValue("grant_type") == "password" {
		var response common.UserLoginResponse
		// User validation
		lower1 := strings.ToLower(c.Request.FormValue("username"))
		lefttrim := strings.TrimLeft(c.Request.FormValue("username"), "0")
		lower2 := strings.ToLower(lefttrim)
		db.Raw("select u.id as user_id, password_hash,r.name as role, version, salt_stored from public.user u join role r on r.id = u.role_id where (USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ? OR USER_NAME = ? OR NATIONAL_NUMBER = ? OR PHONE_NUMBER = ?) and (phone_number_confirmed='true' or email_confirmed='true')", lower1, strings.ToLower(c.Request.FormValue("username")), strings.ToLower(c.Request.FormValue("username")), ("+" + c.Request.FormValue("username")), ("+" + c.Request.FormValue("username")), ("+" + c.Request.FormValue("username")), lower2, strings.TrimLeft(c.Request.FormValue("username"), "0"), strings.TrimLeft(c.Request.FormValue("username"), "0")).Scan(&response)
		if response.UserId != "" {
			validPassword := common.VerifyHashPassword(response.PasswordHash, c.Request.FormValue("password"), response.Version, response.SaltStored)
			if validPassword {
				if response.Role == "User" {
					deviceId := c.Request.FormValue("deviceId")
					DeviceID := c.Request.FormValue("DeviceId")
					if deviceId == "" && DeviceID == "" {
						type DeviceID struct {
							Code        string `json:"code,omitempty"`
							Description string `json:"description,omitempty"`
						}
						type Invalid struct {
							DeviceId *DeviceID `json:"deviceId,omitempty"`
						}
						deviceIdError := DeviceID{"NotEmptyValidator", "Device Id' should not be empty."}
						invalid := Invalid{DeviceId: &deviceIdError}
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
						return
					}
					deviceName := c.Request.FormValue("deviceName")
					DeviceName := c.Request.FormValue("DeviceName")
					if deviceName == "" && DeviceName == "" {
						type DeviceName struct {
							Code        string `json:"code,omitempty"`
							Description string `json:"description,omitempty"`
						}
						type Invalid struct {
							DeviceName *DeviceName `json:"devicePlatform,omitempty"`
						}
						deviceNameError := DeviceName{"NotEmptyValidator", "Device Name' should not be empty."}
						invalid := Invalid{DeviceName: &deviceNameError}
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
						return
					}
					devicePlatform := c.Request.FormValue("devicePlatform")
					DevicePlatform := c.Request.FormValue("DevicePlatform")
					if devicePlatform == "" && DevicePlatform == "" {
						type DevicePlatform struct {
							Code        string `json:"code,omitempty"`
							Description string `json:"description,omitempty"`
						}
						type Invalid struct {
							DevicePlatform *DevicePlatform `json:"devicePlatform,omitempty"`
						}
						devicePlatformError := DevicePlatform{"NotEmptyValidator", "Device Platform' should not be empty."}
						invalid := Invalid{DevicePlatform: &devicePlatformError}
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
						return
					}
					// checking for device limit
					if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
						if deviceId != "" && deviceName != "" && devicePlatform != "" || DeviceID != "" && DeviceName != "" && DevicePlatform != "" {
							var devicelimit DeviceLimit
							fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
							var devicecount string
							db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", response.UserId).Count(&devicecount)
							if devicecount >= devicelimit.Value {
								c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
								return
							}
						}
					}
					//storing device info
					if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
						if deviceId != "" && deviceName != "" && devicePlatform != "" {
							res := strings.ToLower(devicePlatform)
							db.Raw("SELECT platform_id FROM platform WHERE name = ?", res).Find(&platform)
							platform := platform.PlatformId
							db.Table("user_device").Where("user_id = ? and device_id = ?", response.UserId, deviceId).Count(&count)
							device = common.Device{DeviceId: deviceId, Name: deviceName, Platform: platform, CreatedAt: time.Now()}
							userDevice = common.UserDevice{UserId: response.UserId, DeviceId: deviceId, Token: ""}
							if count < 1 {
								db.Create(&device)
								db.Create(&userDevice)
							} else {
								db.Table("device").Where("device_id=(?)", deviceId).Update(common.Device{CreatedAt: time.Now()})
							}
						} else {
							res := strings.ToLower(DevicePlatform)
							db.Raw("SELECT platform_id FROM platform WHERE name = ?", res).Find(&platform)
							platform := platform.PlatformId
							db.Table("user_device").Where("user_id = ? and device_id = ?", response.UserId, DeviceID).Count(&count)
							device = common.Device{DeviceId: DeviceID, Name: DeviceName, Platform: platform, CreatedAt: time.Now()}
							userDevice = common.UserDevice{UserId: response.UserId, DeviceId: DeviceID, Token: ""}
							if count < 1 {
								db.Create(&device)
								db.Create(&userDevice)
							} else {
								db.Table("device").Where("device_id=(?)", DeviceID).Update(common.Device{CreatedAt: time.Now()})

							}
						}
					}
					GenerateTokenWithUserId(c, response.UserId)
					db.Table("user").Where("id=?", response.UserId).Update("last_activity_at", time.Now())
					//db.Table("user").Where("id=?", response.UserId).Update(map[string]interface{}{"last_activity_at": time.Now(), "intiated_deleted_at": nil})
				} else if response.Role == "Admin" {
					fmt.Println("Admin function...")
					GenerateTokenWithUserId(c, response.UserId)
					db.Table("user").Where("id=?", response.UserId).Update("last_activity_at", time.Now())
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32), "Status code": http.StatusInternalServerError})
					return
				}
			} else {
				fmt.Println("Password wrong....")
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "description": "The username or password is incorrect.", "code": "error_user_invalid_credentials", "requestId": randstr.String(32), "Status code": http.StatusBadRequest})
				return
			}
		} else {
			fmt.Println("Not a Valid user")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "description": "The username or password is incorrect.", "code": "error_user_invalid_credentials", "requestId": randstr.String(32), "Status code": http.StatusBadRequest})
			return
		}
	} else if c.Request.FormValue("grant_type") == "refresh_token" {
		// refresh token implementation
		GenerateToken(c)
	} else if c.Request.FormValue("grant_type") == "device_code" {
		// device code implementation
		deviceCode := c.Request.FormValue("device_code")
		type PairingCode struct {
			DeviceId         string    `json:"device_id"`
			DeviceCode       string    `json:"device_code"`
			UserCode         string    `json:"user_code"`
			CreatedAt        time.Time `json:"created_at"`
			ExpiresAt        time.Time `json:"expires_at"`
			SubscriptionDate time.Time `json:"subscription_date"`
		}
		var pairingCode PairingCode
		if deviceCode != "" {
			res := db.Raw(`SELECT * FROM pairing_code WHERE device_code=? and device_id=? and user_id is not null`, deviceCode, c.Request.FormValue("deviceId")).Find(&pairingCode)
			if res.RowsAffected > 0 {
				GenerateToken(c)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "not_found", "description": "Not found.", "code": "", "requestId": randstr.String(32), "Status code": http.StatusNotFound})
				return
			}
		}
	} else if c.Request.FormValue("grant_type") == "facebook_token" {
		deviceId := c.Request.FormValue("deviceId")
		DeviceID := c.Request.FormValue("DeviceId")
		if deviceId == "" && DeviceID == "" {
			type DeviceID struct {
				Code        string `json:"code,omitempty"`
				Description string `json:"description,omitempty"`
			}
			type Invalid struct {
				DeviceId *DeviceID `json:"deviceId,omitempty"`
			}
			deviceIdError := DeviceID{"NotEmptyValidator", "Device Id' should not be empty."}
			invalid := Invalid{DeviceId: &deviceIdError}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
			return
		}
		deviceName := c.Request.FormValue("deviceName")
		DeviceName := c.Request.FormValue("DeviceName")
		if deviceName == "" && DeviceName == "" {
			type DeviceName struct {
				Code        string `json:"code,omitempty"`
				Description string `json:"description,omitempty"`
			}
			type Invalid struct {
				DeviceName *DeviceName `json:"devicePlatform,omitempty"`
			}
			deviceNameError := DeviceName{"NotEmptyValidator", "Device Name' should not be empty."}
			invalid := Invalid{DeviceName: &deviceNameError}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
			return
		}
		devicePlatform := c.Request.FormValue("devicePlatform")
		DevicePlatform := c.Request.FormValue("DevicePlatform")
		if devicePlatform == "" && DevicePlatform == "" {
			type DevicePlatform struct {
				Code        string `json:"code,omitempty"`
				Description string `json:"description,omitempty"`
			}
			type Invalid struct {
				DevicePlatform *DevicePlatform `json:"devicePlatform,omitempty"`
			}
			devicePlatformError := DevicePlatform{"NotEmptyValidator", "Device Platform' should not be empty."}
			invalid := Invalid{DevicePlatform: &devicePlatformError}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "invalid": invalid})
			return
		}
		// Facebook login Implementation
		facebookToken := c.Request.FormValue("token")
		if facebookToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		type FacebookResponse struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		response := common.GetCurlCall("https://graph.facebook.com/v8.0/me?access_token=" + string(facebookToken) + "&fields=id%2Cname%2Cemail%2Cpicture&locale=en_US&method=get&pretty=0&sdk=joey&suppress_http_code=1")

		var DataResponse FacebookResponse
		json.Unmarshal(response, &DataResponse)
		if DataResponse.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error_access_token_invalid", "description": "Invalid access token.", "code": "error_access_token_invalid", "requestId": randstr.String(32)})
			return
		}
		nameSplit := strings.Split(DataResponse.Name, " ")
		var firstName, lastName string
		if len(nameSplit) > 0 {
			firstName = nameSplit[0]
			lastName = nameSplit[len(nameSplit)-1]
		}
		type UserDetails struct {
			UserName   string `json:"username"`
			LanguageId string `json:"languageId"`
			UserId     string `json:"userId"`
			Role       string `json:"role"`
		}
		var user UserDetails
		// create record in user if user doesn't exists
		if DataResponse.Email != "" {
			res := db.Raw(`SELECT user_name,language_id,u.id as user_id,r.name as role FROM public.user u left join role r on r.id=u.role_id WHERE lower(u.user_name)=? `, strings.ToLower(DataResponse.Email)).Find(&user)
			//db.Table("user").Where("id=?", user.UserId).Update("intiated_deleted_at", nil)
			if res.RowsAffected <= 0 {
				db.Exec("INSERT INTO public.user(is_back_office_user, language_id, newsletters_enabled, promotions_enabled, registered_at, email, email_confirmed, user_name, is_deleted, phone_number_confirmed, paycmsstatus, is_adult, privacy_policy, is_recommend, performance, google_analytics, firebase, app_flyer, advertising, aique, google_ads, facebook_ads, is_gdpr_accepted, clever_tap, role_id,first_name,last_name,registration_source) VALUES ('false',1,'false','false',now(),?,?,?,'false','false','true','true','true','true','true','true','true','true','true','true','true','true','true','false','91f15b92-97fd-e611-814f-0af7afba4acb',?,?,3);", DataResponse.Email, "true", DataResponse.Email, firstName, lastName)
			}
			// checking for device limit
			if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
				if deviceId != "" && deviceName != "" && devicePlatform != "" || DeviceID != "" && DeviceName != "" && DevicePlatform != "" {
					var devicelimit DeviceLimit
					fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
					var devicecount string
					db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", user.UserId).Count(&devicecount)
					if devicecount >= devicelimit.Value {
						c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
						return
					}
				}
			}
			//storing device info
			if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
				if deviceId != "" && deviceName != "" && devicePlatform != "" {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", devicePlatform).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, deviceId).Count(&count)
					device = common.Device{DeviceId: deviceId, Name: deviceName, Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: deviceId, Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", deviceId).Update(common.Device{CreatedAt: time.Now()})
					}
				} else {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", DevicePlatform).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, DeviceID).Count(&count)
					device = common.Device{DeviceId: DeviceID, Name: DeviceName, Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: DeviceID, Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", DeviceID).Update(common.Device{CreatedAt: time.Now()})

					}
				}
			}
			GenerateToken(c)
		} else if DataResponse.Email == "" {
			res := db.Raw(`SELECT user_name,language_id,u.id as user_id,r.name as role FROM public.user u left join role r on r.id=u.role_id WHERE u.user_name=? `, DataResponse.ID).Find(&user)
			if res.RowsAffected <= 0 {
				db.Exec("INSERT INTO public.user(is_back_office_user, language_id, newsletters_enabled, promotions_enabled, registered_at, email, email_confirmed, user_name, is_deleted, phone_number_confirmed, paycmsstatus, is_adult, privacy_policy, is_recommend, performance, google_analytics, firebase, app_flyer, advertising, aique, google_ads, facebook_ads, is_gdpr_accepted, clever_tap, role_id,first_name,last_name,registration_source) VALUES ('false',1,'false','false',now(),?,?,?,'false','false','true','true','true','true','true','true','true','true','true','true','true','true','true','false','91f15b92-97fd-e611-814f-0af7afba4acb',?,?,3);", DataResponse.Email, "true", DataResponse.ID, firstName, lastName)
			}
			// checking for device limit
			if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
				if deviceId != "" && deviceName != "" && devicePlatform != "" || DeviceID != "" && DeviceName != "" && DevicePlatform != "" {
					var devicelimit DeviceLimit
					fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
					var devicecount string
					db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", user.UserId).Count(&devicecount)
					if devicecount >= devicelimit.Value {
						c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
						return
					}
				}
			}
			//storing device info
			if deviceName != "web_app" && deviceName != "website" && DeviceName != "web_app" && DeviceName != "website" {
				if deviceId != "" && deviceName != "" && devicePlatform != "" {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", devicePlatform).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, deviceId).Count(&count)
					device = common.Device{DeviceId: deviceId, Name: deviceName, Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: deviceId, Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", deviceId).Update(common.Device{CreatedAt: time.Now()})
					}
				} else {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", DevicePlatform).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, DeviceID).Count(&count)
					device = common.Device{DeviceId: DeviceID, Name: DeviceName, Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: DeviceID, Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", DeviceID).Update(common.Device{CreatedAt: time.Now()})
					}
				}
			}
			GenerateToken(c)
		}
	} else if c.Request.FormValue("grant_type") == "twitter_token" {
		// Twitter login Implementation
		twitterToken := c.Request.FormValue("token")
		twitterTokenSecret := c.Request.FormValue("tokensecret")
		// type TwitterSecret struct {
		// 	OauthTokenSecret string
		// }
		// var twitterSecret TwitterSecret
		// rows := db.Table("public.twitter_request_token").Select("oauth_token_secret").Where("oauth_token=? and oauth_token_secret=?", twitterToken, twitterTokenSecret).Find(&twitterSecret)
		// if rows.RowsAffected > 0 {
		//Generate token
		config := oauth1.NewConfig("z4CAfGU3ToZlln6v440wRoA4x", "Oip8cfbwSDykJ2OdfY5iWMjbpxTGK4t1UqHSf5SVCLQuXgUGmp")
		token := oauth1.NewToken(twitterToken, twitterTokenSecret)

		// httpClient will automatically authorize http.Request's
		httpClient := config.Client(oauth1.NoContext, token)
		type TwitterResponse struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		// example Twitter API request
		path := "https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true"
		resp, _ := httpClient.Get(path)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Raw Response Body:\n%v\n", string(body))
		var dataResponse TwitterResponse
		json.Unmarshal(body, &dataResponse)
		type UserDetails struct {
			UserName   string `json:"username"`
			LanguageId string `json:"languageId"`
			UserId     string `json:"userId"`
			Role       string `json:"role"`
		}
		var user UserDetails
		// create record in user if user doesn't exists
		if dataResponse.Email != "" {
			res := db.Raw(`SELECT user_name,language_id,u.id as user_id,r.name as role FROM public.user u left join role r on r.id=u.role_id WHERE lower(u.user_name)=?`, strings.ToLower(dataResponse.Email)).Find(&user)
			//db.Table("user").Where("id=?", user.UserId).Update("intiated_deleted_at", nil)
			if res.RowsAffected <= 0 {
				db.Exec("INSERT INTO public.user(is_back_office_user, language_id, newsletters_enabled, promotions_enabled, registered_at, email, email_confirmed, user_name, is_deleted, phone_number_confirmed, paycmsstatus, is_adult, privacy_policy, is_recommend, performance, google_analytics, firebase, app_flyer, advertising, aique, google_ads, facebook_ads, is_gdpr_accepted, clever_tap, role_id,first_name,registration_source) VALUES ('false',1,'false','false',now(),?,?,?,'false','false','true','true','true','true','true','true','true','true','true','true','true','true','true','false','91f15b92-97fd-e611-814f-0af7afba4acb',?,3);", dataResponse.Email, "true", dataResponse.Email, dataResponse.Name)
			}
			// checking for device limit
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" || c.Request.FormValue("DeviceId") != "" && c.Request.FormValue("DeviceName") != "" && c.Request.FormValue("DevicePlatform") != "" {
					var devicelimit DeviceLimit
					fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
					var devicecount string
					db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", user.UserId).Count(&devicecount)
					if devicecount >= devicelimit.Value {
						c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
						return
					}
				}
			}
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("devicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("deviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("deviceId"), Name: c.Request.FormValue("deviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("deviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("deviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				} else {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("DevicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("DeviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("DeviceId"), Name: c.Request.FormValue("DeviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("DeviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("DeviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				}
			}
			GenerateToken(c)
			// }
		} else if dataResponse.Email == "" {
			res := db.Raw(`SELECT user_name,language_id,u.id as user_id,r.name as role FROM public.user u left join role r on r.id=u.role_id WHERE u.user_name=?`, dataResponse.ID).Find(&user)
			if res.RowsAffected <= 0 {
				db.Exec("INSERT INTO public.user(is_back_office_user, language_id, newsletters_enabled, promotions_enabled, registered_at, email, email_confirmed, user_name, is_deleted, phone_number_confirmed, paycmsstatus, is_adult, privacy_policy, is_recommend, performance, google_analytics, firebase, app_flyer, advertising, aique, google_ads, facebook_ads, is_gdpr_accepted, clever_tap, role_id,first_name,registration_source) VALUES ('false',1,'false','false',now(),?,?,?,'false','false','true','true','true','true','true','true','true','true','true','true','true','true','true','false','91f15b92-97fd-e611-814f-0af7afba4acb',?,3);", dataResponse.Email, "true", dataResponse.ID, dataResponse.Name)
			}
			// checking for device limit
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" || c.Request.FormValue("DeviceId") != "" && c.Request.FormValue("DeviceName") != "" && c.Request.FormValue("DevicePlatform") != "" {
					var devicelimit DeviceLimit
					fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
					var devicecount string
					db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", user.UserId).Count(&devicecount)
					if devicecount >= devicelimit.Value {
						c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
						return
					}
				}
			}
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("devicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("deviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("deviceId"), Name: c.Request.FormValue("deviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("deviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("deviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				} else {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("DevicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("DeviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("DeviceId"), Name: c.Request.FormValue("DeviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("DeviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("DeviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				}
			}
			GenerateToken(c)
		}
	} else if c.Request.FormValue("grant_type") == "apple_token" {
		// Apple login implementation
		appleToken := c.Request.FormValue("token")
		tokenString := appleToken

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")

			return hmacSampleSecret, nil
		})
		fmt.Println(token.Claims.(jwt.MapClaims))
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims["email"])
			// payload := strings.Split(appleToken, ".")
			// decodedData, _ := base64.StdEncoding.DecodeString(payload[1])
			// type AppleResponse struct {
			// 	Email         string `json:"email"`
			// 	EmailVerified bool   `json:"email_verified"`
			// }
			// var appleReponse AppleResponse
			// json.Unmarshal(decodedData, &appleReponse)
			type UserDetails struct {
				UserName   string `json:"username"`
				LanguageId string `json:"languageId"`
				UserId     string `json:"userId"`
				Role       string `json:"role"`
			}
			var user UserDetails
			// create record in user if user doesn't exists
			res := db.Raw(`SELECT user_name,language_id,u.id as user_id,r.name as role FROM public.user u left join role r on r.id=u.role_id WHERE u.user_name=?`, claims["email"]).Find(&user)
			//db.Table("user").Where("id=?", user.UserId).Update("intiated_deleted_at", nil)
			if res.RowsAffected <= 0 {
				// create record in user if user doesn't exists
				db.Exec("INSERT INTO public.user(is_back_office_user, language_id, newsletters_enabled, promotions_enabled, registration_source, registered_at, email, email_confirmed, user_name, is_deleted, phone_number_confirmed, paycmsstatus, is_adult, privacy_policy, is_recommend, performance, google_analytics, firebase, app_flyer, advertising, aique, google_ads, facebook_ads, is_gdpr_accepted, clever_tap, role_id) VALUES ('false',1,'false','false',5,now(),?,?,?,'false','false','true','true','true','true','true','true','true','true','true','true','true','true','true','false','91f15b92-97fd-e611-814f-0af7afba4acb');", claims["email"], claims["email_verified"], claims["email"])
			}
			// checking for device limit
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" || c.Request.FormValue("DeviceId") != "" && c.Request.FormValue("DeviceName") != "" && c.Request.FormValue("DevicePlatform") != "" {
					var devicelimit DeviceLimit
					fdb.Table("application_setting").Select("value").Where("name = 'UserDevicesLimit'").Find(&devicelimit)
					var devicecount string
					db.Table("user_device").Select("device_id").Joins("inner join device on user_device.device_id = device.device_id").Where("user_id = ? and token is not null and token != '' ", user.UserId).Count(&devicecount)
					if devicecount >= devicelimit.Value {
						c.JSON(http.StatusBadRequest, FinalResponses{"invalid_grant", "Maximum number of allowed devices was reached.", "error_device_limit_reached", randstr.String(32)})
						return
					}
				}
			}
			if c.Request.FormValue("deviceName") != "web_app" && c.Request.FormValue("deviceName") != "website" && c.Request.FormValue("DeviceName") != "web_app" && c.Request.FormValue("DeviceName") != "website" {
				if c.Request.FormValue("deviceId") != "" && c.Request.FormValue("deviceName") != "" && c.Request.FormValue("devicePlatform") != "" {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("devicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("deviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("deviceId"), Name: c.Request.FormValue("deviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("deviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("deviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				} else {
					db.Raw("SELECT platform_id FROM platform WHERE name = ?", c.Request.FormValue("DevicePlatform")).Find(&platform)
					platform := platform.PlatformId
					db.Table("user_device").Where("user_id = ? and device_id = ?", user.UserId, c.Request.FormValue("DeviceId")).Count(&count)
					device = common.Device{DeviceId: c.Request.FormValue("DeviceId"), Name: c.Request.FormValue("DeviceName"), Platform: platform, CreatedAt: time.Now()}
					userDevice = common.UserDevice{UserId: user.UserId, DeviceId: c.Request.FormValue("DeviceId"), Token: ""}
					if count < 1 {
						db.Create(&device)
						db.Create(&userDevice)
					} else {
						db.Table("device").Where("device_id=(?)", c.Request.FormValue("DeviceId")).Update(common.Device{CreatedAt: time.Now()})
					}
				}
			}
			//Generate token
			GenerateToken(c)

		} else {
			fmt.Println("hhhhhhhhhhh")
			fmt.Println(err)
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
	}
	// if err := db.Table("user").Where("id=?", c.MustGet("userid")).Update("last_activity_at", time.Now()).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
	// 	return
	// }
}

func healthsvc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": health()})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3003")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
