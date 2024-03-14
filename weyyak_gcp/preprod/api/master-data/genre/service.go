package genre

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
	qrg := r.Group("/api/genre")
	qrg.GET("/:genre_id", hs.GetGenreById)
}

// GetGenreById -  fetches all genre
// GET /api/genre
// @Summary Show a list of all genre
// @Description get list of all genre
// @Tags Genre
// @Accept  json
// @Produce  json
// @Param castid path string true "Genre Id"
// @Success 200 {array} Genre
// @Router /api/genre/{genre_id} [get]
func (hs *HandlerService) GetGenreById(c *gin.Context) {
	var genre Genre
	genre_id := c.Param("genre_id")
	db := c.MustGet("DB").(*gorm.DB)
	db.Where("id=?", genre_id).Find(&genre)
	c.JSON(http.StatusOK, genre)
}
