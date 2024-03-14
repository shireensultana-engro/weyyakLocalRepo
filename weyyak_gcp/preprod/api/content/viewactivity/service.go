package viewactivity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.GET("/api/viewactivities/filters", hs.ViewActivitiesFilters)
}

// Get ViewActivitiesFilters -  fetches ViewActivitiesFilters
// GET /api/api/viewactivities/filters
// @Summary Show ViewActivitiesFilters
// @Description get ViewActivitiesFilters
// @Tags ViewActivities
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/viewactivities/filters [get]
func (hs *HandlerService) ViewActivitiesFilters(c *gin.Context) {
	newarr := [2]string{"Movie", "Episode"}
	c.JSON(http.StatusOK, gin.H{"contentTypes": newarr})
}
