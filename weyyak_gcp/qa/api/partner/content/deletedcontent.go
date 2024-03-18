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

// package content

// import (
// 	"log"
// 	common "masterdata/common"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/jinzhu/gorm"
// )

// type ContentData struct {
// 	ContentKey int       `json:"content_key"`
// 	ContentId  uuid.UUID `json:"content_id"`
// }

// type SeasonData struct {
// 	SeasonKey int       `json:"season_key"`
// 	SeasonId  uuid.UUID `json:"season_id"`
// }

// type EpisodeData struct {
// 	EpisodeKey int `json:"episode_key"`
// }

// func (hs *HandlerService) GetDeletedContentDetails(c *gin.Context) {
// 	if c.MustGet("AuthorizationRequired") == 1 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
// 		return
// 	}

// 	UserId := c.MustGet("userid")
// 	if UserId == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
// 		return
// 	}

// 	db := c.MustGet("DB").(*gorm.DB)

// 	currentTime := time.Now()
// 	formattedTime := currentTime.Format("2006-01-02")

// 	startdateStr := c.Query("from")
// 	enddateStr := c.Query("to")

// 	var startdate time.Time
// 	var enddate time.Time
// 	var err error

// 	if startdateStr == "" && enddateStr == "" {
// 		startdate, _ = time.Parse("2006-01-02", formattedTime)
// 		enddate, _ = time.Parse("2006-01-02", formattedTime)

// 	} else if startdateStr != "" && enddateStr == "" {
// 		startdate, err = time.Parse("2006-01-02", startdateStr)
// 		if err != nil {
// 			log.Println("Error in parisng startdate: Invalid start date")
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
// 		}
// 		enddate, _ = time.Parse("2006-01-02", formattedTime)

// 	} else if startdateStr == "" && enddateStr != "" {
// 		startdate, _ = time.Parse("2006-01-02", formattedTime)
// 		startdate = startdate.AddDate(0, 0, -7)
// 		enddate, err = time.Parse("2006-01-02", enddateStr)
// 		if err != nil {
// 			log.Println("Error in parisng enddate: Invalid end date")
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
// 		}

// 	} else if startdateStr != "" && enddateStr != "" {
// 		startdate, err = time.Parse("2006-01-02", startdateStr)
// 		if err != nil {
// 			log.Println("Error in parisng startdate: Invalid start date")
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
// 		}

// 		enddate, err = time.Parse("2006-01-02", enddateStr)
// 		if err != nil {
// 			log.Println("Error in parisng enddate: Invalid end date")
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
// 		}
// 	}

// 	enddate = enddate.AddDate(0, 0, 1)

// 	days := daysBetween(startdate, enddate)

// 	if days > 7 {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide dates within 7 days range"})

// 	} else if (startdate.After(time.Now())) || (startdate.After(time.Now()) && enddate.After(time.Now().AddDate(0, 0, 1))) || (enddate.After(time.Now().AddDate(0, 0, 1))) {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide dates less than or equal to today's date"})

// 	} else {
// 		finalData := make(map[string]interface{}, 0)

// 		var CountryResult int32

// 		var country string
// 		if c.Request.URL.Query()["Country"] != nil {
// 			country = c.Request.URL.Query()["Country"][0]
// 		}

// 		CountryResult = common.Countrys(country)
// 		serverError := common.ServerErrorResponse()

// 		oneTier := []string{"movie", "series", "season", "episode"}

// 		var allMultiContents []ContentData
// 		var seasons []SeasonData
// 		var episodes []EpisodeData

// 		for _, types := range oneTier {
// 			var totalCount int

// 			if types == "movie" {
// 				var allOneContents []ContentData
// 				var arr []int
// 				if UserId == os.Getenv("WATCH_NOW") {
// 					db.Debug().Table("content c").Select(`c.content_key`).
// 						Where(`c.watch_now_supplier = true AND
// 								c.content_tier = 1 AND
// 								c.deleted_by_user_id is not null AND
// 								c.modified_at >= ? AND c.modified_at <= ? AND
// 								c.id IS NOT NULL`, startdate, enddate).Order("c.modified_at desc").Find(&allOneContents)
// 					if CountryResult != 0 {
// 						if err := db.Debug().Table("content c").
// 							Joins("join content_variance cv on cv.content_id =c.id").
// 							Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 							Joins("join content_rights cr on cr.id =pi2.rights_id").
// 							Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 							Where("c.watch_now_supplier = true and c.content_tier = 1 and c.deleted_by_user_id is not null and crc.country_id = ? and c.modified_at >= ? AND c.modified_at <= ?", CountryResult, startdate, enddate).Count(&totalCount).Error; err != nil {
// 							c.JSON(http.StatusInternalServerError, serverError)
// 							return
// 						}
// 					} else if country == "" || country == "all" {
// 						if err := db.Debug().Table("content c").
// 							Where("c.watch_now_supplier = true and c.content_tier = 1 and c.deleted_by_user_id is not null and c.modified_at >= ? AND c.modified_at <= ?", startdate, enddate).Count(&totalCount).Error; err != nil {
// 							c.JSON(http.StatusInternalServerError, serverError)
// 							return
// 						}
// 					}
// 				} else {
// 					db.Debug().Table("content c").Select(`c.content_key`).
// 						Where(`
// 					c.content_tier = 1 AND
// 					c.deleted_by_user_id is not null AND
// 					c.modified_at >= ? AND c.modified_at <= ? AND
// 					c.id IS NOT NULL`, startdate, enddate).Order("c.modified_at desc").Find(&allOneContents)

