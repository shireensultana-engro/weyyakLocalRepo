package country

import (
	_ "fmt"
	"frontend_config/common"
	l "frontend_config/logger"
	"net/http"
	"os"
	_ "os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api")
	qrg.GET("/:langcode/countries", hs.GetAllCountries)
	qrg.GET("/countries/:langcode", hs.GetAllCountries)
	// qrg.POST("/countries", hs.CreateNewCountry)
	// qrg.PUT("/:countryid", hs.UpdateCountry)
	// qrg.DELETE("/:countryid", hs.DeleteCountry)

	srg := r.Group("")
	srg.Use(common.ValidateToken())
	srg.GET("/api/countries", hs.GetListOfCountries)
	srg.GET("/v1/countries/:langcode", hs.GetCountriesCallingCode)

	//Error_code check for valid Method
	srg.PUT("/v1/countries/:langcode", hs.GetCountriesCallingCode)
	srg.POST("/v1/countries/:langcode", hs.GetCountriesCallingCode)
	srg.DELETE("/v1/countries/:langcode", hs.GetCountriesCallingCode)

	srg.PUT("/api/countries", hs.GetListOfCountries)
	srg.POST("/api/countries", hs.GetListOfCountries)
	srg.DELETE("/api/countries", hs.GetListOfCountries)

}

// GetAllCountries -  fetches all countries
// GET /api/readonlydata
// @Summary Show a list of all country's
// @Description get list of all country's
// @Tags Country
// @Accept  json
// @Produce  json
// @Success 200 {array} Country
// @Router /api/{langcode}/countries [get]
func (hs *HandlerService) GetAllCountries(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	langCode := c.Param("langcode")
	var countryEn []CountryEN
	var countryAr []CountryAR
	if langCode == "en" {
		if data := db.Debug().Raw("SELECT id, english_name, calling_code FROM country").Limit(os.Getenv("DEFAULT_PAGE_SIZE")).Scan(&countryEn).Error; data != nil {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": data.Error(), "Status": http.StatusBadRequest})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": countryEn})
	} else {
		if data := db.Debug().Raw("SELECT id, arabic_name, calling_code FROM country").Limit(os.Getenv("DEFAULT_PAGE_SIZE")).Scan(&countryAr).Error; data != nil {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": data.Error(), "Status": http.StatusBadRequest})
			return
		}
		l.JSON(c, http.StatusOK, gin.H{"data": countryAr})
	}
}

