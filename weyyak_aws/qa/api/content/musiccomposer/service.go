package musiccomposer

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
	qrg.GET("/musiccomposers", hs.GetAllMusicComposer)
	qrg.PUT("/musiccomposers", hs.CreateMusicComposer)
}

// GetAllMusicComposer -  fetches all MusicComposer's
// GET /api/musiccomposers
// @Summary Show a list of all MusicComposer's
// @Description get list of all MusicComposer's
// @Tags MusicComposers
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} MusicComposer
// @Router /api/musiccomposers [get]
func (hs *HandlerService) GetAllMusicComposer(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var MusicComposer []MusicComposer
	var errorresponse = common.ServerErrorResponse()
	if err := db.Debug().Find(&MusicComposer).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, errorresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": MusicComposer})
}

// CreateMusicComposer -  create new MusicComposer
// PUT /api/musiccomposers
// @Summary Create New MusicComposer
// @Description Create New MusicComposer
// @Tags MusicComposers
// @Accept  json
// @Produce  json
// @Param body body MusicComposer true "Raw JSON string"
// @Success 200 {array} MusicComposer
// @Router /api/musiccomposers [put]
func (hs *HandlerService) CreateMusicComposer(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var musicComposer MusicComposer
	if err := c.ShouldBindJSON(&musicComposer); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var MusicComposerReq MusicComposer
	db.Debug().Table("music_composer").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", musicComposer.EnglishName, musicComposer.ArabicName).Find(&MusicComposerReq)
	if MusicComposerReq.EnglishName == musicComposer.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_musiccomposer_englishname_already_exists", Description: "Music composer with specified 'English Name' of '" + musicComposer.EnglishName + "' already exists."}
	}
	if musicComposer.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if MusicComposerReq.ArabicName == musicComposer.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_musiccomposer_arabicname_already_exists", Description: "Music composer with specified 'Arabic Name' of '" + musicComposer.ArabicName + "' already exists."}
	}
	if musicComposer.ArabicName == "" {
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
	if err := db.Debug().Create(&musicComposer).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err, "Status": http.StatusBadRequest})
		return
	} else {

		l.JSON(c, http.StatusOK, gin.H{"message": "MusicComposer Created Successfully.", "Status": http.StatusOK})
		return
	}
}
