package config

import (
	"encoding/json"
	"fmt"

	"frontend_service/common"
	l "frontend_service/logger"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	r.GET("/config", hs.GetAllConfig)
	r.POST("/setconfig", hs.SetConfig)
	r.GET("/getconfig", hs.GetConfig)
}

// GetAllConfig -  fetches all config
// GET /config
// @Summary Show a list of all country's
// @Description get list of all country's
// @Tags Config
// @Accept  json
// @Produce  json
// @Success 200 {array} ConfigurationDetails
// @Router /config [get]
func (hs *HandlerService) GetAllConfig(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var config ApplicationSetting
	if err := db.Where("name='ConfigEndpointResponseBody'").Find(&config).Error; err != nil {
		serverError := common.ServerErrorResponse("ar")
		l.JSON(c, http.StatusInternalServerError, gin.H{"Message": serverError})
		return
	}

	var formated ConfigurationDetails
	if config.Value != "" {
		data := fmt.Sprintf("%v", config.Value)
		if err := json.Unmarshal([]byte(data), &formated); err != nil {
			serverError := common.ServerErrorResponse("ar")
			l.JSON(c, http.StatusInternalServerError, gin.H{"Message": serverError})
			return
		}
	}
	l.JSON(c, http.StatusOK, formated)
	return
}

func (hs *HandlerService) SetConfig(c *gin.Context) {
	var request interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		l.JSON(c, http.StatusBadRequest, "err in binding")
		return
	}
	jsonString, _ := json.Marshal(&request)
	filename := "configqa.json"
	// storing to s3
	ioutil.WriteFile(filename, jsonString, os.ModePerm)
	common.UploadFileToS3(jsonString, filename)
	//storing to redis
	var redisRequest RedisCacheRequest
	url := os.Getenv("REDIS_CACHE_URL")
	redisRequest.Key = "configqa"
	redisRequest.Value = string(jsonString)
	_, err := common.PostCurlCall("POST", url, redisRequest)
	if err != nil {
		fmt.Println(err)
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	//storing in DB
	db := c.MustGet("FCDB").(*gorm.DB)
	var applicationSetting ApplicationSetting
	applicationSetting.Value = string(jsonString)
	if err := db.Model(&applicationSetting).Where("name='ConfigEndpointResponseBody'").Update(&applicationSetting).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, err)
		return
	}
}

func (hs *HandlerService) GetConfig(c *gin.Context) {
	var result interface{}
	url, errfroms3 := common.GetCurlCall(os.Getenv("S3_URLFORCONFIG"))
	unmarshallerr := json.Unmarshal(url, &result)
	if result == nil || errfroms3 != nil || unmarshallerr != nil {
		// going to redis if any err or no result in s3
		var unmarshallerrredisresponse error
		url := os.Getenv("REDIS_CACHE_URL") + "/" + os.Getenv("CONFIG_KEY")
		response, errcurl := common.GetCurlCall(url)
		var RedisResponse RedisCacheResponse
		errinredis := json.Unmarshal(response, &RedisResponse)
		if RedisResponse.Value != "" {
			unmarshallerrredisresponse = json.Unmarshal([]byte(RedisResponse.Value), &result)
		} else if result == nil || unmarshallerrredisresponse != nil || errinredis != nil || errcurl != nil {
			// going to DB in if data not there or any err in redis
			db := c.MustGet("FCDB").(*gorm.DB)
			var applicationSetting ApplicationSetting
			if err := db.Table("application_setting").Where("name='ConfigEndpointResponseBody'").Find(&applicationSetting).Error; err != nil {
				l.JSON(c, http.StatusInternalServerError, err)
				return
			}
			if applicationSetting.Value != "" {
				if err := json.Unmarshal([]byte(applicationSetting.Value), &result); err != nil {
					l.JSON(c, http.StatusInternalServerError, gin.H{"message": err.Error(), "Status": http.StatusInternalServerError})
				}
			}
		}
	}
	l.JSON(c, http.StatusOK, result)
}