// GetAllCountries -  Get all countries
// GET /api/readonlydata
// @Summary Show Get All countries
// @Description Get All countries
// @Tags Country
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} Country
// @Router /api/countries [get]
func (hs *HandlerService) GetListOfCountries(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var countries []Countries
	serverError := common.ServerErrorResponse()
	if data := db.Debug().Raw("SELECT english_name, arabic_name, id FROM country").Order("english_name ASC").Scan(&countries).Error; data != nil {
		l.JSON(c, http.StatusInternalServerError, serverError)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": countries})
}

// // CreateNewCountry -  create new Country
// // POST /api/readonlydata
// // @Summary Create New Country
// // @Description Create New Country
// // @Tags Country
// // @Accept  json
// // @Produce  json
// // @Param body body country.CountryInput true "Raw JSON string"
// // @Success 200 {array} ContentGenresInput
// // @Router /api/countries-create [post]
// func (appsvc *HandlerService) CreateNewCountry(c *gin.Context) {
// 	// Validate input
// 	var country Country
// 	if err := c.ShouldBindJSON(&country); err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
// 		return
// 	}
// 	db := c.MustGet("DB").(*gorm.DB)
// 	// country.english_name = input.english_name
// 	// country.arabic_name = input.arabic_name
// 	// country.region_id = input.region_id
// 	// country.calling_code = input.calling_code
// 	// country.alpha2code = input.alpha2code
// 	if err := db.Create(&country).Error; err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
// 		return
// 	} else {
// 		l.JSON(c, http.StatusOK, gin.H{"message": "Created Successfully.", "Status": http.StatusOK})
// 		return
// 	}
// }

// // UpdateCountry -  Update Country
// // POST /api/readonlydata/:countryid
// // @Summary Update Country details
// // @Description Update Country details by Country id
// // @Tags Pages
// // @Accept  json
// // @Produce  json
// // @Param countryid path string true "country Id"
// // @Param body body country.CountryInput true "Raw JSON string"
// // @Success 200 {array} country.PostSuccessResponse
// // @Router /api/countries/{countryid} [PUT]
// func (hs *HandlerService) UpdateCountry(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	var country Country

// 	countryid := c.Param("countryid")
// 	if err := db.Where("id=?", countryid).Find(&country).Error; err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Record does not exist.Please provide valid country Id.", "Status": http.StatusBadRequest})
// 		return
// 	}

// 	//var input CountryInput
// 	if err := c.ShouldBindJSON(&country); err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
// 		return
// 	}
// 	//country.english_name = input.english_name
// 	//country.arabic_name = input.arabic_name
// 	//country.region_id = input.region_id
// 	//country.calling_code = input.calling_code
// 	//country.alpha2code = input.alpha2code

// 	if err := db.Where("id=?", countryid).Update(&country).Error; err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
// 		return
// 	} else {
// 		//var response common.PostSuccessResponse
// 		//response.Id = countryid
// 		l.JSON(c, http.StatusOK, gin.H{"message": "Record Updated Successfully.", "Status": http.StatusOK})
// 		return
// 	}

// }

// // DeleteCountry -  Delete Country
// // POST /api/readonlydata/:countryid
// // @Summary Delete Country details
// // @Description Delete Country details by Country id
// // @Tags Pages
// // @Accept  json
// // @Produce  json
// // @Param countryid path string true "country Id"
// // @Param body body country.CountryInput true "Raw JSON string"
// // @Success 200 {array} country.PostSuccessResponse
// // @Router /api/countries/{countryid} [post]
// func (hs *HandlerService) DeleteCountry(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	var country Country

// 	countryid := c.Param("countryid")
// 	if err := db.Where("id=?", countryid).Find(&country).Error; err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Record not exist.Please provide valid country Id.", "Status": http.StatusBadRequest})
// 		return
// 	}

// 	if err := db.Where("id=?", countryid).Delete(&country).Error; err != nil {
// 		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
// 		return
// 	} else {
// 		//var response common.PostSuccessResponse
// 		//response.Id = countryid
// 		//l.JSON(c, http.StatusOK, gin.H{"data": response})
// 		l.JSON(c, http.StatusOK, gin.H{"message": "Record Deleted Successfully.", "Status": http.StatusOK})
// 		return
// 		return
// 	}
// }

// GetCountriesCallingCode -  Fetches all countries Mobile calling Code
// GET /countries/{langcode}
// @Summary Show all countries Mobile calling Code
// @Description get list of countries Mobile calling Code
// @Tags Country
// @Accept  json
// @Produce  json
// @Param langcode path string true "langcode"
// @Success 200 {array} object c.JSON
// @Router /countries/{langcode} [get]
func (hs *HandlerService) GetCountriesCallingCode(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		l.JSON(c, http.StatusMethodNotAllowed, gin.H{"message": "The requested resource does not support http method '" + c.Request.Method + "'."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	langCode := c.Param("langcode")
	var countriescallingcode []CountriesCallingCode
	if langCode == "en" {
		fields := "english_name as name, calling_code as code"
		if data := db.Debug().Table("country").Select(fields).Scan(&countriescallingcode).Error; data != nil {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": data.Error(), "Status": http.StatusBadRequest})
			return
		}
	} else {
		fields := "arabic_name as name, calling_code as code"
		if data := db.Debug().Table("country").Select(fields).Scan(&countriescallingcode).Error; data != nil {
			l.JSON(c, http.StatusBadRequest, gin.H{"message": data.Error(), "Status": http.StatusBadRequest})
			return
		}
	}
	l.JSON(c, http.StatusOK, countriescallingcode)
}
