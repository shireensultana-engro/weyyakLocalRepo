package channel

import (
	// "content/common"

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
	qrg.POST("/channel", hs.CreateProgamList)
	qrg.GET("/channel/:chaneelname", hs.GetProgramlistByChannelName)
	qrg.GET("/channel", hs.GetProgamList)
	qrg.POST("/channel/:id", hs.UpdateProgramById)
	qrg.DELETE("/channel/:id", hs.DeleteProgramById)

}

/*create Program list for channel */
func (hs *HandlerService) CreateProgamList(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var createChannelProgramList CreateChannelProgramList
	if err := c.ShouldBindJSON(&createChannelProgramList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var logo bool
	if createChannelProgramList.Logo != "" {
		logo = true
	} else {
		logo = false
	}
	Channel := ChannelProgramList{Name: createChannelProgramList.Name, Url: createChannelProgramList.Url, StartTime: createChannelProgramList.StartTime, EndTime: createChannelProgramList.EndTime, Duration: createChannelProgramList.Duration, Channel: createChannelProgramList.Channel, Site: createChannelProgramList.Site, Lang: createChannelProgramList.Lang, Logo: logo}
	if err := db.Debug().Create(&Channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	} else {
		go programLogoUpload(createChannelProgramList.Logo, createChannelProgramList.Channel, createChannelProgramList.Name)
		c.JSON(http.StatusOK, gin.H{"message": "program sheduled Successfully."})
		return
	}
}

/*Get all list of Programs*/
func (hs *HandlerService) GetProgamList(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var channelProgramList []ChannelProgramList
	if err := db.Debug().Where("has_deleted is False").Find(&channelProgramList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": channelProgramList})
}

/* Get programs list by channel name */
func (hs *HandlerService) GetProgramlistByChannelName(c *gin.Context) {
	var channelProgramList []ChannelProgramList
	chaneelname := c.Param("chaneelname")
	db := c.MustGet("DB").(*gorm.DB)
	if err := db.Debug().Where("channel=? and has_deleted is False", chaneelname).Find(&channelProgramList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, channelProgramList)
}

/* Update Program information by program Id */
func (hs *HandlerService) UpdateProgramById(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	chaneelId := c.Param("id")
	var createChannelProgramList CreateChannelProgramList
	if err := c.ShouldBindJSON(&createChannelProgramList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	var logo bool
	if createChannelProgramList.Logo != "" {
		logo = true
	} else {
		logo = false
	}
	Channel := ChannelProgramList{Name: createChannelProgramList.Name, Url: createChannelProgramList.Url, StartTime: createChannelProgramList.StartTime, EndTime: createChannelProgramList.EndTime, Duration: createChannelProgramList.Duration, Channel: createChannelProgramList.Channel, Site: createChannelProgramList.Site, Lang: createChannelProgramList.Lang, Logo: logo}
	if err := db.Debug().Table("channel_program_list").Where("id=?", chaneelId).Update(&Channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Record updated Successfully.", "Status": http.StatusOK})
		return
	}
}

/* Delete Program by program Id */
func (hs *HandlerService) DeleteProgramById(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	chaneelId := c.Param("id")
	var channelProgramList ChannelProgramList
	if err := db.Debug().Where("id=? and has_deleted is False ", chaneelId).Find(&channelProgramList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}
	channelProgramList.HasDeleted = true
	if err := db.Debug().Table("channel_program_list").Where("id=? ", chaneelId).Update(&channelProgramList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Record Deleted Successfully.", "Status": http.StatusOK})
		return
	}
}
