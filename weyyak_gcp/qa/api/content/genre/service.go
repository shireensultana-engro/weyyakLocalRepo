package genre

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
	qrg := r.Group("/api/genre")
	qrg.GET("/:genre_id", hs.GetGenreById)
	srg := r.Group("/api")
	srg.Use(common.ValidateToken())
	srg.GET("/genres", hs.GetGenreList)
	srg.GET("/genres/:genre_id/subgenres", hs.GetSubGenereBasedOnGenereID)

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
	db.Debug().Where("id=?", genre_id).Find(&genre)
	l.JSON(c, http.StatusOK, genre)
}

// GetGenreList -  Get all genres
// GET /genres
// @Summary Get all genres list
// @Description get all genres list
// @Tags Genre
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} object c.JSON
// @Router /genres [get]
func (hs *HandlerService) GetGenreList(c *gin.Context) {
	var genre []GenreList
	db := c.MustGet("DB").(*gorm.DB)
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	var errorresponse = common.ServerErrorResponse()
	fields := "english_name,arabic_name,id"
	if data := db.Debug().Table("genre").Select(fields).Scan(&genre).Error; data != nil {
		l.JSON(c, http.StatusInternalServerError, errorresponse)
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": genre})
}

// GetGenreById -  Get Sub Genere Based on Genere ID
// GET /api/genres/:genre_id/subgenres
// @Summary Show a list of all Sub Genere Based on Genere ID
// @Description get list of all Sub Genere Based on Genere ID
// @Tags Genre
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param castid path string true "Genre Id"
// @Success 200 {array} Genre
// @Router /api/genres/{genre_id}/subgenres [get]
func (hs *HandlerService) GetSubGenereBasedOnGenereID(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	var subgenre []SubGenre
	subgenreId := c.Param("genre_id")
	fields := "english_name,arabic_name,id"
	if data := db.Debug().Table("subgenre").Select(fields).Where("genre_id=?", subgenreId).Find(&subgenre).Error; data != nil {
		l.JSON(c, http.StatusInternalServerError, common.ServerErrorResponse())
		return
	}
	l.JSON(c, http.StatusOK, gin.H{"data": subgenre})
}
