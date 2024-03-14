package songwriter

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
	qrg.GET("/songwriters", hs.GetAllSongWriter)
	qrg.PUT("/songwriters", hs.CreateSongWriter)
}

// GetAllSongWriter -  fetches all SongWriter's
// GET /api/songwriters
// @Summary Show a list of all SongWriter's
// @Description get list of all SongWriter's
// @Tags SongWriters
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} SongWriter
// @Router /api/songwriters [get]
func (hs *HandlerService) GetAllSongWriter(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	var erroresponse = common.ServerErrorResponse()
	var SongWriter []SongWriter
	if err := db.Debug().Find(&SongWriter).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, erroresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": SongWriter})
}

// CreateSongWriter -  create new SongWriter
// PUT /api/songwriters
// @Summary Create New SongWriter
// @Description Create New SongWriter
// @Tags SongWriters
// @Accept  json
// @Produce  json
// @Param body body SongWriter true "Raw JSON string"
// @Success 200 {array} SongWriter
// @Router /api/songwriters [put]
func (hs *HandlerService) CreateSongWriter(c *gin.Context) {
	/*Authorization*/
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	db := c.MustGet("DB").(*gorm.DB)
	var songWriter SongWriter
	if err := c.ShouldBindJSON(&songWriter); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var songWriterreq SongWriter
	db.Debug().Table("song_writer").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", songWriter.EnglishName, songWriter.ArabicName).Find(&songWriterreq)
	if songWriterreq.EnglishName == songWriter.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_songwriter_englishname_already_exists", Description: "Song writer with specified 'English Name' of '" + songWriter.EnglishName + "' already exists."}
	}
	if songWriter.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if songWriterreq.ArabicName == songWriter.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_songwriter_arabicname_already_exists", Description: "Song writer with specified 'Arabic Name' of '" + songWriter.ArabicName + "' already exists."}
	}
	if songWriter.ArabicName == "" {
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
	if err := db.Debug().Create(&songWriter).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err, "Status": http.StatusBadRequest})
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "SongWriter Created Successfully.", "Status": http.StatusOK})
		return
	}
}
