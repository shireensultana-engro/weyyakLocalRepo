package anchor

import (
	// "content/common"

	l "content/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

/* All the services should be protected by auth token */
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	/* Setup Routes */
	qrg := r.Group("/api")
	qrg.POST("/anchor", hs.CreateAnchor)
	qrg.GET("/anchor/:id", hs.GetAnchorById)
	qrg.GET("/anchor", hs.GetAllAnchors)
	qrg.POST("/anchor/:id", hs.UpdateAnchorById)
	qrg.DELETE("/anchor/:id", hs.DeleteAnchorById)

}

/*create Anchor */
func (hs *HandlerService) CreateAnchor(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var createAnchor CreateAnchor
	if err := c.ShouldBindJSON(&createAnchor); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var photo bool
	if createAnchor.Photo != "" {
		photo = true
	} else {
		photo = false
	}
	anchor := Anchor{Name: createAnchor.Name, AbouTheAnchor: createAnchor.AbouTheAnchor, Status: createAnchor.Status, Photo: photo}
	if err := db.Debug().Create(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, "internal server error")
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "Anchor created Successfully."})
		return
	}
}

/* Get Anchor by id */
func (hs *HandlerService) GetAnchorById(c *gin.Context) {
	var anchor []Anchor
	id := c.Param("id")
	db := c.MustGet("DB").(*gorm.DB)
	if err := db.Debug().Where("id=? and has_deleted is False", id).Find(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	l.JSON(c, http.StatusOK, anchor)
}

// {
// 	"id": "a1515a21-fbd2-4a16-abf9-5a13afe526f2",
// 	"name": "adhi",
// 	“Description: "test about the achor dteilas",
// 	"status": true,
// 	"photo_url”: “”https://s3.,
// 	“Email”:””,
// 	“shows”:[
// 	{ “show_name”:””,
// 	“timing”:””
// 	}

// {
// 	"show_name": "Test"
// 	"timing": "12:00:00 00:00"
// }

/*Get all anchors*/
func (hs *HandlerService) GetAllAnchors(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var anchor, finalanchor []Anchor

	if err := db.Table("anchor").Where("has_deleted is False").Find(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	for _, anchorVal := range anchor {
		var anchorShows []AnchorShows
		db.Table("anchor_shows").Where("anchor_id=?", anchorVal.Id).Find(&anchorShows)
		anchorVal.Shows = anchorShows
		finalanchor = append(finalanchor, anchorVal)
	}

	l.JSON(c, http.StatusOK, gin.H{"data": finalanchor})
}

/* Update anchor by id*/
func (hs *HandlerService) UpdateAnchorById(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	anchorId := c.Param("id")
	var createAnchor CreateAnchor
	if err := c.ShouldBindJSON(&createAnchor); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var photo bool
	if createAnchor.Photo != "" {
		photo = true
	} else {
		photo = false
	}
	anchor := Anchor{Name: createAnchor.Name, AbouTheAnchor: createAnchor.AbouTheAnchor, Status: createAnchor.Status, Photo: photo}
	if err := db.Debug().Table("anchor").Where("id=?", anchorId).Update(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "anchor updated Successfully.", "Status": http.StatusOK})
		return
	}
}

/* Delete Program by program Id */
func (hs *HandlerService) DeleteAnchorById(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	chaneelId := c.Param("id")
	var anchor Anchor
	if err := db.Debug().Where("id=? and has_deleted is False ", chaneelId).Find(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	anchor.HasDeleted = true
	if err := db.Debug().Table("anchor").Where("id=? ", chaneelId).Update(&anchor).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	} else {
		l.JSON(c, http.StatusOK, gin.H{"message": "Anchor record Deleted Successfully.", "Status": http.StatusOK})
		return
	}
}
