package user

import (
	"net/http"

	l "frontend_config/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// r.GET("/api/ratings/filters", hs.Menu)
	// brg := r.Group("/api/users")
	// brg.Use(common.ValidateToken())
	// brg.GET("/:id/ratings", hs.GetUserRatingsDetailsWithSearchText)
	// r.Use(common.ValidateToken())
	// r.GET("api/viewactivities/:id/watchingissues", hs.GetUserWAtchingIssues)

}

// GetRatingfilters -  fetches filters list
// GET /api/ratings/filters
// @Summary Show rating filters
// @Description get rating filters
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/ratings/filters [get]
func (hs *HandlerService) Menu(c *gin.Context) {
	var newmenu []NewMenu
	db := c.MustGet("DB").(*gorm.DB)
	if err := db.Debug().Table("menu").Select("distinct(url)").Find(&newmenu).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	newarr := []string{}
	contenttype := make(map[string][]string)
	for _, element := range newmenu {
		newarr = append(newarr, strings.Title(element.Url))
	}
	contenttype["contentTypes"] = newarr
	l.JSON(c, http.StatusOK, contenttype)
}
