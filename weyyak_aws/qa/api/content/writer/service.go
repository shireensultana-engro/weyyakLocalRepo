package writer

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
	qrg.GET("/writers", hs.GetAllWriter)
	qrg.PUT("/writers", hs.CreateWriter)
}

// GetAllWriter -  fetches all Writer's
// GET /api/writers
// @Summary Show a list of all Writer's
// @Description get list of all Writer's
// @Tags Writers
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} Writer
// @Router /api/writers [get]
func (hs *HandlerService) GetAllWriter(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var Writer []Writer
	var erroresponse = common.ServerErrorResponse()
	if err := db.Debug().Find(&Writer).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": Writer})
}

// CreateWriter -  create new Writer
// PUT /api/writers
// @Summary Create New Writer
// @Description Create New Writer
// @Tags Writers
// @Accept  json
// @Produce  json
// @Param body body Writer true "Raw JSON string"
// @Success 200 {array} Writer
// @Router /api/writers [put]
func (hs *HandlerService) CreateWriter(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var writer Writer
	if err := c.ShouldBindJSON(&writer); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var writerReq Writer
	db.Debug().Table("writer").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", writer.EnglishName, writer.ArabicName).Find(&writerReq)
	if writerReq.EnglishName == writer.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_writer_englishname_already_exists", Description: "Writer with specified 'English Name' of '" + writer.EnglishName + "' already exists."}
	}
	if writer.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if writerReq.ArabicName == writer.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_writer_arabicname_already_exists", Description: "Writer with specified 'Arabic Name' of '" + writer.ArabicName + "' already exists."}
	}
	if writer.ArabicName == "" {
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
	if err := db.Debug().Create(&writer).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err, "Status": http.StatusBadRequest})
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "Writer Created Successfully.", "Status": http.StatusOK})
		return
	}
}
