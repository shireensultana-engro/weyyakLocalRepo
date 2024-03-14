package dashboard

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	srg := r.Group("/api")
	srg.GET("/dashboard", hs.GetDashBoardDetails)
	srg.GET("/shows", hs.GetShowDetails)

}

// GetDashBoardDetails -  Get Dash Board Details
// GET /api/dashboard
// @Summary Get Dash Board Details
// @Description Get Dash Board Details
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {array} Dashboard
// @Router /api/dashboard [get]
func (hs *HandlerService) GetDashBoardDetails(c *gin.Context) {
	// db := c.MustGet("DB").(*gorm.DB)
	// udb := c.MustGet("UDB").(*gorm.DB)
	// fcdb := c.MustGet("FCDB").(*gorm.DB)

	var dashboard Dashboard

	Data := []byte(`{
		"WatchedContent": [
			{
				"contentTitle": "Noqta Ala Al Sater S4",
				"watchedCount": 11
			},
			{
				"contentTitle": "Light of Hope CK1",
				"watchedCount": 11
			},
			{
				"contentTitle": "Makanak Fee Al Qalb 7",
				"watchedCount": 10
			},
			{
				"contentTitle": "The Tashkent Files HBC",
				"watchedCount": 8
			},
			{
				"contentTitle": "Bent Bnoot",
				"watchedCount": 7
			},
			{
				"contentTitle": "Sharaf s75",
				"watchedCount": 4
			},
			{
				"contentTitle": "Jo Jeeta Wahi Sikandar 1",
				"watchedCount": 4
			},
			{
				"contentTitle": "Hayat Qalbi 3",
				"watchedCount": 4
			},
			{
				"contentTitle": "Asia",
				"watchedCount": 4
			},
			{
				"contentTitle": "PlayNew3_EN",
				"watchedCount": 4
			},
			{
				"contentTitle": "Hawazik",
				"watchedCount": 3
			},
			{
				"contentTitle": "Ramadan 2022",
				"watchedCount": 3
			},
			{
				"contentTitle": "Third Party Movie - EN",
				"watchedCount": 3
			},
			{
				"contentTitle": "AlMohakek B Aqleyt Mojrm2",
				"watchedCount": 3
			},
			{
				"contentTitle": "Eh Gabak And Garak 14",
				"watchedCount": 3
			},
			{
				"contentTitle": "Mom",
				"watchedCount": 3
			},
			{
				"contentTitle": "Ebnat Al Safeer 2",
				"watchedCount": 3
			},
			{
				"contentTitle": "Jurm Jurn 3",
				"watchedCount": 3
			},
			{
				"contentTitle": "Pyaar Kiya toh Darna Kya1",
				"watchedCount": 2
			},
			{
				"contentTitle": "A Hzan Maryam",
				"watchedCount": 2
			},
			{
				"contentTitle": "Ala Kayd Al Hob",
				"watchedCount": 2
			},
			{
				"contentTitle": "Al Bait Al Kabeer 2",
				"watchedCount": 2
			},
			{
				"contentTitle": "Al Mohakek B Aqleyt Mojrm",
				"watchedCount": 2
			},
			{
				"contentTitle": "Amir Ahlame",
				"watchedCount": 2
			},
			{
				"contentTitle": "Driving Dirty",
				"watchedCount": 2
			},
			{
				"contentTitle": "Fatet Lebet",
				"watchedCount": 2
			},
			{
				"contentTitle": "Go Goa Gone",
				"watchedCount": 2
			},
			{
				"contentTitle": "Hadath Fe Demashq",
				"watchedCount": 2
			},
			{
				"contentTitle": "Haider",
				"watchedCount": 2
			},
			{
				"contentTitle": "Hob Marfoud",
				"watchedCount": 2
			},
			{
				"contentTitle": "Hob Yatakhatta Al Zaman 2",
				"watchedCount": 2
			},
			{
				"contentTitle": "Kedarnath",
				"watchedCount": 2
			},
			{
				"contentTitle": "Mann Ana S2",
				"watchedCount": 2
			},
			{
				"contentTitle": "Phobia",
				"watchedCount": 2
			},
			{
				"contentTitle": "programsNew2_EN",
				"watchedCount": 2
			},
			{
				"contentTitle": "Provoked",
				"watchedCount": 2
			},
			{
				"contentTitle": "Adowi...Sadiki",
				"watchedCount": 2
			},
			{
				"contentTitle": "Simmba",
				"watchedCount": 2
			},
			{
				"contentTitle": "Simply Beautiful",
				"watchedCount": 2
			},
			{
				"contentTitle": "Telling Tales 2",
				"watchedCount": 2
			},
			{
				"contentTitle": "The Blind Spot",
				"watchedCount": 2
			},
			{
				"contentTitle": "Total Siyapaa",
				"watchedCount": 2
			},
			{
				"contentTitle": "Boss Wants a Happy Ending",
				"watchedCount": 1
			},
			{
				"contentTitle": "Manikarnika",
				"watchedCount": 1
			},
			{
				"contentTitle": "Tutak Tutak Tutiya",
				"watchedCount": 1
			},
			{
				"contentTitle": "Men Nazret Hob",
				"watchedCount": 1
			},
			{
				"contentTitle": "Mobster's Guru 2",
				"watchedCount": 1
			},
			{
				"contentTitle": "Book Box",
				"watchedCount": 1
			},
			{
				"contentTitle": "BO2_EN",
				"watchedCount": 1
			},
			{
				"contentTitle": "Noqta Al Sater S3",
				"watchedCount": 1
			}
		],
		"UserDevicesresponse": [
			{
				"lable": "android",
				"value": 436078
			},
			{
				"lable": "apple_tv",
				"value": 8
			},
			{
				"lable": "ios",
				"value": 45435
			},
			{
				"lable": "roku",
				"value": 21
			},
			{
				"lable": "smart_tv",
				"value": 65489
			}
		],
		"UserByRegion": [
			{
				"regionName": "Afghanistan",
				"userCount": 8
			},
			{
				"regionName": "Albania",
				"userCount": 2
			},
			{
				"regionName": "Algeria",
				"userCount": 1
			},
			{
				"regionName": "Angola",
				"userCount": 2
			},
			{
				"regionName": "Australia",
				"userCount": 2
			},
			{
				"regionName": "Austria",
				"userCount": 3
			},
			{
				"regionName": "Bahrain",
				"userCount": 12
			},
			{
				"regionName": "Botswana",
				"userCount": 1
			},
			{
				"regionName": "Brazil",
				"userCount": 1
			},
			{
				"regionName": "Central African Republic",
				"userCount": 1
			},
			{
				"regionName": "Ethiopia",
				"userCount": 1
			},
			{
				"regionName": "Eritrea",
				"userCount": 1
			},
			{
				"regionName": "Åland Islands",
				"userCount": 4
			},
			{
				"regionName": "Palestine",
				"userCount": 2
			},
			{
				"regionName": "Germany",
				"userCount": 1
			},
			{
				"regionName": "India",
				"userCount": 317
			},
			{
				"regionName": "Iraq",
				"userCount": 2
			},
			{
				"regionName": "Israel",
				"userCount": 1
			},
			{
				"regionName": "Jordan",
				"userCount": 1
			},
			{
				"regionName": "Kuwait",
				"userCount": 4
			},
			{
				"regionName": "Lebanon",
				"userCount": 1
			},
			{
				"regionName": "Martinique",
				"userCount": 1
			},
			{
				"regionName": "Qatar",
				"userCount": 1
			},
			{
				"regionName": "Saudi Arabia",
				"userCount": 148
			},
			{
				"regionName": "Zimbabwe",
				"userCount": 2
			},
			{
				"regionName": "Syrian Arab Republic",
				"userCount": 1
			},
			{
				"regionName": "United Arab Emirates",
				"userCount": 73
			},
			{
				"regionName": "Turkey",
				"userCount": 2
			},
			{
				"regionName": "Egypt",
				"userCount": 2
			},
			{
				"regionName": "United Kingdom",
				"userCount": 833
			},
			{
				"regionName": "United States of America",
				"userCount": 10
			}
		],
		"ActiveUsers": {
			"activeUserCount": "10"
		}
	}`)

	json.Unmarshal(Data, &dashboard)

	c.JSON(http.StatusOK, gin.H{"Dashboard": dashboard})

	// var watchedcontent []WatchedContent
	// var userdevice []UserDevicesresponse
	// var userbyregion []UserByRegion
	// var activeusercount ActiveUsers
	// var applicationsetting ApplicationSetting

	// redisKey := os.Getenv("REDIS_CONTENT_KEY") + "_GetDashBoardDetails"

	// // val, err := common.GetRedisDataWithKey(redisKey)

	// // fmt.Println("rrrrrrrrrrr", err)
	// // if err == nil {
	// // 	var (
	// // 		redisResponse Dashboard
	// // 	)
	// // 	Data := []byte(val)
	// // 	json.Unmarshal(Data, &redisResponse)
	// // 	c.JSON(http.StatusOK, gin.H{"Dashboard": redisResponse})
	// // 	return
	// // }

	// db.Debug().Raw("select  cpi.transliterated_title as content_title ,count(va.content_id) as watched_count from view_activity va join content c on c.id = va.content_id join content_primary_info cpi on cpi.id = c.primary_info_id where viewed_at between now() - interval '50 DAYS' and now() group by cpi.transliterated_title order by count(va.content_id)  desc limit 50").Find(&watchedcontent)

	// udb.Debug().Raw("select p.name as lable,count(d.platform) as value from platform p join device d on d.platform  = p.platform_id group by p.name").Find(&userdevice)

	// udb.Debug().Raw("select c.english_name as region_name,count(u.id) as user_count from public.user u join public.country c on c.id = u.country where u.last_activity_at between now() - interval '100 DAYS' and now() and u.country is not null and u.country != 0 group by u.country,c.english_name").Find(&userbyregion)

	// fcdb.Raw("select as2.value  from application_setting as2 where name ='ActiveTime'").Find(&applicationsetting)

	// udb.Debug().Raw("select count(id) as active_user_count from public.user where extract(EPOCH from (now() - last_activity_at))< ? ", applicationsetting.Value).Find(&activeusercount)

	// dashboard.WatchedContent = watchedcontent

	// dashboard.UserDevicesresponse = userdevice

	// dashboard.UserByRegion = userbyregion

	// dashboard.ActiveUsers = activeusercount

	// m, _ := json.Marshal(dashboard)
	// err := common.PostRedisDataWithKey(redisKey, m)
	// if err != nil {
	// 	fmt.Println("Redis Value Not Updated")
	// }

}

