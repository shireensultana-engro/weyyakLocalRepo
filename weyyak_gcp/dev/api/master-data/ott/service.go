package ott

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.GET("/v1/:lang/contenttype", hs.GetTopMenuList)
	r.GET("/config", hs.GetConfigDetails)
}

// GetTopMenuList -  Get topmenu list
// GET /v1/:lang/contenttype
// @Summary Get topmenu details list
// @Description Get topmenu details list
// @Tags OTT
// @Accept  json
// @Produce  json
// @Param lang path string true "Language Code"
// @Param device query string true "Device"
// @Success 200 {array} object c.JSON
// @Router /v1/{lang}/contenttype [get]
func (hs *HandlerService) GetTopMenuList(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var menu []TopMenuDetails
	if c.Request.URL.Query()["device"] == nil || c.Request.URL.Query()["device"][0] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide device name.",
			"Status": http.StatusBadRequest})
		return
	}
	deviceName := strings.ToLower(c.Request.URL.Query()["device"][0])
	language := strings.ToLower(c.Param("lang"))
	fields := "device,menu_type as menuType, slider_key as sliderKey,url,menu.order"
	if language == "en" {
		fields += ", menu_english_name as title"
	} else {
		fields += ", menu_arabic_name as title"
	}
	if err := db.Table("menu").Select(fields).
		Where("device=? and is_published=?", deviceName, true).
		Order("menu.order").Find(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(),
			"status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, menu)
	return
}

// GetConfigDetails -  Get application config details
// GET /config
// @Summary Get application config details
// @Description Get application config details
// @Tags OTT
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /config [get]
func (hs *HandlerService) GetConfigDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var app ApplicationSetting
	if err := db.Where("name='ConfigEndpointResponseBody'").Find(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(),
			"status": http.StatusInternalServerError})
	}
	var raw ConfigurationDetails

	if app.Value != "" {
		data := fmt.Sprintf("%v", app.Value)
		if err := json.Unmarshal([]byte(data), &raw); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(),
				"Status": http.StatusInternalServerError})
			return

		}
	}
	c.JSON(http.StatusOK, raw)
	return
}
