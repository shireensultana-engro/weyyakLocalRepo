package config

import (
	"encoding/json"
	"fmt"
	"frontend_service/common"
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
// @Success 200 {array} object ConfigurationDetails
// @Router /config [get]
func (hs *HandlerService) GetAllConfig(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var config ApplicationSetting
	if err := db.Where("name='ConfigEndpointResponseBody'").Find(&config).Error; err != nil {
		serverError := common.ServerErrorResponse("ar")
		c.JSON(http.StatusInternalServerError, gin.H{"Message": serverError})
		return
	}

	var formated ConfigurationDetails
	if config.Value != "" {
		data := fmt.Sprintf("%v", config.Value)
		if err := json.Unmarshal([]byte(data), &formated); err != nil {
			serverError := common.ServerErrorResponse("ar")
			c.JSON(http.StatusInternalServerError, gin.H{"Message": serverError})
			return
		}
	}
	c.JSON(http.StatusOK, formated)
	return
}

func (hs *HandlerService) SetConfig(c *gin.Context) {
	var request interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "err in binding")
		return
	}
	jsonString, _ := json.Marshal(&request)
	filename := "configuat.json"
	ioutil.WriteFile(filename, jsonString, os.ModePerm)
	common.UploadFileToS3(jsonString, filename)
}

func (hs *HandlerService) GetConfig(c *gin.Context) {
	var result interface{}
	// url changed to env variable
	// url, err := common.GetCurlCall("https://z5content-uat.s3.ap-south-1.amazonaws.com/configuat.json")
	url, err := common.GetCurlCall(os.Getenv("S3_URLFORCONFIG"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	err = json.Unmarshal(url, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, result)
}
