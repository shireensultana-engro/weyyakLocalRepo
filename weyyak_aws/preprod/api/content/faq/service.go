package faq

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	srg := r.Group("/api")
	srg.GET("/faq", hs.GetListofFAQ)
	srg.POST("/faq", hs.CreateFAQ)
	srg.GET("/faq/:id", hs.GetFaqByID)
	srg.DELETE("/delete/:id", hs.DeleteFAQ)
	srg.PATCH("/updatevalue/:id", hs.UpdateFAQ)

}

//GET
//getting all the list of FAQ
func (hs *HandlerService) GetListofFAQ(c *gin.Context) {
	var finalResult []Faq
	db := c.MustGet("DB").(*gorm.DB)
	db.Find(&finalResult)
	c.JSON(http.StatusOK, gin.H{"Data": finalResult})
}

//POST
//creating FAQ
func (hs *HandlerService) CreateFAQ(c *gin.Context) {
	var req Faq
	c.ShouldBindJSON(&req)
	db := c.MustGet("DB").(*gorm.DB)
	var faqs = Faq{Question: req.Question, Description: req.Description}
	db.Debug().Create(&faqs)

	c.JSON(http.StatusOK, gin.H{"View": faqs})
}

// GET
// route /FAQ/:id
// Find  FAQ
//@Router/api/faq/{id}
func (hs *HandlerService) GetFaqByID(c *gin.Context) { // Get model if exist
	var faq Faq
	db := c.MustGet("DB").(*gorm.DB)

	if err := db.Where("id = ?", c.Param("id")).Find(&faq).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": faq})
}

//Delete
// route /deletevalue
func (hs *HandlerService) DeleteFAQ(c *gin.Context) {

	var req Faq
	id := c.Params.ByName("id")
	//object
	c.ShouldBindJSON(&req)
	db := c.MustGet("DB").(*gorm.DB)
	if err := db.Where("id = ?", id).Delete(&req).Error; err != nil {

		c.JSON(http.StatusOK, gin.H{"View" + id: "deleted"})
	}
}

//PATCH
// route /updatevalue/:id
// changing the values with respect to id
func (hs *HandlerService) UpdateFAQ(c *gin.Context) {
	var req Faq
	// c.ShouldBindJSON(&req)
	db := c.MustGet("DB").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&req).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	//validate input
	var faq Faq
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Model(&req).Updates(faq)
	c.JSON(http.StatusOK, gin.H{"View": req})
}
