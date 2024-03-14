package ageRatings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	qrg := r.Group("/api/age-ratings")
	qrg.GET("/:rating_id", hs.GetAgeRatingById)
}

// GetAgeRatingById -  fetches all age ratings
// GET /api/age-ratings
// @Summary Show a list of all Age Ratings
// @Description get list of all Age Ratings
// @Tags Age Ratings
// @Accept  json
// @Produce  json
// @Param castid path string true "AgeRating Id"
// @Success 200 {array} AgeRatings
// @Router /api/age-ratings/{rating_id} [get]
func (hs *HandlerService) GetAgeRatingById(c *gin.Context) {
	var AgeRatings AgeRatings
	rating_id := c.Param("rating_id")
	db := c.MustGet("DB").(*gorm.DB)
	db.Where("id=?", rating_id).Find(&AgeRatings)
	c.JSON(http.StatusOK, AgeRatings)
}
