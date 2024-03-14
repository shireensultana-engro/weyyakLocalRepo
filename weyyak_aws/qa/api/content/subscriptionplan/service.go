package subscriptionplan

import (
	"content/common"
	l "content/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.Use(common.ValidateToken())
	r.GET("/api/subscription/plan", hs.GetSubscriptionPlan)
}

// Get GetSubscriptionPlan -  fetches GetSubscriptionPlan
// GET /api/subscription/plan
// @Summary Show GetSubscriptionPlan
// @Description get GetSubscriptionPlan
// @Tags SubscriptionPlan
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/subscription/plan [get]
func (hs *HandlerService) GetSubscriptionPlan(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var response []SubscriptionPlan
	var errorresponse = common.ServerErrorResponse()
	if planresult := db.Debug().Table("subscription_plan").Select("id,name").Find(&response).Error; planresult != nil {
		l.JSON(c, http.StatusInternalServerError, errorresponse)
		return
	}
	l.JSON(c, http.StatusOK, response)
}
