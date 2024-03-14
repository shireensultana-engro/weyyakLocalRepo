package contentType

import (
	"frontend_config/common"
	l "frontend_config/logger"
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
	qrg.Use(common.ValidateToken())
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
// @security Authorization
// @Success 200 {array} ContentOnetierTypes
// @Router /api/contenttypes/onetier [get]
func GetAllConentOneTier(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var contentonetier []ContentOnetierTypes
	if err := db.Debug().Find(&contentonetier).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": contentonetier})
}

// GetAllConentMultiTier -  fetches all multi tier's data
// GET /api/readonlydata
// @Summary Show a list of all multi tier's data
// @Description get list of all multi tier's data
// @Tags Content Types
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} ContentMultitierTypes
// @Router /api/contenttypes/multitier [get]
func (hs *HandlerService) GetAllConentMultiTier(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var contentmultitier []ContentMultitierTypes
	if err := db.Debug().Find(&contentmultitier).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": contentmultitier})
}
