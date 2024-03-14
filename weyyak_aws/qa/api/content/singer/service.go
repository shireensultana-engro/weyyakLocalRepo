package singer

import (
	"content/common"
	l "content/logger"
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
	qrg.GET("/singers", hs.GetAllSinger)
	qrg.PUT("/singers", hs.CreateSinger)
}

// GetAllSinger -  fetches all Singer's
// GET /api/singers
// @Summary Show a list of all Singer's
// @Description get list of all Singer's
// @Tags Singers
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} Singer
// @Router /api/singers [get]
func (hs *HandlerService) GetAllSinger(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var Singer []Singer
	if err := db.Debug().Find(&Singer).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": Singer})
}

// CreateSinger -  create new Singer
// PUT /api/singers
// @Summary Create New Singer
// @Description Create New Singer
// @Tags Singers
// @Accept  json
// @Produce  json
// @Param body body Singer true "Raw JSON string"
// @Success 200 {array} Singer
// @Router /api/singers [put]
func (hs *HandlerService) CreateSinger(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var singer Singer
	if err := c.ShouldBindJSON(&singer); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var singerreq Singer
	db.Debug().Table("singer").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", singer.EnglishName, singer.ArabicName).Find(&singerreq)
	if singerreq.EnglishName == singer.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_singer_englishname_already_exists", Description: "Singer with specified 'English Name' of '" + singer.EnglishName + "' already exists."}
	}
	if singer.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if singerreq.ArabicName == singer.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_singer_arabicname_already_exists", Description: "Singer with specified 'Arabic Name' of '" + singer.ArabicName + "' already exists."}
	}
	if singer.ArabicName == "" {
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
		l.JSON(c, http.StatusBadRequest, inputErrors)
		return
	}
	/*End Of Input Validations*/
	if err := db.Debug().Create(&singer).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err, "Status": http.StatusBadRequest})
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "Singer Created Successfully.", "Status": http.StatusOK})
		return
	}
}
