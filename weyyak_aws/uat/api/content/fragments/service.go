package fragments

import (
	"content/common"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.POST("/contentfrag", hs.CreateContentFragments)
	r.POST("/api/contentfragment", hs.CreateContentFragmentsThroughJsonObject)
	r.GET("/api/contentsdigitalrightsexceeded", hs.ContentsPagesPlaylistsSlidersDigitalRightsExceeded)
}

func (hs *HandlerService) ContentsPagesPlaylistsSlidersDigitalRightsExceeded(c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	fcdb := c.MustGet("FCDB").(*gorm.DB)
	type ContentIds struct {
		ContentId string `json:"content_id"`
	}
	var contentids []ContentIds
	if err := fdb.Debug().Raw("select distinct cf.content_id from content_fragment cf where cf.rights_end_date  between now()  - interval  '1 DAY' and now() ").Find(&contentids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}
	if len(contentids) > 0 {
		// For movies,series and livetv
		rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
		searchPattern := os.Getenv("REDIS_CONTENT_KEY") + "*"
		if len(os.Args) > 1 {
			searchPattern = os.Args[1]
		}
		iter := rdb.Keys(searchPattern)
		for _, Rkey := range iter.Val() {
			rdb.Del(Rkey)
			fmt.Println("Redis key " + Rkey + "is deleted")
		}
	}
	for _, val := range contentids {
		// For pages(play lists and sliders),sidemenu where the content is involved
		common.ContentSynching(val.ContentId, c)
	}
	var playlistids []Ids
	var playlistidstogetsliders []string
	// to fetch sliders which are related to that playlists
	var slideridsofplaylists []Ids
	var sliderids []Ids
	fcdb.Debug().Raw("select id from playlist p where p.scheduling_start_date = current_date or p.scheduling_end_date = current_date").Find(&playlistids)
	// fetching slider ids where the playlist is involved
	for _, val := range playlistids {
		playlistidstogetsliders = append(playlistidstogetsliders, val.Id)
	}
	fcdb.Debug().Raw("select s.id from slider s where s.green_area_playlist_id in (?) or s.red_area_playlist_id in (?) or s.black_area_playlist_id in (?)", playlistidstogetsliders, playlistidstogetsliders, playlistidstogetsliders).Find(&slideridsofplaylists)
	fcdb.Debug().Raw("select id from slider s where s.scheduling_start_date = current_date or s.scheduling_end_date = current_date").Find(&sliderids)
	sliderids = append(sliderids, slideridsofplaylists...)
	// calling related methods here
	for _, value := range playlistids {
		common.PlaylistSynching(value.Id, c)
	}
	for _, value := range sliderids {
		common.SliderSynching(value.Id, c)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func (hs *HandlerService) CreateContentFragments(c *gin.Context) {
	var request RequestIds
	// c.ShouldBindJSON(&request)
	cdb := c.MustGet("DB").(*gorm.DB)
	cdb.Debug().Raw("select id as ids from content where status=1").Find(&request)
	fmt.Println("time start...........", time.Now())
	for _, i := range request.Ids {
		// fmt.Println(i)
		CreateContentFragment(i, c)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("time end...........", time.Now())
	// fdb := c.MustGet("FDB").(*gorm.DB)
	// cdb := c.MustGet("DB").(*gorm.DB)
	// ctx := context.Background()
	// tx := fdb.BeginTx(ctx, nil)

	// var season []Seasons
	// var contentDetails []ContentDetails
	// var contentType ContentType
	// var regionsAndPlatforms RegionsAndPlatforms
	// contentId := c.Param("contentid")

	// if contentTypeErr := cdb.Debug().Table("content").Select("content_tier").Where("id=?", contentId).Find(&contentType).Error; contentTypeErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 	return
	// }
	// fmt.Println(contentType.ContentTier, "????????????????????????????????????????")
	// if contentType.ContentTier == 1 {
	// 	if onetiererr := cdb.Debug().Raw("select c.id as content_id,c.content_key as id,json_agg(distinct aaa.english_name)::varchar as en_cast,json_agg(distinct aaa.arabic_name)::varchar as ar_cast,cpi.original_title as title,json_agg(distinct g.english_name)::varchar as en_genres,json_agg(distinct g.arabic_name)::varchar as ar_genres,pi2.duration as length,cpi.original_title as en_seo_title,cpi.arabic_title as ar_seo_title,cpi.transliterated_title as friendly_url,cpi.original_title as title,ac.production_year,ac.english_synopsis ,ac.arabic_synopsis ,c.content_type ,cv.has_all_rights,cv.id as content_variance_id,"+ /*json_agg(distinct crc.country_id)::varchar country,*/ "c.modified_at ,c.has_poster_image,c.created_at as inserted_at,ac.english_synopsis as english_seo_description,ac.arabic_synopsis as arabic_seo_description,ar.english_name as en_age_rating,ar.arabic_name as ar_age_rating,a.english_name as en_main_actor,a.arabic_name as ar_main_actor,aa.english_name as en_main_actress,aa.arabic_name as ar_main_actress,pi2.video_content_id as video_id,json_agg(distinct tdt.name)::varchar as tags,"+ /*json_agg(distinct pitp.target_platform)::varchar as platforms,*/ "false as geoblock,cr.digital_rights_type as digital_right_type,cr.digital_rights_start_date,cr.digital_rights_end_date from content c join about_the_content_info ac on ac.id = c.about_the_content_info_id join content_primary_info cpi on cpi.id = c.primary_info_id join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id "+ /*join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id */ " join content_rights cr on cr.id = pi2.rights_id "+ /*join content_rights_country crc on crc.content_rights_id = cr.id*/ " join age_ratings ar on ac.age_group = ar.id join content_cast cc on cc.id = c.cast_id left join actor a on a.id = cc.main_actor_id left join actor aa on aa.id = cc.main_actress_id left join content_actor ca on ca.cast_id = c.cast_id left join actor aaa on aaa.id = ca.actor_id join content_genre cg on cg.content_id = c.id join genre g on g.id = cg.genre_id left join content_tag ct on ct.tag_info_id = c.tag_info_id left join textual_data_tag tdt on tdt.id = ct.textual_data_tag_id where c.id =? group by c.id,cpi.transliterated_title,ac.production_year,cv.id,ar.english_name,cpi.arabic_title,ar.arabic_name,a.english_name,a.arabic_name,aa.english_name,aa.arabic_name,pi2.video_content_id,pi2.duration,cr.digital_rights_type,cpi.original_title,ac.english_synopsis ,ac.arabic_synopsis,cr.digital_rights_start_date,cr.digital_rights_end_date,cr.digital_rights_type", contentId).Find(&contentDetails).Error; onetiererr != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 		return
	// 	}
	// 	cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country,json_agg(distinct pitp.target_platform)::varchar as platforms from content c join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where c.id =?", contentId).Find(&regionsAndPlatforms)

	// } else {
	// 	if multitiererr := cdb.Debug().Raw("select c.id as content_id,s.id as season_id,c.content_key as id,json_agg(distinct aaa.english_name)::varchar as en_cast,json_agg(distinct aaa.arabic_name)::varchar as ar_cast,cpi.original_title as title,json_agg(distinct g.english_name)::varchar as en_genres,json_agg(distinct g.arabic_name)::varchar as ar_genres,sum(pi2.duration) as length,cpi.original_title as en_seo_title,cpi.arabic_title as ar_seo_title,cpi.transliterated_title as friendly_url,cpi.original_title as title,ac.production_year,ac.english_synopsis ,ac.arabic_synopsis ,c.content_type ,s.has_all_rights,s.id as content_variance_id,"+ /*json_agg(distinct crc.country_id)::varchar country,*/ "c.modified_at ,c.has_poster_image,c.created_at as inserted_at,ac.english_synopsis as english_seo_description,ac.arabic_synopsis as arabic_seo_description,ar.english_name as en_age_rating,ar.arabic_name as ar_age_rating,a.english_name as en_main_actor,a.arabic_name as ar_main_actor,aa.english_name as en_main_actress,aa.arabic_name as ar_main_actress,min(pi2.video_content_id) as video_id,json_agg(distinct tdt.name)::varchar as tags,"+ /*json_agg(distinct pitp.target_platform)::varchar as platforms,*/ "false as geoblock,cr.digital_rights_type as digital_right_type,s.season_key as season_id,cr.digital_rights_start_date ,cr.digital_rights_end_date from content c join season s on s.content_id = c.id join episode e on e.season_id = s.id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info ac on ac.id = s.about_the_content_info_id join age_ratings ar on ac.age_group = ar.id join playback_item pi2 on pi2.id = e.playback_item_id "+ /*join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id */ "join content_rights cr on cr.id = pi2.rights_id "+ /*join content_rights_country crc on crc.content_rights_id = cr.id*/ " join content_cast cc on cc.id = s.cast_id left join actor a on a.id= cc.main_actor_id left join actor aa on aa.id = cc.main_actress_id left join content_actor ca on ca.cast_id = c.cast_id left join actor aaa on aaa.id= ca.actor_id left join season_genre sg on sg.season_id =s.id left join genre g on g.id = sg.genre_id left join content_tag_info cti on cti.id = s.tag_info_id left join content_tag ct on ct.tag_info_id =s.tag_info_id left join textual_data_tag tdt on tdt.id=ct.textual_data_tag_id where c.id =? group by c.id,cpi.id,ac.id,cr.digital_rights_type,s.id ,ar.english_name,ar.arabic_name,a.id,aa.id,aaa.id,cr.digital_rights_start_date,cr.digital_rights_end_date", contentId).Find(&contentDetails).Error; multitiererr != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 		return
	// 	}
	// 	cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country,json_agg(distinct pitp.target_platform)::varchar as platforms from content c join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi2 on pi2.id = e.playback_item_id join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where c.id = ?", contentId).Find(&regionsAndPlatforms)

	// }

	// var contentFragment ContentFragment
	// var contentFragmentDetails ContentFragmentDetails
	// var movie Movie
	// var movies []Movie
	// var contentFragments []ContentFragment
	// var varianceIds []string
	// lang := "ar"
	// for _, value := range contentDetails {
	// 	/* movie details */
	// 	movie.ID = 0
	// 	movie.Title = value.Title
	// 	movie.Geoblock = false
	// 	movie.DigitalRightType = value.DigitalRightType
	// 	//countries, _ := JsonStringToIntSliceOrMap(value.Country)
	// 	movie.DigitalRightsRegions = nil //countries  current system is not utilizing this value
	// 	subscriptionPlans, _ := JsonStringToIntSliceOrMap(value.SubscriptionPlans)
	// 	movie.SubscriptiontPlans = subscriptionPlans
	// 	if len(subscriptionPlans) < 1 {
	// 		buffer := make([]int, 0)
	// 		movie.SubscriptiontPlans = buffer
	// 	}
	// 	movie.InsertedAt = value.InsertedAt
	// 	movies = append(movies, movie)
	// 	/* Imaginery Details */
	// 	var Imagery ContentImageryDetails

	// 	if value.HasPosterImage == true {
	// 		if contentType.ContentTier == 1 {
	// 			Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/poster-image"
	// 			Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/details-background"
	// 			Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/mobile-details-background"
	// 			Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/poster-image"
	// 			Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + value.ContentVarianceId + "/overlay-poster-image"
	// 		} else {
	// 			Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/poster-image"
	// 			Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/details-background"
	// 			Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/mobile-details-background"
	// 			Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/poster-image"
	// 			Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/overlay-poster-image"
	// 		}
	// 	} else {
	// 		Imagery.Thumbnail = ""
	// 		Imagery.Backdrop = ""
	// 		Imagery.MobileImg = ""
	// 		Imagery.FeaturedImg = ""
	// 		Imagery.Banner = ""
	// 	}
	// 	/* preaparing an object */
	// 	contentFragmentDetails.Id = value.Id
	// 	tags, _ := JsonStringToStringSliceOrMap(value.Tags)
	// 	contentFragmentDetails.Tags = tags
	// 	if len(tags) < 1 {
	// 		buffer := make([]string, 0)
	// 		contentFragmentDetails.Tags = buffer
	// 	}
	// 	contentFragmentDetails.Title = value.Title
	// 	contentFragmentDetails.Length = value.Length
	// 	if contentType.ContentTier == 1 || strings.ToLower(value.ContentType) == "movie" {
	// 		contentFragmentDetails.Movies = movies
	// 	}
	// 	contentFragmentDetails.Imagery = Imagery
	// 	contentFragmentDetails.Geoblock = value.Geoblock
	// 	contentFragmentDetails.VideoId = value.VideoId
	// 	contentFragmentDetails.ContentId = value.ContentId
	// 	contentFragmentDetails.InsertedAt = value.InsertedAt.String()
	// 	contentFragmentDetails.ModifiedAt = value.ModifiedAt.String()
	// 	contentFragmentDetails.ContentType = strings.ToLower(value.ContentType)
	// 	contentFragmentDetails.FriendlyUrl = strings.ToLower(value.FriendlyUrl)
	// 	contentFragmentDetails.ProductionYear = value.ProductionYear
	// 	for i := 1; i <= 2; i++ {
	// 		fmt.Println(i)
	// 		if i == 1 {
	// 			lang = "en"
	// 			cast, _ := JsonStringToStringSliceOrMap(value.EnCast)
	// 			contentFragmentDetails.Cast = cast
	// 			if len(cast) < 1 {
	// 				buffer := make([]string, 0)
	// 				contentFragmentDetails.Cast = buffer
	// 			}
	// 			genres, _ := JsonStringToStringSliceOrMap(value.EnGenres)
	// 			contentFragmentDetails.Genres = genres
	// 			if len(genres) < 1 {
	// 				buffer := make([]string, 0)
	// 				contentFragmentDetails.Genres = buffer
	// 			}
	// 			contentFragmentDetails.Synopsis = value.EnglishSynopsis
	// 			contentFragmentDetails.SeoTitle = value.EnSeoTitle
	// 			contentFragmentDetails.AgeRating = value.EnAgeRating
	// 			contentFragmentDetails.MainActor = value.EnMainActor
	// 			contentFragmentDetails.MainActress = value.EnMainActress
	// 			contentFragmentDetails.SeoDescription = value.EnglishSeoDescription
	// 			contentFragmentDetails.TranslatedTitle = value.ArSeoTitle
	// 			if contentType.ContentTier == 2 {
	// 				if strings.ToLower(value.ContentType) == "series" {
	// 					if err := cdb.Debug().Raw("select cr.digital_rights_type as digital_righttype ,case when ct.language_type = 2 then true else false end as dubbed,false as geoblockgeoblock,s.season_key as id,s.number as season_number,cpi.original_title as seo_title,ac.english_synopsis as seo_description,jsonb_agg(distinct crp.subscription_plan_id)as subscriptiont_plans,cpi.original_title as title from season s join content_rights cr on cr.id =s.rights_id  join content_rights_country crc on crc.content_rights_id =s.rights_id  join content_translation ct on ct.id =s.translation_id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info ac on ac.id = s.about_the_content_info_id  left join content_rights_plan crp on crp.rights_id =s.rights_id  where s.content_id =? group by cr.digital_rights_type,ct.language_type,s.season_key,s.number,crp.subscription_plan_id,cpi.original_title,ac.english_synopsis,cpi.original_title", contentId).Find(&season).Error; err != nil {
	// 						c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 						return
	// 					}
	// 					contentFragmentDetails.Seasons = season
	// 				}
	// 			}
	// 		} else {
	// 			lang = "ar"
	// 			cast, _ := JsonStringToStringSliceOrMap(value.ArCast)
	// 			contentFragmentDetails.Cast = cast
	// 			if len(cast) < 1 {
	// 				buffer := make([]string, 0)
	// 				contentFragmentDetails.Cast = buffer
	// 			}
	// 			genres, _ := JsonStringToStringSliceOrMap(value.ArGenres)
	// 			contentFragmentDetails.Genres = genres
	// 			if len(genres) < 1 {
	// 				buffer := make([]string, 0)
	// 				contentFragmentDetails.Genres = buffer
	// 			}
	// 			contentFragmentDetails.Synopsis = value.ArabicSynopsis
	// 			contentFragmentDetails.SeoTitle = value.ArSeoTitle
	// 			contentFragmentDetails.AgeRating = value.ArAgeRating
	// 			contentFragmentDetails.MainActor = value.ArMainActor
	// 			contentFragmentDetails.MainActress = value.ArMainActress
	// 			contentFragmentDetails.SeoDescription = value.ArabicSeoDescription
	// 			contentFragmentDetails.TranslatedTitle = value.EnSeoTitle
	// 			if contentType.ContentTier == 2 {
	// 				if strings.ToLower(value.ContentType) == "series" {
	// 					if err := cdb.Debug().Raw("select cr.digital_rights_type as digital_righttype , case when ct.language_type = 2 then true else false end as dubbed, false as geoblockgeoblock,s.season_key as id,s.number as season_number,cpi.arabic_title as seo_title,ac.arabic_synopsis as seo_description,jsonb_agg(distinct crp.subscription_plan_id)as subscriptiont_plans,cpi.original_title as title from season s join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = s.rights_id join content_translation ct on ct.id = s.translation_id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info ac on ac.id = s.about_the_content_info_id left join content_rights_plan crp on crp.rights_id = s.rights_id  where s.content_id =? group by cr.digital_rights_type,ct.language_type,s.season_key,s.number,crp.subscription_plan_id,cpi.arabic_title,ac.arabic_synopsis,cpi.original_title", contentId).Find(&season).Error; err != nil {
	// 						c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 						return
	// 					}
	// 					contentFragmentDetails.Seasons = season
	// 				}
	// 			}
	// 		}

	// 		/*mapping content fragment*/
	// 		Response := make(map[string]interface{})
	// 		Response["response_data"] = contentFragmentDetails
	// 		data, _ := json.Marshal(Response)
	// 		type PageResponse struct {
	// 			ResponseData postgres.Jsonb `json:"response_data"`
	// 		}
	// 		var pr PageResponse
	// 		json.Unmarshal(data, &pr)
	// 		contentFragment.ContentId = value.ContentId
	// 		contentFragment.ContentVarianceId = value.ContentVarianceId
	// 		fmt.Println(value.ContentVarianceId)

	// 		varianceIds = append(varianceIds, value.ContentVarianceId)
	// 		contentFragment.Details = pr.ResponseData
	// 		contentFragment.Language = lang
	// 		contentFragment.Country = regionsAndPlatforms.Country
	// 		contentFragment.Platform = regionsAndPlatforms.Platforms
	// 		contentFragment.ContentType = value.ContentType
	// 		contentFragment.ContentKey = value.Id
	// 		contentFragment.RightsStartDate = value.DigitalRightsStartDate
	// 		contentFragment.RightsEndDate = value.DigitalRightsEndDate

	// 		contentFragments = append(contentFragments, contentFragment)
	// 	}
	// }
	// var InsercontentFragment []interface{}
	// for _, content := range contentFragments {
	// 	InsercontentFragment = append(InsercontentFragment, content)
	// }
	// /* Delete All Records With content_variance_id */
	// if err := tx.Debug().Where("content_variance_id in (?)", varianceIds).Delete(&contentFragment).Error; err != nil {
	// 	// if err := tx.Debug().Raw("delete from content_fragment where content_variance_id in (?)",varianceIds).Error; err != nil {
	// 	fmt.Println(err, ">>>>>>")
	// 	tx.Rollback()
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 	return
	// }
	// /* Insertion */
	// if err := gormbulk.BulkInsert(tx.Debug(), InsercontentFragment, 3000); err != nil {
	// 	tx.Rollback()
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
	// 	return
	// }

	// /*Commit Changes*/
	// if err := tx.Commit().Error; err != nil {
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
	// 	return
	// }
}

func CreateContentFragment(contentId string, c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	cdb := c.MustGet("DB").(*gorm.DB)
	ctx := context.Background()
	tx := fdb.Debug().BeginTx(ctx, nil)

	var season []Seasons
	var contentDetails []ContentDetails
	var contentType ContentType
	var regionsAndPlatforms RegionsAndPlatforms
	var missingRights MissingRights

	if contentTypeErr := cdb.Debug().Table("content").Select("content_tier").Where("id=?", contentId).Find(&contentType).Error; contentTypeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}
	fmt.Println(contentType.ContentTier, "????????????????????????????????????????")
	if contentType.ContentTier == 1 {
		if onetiererr := cdb.Debug().Raw("select c.id as content_id,c.content_key as id,json_agg(distinct aaa.english_name)::varchar as en_cast,json_agg(distinct aaa.arabic_name)::varchar as ar_cast,cpi.transliterated_title as title,cpi.arabic_title as arabic_titlen,json_agg(distinct g.english_name)::varchar as en_genres,json_agg(distinct g.arabic_name)::varchar as ar_genres,pi2.duration as length,cpi.original_title as en_seo_title,cpi.arabic_title as ar_seo_title,cpi.transliterated_title as friendly_url,cpi.original_title as title,ac.production_year,ac.english_synopsis ,ac.arabic_synopsis ,c.content_type ,cv.has_all_rights,cv.id as content_variance_id,"+ /*json_agg(distinct crc.country_id)::varchar country,*/ "c.modified_at ,c.has_poster_image,c.created_at as inserted_at,ac.english_synopsis as english_seo_description,ac.arabic_synopsis as arabic_seo_description,ar.english_name as en_age_rating,ar.arabic_name as ar_age_rating,a.english_name as en_main_actor,a.arabic_name as ar_main_actor,aa.english_name as en_main_actress,aa.arabic_name as ar_main_actress,pi2.video_content_id as video_id,json_agg(distinct tdt.name)::varchar as tags,"+ /*json_agg(distinct pitp.target_platform)::varchar as platforms,*/ "false as geoblock,cr.digital_rights_type as digital_right_type,cr.digital_rights_start_date,cr.digital_rights_end_date from content c left join about_the_content_info ac on ac.id = c.about_the_content_info_id left join content_primary_info cpi on cpi.id = c.primary_info_id left join content_variance cv on cv.content_id = c.id left join playback_item pi2 on pi2.id = cv.playback_item_id "+ /*join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id */ "left join content_rights cr on cr.id = pi2.rights_id "+ /*join content_rights_country crc on crc.content_rights_id = cr.id*/ " left join age_ratings ar on ac.age_group = ar.id left join content_cast cc on cc.id = c.cast_id left join actor a on a.id = cc.main_actor_id left join actor aa on aa.id = cc.main_actress_id left join content_actor ca on ca.cast_id = c.cast_id left join actor aaa on aaa.id = ca.actor_id left join content_genre cg on cg.content_id = c.id left join genre g on g.id = cg.genre_id left join content_tag ct on ct.tag_info_id = c.tag_info_id left join textual_data_tag tdt on tdt.id = ct.textual_data_tag_id where c.id =? group by c.id,cpi.transliterated_title,ac.production_year,cv.id,ar.english_name,cpi.arabic_title,ar.arabic_name,a.english_name,a.arabic_name,aa.english_name,aa.arabic_name,pi2.video_content_id,pi2.duration,cr.digital_rights_type,cpi.original_title,ac.english_synopsis ,ac.arabic_synopsis,cr.digital_rights_start_date,cr.digital_rights_end_date,cr.digital_rights_type", contentId).Find(&contentDetails).Error; onetiererr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
			return
		}

		// cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country,json_agg(distinct pitp.target_platform)::varchar as platforms from content c join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where c.id =?", contentId).Find(&regionsAndPlatforms)

	} else {
		if multitiererr := cdb.Debug().Raw("select c.id as content_id,s.id as season_id,c.content_key as id,json_agg(distinct aaa.english_name)::varchar as en_cast,json_agg(distinct aaa.arabic_name)::varchar as ar_cast,cpi2.transliterated_title as title,cpi2.arabic_title as arabic_titlen,json_agg(distinct g.english_name)::varchar as en_genres,json_agg(distinct g.arabic_name)::varchar as ar_genres,sum(pi2.duration) as length,cpi.original_title as en_seo_title,cpi.arabic_title as ar_seo_title,cpi.transliterated_title as friendly_url,cpi.original_title as title,ac.production_year,ac.english_synopsis ,ac.arabic_synopsis ,c.content_type ,s.has_all_rights,s.id as content_variance_id,"+ /*json_agg(distinct crc.country_id)::varchar country,*/ "c.modified_at ,c.has_poster_image,c.created_at as inserted_at,ac.english_synopsis as english_seo_description,ac.arabic_synopsis as arabic_seo_description,ar.english_name as en_age_rating,ar.arabic_name as ar_age_rating,a.english_name as en_main_actor,a.arabic_name as ar_main_actor,aa.english_name as en_main_actress,aa.arabic_name as ar_main_actress,min(pi2.video_content_id) as video_id,json_agg(distinct tdt.name)::varchar as tags,"+ /*json_agg(distinct pitp.target_platform)::varchar as platforms,*/ "false as geoblock,cr.digital_rights_type as digital_right_type,s.season_key as season_id,cr.digital_rights_start_date ,cr.digital_rights_end_date from content c join content_primary_info cpi2 on cpi2.id = c.primary_info_id join season s on s.content_id = c.id left join episode e on e.season_id = s.id left join content_primary_info cpi on cpi.id = s.primary_info_id left join about_the_content_info ac on ac.id = s.about_the_content_info_id left join age_ratings ar on ac.age_group = ar.id left join playback_item pi2 on pi2.id = e.playback_item_id "+ /*join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id */ " left join content_rights cr on cr.id = pi2.rights_id "+ /*join content_rights_country crc on crc.content_rights_id = cr.id*/ " left join content_cast cc on cc.id = s.cast_id left join actor a on a.id= cc.main_actor_id left join actor aa on aa.id = cc.main_actress_id left join content_actor ca on ca.cast_id = c.cast_id left join actor aaa on aaa.id= ca.actor_id left join season_genre sg on sg.season_id =s.id left join genre g on g.id = sg.genre_id left join content_tag_info cti on cti.id = s.tag_info_id left join content_tag ct on ct.tag_info_id =s.tag_info_id left join textual_data_tag tdt on tdt.id=ct.textual_data_tag_id where c.id =? group by c.id,cpi.id,ac.id,cr.digital_rights_type,s.id ,ar.english_name,ar.arabic_name,a.id,aa.id,aaa.id,cr.digital_rights_start_date,cr.digital_rights_end_date,cpi2.original_title,cpi2.arabic_title,cpi2.transliterated_title order by length desc limit 1", contentId).Find(&contentDetails).Error; multitiererr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
			return
		}
		cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country,json_agg(distinct pitp.target_platform)::varchar as platforms from content c join season s on s.content_id = c.id join episode e on e.season_id = s.id join playback_item pi2 on pi2.id = e.playback_item_id join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where c.id = ?", contentId).Find(&regionsAndPlatforms)
		if regionsAndPlatforms.Platforms == "" || regionsAndPlatforms.Country == "" {
			cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country, '[0, 1, 2, 3, 4, 5, 6, 7, 9, 10]' as platform from content c join season s on s.content_id = c.id join content_rights_country crc on crc.content_rights_id = s.rights_id  where c.id = ?", contentId).Find(&regionsAndPlatforms)
			regionsAndPlatforms.Platforms = "[0, 1, 2, 3, 4, 5, 6, 7, 9, 10]"
		}
	}

	var contentFragment ContentFragment
	var contentFragmentDetails ContentFragmentDetails
	var movie Movie
	var movies []Movie
	var contentFragments []ContentFragment
	var varianceIds []string
	lang := "ar"
	for _, value := range contentDetails {
		/* movie details */
		movie.ID = 0
		movie.Title = value.Title
		movie.Geoblock = false
		movie.DigitalRightType = value.DigitalRightType
		if value.DigitalRightType == 0 {
			cdb.Debug().Raw("select cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date from content c join season s on s.content_id = c.id join content_rights cr on cr.id = s.rights_id where c.id = ?", contentId).Find(&missingRights)
			movie.DigitalRightType = missingRights.DigitalRightType
		}
		//countries, _ := JsonStringToIntSliceOrMap(value.Country)
		movie.DigitalRightsRegions = nil //countries  current system is not utilizing this value
		subscriptionPlans, _ := JsonStringToIntSliceOrMap(value.SubscriptionPlans)
		movie.SubscriptiontPlans = subscriptionPlans
		if len(subscriptionPlans) < 1 {
			buffer := make([]int, 0)
			movie.SubscriptiontPlans = buffer
		}
		movie.InsertedAt = value.InsertedAt
		movies = append(movies, movie)
		/* Imaginery Details */
		var Imagery ContentImageryDetails

		// if value.HasPosterImage == true {
		if contentType.ContentTier == 1 {
			Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + value.ContentVarianceId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		} else {
			Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
			Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + value.SeasonId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
		}
		// } else {
		// 	Imagery.Thumbnail = ""
		// 	Imagery.Backdrop = ""
		// 	Imagery.MobileImg = ""
		// 	Imagery.FeaturedImg = ""
		// 	Imagery.Banner = ""
		// }
		/* preaparing an object */
		contentFragmentDetails.Id = value.Id
		tags, _ := JsonStringToStringSliceOrMap(value.Tags)
		contentFragmentDetails.Tags = tags
		if len(tags) < 1 {
			buffer := make([]string, 0)
			contentFragmentDetails.Tags = buffer
		}
		contentFragmentDetails.Title = value.Title
		contentFragmentDetails.Length = value.Length
		if contentType.ContentTier == 1 || strings.ToLower(value.ContentType) == "movie" {
			contentFragmentDetails.Movies = movies
		}
		contentFragmentDetails.Imagery = Imagery
		contentFragmentDetails.Geoblock = value.Geoblock
		contentFragmentDetails.VideoId = value.VideoId
		contentFragmentDetails.ContentId = value.ContentId
		contentFragmentDetails.InsertedAt = value.InsertedAt.String()
		contentFragmentDetails.ModifiedAt = value.ModifiedAt.String()
		contentFragmentDetails.ContentType = strings.ToLower(value.ContentType)
		contentFragmentDetails.FriendlyUrl = strings.ToLower(value.FriendlyUrl)
		contentFragmentDetails.ProductionYear = value.ProductionYear
		for i := 1; i <= 2; i++ {
			fmt.Println(i)
			if i == 1 {
				lang = "en"
				cast, _ := JsonStringToStringSliceOrMap(value.EnCast)
				contentFragmentDetails.Cast = cast
				if len(cast) < 1 {
					buffer := make([]string, 0)
					contentFragmentDetails.Cast = buffer
				}
				genres, _ := JsonStringToStringSliceOrMap(value.EnGenres)
				contentFragmentDetails.Genres = genres
				if len(genres) < 1 {
					buffer := make([]string, 0)
					contentFragmentDetails.Genres = buffer
				}
				contentFragmentDetails.Synopsis = value.EnglishSynopsis
				contentFragmentDetails.SeoTitle = value.EnSeoTitle
				contentFragmentDetails.AgeRating = value.EnAgeRating
				contentFragmentDetails.MainActor = value.EnMainActor
				contentFragmentDetails.MainActress = value.EnMainActress
				contentFragmentDetails.SeoDescription = value.EnglishSeoDescription
				contentFragmentDetails.TranslatedTitle = value.ArSeoTitle
				if contentType.ContentTier == 2 {
					if strings.ToLower(value.ContentType) == "series" {
						if err := cdb.Debug().Raw("select cr.digital_rights_type as digital_righttype ,case when ct.language_type = 2 then true else false end as dubbed,false as geoblockgeoblock,s.season_key as id,s.number as season_number,cpi.original_title as seo_title,ac.english_synopsis as seo_description,jsonb_agg(distinct crp.subscription_plan_id)as subscriptiont_plans,cpi.original_title as title from season s join content_rights cr on cr.id =s.rights_id  join content_rights_country crc on crc.content_rights_id =s.rights_id  join content_translation ct on ct.id =s.translation_id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info ac on ac.id = s.about_the_content_info_id  left join content_rights_plan crp on crp.rights_id =s.rights_id  where s.content_id =? group by cr.digital_rights_type,ct.language_type,s.season_key,s.number,crp.subscription_plan_id,cpi.original_title,ac.english_synopsis,cpi.original_title", contentId).Find(&season).Error; err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
							return
						}
						contentFragmentDetails.Seasons = season
					}
				}
			} else {
				lang = "ar"
				contentFragmentDetails.Title = value.ArabicTitlen
				cast, _ := JsonStringToStringSliceOrMap(value.ArCast)
				contentFragmentDetails.Cast = cast
				if len(cast) < 1 {
					buffer := make([]string, 0)
					contentFragmentDetails.Cast = buffer
				}
				genres, _ := JsonStringToStringSliceOrMap(value.ArGenres)
				contentFragmentDetails.Genres = genres
				if len(genres) < 1 {
					buffer := make([]string, 0)
					contentFragmentDetails.Genres = buffer
				}
				contentFragmentDetails.Synopsis = value.ArabicSynopsis
				contentFragmentDetails.SeoTitle = value.ArSeoTitle
				contentFragmentDetails.AgeRating = value.ArAgeRating
				contentFragmentDetails.MainActor = value.ArMainActor
				contentFragmentDetails.MainActress = value.ArMainActress
				contentFragmentDetails.SeoDescription = value.ArabicSeoDescription
				contentFragmentDetails.TranslatedTitle = value.EnSeoTitle
				if contentType.ContentTier == 2 {
					if strings.ToLower(value.ContentType) == "series" {
						if err := cdb.Debug().Raw("select cr.digital_rights_type as digital_righttype , case when ct.language_type = 2 then true else false end as dubbed, false as geoblockgeoblock,s.season_key as id,s.number as season_number,cpi.arabic_title as seo_title,ac.arabic_synopsis as seo_description,jsonb_agg(distinct crp.subscription_plan_id)as subscriptiont_plans,cpi.original_title as title from season s join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id = s.rights_id join content_translation ct on ct.id = s.translation_id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info ac on ac.id = s.about_the_content_info_id left join content_rights_plan crp on crp.rights_id = s.rights_id  where s.content_id =? group by cr.digital_rights_type,ct.language_type,s.season_key,s.number,crp.subscription_plan_id,cpi.arabic_title,ac.arabic_synopsis,cpi.original_title", contentId).Find(&season).Error; err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
							return
						}
						contentFragmentDetails.Seasons = season
					}
				}
			}

			/*mapping content fragment*/
			Response := make(map[string]interface{})
			Response["response_data"] = contentFragmentDetails
			data, _ := json.Marshal(Response)
			type PageResponse struct {
				ResponseData postgres.Jsonb `json:"response_data"`
			}
			var pr PageResponse
			json.Unmarshal(data, &pr)
			contentFragment.ContentId = value.ContentId
			contentFragment.ContentVarianceId = value.ContentVarianceId
			fmt.Println(value.ContentVarianceId)

			varianceIds = append(varianceIds, value.ContentVarianceId)
			contentFragment.Details = pr.ResponseData
			contentFragment.Language = lang
			if contentType.ContentTier == 1 {
				cdb.Debug().Raw("select json_agg(distinct crc.country_id)::varchar country,json_agg(distinct pitp.target_platform)::varchar as platforms from content c join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id join playback_item_target_platform pitp on pitp.playback_item_id = pi2.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where c.id =?and cv.id = ?", contentId, value.ContentVarianceId).Find(&regionsAndPlatforms)
			}
			contentFragment.Country = regionsAndPlatforms.Country
			contentFragment.Platform = regionsAndPlatforms.Platforms
			contentFragment.ContentType = value.ContentType
			contentFragment.ContentKey = value.Id
			contentFragment.RightsStartDate = value.DigitalRightsStartDate
			contentFragment.RightsEndDate = value.DigitalRightsEndDate
			startdate := value.DigitalRightsStartDate.IsZero()
			enddate := value.DigitalRightsEndDate.IsZero()
			if startdate == true || enddate == true {
				cdb.Debug().Raw("select cr.digital_rights_type,cr.digital_rights_start_date,cr.digital_rights_end_date from content c join season s on s.content_id = c.id join content_rights cr on cr.id = s.rights_id where c.id = ?", contentId).Find(&missingRights)
				contentFragment.RightsStartDate = missingRights.DigitalRightsStartDate
				contentFragment.RightsEndDate = missingRights.DigitalRightsEndDate
			}

			contentFragments = append(contentFragments, contentFragment)
		}
	}
	var InsercontentFragment []interface{}
	for _, content := range contentFragments {
		InsercontentFragment = append(InsercontentFragment, content)
	}
	/* Delete All Records With content_variance_id */
	if err := tx.Debug().Where("content_variance_id in (?)", varianceIds).Delete(&contentFragment).Error; err != nil {
		// if err := tx.Debug().Raw("delete from content_fragment where content_variance_id in (?)",varianceIds).Error; err != nil {
		fmt.Println(err, ">>>>>>")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}
	/* Insertion */
	if err := gormbulk.BulkInsert(tx, InsercontentFragment, 3000); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}

	/*Commit Changes*/
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server-error"})
		return
	}
}

