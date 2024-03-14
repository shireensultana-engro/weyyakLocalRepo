package digitalRights

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
	qrg := r.Group("/api/content")
	qrg.Use(common.ValidateToken())
	qrg.GET("/digitalrights/types", hs.GetAllDigitalTypes)
	qrg.GET("/displaystatuses", hs.GetAllDisplayStatus)
	qrg.GET("/digitalrights/regions/all", hs.GetDigitalRightsRegions)
}

// GetAllDigitalTypes -  fetches all digital types
// GET /api/content/digitalrights/types
// @Summary Show a list of all digital types
// @Description get list of all digital types
// @Tags Digital
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} ContentRights
// @Router /api/content/digitalrights/types [get]
func (hs *HandlerService) GetAllDigitalTypes(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	fcdb := c.MustGet("FCDB").(*gorm.DB)
	var errorresponse = common.ServerErrorResponse()
	var digitalrights []ContentRights
	if err := fcdb.Debug().Table("digital_rights_types").Find(&digitalrights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": digitalrights})
}

// GetAllDisplayStatus -  fetches all display status
// GET /api/content/digitalrights
// @Summary Show a list of all display status
// @Description get list of all display status
// @Tags Digital
// @Accept  json
// @Produce  json
// @security Authorization
// @Success 200 {array} DisplayStatus
// @Router /api/content/displaystatuses [get]
func (hs *HandlerService) GetAllDisplayStatus(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var erroresponse = common.ServerErrorResponse()
	var displayStatus []DisplayStatus
	if err := db.Debug().Find(&displayStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": displayStatus})
}

// GetDigitalRightsRegions -  Get digital rights regions
// GET /api/content/digitalrights/regions/all
// @Summary show list of digital rights regions
// @Description Get digital rights regions
// @Tags Digital
// @Security Authorization
// @Accept  json
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/content/digitalrights/regions/all [get]
func (hs *HandlerService) GetDigitalRightsRegions(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	// db := c.MustGet("DB").(*gorm.DB)
	db := c.MustGet("FCDB").(*gorm.DB)
	var continents []Continent
	var regions []Region
	var countries []CountrysResponse
	var region RegionsResponse
	var resp DigitalrightsResponse
	var response []DigitalrightsResponse
	var errorresponse = common.ServerErrorResponse()
	// For Continent
	fields := "id,name"
	if data := db.Debug().Select(fields).Order("name ASC").Find(&continents).Error; data != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	if continents != nil {
		for _, con := range continents {
			// For Regions
			fields := "id,name ,continent_id"
			if data := db.Debug().Select(fields).Where("continent_id = ?", con.Id).Order("name ASC").Find(&regions).Error; data != nil {
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
			if regions != nil {
				var regionResponse []RegionsResponse
				for _, reg := range regions {

					// For Country
					// fields := "id,english_name as name,calling_code"
					fields := "id,english_name as name"
					if data := db.Debug().Table("country").Select(fields).Where("region_id = ?", reg.Id).Order("english_name ASC").Find(&countries).Error; data != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}

					region.Name = reg.Name
					region.Countries = countries
					regionResponse = append(regionResponse, region)
				} // for each loop closed here
				resp.Name = con.Name
				resp.Regions = regionResponse
				resp.Countries = nil
				response = append(response, resp)
			} // If Closed

		} // for each loop closed here
	} // If Closed

	c.JSON(http.StatusOK, gin.H{"data": response})
	return
}
