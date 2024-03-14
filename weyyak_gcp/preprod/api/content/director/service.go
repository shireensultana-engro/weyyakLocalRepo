package director

import (
	"content/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api")
	qrg.Use(common.ValidateToken())
	qrg.GET("/directors", hs.GetAllDirector)
	qrg.PUT("/directors", hs.CreateDirector)
}

// GetAllDirector -  fetches all Director's
// GET /api/directors
// @Summary Show a list of all Director's
// @Description get list of all Director's
// @Tags Directors
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} Director
// @Router /api/directors [get]
func (hs *HandlerService) GetAllDirector(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var Director []Director
	if err := db.Debug().Find(&Director).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": Director})
}

// CreateDirector -  create new Director
// PUT /api/directors
// @Summary Create New Director
// @Description Create New Director
// @Tags Directors
// @Accept  json
// @Produce  json
// @Param body body Director true "Raw JSON string"
// @Success 200 {array} Director
// @Router /api/directors [put]
func (hs *HandlerService) CreateDirector(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var director Director
	if err := c.ShouldBindJSON(&director); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var directorReq Director
	db.Debug().Table("director").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", director.EnglishName, director.ArabicName).Find(&directorReq)
	if directorReq.EnglishName == director.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_director_englishname_already_exists", Description: "Director with specified 'English Name' of '" + director.EnglishName + "' already exists."}
	}
	if director.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if directorReq.ArabicName == director.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_director_arabicname_already_exists", Description: "Director with specified 'Arabic Name' of '" + director.ArabicName + "' already exists."}
	}
	if director.ArabicName == "" {
		errorFlag = true
		arabicName = common.ArabicName{Code: "NotEmptyValidator", Description: "'Arabic Name' should not be empty."}
	}
	var invalid common.Invalid
	if arabicName.Code != "" {
		invalid.ArabicName = arabicName
	}
	if englishName.Code != "" {
		invalid.EnglishName = englishName
	}
	if errorFlag {
		inputErrors := common.InpurError{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: common.GenerateRandomString(32), Invalid: invalid}
		c.JSON(http.StatusBadRequest, inputErrors)
		return
	}
	/*End Of Input Validations*/
	if err := db.Debug().Create(&director).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.ServerErrorResponse())
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Director Created Successfully.", "Status": http.StatusOK})
		return
	}
}
