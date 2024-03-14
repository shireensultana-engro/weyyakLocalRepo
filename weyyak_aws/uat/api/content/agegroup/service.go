package agegroup

import (
	"content/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.Use(common.ValidateToken())
	r.GET("/api/agegroups", hs.AgeGroup)
}

// Getagegroups -  fetches agegroups
// GET /api/agegroups
// @Summary Show agegroups
// @Description get agegroups
// @Tags Agegroup
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} AgeGroup
// @Router /api/agegroups [get]
func (hs *HandlerService) AgeGroup(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var agegroup []AgeGroup
	fields := "english_name,arabic_name,id"
	if err := db.Debug().Table("age_ratings").Select(fields).Find(&agegroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": agegroup})
}
