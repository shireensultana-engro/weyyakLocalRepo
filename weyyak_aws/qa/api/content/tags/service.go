package tags

import (
	"content/common"
	l "content/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api/")
	qrg.Use(common.ValidateToken())
	qrg.GET("/textualdatatags", hs.GetAllTags)
	qrg.PUT("/textualdatatags", hs.CreateTags)
}

// GetAllTags -  fetches all Tags's
// GET /api/textualdatatags
// @Summary Show a list of all Tags's
// @Description get list of all Tags's
// @Tags Tags
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} TextualDataTag
// @Router /api/textualdatatags [get]
func (hs *HandlerService) GetAllTags(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	var Tags []TextualDataTag
	if err := db.Debug().Select("distinct on (name) name,id").Order("name Asc").Find(&Tags).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, errorresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": Tags})
}

// CreateTags -  create new TextualDataTag
// POST /api/textualdatatags
// @Summary Create New TextualDataTag
// @Description Create New TextualDataTag
// @Tags Tags
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param body body TextualDataTag true "Raw JSON string"
// @Success 200 {array} TextualDataTag
// @Router /api/textualdatatags [post]
func (hs *HandlerService) CreateTags(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var Tags TextualDataTag
	if err := c.ShouldBindJSON(&Tags); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	if Tags.Name == "" {
		name := Name{"NotEmptyValidator", "'Name' should not be empty."}
		l.JSON(c, http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "name": name})
		return
	}
	if err := db.Debug().Where("name=?", Tags.Name).Find(&Tags).RowsAffected; err != 0 {
		name := Name{"error_slider_name_not_unique", "Slider with specified 'Name' of '" + Tags.Name + "' already exists."}
		l.JSON(c, http.StatusBadRequest, gin.H{"error": "invalid_request", "description": "Validation failed.", "code": "error_validation_failed", "requestId": randstr.String(32), "name": name})
		return
	} else {
		if err := db.Debug().Create(&Tags).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		} else {
			l.JSON(c, http.StatusOK, gin.H{"message": "Tags Created Successfully.", "Status": http.StatusOK})
			return
		}
	}
}
