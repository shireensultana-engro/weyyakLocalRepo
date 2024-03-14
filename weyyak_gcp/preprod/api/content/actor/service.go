package actor

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
	qrg.GET("/Actors", hs.GetAllActor)
	qrg.PUT("/Actors", hs.CreateActor)
	qrg.GET("/actors", hs.GetAllActor)
	qrg.PUT("/actors", hs.CreateActor)
}

// GetAllActor -  fetches all actor's
// GET /api/Actors
// @Summary Show a list of all actor's
// @Description get list of all actor's
// @Tags Actors
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} Actor
// @Router /api/Actors [get]
func (hs *HandlerService) GetAllActor(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var erroresponse = common.ServerErrorResponse()
	var actor []Actor
	if err := db.Debug().Find(&actor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": actor})
}

// CreateActor -  create new actor
// PUT /api/actors
// @Summary Create New Actor
// @Description Create New Actor
// @Tags Actors
// @Accept  json
// @security Authorization
// @Produce  json
// @Param body body Actor true "Raw JSON string"
// @Success 200 {array} Actor
// @Router /api/actors [put]
func (hs *HandlerService) CreateActor(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var actor Actor
	var erroresponse = common.ServerErrorResponse()
	if err := c.ShouldBindJSON(&actor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	/*Input Validations*/
	var errorFlag bool
	errorFlag = false
	var englishName common.EnglishName
	var actorRequest Actor
	db.Debug().Table("actor").Select("english_name,arabic_name").Where("english_name=? or arabic_name=?", actor.EnglishName, actor.ArabicName).Find(&actorRequest)
	if actorRequest.EnglishName == actor.EnglishName {
		errorFlag = true
		englishName = common.EnglishName{Code: "error_actor_englishname_already_exists", Description: "Actor with specified 'English Name' of '" + actor.EnglishName + "' already exists."}
	}
	if actor.EnglishName == "" {
		errorFlag = true
		englishName = common.EnglishName{Code: "NotEmptyValidator", Description: "'English Name' should not be empty."}
	}
	var arabicName common.ArabicName
	if actorRequest.ArabicName == actor.ArabicName {
		errorFlag = true
		arabicName = common.ArabicName{Code: "error_actor_arabicname_already_exists", Description: "Actor with specified 'Arabic Name' of '" + actor.ArabicName + "' already exists."}
	}
	if actor.ArabicName == "" {
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
	if err := db.Debug().Create(&actor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Actor Created Successfully."})
		return
	}
}