// GetShowDetails -  Get Show Details
// GET /api/shows
// @Summary Get Show Details
// @Description Get Show Details
// @Tags Shows
// @Accept  json
// @Produce  json
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} Content
// @Router /api/shows [get]
func (hs *HandlerService) GetShowDetails(c *gin.Context) {
	cdb := c.MustGet("DB").(*gorm.DB)
	var showkeys []Content
	var seasonkeys []Content
	var episodekeys []Content
	var finalkeys []Content
	var showids []int
	var seasonids []int
	//for creating view
	/*cdb.Debug().Raw("select (case when show_id != 0 then show_id when season_id != 0 then season_id when episode_id != 0 then episode_id end) as id,parent,droppable,text,content_type from (select content_key as show_id,0 as season_id, 0 as episode_id, 0 as parent,true as droppable,0 as content_type,cpi.transliterated_title as text from content c join content_primary_info cpi on cpi.id = c.primary_info_id where  c.content_type = 'Series' and c.deleted_by_user_id is null union select 0 as show_id,s.season_key  as season_id, 0 as episode_id, c.content_key  as parent,true as droppable,1 as content_type,cpi2.transliterated_title as text from season s join content_primary_info cpi2 on cpi2.id = s.primary_info_id join content c on c.id  = s.content_id where s.deleted_by_user_id is null union select 0 as show_id,0 as season_id,e.episode_key  as episode_id, s2.season_key  as parent,false as droppable,2 as content_type,cpi3.transliterated_title as text from episode e join content_primary_info cpi3 on cpi3.id =e.primary_info_id join season s2 on s2.id = e.season_id where e.deleted_by_user_id is null) ss").Find(&keys)*/
	//imagery
	// for _, k := range keys {
	// 	image := os.Getenv("IMAGERY_URL") + k.ShowId + "/" + k.SeasonId + "/poster-image"
	// 	image := os.Getenv("IMAGERY_URL") + k.ShowId + "/" + k.SeasonId + "/" + k.EpisodeId + "/poster-image"
	// }
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	cdb.Debug().Raw("select *  from public.shows s2 where content_type = 0 limit ? offset ?", limit, offset).Find(&showkeys)
	for _, k := range showkeys {
		showids = append(showids, k.Id)
	}
	cdb.Debug().Raw("select * from public.shows s where parent in (?)", showids).Find(&seasonkeys)
	for _, k := range seasonkeys {
		seasonids = append(seasonids, k.Id)
	}
	cdb.Debug().Raw("select * from public.shows s3 where parent in (?)", seasonids).Find(&episodekeys)
	finalkeys = append(finalkeys, showkeys...)
	finalkeys = append(finalkeys, seasonkeys...)
	finalkeys = append(finalkeys, episodekeys...)
	c.JSON(http.StatusOK, gin.H{"keys": finalkeys})
}
