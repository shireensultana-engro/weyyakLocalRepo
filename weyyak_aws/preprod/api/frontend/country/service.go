package country

import (
	_ "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/thanhpk/randstr"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/v1")
	qrg.GET("/:lang/countries", hs.GetAllCountries)
	qrg.GET("/countries/:lang", hs.GetAllCountries)
}

// GetAllCountries -  Fetches all countries Mobile calling Code
// GET /countries/{langcode}
// @Summary Show all countries Mobile calling Code
// @Description get list of countries Mobile calling Code
// @Tags Country
// @Accept  json
// @Produce  json
// @Param langcode path string true "langcode"
// @Success 200 {array} object c.JSON
// @Router /countries/{langcode} [get]
func (hs *HandlerService) GetAllCountries(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	langCode := c.Param("lang")
	var countries []CountryDetails
	var fields string
	if langCode == "en" {
		fields = "english_name as name, calling_code as code"
	} else {
		fields = "arabic_name as name, calling_code as code"
	}
	if data := db.Table("country").Select(fields).Scan(&countries).Error; data != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "error": "server_error", "description": "حدث خطأ ما", "code": "error_server_error", "requestId": randstr.String(32) })
		return
	}
	c.JSON(http.StatusOK, countries)
}
