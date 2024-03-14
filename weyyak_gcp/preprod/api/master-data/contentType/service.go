package contentType

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
	qrg := r.Group("/api/contenttypes")
	qrg.GET("/onetier", GetAllConentOneTier)
	qrg.GET("/multitier", hs.GetAllConentMultiTier)
}

// GetAllConentOneTier -  fetches all one tier's data
// GET /api/readonlydata
// @Summary Show a list of all one tier's data
// @Description get list of all one tier's data
// @Tags Content Types
// @Accept  json
// @Produce  json
// @Success 200 {array} ContentOneTier
// @Router /api/contenttypes/onetier [get]
func GetAllConentOneTier(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var contentonetier []ContentOneTier
	db.Find(&contentonetier)
	c.JSON(http.StatusOK, gin.H{"data": contentonetier})
}

// GetAllConentMultiTier -  fetches all multi tier's data
// GET /api/readonlydata
// @Summary Show a list of all multi tier's data
// @Description get list of all multi tier's data
// @Tags Content Types
// @Accept  json
// @Produce  json
// @Success 200 {array} ContentMultiTier
// @Router /api/contenttypes/multitier [get]
func (hs *HandlerService) GetAllConentMultiTier(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var contentmultitier []ContentMultiTier
	db.Find(&contentmultitier)
	c.JSON(http.StatusOK, gin.H{"data": contentmultitier})
}
