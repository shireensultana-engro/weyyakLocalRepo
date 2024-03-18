package content

import (
	"log"

	// common "masterdata/common"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ContentData struct {
	ContentKey int `json:"content_key"`
}

type SeasonData struct {
	SeasonKey int `json:"season_key"`
}

type EpisodeData struct {
	EpisodeKey int `json:"episode_key"`
}

func (hs *HandlerService) GetDeletedContentDetails(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}

	db := c.MustGet("DB").(*gorm.DB)

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02")

	startdateStr := c.Query("from")
	enddateStr := c.Query("to")

	var startdate time.Time
	var enddate time.Time
	var err error

	if startdateStr == "" && enddateStr == "" {
		startdate, _ = time.Parse("2006-01-02", formattedTime)
		enddate, _ = time.Parse("2006-01-02", formattedTime)

	} else if startdateStr != "" && enddateStr == "" {
		startdate, err = time.Parse("2006-01-02", startdateStr)
		if err != nil {
			log.Println("Error in parisng startdate: Invalid start date")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
		}
		enddate, _ = time.Parse("2006-01-02", formattedTime)

	} else if startdateStr == "" && enddateStr != "" {
		startdate, _ = time.Parse("2006-01-02", formattedTime)
		startdate = startdate.AddDate(0, 0, -7)
		enddate, err = time.Parse("2006-01-02", enddateStr)
		if err != nil {
			log.Println("Error in parisng enddate: Invalid end date")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
		}

	} else if startdateStr != "" && enddateStr != "" {
		startdate, err = time.Parse("2006-01-02", startdateStr)
		if err != nil {
			log.Println("Error in parisng startdate: Invalid start date")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
		}

		enddate, err = time.Parse("2006-01-02", enddateStr)
		if err != nil {
			log.Println("Error in parisng enddate: Invalid end date")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
		}
	}

	enddate = enddate.AddDate(0, 0, 1)

	days := daysBetween(startdate, enddate)
	finalData := make(map[string]interface{}, 0)

	if days > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide dates within 7 days range"})
		return

	} else if (startdate.After(time.Now())) || (startdate.After(time.Now()) && enddate.After(time.Now().AddDate(0, 0, 1))) || (enddate.After(time.Now().AddDate(0, 0, 1))) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide dates less than or equal to today's date"})
		return

	} else {
		var allOneContents []ContentData
		var allMultiContents []ContentData
		var seasons []SeasonData
		var episodes []EpisodeData

		if UserId == os.Getenv("WATCH_NOW") {
			db.Table("content").Select(`content_key`).
				Where(`watch_now_supplier = true AND
								content_tier = 1 AND 
								deleted_by_user_id is not null AND 
								modified_at >= ? AND modified_at <= ? AND
								id IS NOT NULL`, startdate, enddate).Order("modified_at desc").Find(&allOneContents)
		} else {
			db.Table("content").Select(`content_key`).
				Where(`
					content_tier = 1 AND 
					deleted_by_user_id is not null AND 
					modified_at >= ? AND modified_at <= ? AND
					id IS NOT NULL`, startdate, enddate).Order("modified_at desc").Find(&allOneContents)
		}
		finalData["movie"] = allOneContents

		if UserId == os.Getenv("WATCH_NOW") {
			db.Raw(`
								select content_key
								from content
								where watch_now_supplier = 'true'
									and content_tier = 2
									and deleted_by_user_id is not null
									and modified_at >= ? AND modified_at <= ?										
							
					`, startdate, enddate).Order("modified_at desc").Find(&allMultiContents)
		} else {
			db.Raw(`
								select content_key
								from content
								where content_tier = 2						
									and deleted_by_user_id is not null
									and modified_at >= ? AND modified_at <= ?
							
						`, startdate, enddate).Order("modified_at desc").Find(&allMultiContents)
		}
		finalData["series"] = allMultiContents

		db.Raw(`
					select
						season_key
					from
						season
					where
						deleted_by_user_id is not null
						and modified_at >= ? AND modified_at <= ?
		`, startdate, enddate).Order("modified_at desc").Find(&seasons)
		finalData["season"] = seasons

		db.Raw(`
					select
						episode_key
					from
						episode
					where
						deleted_by_user_id is not null
						and modified_at >= ? AND modified_at <= ?
				`, startdate, enddate).Order("modified_at desc").Find(&episodes)
		finalData["episode"] = episodes

	}
	c.JSON(http.StatusOK, gin.H{"data": finalData})
	return
}
