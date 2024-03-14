package language

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
	langqrg := r.Group("/api/languages")
	langqrg.Use(common.ValidateToken())
	langqrg.GET("/", hs.GetAvailableLanguages)
	langqrg.GET("", hs.GetAvailableLanguages)
	langqrg.GET("/origintypes", hs.LanguageOriginType)
	langqrg.GET("/ar/dialects", hs.GetAllLanuageDialects)
	langqrg.GET("/dubbing", hs.GetAllLanuageDubbing)
	langqrg.GET("/subtitling", hs.GetAllLanuageSubtitles)
	langqrg.GET("/original", hs.GetAllLanuageOriginal)
}

// LanguageOriginType -  fetches all Languages origin types
// GET /languages/origintypes
// @Summary Show a list of all Languages  origin types
// @Description get list of all Languages  origin types
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} LanguageOriginType
// @Router /languages/origintypes [get]
func (hs *HandlerService) LanguageOriginType(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var languageOriginTypes []LanguageOriginType
	if err := db.Debug().Find(&languageOriginTypes).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": languageOriginTypes})
}

// GetAllLanuageDialects -  fetches all Languages dialects
// GET /languages/ar/dialects
// @Summary Show a list of all Languages dialects
// @Description get list of all Languages dialects
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} LanguageDialect
// @Router /languages/ar/dialects [get]
func (hs *HandlerService) GetAllLanuageDialects(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var languageDialects []LanguageDialect
	if err := db.Debug().Find(&languageDialects).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": languageDialects})
}

// GetAllLanuageDubbing -  fetches all Languages dubbing
// GET /languages/dubbing
// @Summary Show a list of all Languages dubbing
// @Description get list of all Languages dubbing
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} LanguageDubbing
// @Router /languages/dubbing [get]
func (hs *HandlerService) GetAllLanuageDubbing(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var languageDubbing []LanguageDubbing
	if err := db.Debug().Find(&languageDubbing).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": languageDubbing})
}

// GetAllLanuageSubtitles -  fetches all Languages subtitle
// GET /languages/subtitling
// @Summary Show a list of all Languages subtitle
// @Description get list of all Languages subtitle
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} LanguageSubtitles
// @Router /languages/subtitling [get]
func (hs *HandlerService) GetAllLanuageSubtitles(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var languageSubtitle []LanguageSubtitles
	if err := db.Debug().Find(&languageSubtitle).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": languageSubtitle})
}

// GetAllLanuageOriginal -  fetches all Languages original
// GET /languages/original
// @Summary Show a list of all Languages original
// @Description get list of all Languages original
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} Language
// @Router /languages/original [get]
func (hs *HandlerService) GetAllLanuageOriginal(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var languageOriginal []Language
	if err := db.Debug().Find(&languageOriginal).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": languageOriginal})
}

// GetAvailableLanguages -  fetches all available languages
// GET /languages/
// @Summary Show a list of all Languages available
// @Description get list of all Languages available
// @Tags Language
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} LanguageAvailable
// @Router /languages/ [get]
func (hs *HandlerService) GetAvailableLanguages(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var language []LanguageAvailable
	if data := db.Debug().Raw("SELECT id, english_name as name FROM language order by id asc").Limit(2).Scan(&language).Error; data != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	// db.Find(&language)
	l.JSON(c, http.StatusOK, gin.H{"data": language})
}