/* Split String to slice(string) */
func JsonStringToStringSliceOrMap(data string) ([]string, error) {
	output := make([]string, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	sort.Strings(output)
	return output, nil
}

/* Split String to slice(int) */
func JsonStringToIntSliceOrMap(data string) ([]int, error) {
	output := make([]int, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	sort.Ints(output)
	return output, nil
}

func OnetierImagery(contentId, varianceId string, ImageryDetails chan ContentImageryDetails) {
	var Imagery ContentImageryDetails
	Imagery.Thumbnail = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Backdrop = os.Getenv("IMAGES") + contentId + "/details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.MobileImg = os.Getenv("IMAGES") + contentId + "/mobile-details-background" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.FeaturedImg = os.Getenv("IMAGES") + contentId + "/poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	Imagery.Banner = os.Getenv("IMAGES") + contentId + "/" + varianceId + "/overlay-poster-image" + "?c=" + strconv.Itoa(rand.Intn(200000))
	ImageryDetails <- Imagery
}

func (hs *HandlerService) CreateContentFragmentsThroughJsonObject(c *gin.Context) {
	var request RequestIds
	c.ShouldBindJSON(&request)
	for _, i := range request.Ids {
		CreateContentFragment(i, c)
		time.Sleep(2 * time.Second)
	}
	fmt.Println("time end...........", time.Now())
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func RemoveContentFragmentSeason(c *gin.Context, seasonID string) {
	
	db := c.MustGet("FDB").(*gorm.DB)

	type DeleteSeason struct {
		ContentVarianceID string `json:"content_variance_id"`
	}

	var deleteSeason DeleteSeason

	if res := db.Debug().Table("content_fragment").Where("content_variance_id=?", seasonID).Delete(&deleteSeason).Error; res != nil {
		log.Println("Fragment not deleted ", seasonID)
	}
}