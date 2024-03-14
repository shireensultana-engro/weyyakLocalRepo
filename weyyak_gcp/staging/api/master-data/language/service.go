package language

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
	langqrg := r.Group("/api/languages")
	langqrg.GET("/origintypes", hs.GetAllOriginTypes)
	langqrg.GET("/ar/dialects", hs.GetAllLanuageDialects)
	langqrg.GET("/dubbing", hs.GetAllLanuageDubbing)
	langqrg.GET("/subtitling", hs.GetAllLanuageSubtitles)
	langqrg.GET("/original", hs.GetAllLanuageOriginal)
}

// GetAllOriginTypes -  fetches all Languages origin types
// GET /api/languages/origintypes
// @Summary Show a list of all Languages  origin types
// @Description get list of all Languages  origin types
// @Tags Language
// @Accept  json
// @Produce  json
// @Success 200 {array} LanguageOriginTypes
// @Router /api/languages/origintypes [get]
func (hs *HandlerService) GetAllOriginTypes(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var languageOriginTypes []LanguageOriginTypes
	db.Find(&languageOriginTypes)
	c.JSON(http.StatusOK, gin.H{"data": languageOriginTypes})
}

// GetAllLanuageDialects -  fetches all Languages dialects
// GET /api/languages/ar/dialects
// @Summary Show a list of all Languages dialects
// @Description get list of all Languages dialects
// @Tags Language
// @Accept  json
// @Produce  json
// @Success 200 {array} LanguageDialects
// @Router /api/languages/ar/dialects [get]
func (hs *HandlerService) GetAllLanuageDialects(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var languageDialects []LanguageDialects
	db.Find(&languageDialects)
	c.JSON(http.StatusOK, gin.H{"data": languageDialects})
}

// GetAllLanuageDubbing -  fetches all Languages dubbing
// GET /api/languages/dubbing
// @Summary Show a list of all Languages dubbing
// @Description get list of all Languages dubbing
// @Tags Language
// @Accept  json
// @Produce  json
// @Success 200 {array} LanguageDubbing
// @Router /api/languages/dubbing [get]
func (hs *HandlerService) GetAllLanuageDubbing(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var languageDubbing []LanguageDubbing
	db.Find(&languageDubbing)
	c.JSON(http.StatusOK, gin.H{"data": languageDubbing})
}

// GetAllLanuageSubtitles -  fetches all Languages subtitle
// GET /api/languages/subtitling
// @Summary Show a list of all Languages subtitle
// @Description get list of all Languages subtitle
// @Tags Language
// @Accept  json
// @Produce  json
// @Success 200 {array} LanguageSubtitle
// @Router /api/languages/subtitling [get]
func (hs *HandlerService) GetAllLanuageSubtitles(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var languageSubtitle []LanguageSubtitle
	db.Find(&languageSubtitle)
	c.JSON(http.StatusOK, gin.H{"data": languageSubtitle})
}

// GetAllLanuageOriginal -  fetches all Languages original
// GET /api/languages/original
// @Summary Show a list of all Languages original
// @Description get list of all Languages original
// @Tags Language
// @Accept  json
// @Produce  json
// @Success 200 {array} LanguageOriginal
// @Router /api/languages/original [get]
func (hs *HandlerService) GetAllLanuageOriginal(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var languageOriginal []LanguageOriginal
	db.Find(&languageOriginal)
	c.JSON(http.StatusOK, gin.H{"data": languageOriginal})
}