// 					if CountryResult != 0 {
// 						if err := db.Debug().Table("content c").
// 							Joins("join content_variance cv on cv.content_id =c.id").
// 							Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 							Joins("join content_rights cr on cr.id =pi2.rights_id").
// 							Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 							Where("c.content_tier = 1 and c.deleted_by_user_id is not null and crc.country_id = ? and c.modified_at >= ? AND c.modified_at <= ?", CountryResult, startdate, enddate).Count(&totalCount).Error; err != nil {
// 							c.JSON(http.StatusInternalServerError, serverError)
// 							return
// 						}
// 					} else if country == "" || country == "all" {
// 						if err := db.Debug().Table("content c").
// 							Where("c.content_tier = 1 and c.deleted_by_user_id is not null and c.modified_at >= ? AND c.modified_at <= ?", startdate, enddate).Count(&totalCount).Error; err != nil {
// 							c.JSON(http.StatusInternalServerError, serverError)
// 							return
// 						}
// 					}
// 				}
// 				for _, v := range allOneContents {
// 					arr = append(arr, v.ContentKey)
// 				}
// 				finalData[types] = arr
// 			} else if types == "series" {
// 				var arr []int
// 				if UserId == os.Getenv("WATCH_NOW") {
// 					if CountryResult != 0 {
// 						db.Debug().Raw(`
// 								select	c.content_key, c.id as content_id
// 								from content c
// 								join season s on s.content_id = c.id
// 								join content_rights cr on cr.id = s.rights_id
// 								join content_rights_country crc on crc.content_rights_id = cr.id
// 								where
// 								c.watch_now_supplier = 'true'
// 								and c.content_tier = 2
// 									and c.deleted_by_user_id is not null
// 									and c.modified_at >= ? AND c.modified_at <= ?
// 									and crc.country_id = ?

// 					`, startdate, enddate, CountryResult).Order("c.modified_at desc").Find(&allMultiContents)
// 					} else if country == "" || country == "all" {

// 						db.Debug().Raw(`
// 								select	c.content_key, c.id as content_id
// 								from content c
// 								where
// 									c.watch_now_supplier = 'true'
// 									and c.content_tier = 2
// 									and c.deleted_by_user_id is not null
// 									and c.modified_at >= ? AND c.modified_at <= ?
// 					`, startdate, enddate).Order("c.modified_at desc").Find(&allMultiContents)
// 					}
// 				} else {
// 					if CountryResult != 0 {
// 						db.Debug().Raw(`
// 								select
// 								c.content_key, c.id as content_id
// 								from
// 									content c
// 								join season s on s.content_id = c.id
// 								join content_rights cr on cr.id = s.rights_id
// 								join content_rights_country crc on crc.content_rights_id = cr.id

// 								where
// 									c.content_tier = 2
// 									and c.deleted_by_user_id is not null
// 									and c.modified_at >= ? AND c.modified_at <= ?
// 									and crc.country_id = ?

// 					`, startdate, enddate, CountryResult).Order("c.modified_at desc").Find(&allMultiContents)
// 					} else if country == "" || country == "all" {

// 						db.Debug().Raw(`
// 								select
// 									c.content_key, c.id as content_id
// 								from
// 									content c
// 								where
// 									c.content_tier = 2
// 									and c.deleted_by_user_id is not null
// 									and c.modified_at >= ? AND c.modified_at <= ?

// 						`, startdate, enddate).Order("c.modified_at desc").Find(&allMultiContents)
// 					}
// 				}
// 				for _, v := range allMultiContents {
// 					if v.ContentKey != 0 {
// 						arr = append(arr, v.ContentKey)
// 					}
// 				}
// 				finalData[types] = arr
// 			} else if types == "season" {
// 				var season SeasonData
// 				var sarr []int
// 				for _, v := range allMultiContents {
// 					db.Debug().Raw(`
// 						select
// 							s.season_key, s.id as season_id
// 						from
// 							season s
// 						where
// 							s.content_id = ?
// 							and s.deleted_by_user_id is not null
// 							and s.modified_at >= ? AND s.modified_at <= ?
// 					`, v.ContentId, startdate, enddate).Order("s.modified_at desc").Find(&season)

// 					seasons = append(seasons, season)
// 				}
// 				for _, s := range seasons {
// 					if s.SeasonKey != 0 {
// 						sarr = append(sarr, s.SeasonKey)
// 					}
// 				}
// 				finalData[types] = sarr
// 			} else if types == "episode" {
// 				var episode EpisodeData
// 				var earr []int
// 				for _, v := range seasons {
// 					db.Debug().Raw(`
// 						select
// 							e.episode_key
// 						from
// 							episode e
// 						where
// 							e.season_id = ?
// 							and e.deleted_by_user_id is not null
// 							and e.modified_at >= ? AND e.modified_at <= ?
// 					`, v.SeasonId, startdate, enddate).Order("e.modified_at desc").Find(&episode)

// 					episodes = append(episodes, episode)
// 				}
// 				for _, e := range episodes {
// 					if e.EpisodeKey != 0 {
// 						earr = append(earr, e.EpisodeKey)
// 					}
// 				}
// 				finalData[types] = earr
// 			}
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": finalData})
// 	}
// }
