package productname

import (
	"content/common"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.Use(common.ValidateToken())
	r.GET("/api/content/productnames", hs.ProductName)
}

// GetProductNames -  fetches Productnames
// GET /api/content/productnames
// @Summary Show ProductNames
// @Description get product names
// @Tags Product Name
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} ProductName
// @Router /api/content/productnames [get]
func (hs *HandlerService) ProductName(c *gin.Context) {
	var productname []ProductName
	db := c.MustGet("DB").(*gorm.DB)
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	var errorresponse = common.ServerErrorResponse()
	fmt.Println(productname)
	if err := db.Debug().Table("product_name").Select("Id,Name").Order("id").Find(&productname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": productname})
}
