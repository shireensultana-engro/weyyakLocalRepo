package digitalRights

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api/content")
	qrg.GET("/digitalrights/types", hs.GetAllDigitalTypes)
	qrg.GET("/displaystatuses", hs.GetAllDisplayStatus)
}

// GetAllDigitalTypes -  fetches all digital types
// GET /api/content/digitalrights
// @Summary Show a list of all digital types
// @Description get list of all digital types
// @Tags Digital
// @Accept  json
// @Produce  json
// @Success 200 {array} DigitalRights
// @Router /api/content/digitalrights/types [get]
func (hs *HandlerService) GetAllDigitalTypes(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var digitalrights []DigitalRights
	db.Find(&digitalrights)
	c.JSON(http.StatusOK, gin.H{"data": digitalrights})
}

// GetAllDisplayStatus -  fetches all display status
// GET /api/content/digitalrights
// @Summary Show a list of all display status
// @Description get list of all display status
// @Tags Digital
// @Accept  json
// @Produce  json
// @Success 200 {array} DisplayStatus
// @Router /api/content/displaystatuses [get]
func (hs *HandlerService) GetAllDisplayStatus(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var displayStatus []DisplayStatus
	db.Find(&displayStatus)
	c.JSON(http.StatusOK, gin.H{"data": displayStatus})
}
