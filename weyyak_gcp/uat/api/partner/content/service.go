package content

import (
	"encoding/json"
	"fmt"
	"log"
	common "masterdata/common"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const IMAGES = "https://content.weyyak.com/"

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.POST("/oauth2/token", hs.Login)
	srg := r.Group("/v1")
	srg.Use(common.ValidateToken())
	routerOneTier := srg.Group("/contents/onetier")
	routerMultiTier := srg.Group("/contents/multitier")
	routerOneTier.GET("/:contentId", hs.GetOneTierContentDetailsBasedonContentID)
	routerOneTier.GET("/", hs.GetAllOneTierContentDetails)
	routerMultiTier.GET(":contentId", hs.GetMultiTierDetailsBasedonContentID)
	routerMultiTier.GET("/", hs.GetAllMultiTierDetails)
	srg.GET("/episode/:contentId", hs.GetEpisodeDetailsByEpisodeId)
	srg.GET("/get_menu", hs.GetMenuDetails)
	srg.GET("/get_page/:pageId", hs.GetPageDetails)
	srg.GET("/get_info/:videoId", hs.GetVideoDuration)

	srg.GET("/deleted-content", hs.GetDeletedContentDetails)
	srg.GET("/modified-content", hs.GetModifiedContentDetails)
}

// Login -  Login user
// POST /oauth2/token
// @Summary User login with generate token
// @Description User login with generate token
// @Tags Login
// @Accept  multipart/form-data
// @Produce  json
// @Param   username formData string true  "Enter Username"
// @Param   password formData string true  "Enter Password"
// @Success 200 "success"
// @Failure 400 "Bad Request."
// @Failure 401 "Authorization has been denied for this request."
// @Router /oauth2/token [post]
func (hs *HandlerService) Login(c *gin.Context) {
	// GrantType := c.Request.FormValue("grant_type")
	// if needed uncomment below values
	// DeviceID := c.Request.FormValue("deviceId")
	// DeviceName := c.Request.FormValue("deviceName")
	// DevicePlatform := c.Request.FormValue("devicePlatform")
	GrantType := "password"
	DeviceID := "web_app"
	DeviceName := "web_app"
	DevicePlatform := "web_app"
	UserName := c.Request.FormValue("username")
	Password := c.Request.FormValue("password")
	data := url.Values{
		"grant_type":     {GrantType},
		"deviceId":       {DeviceID},
		"deviceName":     {DeviceName},
		"devicePlatform": {DevicePlatform},
		"username":       {UserName},
		"password":       {Password},
	}
	resp, err := http.PostForm(os.Getenv("LOGIN_API"), data)
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetOneTierContentDetailsBasedonContentID
// GET /contents/onetier/:contentId
// @Description Get One Tier Content Details Based on Content ID
// @Tags OneTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param contentId path string true "Content Id."
// @Param Country query string false "Country code of the user."
// @Success 200  {array} OnetierContent
// @Failure 404 "The object was not found."
// @Failure 500  object ErrorResponse "Internal server error."
// @Router /contents/onetier/{contentId} [get]
func (hs *HandlerService) GetOneTierContentDetailsBasedonContentID(c *gin.Context) {
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

	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()

	var finalContentResult []FinalSeasonResultOneTire
	var allContents []AllOnetierContent

	var CountryResult int32

	content_key, _ := strconv.Atoi(c.Param("contentId"))

	var count int
	if err := db.Debug().Table("content").Where("third_party_content_key=?", content_key).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if count < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	}

	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
		fmt.Println(country)
	}

	// type RedisResponse struct {
	// 	Pagination Pagination          `json:"pagination"`
	// 	Data       []AllOnetierContent `json:"data"`
	// }

	// redisKey := os.Getenv("REDIS_CONTENT_KEY") + "_GetAllOneTierContentDetails" + "*offset_" + strconv.FormatInt(offset, 10) + "_limit_" + strconv.FormatInt(limit, 10)

	// val, err := common.GetRedisDataWithKey(redisKey)
	// if err == nil {
	// 	var (
	// 		redisResponse      common.RedisCacheResponse
	// 		finalredisResponse RedisResponse
	// 		redisErrorResponse common.RedisErrorResponse
	// 	)

	// 	Data := []byte(val)
	// 	json.Unmarshal(Data, &redisResponse)
	// 	json.Unmarshal(Data, &redisErrorResponse)
	// 	json.Unmarshal([]byte(redisResponse.Value), &finalredisResponse)
	// 	if redisErrorResponse.Error != "redis: nil" {
	// 		c.JSON(http.StatusOK, gin.H{"pagination": finalredisResponse.Pagination, "data": finalredisResponse.Data})
	// 		return
	// 	}
	// }

	CountryResult = common.Countrys(country)

	if UserId == os.Getenv("WATCH_NOW") {
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`c.watch_now_supplier = 'true'
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.third_party_content_key = ? `, content_key).Find(&finalContentResult)

	} else {
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.third_party_content_key = ? `, content_key).Find(&finalContentResult)

	}

	fmt.Printf("%#v", finalContentResult)

	for _, singleContent := range finalContentResult {

		var contentVariances []ContentVariancesSource

		var ContentVariancesRecord []map[string]interface{}

		if CountryResult != 0 {
			db.Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
						cv.has_subtitling_script, ct.dubbing_language, pi2.rights_id, cv.has_overlay_poster_image from content_variance cv 
						join playback_item pi2 on pi2.id = cv.playback_item_id	
						join content_translation ct on ct.id = pi2.translation_id
						join content_rights cr on cr.id =pi2.rights_id
						join content_rights_country crc on crc.content_rights_id = cr.id
						where 
						( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and
						(cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and
						(cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and 
						crc.country_id = ? and 
						cv.content_id = ?;`, CountryResult, singleContent.Id).Find(&contentVariances)
		} else if country == "" || country == "all" {
			db.Raw(`select cv.id, pi2.duration as length, 
						pi2.video_content_id, ct.language_type, cv.has_dubbing_script, 
						cv.has_subtitling_script, ct.dubbing_dialect_id, ct.dubbing_language, pi2.rights_id, 
						cv.has_overlay_poster_image from content_variance cv 
						join playback_item pi2 on pi2.id = cv.playback_item_id	
						join content_translation ct on ct.id = pi2.translation_id
						join content_rights cr on cr.id =pi2.rights_id
						where 
						( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and
						(cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and
						(cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and 
						cv.content_id = ?;`, singleContent.Id).Find(&contentVariances)
		} else if country != "all" && CountryResult == 0 {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, cv := range contentVariances {
			var dubbingScript string
			var subtitlingScript string
			var OverlayPosterImage string

			if cv.HasDubbingScript {
				dubbingScript = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("DUBBLING_SCRIPT")
			} else {
				dubbingScript = ""
			}

			if cv.HasSubtitlingScript {
				subtitlingScript = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("SUBTITLING_SCRIPT")
			} else {
				subtitlingScript = ""
			}

			var digitalRightsRegions []DigitalRightsRegions
			var trailerInfo []VarianceTrailers

			db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

			var digitalRights []int
			for _, region := range digitalRightsRegions {
				digitalRights = append(digitalRights, region.CountryId)
			}

			db.Table("variance_trailer").Where("content_variance_id=?", cv.Id).Scan(&trailerInfo)

			if cv.HasOverlayPosterImage {
				OverlayPosterImage = os.Getenv("IMAGERY_URL") + singleContent.Id + "/" + cv.Id + "/overlay-poster-image"
			}

			ContentVariancesRecord = append(ContentVariancesRecord, map[string]interface{}{
				"id":                   cv.Id,
				"length":               cv.Length,
				"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
				"languageType":         common.LanguageOriginTypes(cv.LanguageType),
				"dubbingScript":        dubbingScript,
				"subtitlingScript":     subtitlingScript,
				"dubbingLanguage":      cv.DubbingLanguage,
				"dubbingDialectName":   common.DialectIdname(cv.DubbingDialectId, "en"),
				"digitalRightsRegions": digitalRights,
				"trailerInfo":          trailerInfo,
				"overlayPosterImage":   OverlayPosterImage,
			})
		}

		var contentGenres []SeasonGenres
		var finalContentGenre []interface{}
		var newContentGenres NewSeasonGenres
		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
			Joins("join genre g on g.id=cg.genre_id").
			Where("content_id=?", singleContent.Id).Scan(&contentGenres).Error; genreResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, tagInfoIds := range contentGenres {
			var contentSubgenre []SeasonSubgenre
			if subgenreResult := db.Debug().Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
				Joins("join subgenre sg on sg.id=csg.subgenre_id").
				Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			var SubgenreEn []string
			var SubgenreAr []string
			for _, data := range contentSubgenre {
				SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
				SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
				newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
				newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
				newContentGenres.SubGenerEnglish = SubgenreEn
				newContentGenres.SubGenerArabic = SubgenreAr
				newContentGenres.Id = tagInfoIds.Id
			}
			finalContentGenre = append(finalContentGenre, newContentGenres)
		}

		var contentCast Cast
		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
			Joins("left join actor actor on actor.id =cc.main_actor_id").
			Joins("left join actor actress on actress.id =cc.main_actress_id").
			Where("cc.id=?", singleContent.CastId).Scan(&contentCast).Error; castResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var contentActor []ContentActor
		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
			Joins("left join actor a on a.id =ca.actor_id").
			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
			Joins("left join writer w on w.id =cw.writer_id").
			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
			Joins("left join director d on d.id =cd.director_id").
			Where("ca.cast_id=?", singleContent.CastId).Scan(&contentActor).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic, actorId []string
		for _, actorIds := range contentActor {
			actorId = append(actorId, actorIds.ActorId)
			actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
			actorArabic = append(actorArabic, actorIds.ActorArabic)
			writerId = append(writerId, actorIds.WriterId)
			writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
			writerArabic = append(writerArabic, actorIds.WriterArabic)
			directorId = append(directorId, actorIds.DirectorId)
			directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
			directorArabic = append(directorArabic, actorIds.DirectorArabic)
		}

		var contentMusic []ContentMusic
		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
			Joins("left join singer s on s.id =cs.singer_id").
			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
			Where(" cs.music_id=?", singleContent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
		for _, musicIds := range contentMusic {
			singerId = append(singerId, musicIds.SingerIds)
			singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
			singerArabic = append(singerArabic, musicIds.SingersArabic)
			composerId = append(composerId, musicIds.MusicComposerIds)
			composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
			composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
			SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
			SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
			SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
		}

		var contentTagInfo []ContentTag
		db.Table("content_tag ct").Select("tdt.name").
			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
			Where("ct.tag_info_id=?", singleContent.TagInfoId).Scan(&contentTagInfo)
		var tagInfo []string
		for _, tagInfoIds := range contentTagInfo {
			tagInfo = append(tagInfo, tagInfoIds.Name)
		}

		var Tags []string

		Tags = tagInfo
		if len(tagInfo) < 1 {
			buffer := make([]string, 0)
			Tags = buffer
		}

		/*Fetch Production_country*/
		var ProductionCountries []int
		var productionCountry []ProductionCountry
		if productionCountryResult := db.Table("production_country").Select("country_id").Where("about_the_content_info_id=?", singleContent.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var countries []int
		for _, prcountries := range productionCountry {
			countries = append(countries, prcountries.CountryId)
		}
		ProductionCountries = countries
		if len(tagInfo) < 1 {
			buffer := make([]int, 0)
			ProductionCountries = buffer
		}

		/*SeoDetails*/

		var PosterImage string
		var OverlayPosterImage string
		var DetailsBackground string
		var MobileDetailsBackground string

		/*non_textual Data*/
		if singleContent.HasPosterImage {
			PosterImage = IMAGES + singleContent.Id + os.Getenv("POSTER_IMAGE")
		}

		for _, cv := range contentVariances {
			if cv.HasOverlayPosterImage {
				OverlayPosterImage = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("OVERLAY_POSTER_IMAGE")
			}
		}

		if singleContent.HasDetailsBackground {
			DetailsBackground = IMAGES + singleContent.Id + os.Getenv("DETAILS_BACKGROUND")
		}
		if singleContent.HasMobileDetailsBackground {
			MobileDetailsBackground = IMAGES + singleContent.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
		}

		allContents = append(allContents, AllOnetierContent{
			Id:               singleContent.Id,
			CreatedAt:        singleContent.CreatedAt,
			ModifiedAt:       singleContent.ModifiedAt,
			ContentKey:       singleContent.ContentKey,
			ContentVariances: ContentVariancesRecord,
			PrimaryInfo: PrimaryInfo{
				ContentType:         singleContent.ContentType,
				OriginalTitle:       singleContent.OriginalTitle,
				AlternativeTitle:    singleContent.AlternativeTitle,
				ArabicTitle:         singleContent.ArabicTitle,
				TransliteratedTitle: singleContent.TransliteratedTitle,
				Notes:               singleContent.Notes,
			},
			ContentGenres: finalContentGenre,
			Cast: Cast{
				CastId:             singleContent.CastId,
				MainActorId:        contentCast.MainActorId,
				MainActressId:      contentCast.MainActressId,
				MainActorEnglish:   contentCast.MainActorEnglish,
				MainActorArabic:    contentCast.MainActorArabic,
				MainActressEnglish: contentCast.MainActressEnglish,
				MainActressArabic:  contentCast.MainActressArabic,
				ActorIds:           common.RemoveDuplicateValues(actorId),
				ActorEnglish:       common.RemoveDuplicateValues(actorEnglish),
				ActorArabic:        common.RemoveDuplicateValues(actorArabic),
				WriterId:           common.RemoveDuplicateValues(writerId),
				WriterEnglish:      common.RemoveDuplicateValues(writerEnglish),
				WriterArabic:       common.RemoveDuplicateValues(writerArabic),
				DirectorEnglish:    common.RemoveDuplicateValues(directorEnglish),
				DirectorArabic:     common.RemoveDuplicateValues(directorArabic),
				DirectorIds:        common.RemoveDuplicateValues(directorId),
			},
			Music: Music{
				MusicId:               singleContent.MusicId,
				SingerIds:             common.RemoveDuplicateValues(singerId),
				SingersEnglish:        common.RemoveDuplicateValues(singerEnglish),
				SingersArabic:         common.RemoveDuplicateValues(singerArabic),
				MusicComposerIds:      common.RemoveDuplicateValues(composerId),
				MusicComposersEnglish: common.RemoveDuplicateValues(composerEnglish),
				MusicComposersArabic:  common.RemoveDuplicateValues(composerArabic),
				SongWriterIds:         common.RemoveDuplicateValues(SongWriterId),
				SongWritersEnglish:    common.RemoveDuplicateValues(SongWriterEnglish),
				SongWritersArabic:     common.RemoveDuplicateValues(SongWriterArabic),
			},
			TagInfo: TagInfo{
				Tags: Tags,
			},
			AboutTheContent: AboutTheContent{
				OriginalLanguage:      singleContent.OriginalLanguage,
				Supplier:              singleContent.Supplier,
				AcquisitionDepartment: singleContent.AcquisitionDepartment,
				EnglishSynopsis:       singleContent.EnglishSynopsis,
				ArabicSynopsis:        singleContent.ArabicSynopsis,
				ProductionYear:        singleContent.ProductionYear,
				ProductionHouse:       singleContent.ProductionHouse,
				AgeGroup:              common.AgeRatings(singleContent.AgeGroup, "en"),
				ProductionCountries:   ProductionCountries,
			},
			SeoDetails: SeoDetails{
				EnglishMetaTitle:       singleContent.EnglishMetaTitle,
				ArabicMetaTitle:        singleContent.ArabicMetaTitle,
				EnglishMetaDescription: singleContent.EnglishMetaDescription,
				ArabicMetaDescription:  singleContent.ArabicMetaDescription,
			},
			NonTextualData: NonTextualData{
				PosterImage:             PosterImage,
				OverlayPosterImage:      OverlayPosterImage,
				DetailsBackground:       DetailsBackground,
				MobileDetailsBackground: MobileDetailsBackground,
			},
		})

	}

	// pagination := Pagination{
	// 	Size:   totalCount,
	// 	Limit:  int(limit),
	// 	Offset: int(offset),
	// }

	// m, _ := json.Marshal(RedisResponse{
	// 	Pagination: pagination,
	// 	Data:       allContents,
	// })

	// err = common.PostRedisDataWithKey(redisKey, m)
	// if err != nil {
	// 	fmt.Println("Redis Value Not Updated")
	// }

	c.JSON(http.StatusOK, gin.H{"data": allContents})

}

// func (hs *HandlerService) GetOneTierContentDetailsBasedonContentID(c *gin.Context) {
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
// 	var contentResult OnetierContent
// 	serverError := common.ServerErrorResponse()
// 	notFound := common.NotFoundErrorResponse()

// 	var finalContentResults []FinalSeasonResult

// 	var country string
// 	if c.Request.URL.Query()["Country"] != nil {
// 		country = c.Request.URL.Query()["Country"][0]
// 	}
// 	CountryResult := common.Countrys(country)
// 	content_key, _ := strconv.Atoi(c.Param("contentId"))

// 	var count int
// 	if err := db.Debug().Table("content").Where("third_party_content_key=?", content_key).Count(&count).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return
// 	}
// 	if count < 1 {
// 		c.JSON(http.StatusNotFound, notFound)
// 		return
// 	}
// 	/*digital rights*/
// 	if c.Request.URL.Query()["Country"] != nil {
// 		country = c.Request.URL.Query()["Country"][0]
// 	}
// 	CountryResult = common.Countrys(country)
// 	if UserId == os.Getenv("WATCH_NOW") {
// 		if CountryResult != 0 {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where("c.watch_now_supplier = 'true' and c.third_party_content_key = ? and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)and crc.country_id = ?", content_key, CountryResult).Find(&finalContentResults).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country == "" || country == "all" {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Where("c.watch_now_supplier = 'true' and c.third_party_content_key = ? and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) ", content_key).Find(&finalContentResults).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country != "all" && CountryResult == 0 {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 	} else {
// 		if CountryResult != 0 {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where(" c.third_party_content_key = ? and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)and crc.country_id = ?", content_key, CountryResult).Find(&finalContentResults).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country == "" || country == "all" {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Where(" c.third_party_content_key = ? and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) ", content_key).Find(&finalContentResults).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country != "all" && CountryResult == 0 {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 	}
// 	/*content-Data*/
// 	for _, finalContentResult := range finalContentResults {
// 		contentResult.ContentKey = finalContentResult.ContentKey
// 		contentResult.PrimaryInfo.ContentType = finalContentResult.ContentType
// 		contentResult.PrimaryInfo.OriginalTitle = finalContentResult.OriginalTitle
// 		contentResult.PrimaryInfo.AlternativeTitle = finalContentResult.AlternativeTitle
// 		contentResult.PrimaryInfo.ArabicTitle = finalContentResult.ArabicTitle
// 		contentResult.PrimaryInfo.TransliteratedTitle = finalContentResult.TransliteratedTitle
// 		contentResult.PrimaryInfo.Notes = finalContentResult.Notes
// 		/*Fetch content_geners*/
// 		var contentGenres []SeasonGenres
// 		var finalContentGenre []interface{}
// 		var newContentGenres NewSeasonGenres
// 		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
// 			Joins("left join genre g on g.id=cg.genre_id").
// 			Where("content_id=?", finalContentResult.Id).Scan(&contentGenres).Error; genreResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}

// 		for _, tagInfoIds := range contentGenres {
// 			var contentSubgenre []SeasonSubgenre
// 			if subgenreResult := db.Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
// 				Joins("left join subgenre sg on sg.id=csg.subgenre_id").
// 				Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			var SubgenreEn []string
// 			var SubgenreAr []string
// 			for _, data := range contentSubgenre {
// 				SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
// 				SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
// 				newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
// 				newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
// 				newContentGenres.SubGenerEnglish = SubgenreEn
// 				newContentGenres.SubGenerArabic = SubgenreAr
// 				newContentGenres.Id = tagInfoIds.Id
// 			}
// 			finalContentGenre = append(finalContentGenre, newContentGenres)
// 		}
// 		contentResult.ContentGenres = finalContentGenre
// 	}
// 	/*content_variance*/
// 	var ContentVariance []interface{}
// 	var contentVariances ContentVariances

// 	for _, finalContentResult := range finalContentResults {

// 		contentVariances.Length = finalContentResult.Duration
// 		contentUrl := os.Getenv("CONTENT_URL")
// 		contentVariances.VideoContentUrl = contentUrl + finalContentResult.VideoContentId
// 		contentVariances.LanguageType = common.LanguageOriginTypes(finalContentResult.LanguageType)
// 		if finalContentResult.HasDubbingScript {
// 			contentVariances.DubbingScript = IMAGES + finalContentResult.Id + "/" + finalContentResult.VarianceId + os.Getenv("DUBBLING_SCRIPT")
// 		} else {
// 			contentVariances.DubbingScript = ""
// 		}
// 		if finalContentResult.HasSubtitlingScript {
// 			contentVariances.SubtitlingScript = IMAGES + finalContentResult.Id + "/" + finalContentResult.VarianceId + os.Getenv("SUBTITLING_SCRIPT")
// 		} else {
// 			contentVariances.SubtitlingScript = ""
// 		}
// 		contentVariances.DubbingLanguage = finalContentResult.DubbingLanguage
// 		contentVariances.DubbingDialectName = common.DialectIdname(finalContentResult.DubbingDialectId, "en")
// 		/*Fetch Digital_right_Regions*/
// 		var digitalRightsRegions []DigitalRightsRegions
// 		db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalContentResult.RightsId).Scan(&digitalRightsRegions)
// 		var RegionRights []int
// 		for _, idarr := range digitalRightsRegions {
// 			RegionRights = append(RegionRights, idarr.CountryId)
// 		}
// 		/*for digital rights*/
// 		var IsCheck bool
// 		for _, value := range RegionRights {
// 			if CountryResult == int32(value) {
// 				IsCheck = true
// 			}
// 		}
// 		if country == "" || country == "all" {
// 			contentVariances.DigitalRightsRegions = RegionRights
// 		} else if country != "all" {
// 			if IsCheck {
// 				contentVariances.DigitalRightsRegions = nil
// 			} else if !IsCheck {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		}
// 		if len(RegionRights) < 1 {
// 			buffer := make([]int, 0)
// 			contentVariances.DigitalRightsRegions = buffer
// 		}
// 		contentVariances.Id = finalContentResult.VarianceId
// 		contentVariances.DigitalRightsRegions = RegionRights
// 		ContentVariance = append(ContentVariance, contentVariances)
// 	}
// 	contentResult.ContentVariances = ContentVariance
// 	/* Fetch content_cast*/
// 	var contentCast Cast
// 	for _, finalContentResult := range finalContentResults {
// 		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
// 			Joins("left join actor actor on actor.id =cc.main_actor_id").
// 			Joins("left join actor actress on actress.id =cc.main_actress_id").
// 			Where("cc.id=?", finalContentResult.CastId).Scan(&contentCast).Error; castResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		contentResult.Cast.CastId = finalContentResult.CastId
// 		contentResult.Cast.MainActorId = contentCast.MainActorId
// 		contentResult.Cast.MainActressId = contentCast.MainActressId
// 		contentResult.Cast.MainActorEnglish = contentCast.MainActorEnglish
// 		contentResult.Cast.MainActorArabic = contentCast.MainActorArabic
// 		contentResult.Cast.MainActressEnglish = contentCast.MainActressEnglish
// 		contentResult.Cast.MainActressArabic = contentCast.MainActressArabic
// 		var contentActor []ContentActor
// 		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
// 			Joins("left join actor a on a.id =ca.actor_id").
// 			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
// 			Joins("left join writer w on w.id =cw.writer_id").
// 			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
// 			Joins("left join director d on d.id =cd.director_id").
// 			Where("ca.cast_id=?", finalContentResult.CastId).Scan(&contentActor).Error; actorResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var actorId, actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic []string
// 		for _, actorIds := range contentActor {
// 			actorId = append(actorId, actorIds.ActorId)
// 			actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
// 			actorArabic = append(actorArabic, actorIds.ActorArabic)
// 			writerId = append(writerId, actorIds.WriterId)
// 			writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
// 			writerArabic = append(writerArabic, actorIds.WriterArabic)
// 			directorId = append(directorId, actorIds.DirectorId)
// 			directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
// 			directorArabic = append(directorArabic, actorIds.DirectorArabic)
// 		}
// 		contentResult.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
// 		contentResult.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
// 		contentResult.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
// 		contentResult.Cast.WriterId = common.RemoveDuplicateValues(writerId)
// 		contentResult.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
// 		contentResult.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
// 		contentResult.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
// 		contentResult.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
// 		contentResult.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)

// 		/* Fetch content_music*/
// 		var contentMusic []ContentMusic
// 		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
// 			Joins("left join singer s on s.id =cs.singer_id").
// 			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
// 			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
// 			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
// 			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
// 			Where(" cs.music_id=?", finalContentResult.MusicId).Scan(&contentMusic).Error; actorResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}

// 		var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
// 		for _, musicIds := range contentMusic {
// 			singerId = append(singerId, musicIds.SingerIds)
// 			singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
// 			singerArabic = append(singerArabic, musicIds.SingersArabic)
// 			composerId = append(composerId, musicIds.MusicComposerIds)
// 			composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
// 			composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
// 			SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
// 			SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
// 			SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
// 		}
// 		contentResult.Music.MusicId = finalContentResult.MusicId
// 		contentResult.Music.SingerIds = common.RemoveDuplicateValues(singerId)
// 		contentResult.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
// 		contentResult.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
// 		contentResult.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
// 		contentResult.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
// 		contentResult.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
// 		contentResult.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
// 		contentResult.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
// 		contentResult.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
// 		/*Fetch tag_info*/
// 		var contentTagInfo []ContentTag
// 		db.Table("content_tag ct").Select("tdt.name").
// 			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
// 			Where("ct.tag_info_id=?", finalContentResult.TagInfoId).Scan(&contentTagInfo)
// 		var tagInfo []string
// 		for _, tagInfoIds := range contentTagInfo {
// 			tagInfo = append(tagInfo, tagInfoIds.Name)
// 		}
// 		contentResult.TagInfo.Tags = tagInfo
// 		if len(tagInfo) < 1 {
// 			buffer := make([]string, 0)
// 			contentResult.TagInfo.Tags = buffer
// 		}
// 		contentResult.AboutTheContent.OriginalLanguage = finalContentResult.OriginalLanguage
// 		contentResult.AboutTheContent.Supplier = finalContentResult.Supplier
// 		contentResult.AboutTheContent.AcquisitionDepartment = finalContentResult.AcquisitionDepartment
// 		contentResult.AboutTheContent.EnglishSynopsis = finalContentResult.EnglishSynopsis
// 		contentResult.AboutTheContent.ArabicSynopsis = finalContentResult.ArabicSynopsis
// 		contentResult.AboutTheContent.ProductionYear = finalContentResult.ProductionYear
// 		contentResult.AboutTheContent.ProductionHouse = finalContentResult.ProductionHouse
// 		contentResult.AboutTheContent.AgeGroup = common.AgeRatings(finalContentResult.AgeGroup, "en")
// 		/*Fetch Production_country*/
// 		var productionCountry []ProductionCountry
// 		if productionCountryResult := db.Table("production_country").Select("country_id").Where("about_the_content_info_id=?", finalContentResult.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var countries []int
// 		for _, prcountries := range productionCountry {
// 			countries = append(countries, prcountries.CountryId)
// 		}
// 		contentResult.AboutTheContent.ProductionCountries = countries
// 		if len(tagInfo) < 1 {
// 			buffer := make([]int, 0)
// 			contentResult.AboutTheContent.ProductionCountries = buffer
// 		}
// 		/*SeoDetails*/
// 		contentResult.SeoDetails.EnglishMetaTitle = finalContentResult.EnglishMetaTitle
// 		contentResult.SeoDetails.ArabicMetaTitle = finalContentResult.ArabicMetaTitle
// 		contentResult.SeoDetails.EnglishMetaDescription = finalContentResult.EnglishMetaDescription
// 		contentResult.SeoDetails.ArabicMetaDescription = finalContentResult.ArabicMetaDescription
// 		/*non_textual Data*/
// 		if finalContentResult.HasPosterImage {
// 			contentResult.NonTextualData.PosterImage = IMAGES + finalContentResult.Id + os.Getenv("POSTER_IMAGE")
// 		}
// 		if finalContentResult.HasOverlayPosterImage {
// 			contentResult.NonTextualData.OverlayPosterImage = IMAGES + finalContentResult.Id + "/" + finalContentResult.VarianceId + os.Getenv("OVERLAY_POSTER_IMAGE")
// 		}
// 		if finalContentResult.HasDetailsBackground {
// 			contentResult.NonTextualData.DetailsBackground = IMAGES + finalContentResult.Id + os.Getenv("DETAILS_BACKGROUND")
// 		}
// 		if finalContentResult.HasMobileDetailsBackground {
// 			contentResult.NonTextualData.MobileDetailsBackground = IMAGES + finalContentResult.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
// 		}
// 		contentResult.Id = finalContentResult.Id
// 		contentResult.CreatedAt = finalContentResult.CreatedAt
// 		contentResult.ModifiedAt = finalContentResult.ModifiedAt
// 	}
// 	if CountryResult != 0 || (country == "" || country == "all") {
// 		c.JSON(http.StatusOK, gin.H{"data": contentResult})
// 	}
// }

// GetOneTierContentDetailsBasedonContentID - Get All One Tier Contents Details
// GET /contents/onetier
// @Description Get All One Tier Contents Details
// @Tags OneTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param offset query string false "Zero-based index of the first requested item."
// @Param limit query string false "Maximum number of collection items to return for a single request."
// @Param Country query string false "Country code of the user."
// @Success 200  object Response
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /contents/onetier [get]
func (hs *HandlerService) GetAllOneTierContentDetails(c *gin.Context) {
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

	serverError := common.ServerErrorResponse()

	var finalContentResult []FinalSeasonResultOneTire
	var allContents []AllOnetierContent
	var limit, offset int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if offset == 0 {
		offset = 0
	}
	if limit == 0 {
		limit = 5
	}

	var totalCount int
	var CountryResult int32

	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
		fmt.Println(country)
	}

	type RedisResponse struct {
		Pagination Pagination          `json:"pagination"`
		Data       []AllOnetierContent `json:"data"`
	}

	redisKey := os.Getenv("REDIS_CONTENT_KEY") + "_GetAllOneTierContentDetails" + "*offset_" + strconv.FormatInt(offset, 10) + "_limit_" + strconv.FormatInt(limit, 10)

	val, err := common.GetRedisDataWithKey(redisKey)
	if err == nil {
		var (
			redisResponse      common.RedisCacheResponse
			finalredisResponse RedisResponse
			redisErrorResponse common.RedisErrorResponse
		)

		Data := []byte(val)
		json.Unmarshal(Data, &redisResponse)
		json.Unmarshal(Data, &redisErrorResponse)
		json.Unmarshal([]byte(redisResponse.Value), &finalredisResponse)
		if redisErrorResponse.Error != "redis: nil" {
			c.JSON(http.StatusOK, gin.H{"pagination": finalredisResponse.Pagination, "data": finalredisResponse.Data})
			return
		}
	}

	CountryResult = common.Countrys(country)

	if UserId == os.Getenv("WATCH_NOW") {
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`c.watch_now_supplier = 'true'
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.id IS NOT NULL`).Limit(limit).Offset(offset).Find(&finalContentResult)

		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`c.watch_now_supplier = 'true'
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			c.id IS NOT NULL`).Count(&totalCount)
	} else {
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.id IS NOT NULL`).Limit(limit).Offset(offset).Find(&finalContentResult)

		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			c.id IS NOT NULL`).Count(&totalCount)
	}

	for _, singleContent := range finalContentResult {

		var contentVariances []ContentVariancesSource

		var ContentVariancesRecord []map[string]interface{}

		if CountryResult != 0 {
			db.Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
						cv.has_subtitling_script, ct.dubbing_language, pi2.rights_id, cv.has_overlay_poster_image from content_variance cv 
						join playback_item pi2 on pi2.id = cv.playback_item_id	
						join content_translation ct on ct.id = pi2.translation_id
						join content_rights cr on cr.id =pi2.rights_id
						join content_rights_country crc on crc.content_rights_id = cr.id
						where 
						( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and
						(cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and
						(cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and 
						crc.country_id = ? and 
						cv.content_id = ?;`, CountryResult, singleContent.Id).Find(&contentVariances)
		} else if country == "" || country == "all" {
			db.Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
						cv.has_subtitling_script, ct.dubbing_language, pi2.rights_id, cv.has_overlay_poster_image from content_variance cv 
						join playback_item pi2 on pi2.id = cv.playback_item_id	
						join content_translation ct on ct.id = pi2.translation_id
						join content_rights cr on cr.id =pi2.rights_id
						where 
						( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and
						(cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and
						(cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and 
						cv.content_id = ?;`, singleContent.Id).Find(&contentVariances)
		} else if country != "all" && CountryResult == 0 {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, cv := range contentVariances {
			var dubbingScript string
			var subtitlingScript string

			if cv.HasDubbingScript {
				dubbingScript = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("DUBBLING_SCRIPT")
			} else {
				dubbingScript = ""
			}

			if cv.HasSubtitlingScript {
				subtitlingScript = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("SUBTITLING_SCRIPT")
			} else {
				subtitlingScript = ""
			}

			var digitalRightsRegions []DigitalRightsRegions
			var trailerInfo []VarianceTrailers

			db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

			var digitalRights []int
			for _, region := range digitalRightsRegions {
				digitalRights = append(digitalRights, region.CountryId)
			}

			db.Table("variance_trailer").Where("content_variance_id=?", cv.Id).Scan(&trailerInfo)

			ContentVariancesRecord = append(ContentVariancesRecord, map[string]interface{}{
				"id":                   cv.Id,
				"length":               cv.Length,
				"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
				"languageType":         common.LanguageOriginTypes(cv.LanguageType),
				"dubbingScript":        dubbingScript,
				"subtitlingScript":     subtitlingScript,
				"dubbingLanguage":      cv.DubbingLanguage,
				"dubbingDialectName":   common.DialectIdname(cv.DubbingDialectId, "en"),
				"digitalRightsRegions": digitalRights,
				"trailerInfo":          trailerInfo,
			})
		}

		var contentGenres []SeasonGenres
		var finalContentGenre []interface{}
		var newContentGenres NewSeasonGenres
		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
			Joins("join genre g on g.id=cg.genre_id").
			Where("content_id=?", singleContent.Id).Scan(&contentGenres).Error; genreResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, tagInfoIds := range contentGenres {
			var contentSubgenre []SeasonSubgenre
			if subgenreResult := db.Debug().Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
				Joins("join subgenre sg on sg.id=csg.subgenre_id").
				Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			var SubgenreEn []string
			var SubgenreAr []string
			for _, data := range contentSubgenre {
				SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
				SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
				newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
				newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
				newContentGenres.SubGenerEnglish = SubgenreEn
				newContentGenres.SubGenerArabic = SubgenreAr
				newContentGenres.Id = tagInfoIds.Id
			}
			finalContentGenre = append(finalContentGenre, newContentGenres)
		}

		var contentCast Cast
		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
			Joins("left join actor actor on actor.id =cc.main_actor_id").
			Joins("left join actor actress on actress.id =cc.main_actress_id").
			Where("cc.id=?", singleContent.CastId).Scan(&contentCast).Error; castResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var contentActor []ContentActor
		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
			Joins("left join actor a on a.id =ca.actor_id").
			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
			Joins("left join writer w on w.id =cw.writer_id").
			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
			Joins("left join director d on d.id =cd.director_id").
			Where("ca.cast_id=?", singleContent.CastId).Scan(&contentActor).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic, actorId []string
		for _, actorIds := range contentActor {
			actorId = append(actorId, actorIds.ActorId)
			actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
			actorArabic = append(actorArabic, actorIds.ActorArabic)
			writerId = append(writerId, actorIds.WriterId)
			writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
			writerArabic = append(writerArabic, actorIds.WriterArabic)
			directorId = append(directorId, actorIds.DirectorId)
			directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
			directorArabic = append(directorArabic, actorIds.DirectorArabic)
		}

		var contentMusic []ContentMusic
		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
			Joins("left join singer s on s.id =cs.singer_id").
			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
			Where(" cs.music_id=?", singleContent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
		for _, musicIds := range contentMusic {
			singerId = append(singerId, musicIds.SingerIds)
			singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
			singerArabic = append(singerArabic, musicIds.SingersArabic)
			composerId = append(composerId, musicIds.MusicComposerIds)
			composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
			composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
			SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
			SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
			SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
		}

		var contentTagInfo []ContentTag
		db.Table("content_tag ct").Select("tdt.name").
			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
			Where("ct.tag_info_id=?", singleContent.TagInfoId).Scan(&contentTagInfo)
		var tagInfo []string
		for _, tagInfoIds := range contentTagInfo {
			tagInfo = append(tagInfo, tagInfoIds.Name)
		}

		var Tags []string

		Tags = tagInfo
		if len(tagInfo) < 1 {
			buffer := make([]string, 0)
			Tags = buffer
		}

		/*Fetch Production_country*/
		var ProductionCountries []int
		var productionCountry []ProductionCountry
		if productionCountryResult := db.Table("production_country").Select("country_id").Where("about_the_content_info_id=?", singleContent.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var countries []int
		for _, prcountries := range productionCountry {
			countries = append(countries, prcountries.CountryId)
		}
		ProductionCountries = countries
		if len(tagInfo) < 1 {
			buffer := make([]int, 0)
			ProductionCountries = buffer
		}

		/*SeoDetails*/

		var PosterImage string
		var OverlayPosterImage string
		var DetailsBackground string
		var MobileDetailsBackground string

		/*non_textual Data*/
		if singleContent.HasPosterImage {
			PosterImage = IMAGES + singleContent.Id + os.Getenv("POSTER_IMAGE")
		}

		for _, cv := range contentVariances {
			if cv.HasOverlayPosterImage {
				OverlayPosterImage = IMAGES + singleContent.Id + "/" + cv.Id + os.Getenv("OVERLAY_POSTER_IMAGE")
			}
		}

		if singleContent.HasDetailsBackground {
			DetailsBackground = IMAGES + singleContent.Id + os.Getenv("DETAILS_BACKGROUND")
		}
		if singleContent.HasMobileDetailsBackground {
			MobileDetailsBackground = IMAGES + singleContent.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
		}

		allContents = append(allContents, AllOnetierContent{
			Id:               singleContent.Id,
			CreatedAt:        singleContent.CreatedAt,
			ModifiedAt:       singleContent.ModifiedAt,
			ContentKey:       singleContent.ContentKey,
			ContentVariances: ContentVariancesRecord,
			PrimaryInfo: PrimaryInfo{
				ContentType:         singleContent.ContentType,
				OriginalTitle:       singleContent.OriginalTitle,
				AlternativeTitle:    singleContent.AlternativeTitle,
				ArabicTitle:         singleContent.ArabicTitle,
				TransliteratedTitle: singleContent.TransliteratedTitle,
				Notes:               singleContent.Notes,
			},
			ContentGenres: finalContentGenre,
			Cast: Cast{
				CastId:             singleContent.CastId,
				MainActorId:        contentCast.MainActorId,
				MainActressId:      contentCast.MainActressId,
				MainActorEnglish:   contentCast.MainActorEnglish,
				MainActorArabic:    contentCast.MainActorArabic,
				MainActressEnglish: contentCast.MainActressEnglish,
				MainActressArabic:  contentCast.MainActressArabic,
				ActorIds:           common.RemoveDuplicateValues(actorId),
				ActorEnglish:       common.RemoveDuplicateValues(actorEnglish),
				ActorArabic:        common.RemoveDuplicateValues(actorArabic),
				WriterId:           common.RemoveDuplicateValues(writerId),
				WriterEnglish:      common.RemoveDuplicateValues(writerEnglish),
				WriterArabic:       common.RemoveDuplicateValues(writerArabic),
				DirectorEnglish:    common.RemoveDuplicateValues(directorEnglish),
				DirectorArabic:     common.RemoveDuplicateValues(directorArabic),
				DirectorIds:        common.RemoveDuplicateValues(directorId),
			},
			Music: Music{
				MusicId:               singleContent.MusicId,
				SingerIds:             common.RemoveDuplicateValues(singerId),
				SingersEnglish:        common.RemoveDuplicateValues(singerEnglish),
				SingersArabic:         common.RemoveDuplicateValues(singerArabic),
				MusicComposerIds:      common.RemoveDuplicateValues(composerId),
				MusicComposersEnglish: common.RemoveDuplicateValues(composerEnglish),
				MusicComposersArabic:  common.RemoveDuplicateValues(composerArabic),
				SongWriterIds:         common.RemoveDuplicateValues(SongWriterId),
				SongWritersEnglish:    common.RemoveDuplicateValues(SongWriterEnglish),
				SongWritersArabic:     common.RemoveDuplicateValues(SongWriterArabic),
			},
			TagInfo: TagInfo{
				Tags: Tags,
			},
			AboutTheContent: AboutTheContent{
				OriginalLanguage:      singleContent.OriginalLanguage,
				Supplier:              singleContent.Supplier,
				AcquisitionDepartment: singleContent.AcquisitionDepartment,
				EnglishSynopsis:       singleContent.EnglishSynopsis,
				ArabicSynopsis:        singleContent.ArabicSynopsis,
				ProductionYear:        singleContent.ProductionYear,
				ProductionHouse:       singleContent.ProductionHouse,
				AgeGroup:              common.AgeRatings(singleContent.AgeGroup, "en"),
				ProductionCountries:   ProductionCountries,
			},
			SeoDetails: SeoDetails{
				EnglishMetaTitle:       singleContent.EnglishMetaTitle,
				ArabicMetaTitle:        singleContent.ArabicMetaTitle,
				EnglishMetaDescription: singleContent.EnglishMetaDescription,
				ArabicMetaDescription:  singleContent.ArabicMetaDescription,
			},
			NonTextualData: NonTextualData{
				PosterImage:             PosterImage,
				OverlayPosterImage:      OverlayPosterImage,
				DetailsBackground:       DetailsBackground,
				MobileDetailsBackground: MobileDetailsBackground,
			},
		})

	}

	pagination := Pagination{
		Size:   totalCount,
		Limit:  int(limit),
		Offset: int(offset),
	}

	m, _ := json.Marshal(RedisResponse{
		Pagination: pagination,
		Data:       allContents,
	})

	err = common.PostRedisDataWithKey(redisKey, m)
	if err != nil {
		fmt.Println("Redis Value Not Updated")
	}

	c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})

}

// func (hs *HandlerService) GetAllOneTierContentDetailsSS(c *gin.Context) {
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

// 	var contentResult AllOnetierContent
// 	var allContents []AllOnetierContent
// 	serverError := common.ServerErrorResponse()
// 	var finalContentResult []FinalSeasonResult
// 	var limit, offset int64
// 	if c.Request.URL.Query()["limit"] != nil {
// 		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
// 	}
// 	if c.Request.URL.Query()["offset"] != nil {
// 		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
// 	}
// 	if offset == 0 {
// 		offset = 0
// 	}
// 	if limit == 0 {
// 		limit = 5
// 	}

// 	var totalCount int
// 	var CountryResult int32
// 	/*digital rights*/
// 	var country string
// 	if c.Request.URL.Query()["Country"] != nil {
// 		country = c.Request.URL.Query()["Country"][0]
// 		fmt.Println(country)
// 	}
// 	CountryResult = common.Countrys(country)
// 	if UserId == os.Getenv("WATCH_NOW") {
// 		if CountryResult != 0 {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where("c.watch_now_supplier = 'true' and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			if err := db.Debug().Table("content c").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where("c.watch_now_supplier = 'true' and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Count(&totalCount).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country == "" || country == "all" {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
// 				Where("c.watch_now_supplier = 'true' and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			if err := db.Debug().Table("content c").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
// 				Where("c.watch_now_supplier = 'true' and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Count(&totalCount).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country != "all" && CountryResult == 0 {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 	} else {
// 		if CountryResult != 0 {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where(" c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			if err := db.Debug().Table("content c").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
// 				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 				Where(" c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Count(&totalCount).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country == "" || country == "all" {
// 			if err := db.Debug().Table("content c").Select("c.id ,c.third_party_content_key as content_key , c.status,c.content_type ,c.english_meta_title ,c.arabic_meta_title ,c.english_meta_description ,c.arabic_meta_description ,c.created_at,c.modified_at,c.has_poster_image ,c.has_details_background ,c.has_mobile_details_background ,c.cast_id,c.music_id,c.tag_info_id,c.about_the_content_info_id,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.outro_start ,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language,pi2.id as playback_item_id,pi2.video_content_id ,pi2.rights_id,pi2.scheduling_date_time,pi2.duration,cv.id as variance_id,cv.status as variance_status,cv.has_overlay_poster_image ,cv.has_dubbing_script ,cv.has_subtitling_script ,cr.digital_rights_type ,cr.digital_rights_start_date ,cr.digital_rights_end_date ,atci.original_language ,atci.supplier ,atci.acquisition_department ,atci.english_synopsis ,atci.arabic_synopsis ,atci.production_year ,atci.production_house ,atci.age_group ,atci.outro_start as about_outro_start").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
// 				Where("c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			if err := db.Table("content c").
// 				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
// 				Joins("join content_variance cv on cv.content_id =c.id").
// 				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
// 				Joins("join content_translation ct on ct.id =pi2.translation_id").
// 				Joins("join content_rights cr on cr.id =pi2.rights_id").
// 				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
// 				Where(" c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Count(&totalCount).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 		} else if country != "all" && CountryResult == 0 {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 	}

// 	for _, singleContent := range finalContentResult {

// 		/*content-Data*/
// 		contentResult.ContentKey = singleContent.ContentKey
// 		/*Textual-Data*/
// 		contentResult.PrimaryInfo.ContentType = singleContent.ContentType
// 		contentResult.PrimaryInfo.OriginalTitle = singleContent.OriginalTitle
// 		contentResult.PrimaryInfo.AlternativeTitle = singleContent.AlternativeTitle
// 		contentResult.PrimaryInfo.ArabicTitle = singleContent.ArabicTitle
// 		contentResult.PrimaryInfo.TransliteratedTitle = singleContent.TransliteratedTitle
// 		contentResult.PrimaryInfo.Notes = singleContent.Notes
// 		/*Fetch content_geners*/
// 		var contentGenres []SeasonGenres
// 		var finalContentGenre []interface{}
// 		var newContentGenres NewSeasonGenres
// 		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
// 			Joins("join genre g on g.id=cg.genre_id").
// 			Where("content_id=?", singleContent.Id).Scan(&contentGenres).Error; genreResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		for _, tagInfoIds := range contentGenres {
// 			var contentSubgenre []SeasonSubgenre
// 			if subgenreResult := db.Debug().Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
// 				Joins("join subgenre sg on sg.id=csg.subgenre_id").
// 				Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
// 				c.JSON(http.StatusInternalServerError, serverError)
// 				return
// 			}
// 			var SubgenreEn []string
// 			var SubgenreAr []string
// 			for _, data := range contentSubgenre {
// 				SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
// 				SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
// 				newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
// 				newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
// 				newContentGenres.SubGenerEnglish = SubgenreEn
// 				newContentGenres.SubGenerArabic = SubgenreAr
// 				newContentGenres.Id = tagInfoIds.Id

// 			}
// 			finalContentGenre = append(finalContentGenre, newContentGenres)

// 		}
// 		contentResult.ContentGenres = finalContentGenre
// 		/*content_variance*/
// 		var ContentVariance []interface{}
// 		var contentVariances ContentVariances
// 		contentVariances.Length = singleContent.Duration
// 		contentVariances.VideoContentUrl = os.Getenv("CONTENT_URL") + singleContent.VideoContentId
// 		contentVariances.LanguageType = common.LanguageOriginTypes(singleContent.LanguageType)
// 		if singleContent.HasDubbingScript {
// 			contentVariances.DubbingScript = IMAGES + singleContent.Id + "/" + singleContent.VarianceId + os.Getenv("DUBBLING_SCRIPT")
// 		} else {
// 			contentVariances.DubbingScript = ""
// 		}
// 		if singleContent.HasSubtitlingScript {
// 			contentVariances.SubtitlingScript = IMAGES + singleContent.Id + "/" + singleContent.VarianceId + os.Getenv("SUBTITLING_SCRIPT")
// 		} else {
// 			contentVariances.SubtitlingScript = ""
// 		}
// 		contentVariances.DubbingLanguage = singleContent.DubbingLanguage

// 		contentVariances.DubbingDialectName = common.DialectIdname(singleContent.DubbingDialectId, "en")
// 		/*Fetch Digital_right_Regions*/
// 		var digitalRightsRegions []DigitalRightsRegions

// 		db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", singleContent.RightsId).Scan(&digitalRightsRegions)

// 		var digitalRights []int
// 		for _, region := range digitalRightsRegions {
// 			digitalRights = append(digitalRights, region.CountryId)
// 		}
// 		contentVariances.DigitalRightsRegions = digitalRights
// 		contentVariances.Id = singleContent.VarianceId

// 		ContentVariance = append(ContentVariance, contentVariances)
// 		// contentResult.ContentVariances = ContentVariance
// 		/* Fetch content_cast*/
// 		var contentCast Cast
// 		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
// 			Joins("left join actor actor on actor.id =cc.main_actor_id").
// 			Joins("left join actor actress on actress.id =cc.main_actress_id").
// 			Where("cc.id=?", singleContent.CastId).Scan(&contentCast).Error; castResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		contentResult.Cast.CastId = singleContent.CastId
// 		contentResult.Cast.MainActorId = contentCast.MainActorId
// 		contentResult.Cast.MainActressId = contentCast.MainActressId
// 		contentResult.Cast.MainActorEnglish = contentCast.MainActorEnglish
// 		contentResult.Cast.MainActorArabic = contentCast.MainActorArabic
// 		contentResult.Cast.MainActressEnglish = contentCast.MainActressEnglish
// 		contentResult.Cast.MainActressArabic = contentCast.MainActressArabic
// 		var contentActor []ContentActor
// 		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
// 			Joins("left join actor a on a.id =ca.actor_id").
// 			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
// 			Joins("left join writer w on w.id =cw.writer_id").
// 			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
// 			Joins("left join director d on d.id =cd.director_id").
// 			Where("ca.cast_id=?", singleContent.CastId).Scan(&contentActor).Error; actorResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic, actorId []string
// 		for _, actorIds := range contentActor {
// 			actorId = append(actorId, actorIds.ActorId)
// 			actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
// 			actorArabic = append(actorArabic, actorIds.ActorArabic)
// 			writerId = append(writerId, actorIds.WriterId)
// 			writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
// 			writerArabic = append(writerArabic, actorIds.WriterArabic)
// 			directorId = append(directorId, actorIds.DirectorId)
// 			directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
// 			directorArabic = append(directorArabic, actorIds.DirectorArabic)
// 		}

// 		contentResult.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
// 		contentResult.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
// 		contentResult.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
// 		contentResult.Cast.WriterId = common.RemoveDuplicateValues(writerId)
// 		contentResult.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
// 		contentResult.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
// 		contentResult.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
// 		contentResult.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)
// 		contentResult.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)

// 		/* Fetch content_music*/
// 		var contentMusic []ContentMusic
// 		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
// 			Joins("left join singer s on s.id =cs.singer_id").
// 			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
// 			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
// 			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
// 			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
// 			Where(" cs.music_id=?", singleContent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}

// 		var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
// 		for _, musicIds := range contentMusic {
// 			singerId = append(singerId, musicIds.SingerIds)
// 			singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
// 			singerArabic = append(singerArabic, musicIds.SingersArabic)
// 			composerId = append(composerId, musicIds.MusicComposerIds)
// 			composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
// 			composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
// 			SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
// 			SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
// 			SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
// 		}
// 		contentResult.Music.MusicId = singleContent.MusicId
// 		contentResult.Music.SingerIds = common.RemoveDuplicateValues(singerId)
// 		contentResult.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
// 		contentResult.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
// 		contentResult.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
// 		contentResult.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
// 		contentResult.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
// 		contentResult.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
// 		contentResult.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
// 		contentResult.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
// 		/*Fetch tag_info*/
// 		var contentTagInfo []ContentTag
// 		db.Table("content_tag ct").Select("tdt.name").
// 			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
// 			Where("ct.tag_info_id=?", singleContent.TagInfoId).Scan(&contentTagInfo)
// 		var tagInfo []string
// 		for _, tagInfoIds := range contentTagInfo {
// 			tagInfo = append(tagInfo, tagInfoIds.Name)
// 		}
// 		contentResult.TagInfo.Tags = tagInfo
// 		if len(tagInfo) < 1 {
// 			buffer := make([]string, 0)
// 			contentResult.TagInfo.Tags = buffer
// 		}
// 		contentResult.AboutTheContent.OriginalLanguage = singleContent.OriginalLanguage
// 		contentResult.AboutTheContent.Supplier = singleContent.Supplier
// 		contentResult.AboutTheContent.AcquisitionDepartment = singleContent.AcquisitionDepartment
// 		contentResult.AboutTheContent.EnglishSynopsis = singleContent.EnglishSynopsis
// 		contentResult.AboutTheContent.ArabicSynopsis = singleContent.ArabicSynopsis
// 		contentResult.AboutTheContent.ProductionYear = singleContent.ProductionYear
// 		contentResult.AboutTheContent.ProductionHouse = singleContent.ProductionHouse
// 		contentResult.AboutTheContent.AgeGroup = common.AgeRatings(singleContent.AgeGroup, "en")
// 		/*Fetch Production_country*/
// 		var productionCountry []ProductionCountry
// 		if productionCountryResult := db.Table("production_country").Select("country_id").Where("about_the_content_info_id=?", singleContent.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var countries []int
// 		for _, prcountries := range productionCountry {
// 			countries = append(countries, prcountries.CountryId)
// 		}
// 		contentResult.AboutTheContent.ProductionCountries = countries
// 		if len(tagInfo) < 1 {
// 			buffer := make([]int, 0)
// 			contentResult.AboutTheContent.ProductionCountries = buffer
// 		}
// 		/*SeoDetails*/
// 		contentResult.SeoDetails.EnglishMetaTitle = singleContent.EnglishMetaTitle
// 		contentResult.SeoDetails.ArabicMetaTitle = singleContent.ArabicMetaTitle
// 		contentResult.SeoDetails.EnglishMetaDescription = singleContent.EnglishMetaDescription
// 		contentResult.SeoDetails.ArabicMetaDescription = singleContent.ArabicMetaDescription
// 		/*non_textual Data*/
// 		if singleContent.HasPosterImage {
// 			contentResult.NonTextualData.PosterImage = IMAGES + singleContent.Id + os.Getenv("POSTER_IMAGE")
// 		}
// 		if singleContent.HasOverlayPosterImage {
// 			contentResult.NonTextualData.OverlayPosterImage = IMAGES + singleContent.Id + "/" + singleContent.VarianceId + os.Getenv("OVERLAY_POSTER_IMAGE")
// 		}
// 		if singleContent.HasDetailsBackground {
// 			contentResult.NonTextualData.DetailsBackground = IMAGES + singleContent.Id + os.Getenv("DETAILS_BACKGROUND")
// 		}
// 		if singleContent.HasMobileDetailsBackground {
// 			contentResult.NonTextualData.MobileDetailsBackground = IMAGES + singleContent.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
// 		}
// 		contentResult.Id = singleContent.Id
// 		contentResult.CreatedAt = singleContent.CreatedAt
// 		contentResult.ModifiedAt = singleContent.ModifiedAt
// 		allContents = append(allContents, contentResult)
// 	}
// 	/*Pagination*/
// 	var pagination Pagination
// 	pagination.Limit = int(limit)
// 	pagination.Offset = int(offset)
// 	pagination.Size = totalCount
// 	var final Response
// 	final.Pagination = pagination
// 	final.Data = allContents

// 	if CountryResult != 0 || (country == "" || country == "all") {
// 		c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})
// 	}
// }

// GetEpisodeDetailsByEpisodeId - Get Episode Details By Episode Id
// GET /episode/:contentId
// @Description Get Episode Details By Episode Id
// @Tags MultiTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param contentId path string true "Content Id"
// @Param country path string false "Country code of the user"
// @Success 200  object EpisodeResult
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /episode/{contentId} [get]
func (hs *HandlerService) GetEpisodeDetailsByEpisodeId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	episode_key, _ := strconv.Atoi(c.Param("contentId"))
	var count int
	var CountryResult int32
	if err := db.Table("episode").Where("third_party_episode_key=?", episode_key).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if count < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	}
	type ContentId struct {
		ContentId string
	}
	var watchnowsupplier ContentId
	if UserId == os.Getenv("WATCH_NOW") {
		if watchnowResult := db.Debug().Table("content as c").Select("c.id as content_id").Joins("left join season s on s.content_id = c.id").Joins("left join episode e on e.season_id = s.id").Where("c.watch_now_supplier = 'true' and e.third_party_episode_key =?", episode_key).Find(&watchnowsupplier).Error; watchnowResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	var finalContentResult FinalSeasonResult
	var contentResult EpisodeResult
	if episodeResult := db.Debug().Table("episode as e").Select("e.number as episode_number ,s.content_id ,e.third_party_episode_key as episode_key ,pi2.duration,pi2.video_content_id,e.synopsis_english ,e.synopsis_arabic,e.season_id,e.has_poster_image ,e.has_dubbing_script ,e.has_subtitling_script ,e.english_meta_title ,e.arabic_meta_title ,e.english_meta_description ,e.arabic_meta_description ,e.created_at ,e.modified_at ,e.id ,cpi.original_title ,cpi.alternative_title ,cpi.arabic_title ,cpi.transliterated_title ,cpi.notes ,cpi.intro_start ,cpi.outro_start,ct.language_type ,ct.dubbing_language ,ct.dubbing_dialect_id ,ct.subtitling_language ,e.cast_id ,e.music_id ,e.tag_info_id ,s.rights_id").
		Joins("left join season s on e.season_id =s.id").
		Joins("left join playback_item pi2 on pi2.id =e.playback_item_id").
		Joins("left join content_primary_info cpi on cpi.id =e.primary_info_id").
		Joins("left join content_translation ct on ct.id =pi2.translation_id").
		Where("e.deleted_by_user_id is null and e.status = 1 and e.third_party_episode_key =?", episode_key).Find(&finalContentResult).Error; episodeResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	}

	if country == "" || country == "all" {
		CountryResult = 1
	} else if country != "" {
		CountryResult = common.Countrys(country)
	}

	contentResult.EpisodeNumber = finalContentResult.EpisodeNumber
	contentResult.ContentKey = finalContentResult.ContentId
	contentResult.EpisodeKey = finalContentResult.EpisodeKey
	contentResult.Length = finalContentResult.Duration
	contentResult.VideoContentUrl = os.Getenv("CONTENT_URL") + finalContentResult.VideoContentId
	contentResult.SynopsisEnglish = finalContentResult.SynopsisEnglish
	contentResult.SynopsisArabic = finalContentResult.SynopsisArabic
	contentResult.SeasonId = finalContentResult.SeasonId
	/*Fetch Digital_right_Regions*/
	var digitalRightsRegions []DigitalRightsRegions
	if countryError := db.Debug().Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalContentResult.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var RegionRights []int
	for _, idarr := range digitalRightsRegions {
		RegionRights = append(RegionRights, idarr.CountryId)
	}
	/*for digital rights*/
	var IsCheck bool
	for _, value := range RegionRights {
		if CountryResult == int32(value) {
			IsCheck = true
		}
	}
	if country == "" || country == "all" {
		contentResult.DigitalRightsRegions = RegionRights
	} else if country != "all" {
		if IsCheck {
			contentResult.DigitalRightsRegions = nil
		} else if !IsCheck {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	if len(RegionRights) < 1 {
		buffer := make([]int, 0)
		contentResult.DigitalRightsRegions = buffer
	}
	contentResult.PrimaryInfo.Number = finalContentResult.EpisodeNumber
	contentResult.PrimaryInfo.VideoContentId = finalContentResult.VideoContentId
	contentResult.PrimaryInfo.SynopsisEnglish = finalContentResult.SynopsisEnglish
	contentResult.PrimaryInfo.SynopsisArabic = finalContentResult.SynopsisArabic
	contentResult.PrimaryInfo.OriginalTitle = finalContentResult.OriginalTitle
	contentResult.PrimaryInfo.AlternativeTitle = finalContentResult.AlternativeTitle
	contentResult.PrimaryInfo.ArabicTitle = finalContentResult.ArabicTitle
	contentResult.PrimaryInfo.TransliteratedTitle = finalContentResult.TransliteratedTitle
	contentResult.PrimaryInfo.Notes = finalContentResult.Notes
	/* Fetch content_cast*/
	var contentCast Cast
	if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
		Joins("left join actor actor on actor.id =cc.main_actor_id").
		Joins("left join actor actress on actress.id =cc.main_actress_id").
		Where("cc.id=?", finalContentResult.CastId).Scan(&contentCast).Error; castResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	contentResult.Cast.CastId = finalContentResult.CastId
	contentResult.Cast.MainActorId = contentCast.MainActorId
	contentResult.Cast.MainActressId = contentCast.MainActressId
	contentResult.Cast.MainActorEnglish = contentCast.MainActorEnglish
	contentResult.Cast.MainActorArabic = contentCast.MainActorArabic
	contentResult.Cast.MainActressEnglish = contentCast.MainActressEnglish
	contentResult.Cast.MainActressArabic = contentCast.MainActressArabic
	var contentActor []ContentActor
	if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
		Joins("left join actor a on a.id =ca.actor_id").
		Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
		Joins("left join writer w on w.id =cw.writer_id").
		Joins("left join content_director cd on cd.cast_id =ca.cast_id").
		Joins("left join director d on d.id =cd.director_id").
		Where("ca.cast_id=?", finalContentResult.CastId).Scan(&contentActor).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var actorEnglish, actorId, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic []string
	for _, actorIds := range contentActor {
		actorId = append(actorId, actorIds.ActorId)
		actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
		actorArabic = append(actorArabic, actorIds.ActorArabic)
		writerId = append(writerId, actorIds.WriterId)
		writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
		writerArabic = append(writerArabic, actorIds.WriterArabic)
		directorId = append(directorId, actorIds.DirectorId)
		directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
		directorArabic = append(directorArabic, actorIds.DirectorArabic)
	}
	contentResult.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
	contentResult.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
	contentResult.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
	contentResult.Cast.WriterId = common.RemoveDuplicateValues(writerId)
	contentResult.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
	contentResult.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
	contentResult.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
	contentResult.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
	contentResult.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)

	/* Fetch content_music*/
	var contentMusic []ContentMusic
	if actorResult := db.Debug().Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
		Joins("left join singer s on s.id =cs.singer_id").
		Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
		Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
		Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
		Joins("left join song_writer sw on sw.id =csw.song_writer_id").
		Where(" cs.music_id=?", finalContentResult.MusicId).Scan(&contentMusic).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	var singerEnglish, singerId, singerArabic, composerId, composerEnglish, SongWriterId, composerArabic, SongWriterEnglish, SongWriterArabic []string
	for _, musicIds := range contentMusic {
		singerId = append(singerId, musicIds.SingerIds)
		singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
		singerArabic = append(singerArabic, musicIds.SingersArabic)
		composerId = append(composerId, musicIds.MusicComposerIds)
		composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
		composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
		SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
		SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
		SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
	}
	contentResult.Music.MusicId = finalContentResult.MusicId
	contentResult.Music.SingerIds = common.RemoveDuplicateValues(singerId)
	contentResult.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
	contentResult.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
	contentResult.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
	contentResult.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
	contentResult.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
	contentResult.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
	contentResult.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
	contentResult.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
	/*Fetch tag_info*/
	var contentTagInfo []ContentTag
	db.Table("content_tag ct").Select("tdt.name").
		Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
		Where("ct.tag_info_id=?", finalContentResult.TagInfoId).Scan(&contentTagInfo)
	var tagInfo []string
	for _, tagInfoIds := range contentTagInfo {
		tagInfo = append(tagInfo, tagInfoIds.Name)
	}
	contentResult.TagInfo.Tags = tagInfo
	if len(tagInfo) < 1 {
		buffer := make([]string, 0)
		contentResult.TagInfo.Tags = buffer
	}
	/*Non_textual Data*/
	if finalContentResult.HasPosterImage {
		contentResult.NonTextualDataEpisode.PosterImage = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + os.Getenv("POSTER_IMAGE")
	} else {
		contentResult.NonTextualDataEpisode.PosterImage = ""
	}
	if finalContentResult.HasDubbingScript {
		contentResult.NonTextualDataEpisode.DubbingScript = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + os.Getenv("DUBBLING_SCRIPT")
	} else {
		contentResult.NonTextualDataEpisode.DubbingScript = ""
	}
	if finalContentResult.HasSubtitlingScript {
		contentResult.NonTextualDataEpisode.SubtitlingScript = os.Getenv("IMAGE_URL_GO") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + finalContentResult.Id + "/subtitling-script"
	} else {
		contentResult.NonTextualDataEpisode.SubtitlingScript = ""
	}
	/*Translation data*/
	contentResult.Translation.LanguageType = common.LanguageOriginTypes(finalContentResult.LanguageType)
	contentResult.Translation.DubbingLanguage = finalContentResult.DubbingLanguage
	// contentResult.Translation.DubbingDialectId = finalContentResult.DubbingDialectId
	contentResult.Translation.SubtitlingLanguage = finalContentResult.SubtitlingLanguage
	/*SeoDetails*/
	contentResult.SeoDetails.EnglishMetaTitle = finalContentResult.EnglishMetaTitle
	contentResult.SeoDetails.ArabicMetaTitle = finalContentResult.ArabicMetaTitle
	contentResult.SeoDetails.EnglishMetaDescription = finalContentResult.EnglishMetaDescription
	contentResult.SeoDetails.ArabicMetaDescription = finalContentResult.ArabicMetaDescription

	contentResult.CreatedAt = finalContentResult.CreatedAt
	contentResult.ModifiedAt = finalContentResult.ModifiedAt
	contentResult.Id = finalContentResult.Id
	if CountryResult != 0 || IsCheck {
		c.JSON(http.StatusOK, gin.H{"data": contentResult})
	}
}

// GetMenuDetails - Get all menu list details
// GET /get_menu
// @Description Get All menu list details by platform ID
// @Tags Menu
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param device query string true "Device Name"
// @Success 200  object MenuDetails
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /get_menu [get]
func (hs *HandlerService) GetMenuDetails(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	serverError := common.ServerErrorResponse()
	db := c.MustGet("FCDB").(*gorm.DB)
	var pageDetails []PageDetails
	var menu MenuData
	var menuResponse []MenuData
	var response MenuDetails
	var limit, offset, current_page int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
	}
	if limit == 0 {
		limit, _ = strconv.ParseInt(os.Getenv("DEFAULT_PAGE_SIZE"), 10, 64)
	}

	if current_page == 0 {
		current_page = 1
		offset = 0
	} else {
		offset = current_page*limit - limit
	}

	// offset = current_page * limit
	if c.Request.URL.Query()["device"] == nil || c.Request.URL.Query()["device"][0] == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}
	DeviceName := strings.ToLower(c.Request.URL.Query()["device"][0])
	deviceId := common.DeviceIds(DeviceName)
	rows := db.Debug().Table("page p").Select("p.*").Joins("inner join page_target_platform ptp on ptp.page_id=p.id").Where("p.is_disabled=false and p.deleted_by_user_id is null and p.third_party_page_key is not null and ptp.target_platform=?", deviceId).Group("p.id,ptp.page_order_number").Order("ptp.page_order_number asc").Find(&pageDetails).RowsAffected
	if err := db.Table("page p").Select("p.*").Joins("inner join page_target_platform ptp on ptp.page_id=p.id").Where("p.is_disabled=false and p.deleted_by_user_id is null and p.third_party_page_key is not null and ptp.target_platform=?", deviceId).Group("p.id,ptp.page_order_number").Order("ptp.page_order_number asc").Limit(limit).Offset(offset).Find(&pageDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	type PageIds struct {
		Id string `json:"id"`
	}
	var pageids []PageIds
	var ids []string
	if err := db.Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW()").Find(&pageids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if pageids != nil {
		for _, pageid := range pageids {
			ids = append(ids, pageid.Id)
		}
	}
	for _, details := range pageDetails {
		//changed to third party key
		menu.Id = details.ThirdPartyPageKey
		menu.FriendlyUrlEnglish = details.EnglishPageFriendlyUrl
		menu.SeoDescription = details.EnglishMetaDescription
		menu.TitleEnglish = details.EnglishTitle
		menu.FriendlyUrlArabic = details.ArabicPageFriendlyUrl
		menu.TitleArabic = details.ArabicTitle
		menu.Type = common.PageTypes(details.PageType)
		if details.PageType != 16 && details.PageType != 8 {
			exists := common.FindString(ids, details.Id)
			if (details.PageType == 0 && exists == true) || details.PageType == 1 {
				menu.Type = "Home"
			} else {
				menu.Type = "VOD"
			}
		}
		menu.Featured = nil
		menu.Playlists = nil
		menuResponse = append(menuResponse, menu)
	}

	lastPage := rows / limit
	response.Total = int(rows)
	response.PerPage = int(limit)
	response.CurrentPage = int(current_page)
	if lastPage == 0 {
		lastPage = 1
		response.LastPage = 1
	} else {
		response.LastPage = int(lastPage)
	}
	var Host string
	if c.Request.Host == "localhost:3006" {
		Host = "http://" + c.Request.Host
	} else {
		Host = os.Getenv("BASE_URL")
	}
	if current_page < lastPage {
		var NextPageUrl string
		NextPageUrl = Host + "/get_menu?device=" + DeviceName + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page+1, 10)
		response.NextPageUrl = &NextPageUrl
	} else {
		response.NextPageUrl = nil
	}

	if current_page == 1 && lastPage == 1 {
		response.PrevPageUrl = nil

	} else if current_page-1 > 0 || current_page == 1 {
		var PrevPageUrl string
		if current_page == 1 {
			response.PrevPageUrl = nil
		} else {
			PrevPageUrl = Host + "/get_menu?device=" + DeviceName + "&limit=" + strconv.FormatInt(limit, 10) + "&page=" + strconv.FormatInt(current_page-1, 10)
			response.PrevPageUrl = &PrevPageUrl
		}
	} else {
		response.PrevPageUrl = nil
	}

	response.From = int(offset + 1)
	if int(rows) < int(limit) {
		response.To = int(rows)
	} else {
		response.To = int(offset + limit)
	}

	if int(offset+1) > int(rows) {
		c.JSON(http.StatusBadRequest, "{'message':'Not Found'}")
		return
	}

	// response.To = int(offset + rows)
	response.Data = menuResponse
	c.JSON(http.StatusOK, response)
	return
}

// GetPageDetails - Get  page details
// GET /get_page/:pageId
// @Description Get All menu list details by platform ID
// @Tags Menu
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param pageId path string true "Page Id"
// @Param country path string false "Country code of the user."
// @Success 200  object PageDetails
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /get_page/{pageId} [get]
func (hs *HandlerService) GetPageDetails(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	pagekey, _ := strconv.Atoi(c.Param("pageId"))
	var countryCode string
	if c.Request.URL.Query()["Country"] != nil {
		countryCode = strings.ToUpper(c.Request.URL.Query()["Country"][0])
	}
	if len(countryCode) != 2 {
		countryCode = "AE"
	}
	countryId := int(common.Countrys(countryCode))
	serverError := common.ServerErrorResponse()
	db := c.MustGet("FCDB").(*gorm.DB)
	var details PageDetails
	var menu MenuDatas
	if err := db.Table("page p").Select("p.*").Where("p.is_disabled=false and p.deleted_by_user_id is null and p.page_key=?", pagekey).Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	type PageIds struct {
		Id string `json:"id"`
	}
	var pageids []PageIds
	var ids []string
	if err := db.Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW() and p.page_key=?", pagekey).Find(&pageids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if pageids != nil {
		for _, pageid := range pageids {
			ids = append(ids, pageid.Id)
		}
	}
	menu.Id = details.PageKey
	menu.FriendlyUrlEnglish = details.EnglishPageFriendlyUrl
	menu.SeoDescription = details.EnglishMetaDescription
	menu.TitleEnglish = details.EnglishTitle
	menu.FriendlyUrlArabic = details.ArabicPageFriendlyUrl
	menu.TitleArabic = details.ArabicTitle
	menu.Type = common.PageTypes(details.PageType)
	if details.PageType != 16 && details.PageType != 8 {
		exists := common.FindString(ids, details.Id)
		if (details.PageType == 0 && exists == true) || details.PageType == 1 {
			menu.Type = "Home"
		} else {
			menu.Type = "VOD"
		}
	}
	//page slider details
	var blackPlaylistCount, redPlaylistCount, greenPlaylistCount int
	var slider Slider
	if err := db.Select("s.*").Table("slider s").Joins("inner join page_slider ps on ps.slider_id=s.id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and ps.page_id=? and (s.scheduling_start_date <=NOW() or ps.order =0) and (s.scheduling_end_date >=NOW()  or ps.order =0)", details.Id).Limit(1).Order("ps.order desc").Find(&slider).Error; err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if slider.SliderKey != 0 {
		var featuredDetails FeaturedDetails
		var featuredPlaylist FeaturedPlaylists
		var featuredPlaylists []FeaturedPlaylists
		featuredDetails.ID = int64(slider.SliderKey)
		featuredDetails.Type = common.SliderTypes(slider.Type)
		if slider.BlackAreaPlaylistId != "" || slider.RedAreaPlaylistId != "" || slider.GreenAreaPlaylistId != "" {
			playlists, _ := SliderPlaylists(slider.BlackAreaPlaylistId, slider.RedAreaPlaylistId, slider.GreenAreaPlaylistId, c)
			for _, playlist := range playlists {
				featuredPlaylist.ID = int32(playlist.PlaylistKey)
				featuredPlaylist.PlaylistType = playlist.PlaylistType
				contentIds, err := PlaylistItemContents(playlist.ID, c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				var Ids []string
				for _, content := range contentIds {
					Ids = append(Ids, content.ContentId)
				}
				var playlistContents []PlaylistContent
				response := make(chan FunctionResponse)
				go GetContentDetails(Ids, countryId, c, response)
				details := <-response
				if details.Err != nil {
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}
				playlistContents = details.ContentDetails
				if playlistContents == nil {
					continue
				}
				if playlist.PlaylistType == "black_playlist" {
					blackPlaylistCount = len(playlistContents)
				} else if playlist.PlaylistType == "red_playlist" {
					redPlaylistCount = len(playlistContents)
				} else {
					greenPlaylistCount = len(playlistContents)
				}
				var contents []PlaylistContent
				for _, id := range contentIds {
					for _, content := range playlistContents {
						if id.ContentId == content.ContentId {
							contents = append(contents, content)
						}
					}
				}
				featuredPlaylist.Content = contents
				featuredPlaylists = append(featuredPlaylists, featuredPlaylist)
			}
		}
		if details.PageType == 1 && blackPlaylistCount >= common.BlackPlaylistCount && redPlaylistCount == common.RedPlaylistCount && greenPlaylistCount >= common.GreenPlaylistCount {
			featuredDetails.Playlists = featuredPlaylists
		} else if details.PageType != 1 {
			featuredDetails.Playlists = featuredPlaylists
		}
		if len(featuredDetails.Playlists) > 0 {
			menu.Featured = &featuredDetails
		}
	}
	//page palylist details
	var playlists []Playlist
	if err := db.Select("p.id,english_title,arabic_title,p.scheduling_start_date,p.scheduling_end_date,p.deleted_by_user_id,p.is_disabled,p.created_at,p.playlist_key,p.modified_at,p.playlist_type").Table("page_playlist pp").Joins("join playlist p on p.id =pp.playlist_id").Where("p.is_disabled =false and p.deleted_by_user_id is null and pp.page_id =? and (p.scheduling_start_date <=now() or p.scheduling_start_date is null) and (p.scheduling_end_date >=now() or p.scheduling_end_date is null)", details.Id).Order("pp.order asc").Find(&playlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var pagePlaylists []MenuPlaylists
	for _, playlist := range playlists {
		pagePlaylist := MenuPlaylists{}
		pagePlaylist.ID = int32(playlist.PlaylistKey)
		pagePlaylist.Content = []PlaylistContent{}
		pagePlaylist.TitleEnglish = playlist.EnglishTitle
		pagePlaylist.TitleArabic = playlist.ArabicTitle
		contentIds, err := PlaylistItemContents(playlist.ID, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var Ids []string
		for _, content := range contentIds {
			Ids = append(Ids, content.ContentId)
		}
		var playlistContents []PlaylistContent
		response := make(chan FunctionResponse)
		go GetContentDetails(Ids, countryId, c, response)
		details := <-response
		if details.Err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		playlistContents = details.ContentDetails
		if playlistContents == nil {
			continue
		}
		var contents []PlaylistContent
		for _, id := range contentIds {
			for _, content := range playlistContents {
				if id.ContentId == content.ContentId {
					contents = append(contents, content)
				}
			}
		}
		pagePlaylist.Content = contents
		pagePlaylists = append(pagePlaylists, pagePlaylist)
	}
	menu.Playlists = pagePlaylists
	c.JSON(http.StatusOK, menu)
	return
}
func SliderPlaylists(BlackAreaPlaylistId string, RedAreaPlaylistId string, GreenAreaPlaylistId string, c *gin.Context) ([]Playlist, error) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var playlists []Playlist
	playlist := []interface{}{}
	fields := "id,playlist_key,case"
	if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'  when id='" + RedAreaPlaylistId + "' then 'red_playlist' when id='" + GreenAreaPlaylistId + "' then 'green_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, RedAreaPlaylistId, GreenAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'  when id='" + RedAreaPlaylistId + "' then 'red_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, RedAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist' when id='" + GreenAreaPlaylistId + "' then 'green_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId, GreenAreaPlaylistId}
	} else if BlackAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 'black_playlist'"
		playlist = []interface{}{BlackAreaPlaylistId}
	}
	fields += " end as playlist_type,"
	fields += " case"
	if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1  when id='" + RedAreaPlaylistId + "' then 2 when id='" + GreenAreaPlaylistId + "' then 3"
	} else if BlackAreaPlaylistId != "" && RedAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1  when id='" + RedAreaPlaylistId + "' then 2"
	} else if BlackAreaPlaylistId != "" && GreenAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1 when id='" + GreenAreaPlaylistId + "' then 3"
	} else if BlackAreaPlaylistId != "" {
		fields += " when id='" + BlackAreaPlaylistId + "' then 1"
	}
	fields += " else 0 end as playlist_order"
	if err := db.Select(fields).Where("id in(?) and (scheduling_start_date <=now() or scheduling_start_date is null) and (scheduling_end_date >=now() or scheduling_end_date is null)", playlist).Order("playlist_type desc").Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}
func PlaylistItemContents(playlistId string, c *gin.Context) ([]PlaylistContentIds, error) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var contentIds []PlaylistContentIds
	if err := db.Table("playlist_item_content pic").Select("pic.content_id").Joins("inner join playlist_item pi2 on pi2.id=pic.playlist_item_id inner join playlist p on p.id=pi2.playlist_id").Where("p.id =?", playlistId).Order("pi2.order asc").Find(&contentIds).Error; err != nil {
		return nil, err
	}
	return contentIds, nil
}
func GetContentDetails(contentIds []string, country int, c *gin.Context, response chan FunctionResponse) {
	var contentResponse FunctionResponse
	var content PlaylistContent
	var contents []PlaylistContent
	db := c.MustGet("DB").(*gorm.DB)
	var cDetails []ContentDetails
	if err := db.Raw("select c.id,c.content_key,round(CAST(c.average_rating as numeric),1) as average_rating,min(pi2.video_content_id) as video_id,Replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,lower(c.content_type) as content_type,atci.english_synopsis as synopsis_english,atci.arabic_synopsis as synopsis_arabic,c.english_meta_title as seo_title_english,c.arabic_meta_title as seo_title_arabic,c.english_meta_description as seo_description_english,c.arabic_meta_description as seo_description_arabic,min(pi2.duration) as length,cpi.transliterated_title as title_english,cpi.arabic_title as title_arabic,c.english_meta_title as seo_title,c.created_at as inserted_at,c.modified_at,cv.id as varience_id from content c join content_primary_info cpi on cpi.id = c.primary_info_id join about_the_content_info atci on atci.id =c.about_the_content_info_id join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where crc.country_id = ? and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=now() or cr.digital_rights_end_date is null) and c.status =1 and c.deleted_by_user_id is null and c.id in(?) and cv.status =1 and cv.deleted_by_user_id is null group by c.id,c.content_key,c.average_rating,cpi.transliterated_title,c.content_type,atci.english_synopsis,atci.arabic_synopsis,c.english_meta_title,c.arabic_meta_title,c.english_meta_description,c.arabic_meta_description,cpi.transliterated_title,cpi.arabic_title,c.english_meta_title,c.created_at,c.modified_at,cv.id union select c.id,c.content_key,round(CAST(c.average_rating as numeric),1) as average_rating,min(pi2.video_content_id) as video_id,Replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,lower(c.content_type) as content_type,atci.english_synopsis as synopsis_english,atci.arabic_synopsis as synopsis_arabic,s.english_meta_title as seo_title_english,s.arabic_meta_title as seo_title_arabic,s.english_meta_description as seo_description_english,s.arabic_meta_description as seo_description_arabic,min(pi2.duration) as length,cpi.transliterated_title as title_english,cpi.arabic_title as title_arabic,s.english_meta_title as seo_title,c.created_at as inserted_at,c.modified_at,s.id as varience_id from content c join season s on s.content_id=c.id join episode e on e.season_id=s.id join content_primary_info cpi on cpi.id = s.primary_info_id join about_the_content_info atci on atci.id =s.about_the_content_info_id join playback_item pi2 on pi2.id = e.playback_item_id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where crc.country_id = ? and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=now() or cr.digital_rights_end_date is null) and c.status =1 and c.deleted_by_user_id is null and c.id in(?) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null group by c.id,c.content_key,c.average_rating,cpi.transliterated_title,c.content_type,atci.english_synopsis,atci.arabic_synopsis,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,cpi.transliterated_title,cpi.arabic_title,s.english_meta_title,c.created_at,c.modified_at,s.id", country, contentIds, country, contentIds).Find(&cDetails).Error; err != nil {
		contentResponse.ContentDetails = nil
		contentResponse.Err = err
		response <- contentResponse
	}
	for _, details := range cDetails {
		content.ContentId = details.Id
		content.ContentKey = details.ContentKey
		content.AgeRating = details.AgeRating
		content.VideoId = details.VideoId
		content.FriendlyUrl = details.FriendlyUrl
		content.ContentType = details.ContentType
		content.SynopsisEnglish = details.SynopsisEnglish
		content.SynopsisArabic = details.SynopsisArabic
		content.SeoTitleEnglish = details.SeoTitleEnglish
		content.SeoTitleArabic = details.SeoTitleArabic
		content.SeoDescriptionEnglish = details.SeoDescriptionEnglish
		content.SeoDescriptionArabic = details.SeoDescriptionArabic
		content.Length = details.Length
		content.TitleEnglish = details.TitleEnglish
		content.TitleArabic = details.TitleArabic
		content.SeoTitle = details.SeoTitle
		var imagery ContentImageryDetails
		if details.ContentTier == 1 {
			imagery.Thumbnail = os.Getenv("IMAGE_URL") + details.Id + os.Getenv("POSTER_IMAGE")
			imagery.Backdrop = os.Getenv("IMAGE_URL") + details.Id + os.Getenv("DETAILS_BACKGROUND")
			imagery.MobileImg = os.Getenv("IMAGE_URL") + details.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
			imagery.FeaturedImg = os.Getenv("IMAGE_URL") + details.Id + os.Getenv("POSTER_IMAGE")
			imagery.OverlayPoster = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("OVERLAY_POSTER_IMAGE")
		} else {
			imagery.Thumbnail = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("POSTER_IMAGE")
			imagery.Backdrop = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("DETAILS_BACKGROUND")
			imagery.MobileImg = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("MOBILE_DETAILS_BACKGROUND")
			imagery.FeaturedImg = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("POSTER_IMAGE")
			imagery.OverlayPoster = os.Getenv("IMAGE_URL") + details.Id + "/" + details.VarienceId + os.Getenv("OVERLAY_POSTER_IMAGE")
		}
		content.Imagery = imagery
		content.InsertedAt = details.InsertedAt
		content.ModifiedAt = details.ModifiedAt
		contents = append(contents, content)
	}
	contentResponse.ContentDetails = contents
	contentResponse.Err = nil
	response <- contentResponse
}

type VideoDurationInfo struct {
	Duration     int           `json:"duration"`
	Thumbnails   []interface{} `json:"thumbnails"`
	UrlTrickplay string        `json:"url_trickplay"`
	UrlVideo     string        `json:"url_video"`
}

// GetVideoDurationDetails -  Get all menu list details
// GET /get_info/:videoId
// @Description Get all menu list details
// @Tags Video
// @Accept  json
// @Produce  json
// @Param videoId path string true "video Id."
// @Param country path string false "Country code of the user."
// @Success 200  object VideoDurationInfo
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /get_info/{videoId} [get]
func (hs *HandlerService) GetVideoDuration(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	serverError := common.ServerErrorResponse()

	VideoId := c.Param("videoId")
	response, err := common.GetCurlCall(os.Getenv("VIDEO_API") + VideoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var details VideoDurationInfo
	json.Unmarshal(response, &details)
	c.JSON(http.StatusOK, details)
	return
}

// GetMultiTierContentDetailsBasedonContentID
// GET /contents/multitier/:contentId
// @Description Get Multi Tier Content Details Based on Content ID
// @Tags MultiTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param contentId path string true "Content ID"
// @Param Country query string false "Country code of the user."
// @Success 200 {array} MultiTierContent
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /contents/multitier/{contentId} [get]
func (hs *HandlerService) GetMultiTierDetailsBasedonContentID(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}
	var contentResult MultiTierContent
	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()
	var finalContentResult FinalSeasonResult
	/*for country rights*/
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	}
	CountryResult := common.Countrys(country)
	ContentKey, _ := strconv.Atoi(c.Param("contentId"))
	var count int
	if err := db.Table("content").Where("third_party_content_key=?", ContentKey).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if count < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	}
	if UserId == os.Getenv("WATCH_NOW") {
		if CountryResult != 0 {
			if err := db.Debug().Table("content c").Select("distinct s.id,pi2.rights_id,c.third_party_content_key as multi_tier_content_key, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.third_party_season_key as season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background").
				Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
				Joins("join content_genre cg on cg.content_id  = c.id").
				Joins("join season s on s.content_id = c.id").
				Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
				Joins("join content_cast cc  on cc.Id  = s.cast_id ").
				Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id and atci.supplier !='Others'").
				Joins("join content_translation ct on ct.id = s. translation_id").
				Joins("join playback_item pi2 on pi2.translation_id = ct.id").
				Joins("join content_rights cr on cr.id = pi2.rights_id").
				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
				Where("c.watch_now_supplier = 'true' and c.third_party_content_key = ? and c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)	and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null and crc.country_id = ? ", ContentKey, CountryResult).
				Find(&finalContentResult).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country == "" || country == "all" {
			if err := db.Table("content c").Select("distinct s.id,pi2.rights_id,c.third_party_content_key as multi_tier_content_key, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.third_party_season_key as season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background").
				Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
				Joins("join content_genre cg on cg.content_id  = c.id").
				Joins("join season s on s.content_id = c.id").
				Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
				Joins("join content_cast cc  on cc.Id  = s.cast_id ").
				Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id and atci.supplier !='Others'").
				Joins("join content_translation ct on ct.id = s. translation_id").
				Joins("join playback_item pi2 on pi2.translation_id = ct.id").
				Joins("join content_rights cr on cr.id = pi2.rights_id").
				Where("c.watch_now_supplier = 'true' and c.third_party_content_key = ? and c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)	and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null ", ContentKey).
				Find(&finalContentResult).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country != "all" && CountryResult == 0 {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	} else {
		if CountryResult != 0 {
			if err := db.Debug().Table("content c").Select("distinct s.id,pi2.rights_id,c.third_party_content_key as multi_tier_content_key, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.third_party_season_key as season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background").
				Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
				Joins("join content_genre cg on cg.content_id  = c.id").
				Joins("join season s on s.content_id = c.id").
				Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
				Joins("join content_cast cc  on cc.Id  = s.cast_id ").
				Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id and atci.supplier !='Others'").
				Joins("join content_translation ct on ct.id = s. translation_id").
				Joins("join playback_item pi2 on pi2.translation_id = ct.id").
				Joins("join content_rights cr on cr.id = pi2.rights_id").
				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
				Where("	c.third_party_content_key = ? and c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)	and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null and crc.country_id = ? ", ContentKey, CountryResult).
				Find(&finalContentResult).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country == "" || country == "all" {
			if err := db.Debug().Table("content c").Select("distinct s.id,pi2.rights_id,c.third_party_content_key as multi_tier_content_key, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.third_party_season_key as season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background ").
				Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
				Joins("join content_genre cg on cg.content_id  = c.id").
				Joins("join season s on s.content_id = c.id").
				Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
				Joins("join content_cast cc  on cc.Id  = s.cast_id ").
				Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id and atci.supplier !='Others'").
				Joins("join content_translation ct on ct.id = s. translation_id").
				Joins("join playback_item pi2 on pi2.translation_id = ct.id").
				Joins("join content_rights cr on cr.id = pi2.rights_id").
				Where("	c.third_party_content_key = ? and c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)	and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null ", ContentKey).
				Find(&finalContentResult).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country != "all" && CountryResult == 0 {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
	}
	contentResult.ContentKey = finalContentResult.MultiTierContentKey
	contentResult.PrimaryInfo.ContentType = finalContentResult.ContentType
	contentResult.PrimaryInfo.OriginalTitle = finalContentResult.OriginalTitle
	contentResult.PrimaryInfo.AlternativeTitle = finalContentResult.AlternativeTitle
	contentResult.PrimaryInfo.ArabicTitle = finalContentResult.ArabicTitle
	contentResult.PrimaryInfo.TransliteratedTitle = finalContentResult.TransliteratedTitle
	contentResult.PrimaryInfo.Notes = finalContentResult.Notes

	/*Fetch content_geners*/
	var contentGenres []SeasonGenres
	var finalContentGenre []interface{}
	var newContentGenres NewSeasonGenres
	if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
		Joins("left join genre g on g.id=cg.genre_id").
		Where("content_id=?", finalContentResult.ContentId).Scan(&contentGenres).Error; genreResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, tagInfoIds := range contentGenres {
		var contentSubgenre []SeasonSubgenre
		if subgenreResult := db.Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
			Joins("left join subgenre sg on sg.id=csg.subgenre_id").
			Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var SubgenreEn []string
		var SubgenreAr []string
		for _, data := range contentSubgenre {
			SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
			SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
			newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
			newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
			newContentGenres.SubGenerEnglish = SubgenreEn
			newContentGenres.SubGenerArabic = SubgenreAr
			newContentGenres.Id = tagInfoIds.Id
		}
		finalContentGenre = append(finalContentGenre, newContentGenres)
	}
	contentResult.ContentGenres = finalContentGenre
	//content details

	var varianceTrailers []VarianceTrailers
	if varianceTrailersError := db.Debug().Raw("select * from variance_trailer vt where vt.season_id=? order by vt.order asc ", finalContentResult.SeasonId).Find(&varianceTrailers).Error; varianceTrailersError != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var variance VarianceTrailers
	var finalVarianceResult []interface{}

	for _, varianceData := range varianceTrailers {
		variance.Order = varianceData.Order
		variance.VideoTrailerId = varianceData.VideoTrailerId
		variance.EnglishTitle = varianceData.EnglishTitle
		variance.ArabicTitle = varianceData.ArabicTitle
		variance.Duration = varianceData.Duration
		variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
		//const url = "https://z5content-qa.s3.amazonaws.com/"
		variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + varianceData.Id + "/trailer-poster-image"
		variance.Id = varianceData.Id
		finalVarianceResult = append(finalVarianceResult, variance)
	}

	contentResult.ContentSeasons[0].TrailerInfo = finalVarianceResult

	contentResult.ContentSeasons[0].ContentId = finalContentResult.ContentId
	contentResult.ContentSeasons[0].SeasonKey = finalContentResult.SeasonKey
	contentResult.ContentSeasons[0].SeasonNumber = finalContentResult.SeasonNumber
	contentResult.ContentSeasons[0].CreatedAt = finalContentResult.InsertedAt
	contentResult.ContentSeasons[0].ModifiedAt = finalContentResult.ModifiedAt
	//season primary info

	contentResult.ContentSeasons[0].PrimaryInfo.SeasonNumber = finalContentResult.SeasonNumber
	contentResult.ContentSeasons[0].PrimaryInfo.OriginalTitle = finalContentResult.SeasonOriginalTitle
	contentResult.ContentSeasons[0].PrimaryInfo.AlternativeTitle = finalContentResult.SeasonAlternativeTitle
	contentResult.ContentSeasons[0].PrimaryInfo.ArabicTitle = finalContentResult.SeasonArabicTitle
	contentResult.ContentSeasons[0].PrimaryInfo.TransliteratedTitle = finalContentResult.SeasonTransliteratedTitle
	contentResult.ContentSeasons[0].PrimaryInfo.Notes = finalContentResult.SeasonNotes

	/* Fetch content_cast normal*/
	var contentCast Cast
	if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
		Joins("left join actor actor on actor.id =cc.main_actor_id").
		Joins("left join actor actress on actress.id =cc.main_actress_id").
		Where("cc.id=?", finalContentResult.CastId).Scan(&contentCast).Error; castResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	contentResult.ContentSeasons[0].Cast.CastId = finalContentResult.CastId
	contentResult.ContentSeasons[0].Cast.MainActorId = contentCast.MainActorId
	contentResult.ContentSeasons[0].Cast.MainActressId = contentCast.MainActressId
	contentResult.ContentSeasons[0].Cast.MainActorEnglish = contentCast.MainActorEnglish
	contentResult.ContentSeasons[0].Cast.MainActorArabic = contentCast.MainActorArabic
	contentResult.ContentSeasons[0].Cast.MainActressEnglish = contentCast.MainActressEnglish
	contentResult.ContentSeasons[0].Cast.MainActressArabic = contentCast.MainActressArabic
	var contentActor []ContentActor
	if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
		Joins("left join actor a on a.id =ca.actor_id").
		Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
		Joins("left  join writer w on w.id =cw.writer_id").
		Joins("left join content_director cd on cd.cast_id =ca.cast_id").
		Joins("left join director d on d.id =cd.director_id").
		Where("ca.cast_id=?", finalContentResult.CastId).Scan(&contentActor).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var actorId, actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic []string
	for _, actorIds := range contentActor {
		actorId = append(actorId, actorIds.ActorId)
		actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
		actorArabic = append(actorArabic, actorIds.ActorArabic)
		writerId = append(writerId, actorIds.WriterId)
		writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
		writerArabic = append(writerArabic, actorIds.WriterArabic)
		directorId = append(directorId, actorIds.DirectorId)
		directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
		directorArabic = append(directorArabic, actorIds.DirectorArabic)
	}
	contentResult.ContentSeasons[0].Cast.ActorIds = common.RemoveDuplicateValues(actorId)
	contentResult.ContentSeasons[0].Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
	contentResult.ContentSeasons[0].Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
	contentResult.ContentSeasons[0].Cast.WriterId = common.RemoveDuplicateValues(writerId)
	contentResult.ContentSeasons[0].Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
	contentResult.ContentSeasons[0].Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
	contentResult.ContentSeasons[0].Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
	contentResult.ContentSeasons[0].Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
	contentResult.ContentSeasons[0].Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)
	/*fetching music details*/
	var contentMusic []ContentMusic
	if actorResult := db.Debug().Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_composers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
		Joins("left join singer s on s.id =cs.singer_id").
		Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
		Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
		Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
		Joins("left join song_writer sw on sw.id =csw.song_writer_id").
		Where("cs.music_id=?", finalContentResult.MusicId).Scan(&contentMusic).Error; actorResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
	for _, musicIds := range contentMusic {
		singerId = append(singerId, musicIds.SingerIds)
		singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
		singerArabic = append(singerArabic, musicIds.SingersArabic)
		composerId = append(composerId, musicIds.MusicComposerIds)
		composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
		composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
		SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
		SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
		SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
	}
	contentResult.ContentSeasons[0].Music.MusicId = finalContentResult.MusicId
	contentResult.ContentSeasons[0].Music.SingerIds = common.RemoveDuplicateValues(singerId)
	contentResult.ContentSeasons[0].Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
	contentResult.ContentSeasons[0].Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
	contentResult.ContentSeasons[0].Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
	contentResult.ContentSeasons[0].Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
	contentResult.ContentSeasons[0].Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
	contentResult.ContentSeasons[0].Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
	contentResult.ContentSeasons[0].Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
	contentResult.ContentSeasons[0].Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
	/*content tag info*/
	var contentTagInfo []ContentTag
	db.Table("content_tag ct").Select("tdt.name").
		Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
		Where("ct.tag_info_id=?", finalContentResult.TagInfoId).Scan(&contentTagInfo)
	var tagInfo []string
	for _, tagInfoIds := range contentTagInfo {
		tagInfo = append(tagInfo, tagInfoIds.Name)
	}
	contentResult.ContentSeasons[0].TagInfo.Tags = tagInfo
	if len(tagInfo) < 1 {
		buffer := make([]string, 0)
		contentResult.ContentSeasons[0].TagInfo.Tags = buffer
	}
	/*season geners*/
	var seasonGenres []SeasonGenres
	var finalSeasonGenre []interface{}
	var newSeasonGenres NewSeasonGenres
	if genreResult := db.Table("season_genre sg").Select("sg.id ,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
		Joins("left join genre g on g.id=sg.genre_id").
		Where("sg.season_id = ?", finalContentResult.Id).Scan(&seasonGenres).Error; genreResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	for _, tagInfoIds := range seasonGenres {
		var seasonSubgenre []SeasonSubgenre
		if subgenreResult := db.Table("season_subgenre ssg").Select("sg.english_name as sub_gener_english,sg.arabic_name as sub_gener_arabic").
			Joins("left join subgenre sg on sg.id=ssg.subgenre_id").
			Where("ssg.season_genre_id=?", tagInfoIds.Id).Scan(&seasonSubgenre).Error; subgenreResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var SubgenreEn []string
		var SubgenreAr []string
		for _, data := range seasonSubgenre {
			SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
			SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
			newSeasonGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
			newSeasonGenres.GenerArabicName = tagInfoIds.GenerArabicName
			newSeasonGenres.SubGenerEnglish = SubgenreEn
			newSeasonGenres.SubGenerArabic = SubgenreAr
			newSeasonGenres.Id = tagInfoIds.Id

		}
		finalSeasonGenre = append(finalSeasonGenre, newSeasonGenres)
	}
	contentResult.ContentSeasons[0].SeasonGenres = finalSeasonGenre
	/*about the content*/
	contentResult.ContentSeasons[0].AboutTheContent.OriginalLanguage = finalContentResult.OriginalLanguage
	contentResult.ContentSeasons[0].AboutTheContent.Supplier = finalContentResult.Supplier
	contentResult.ContentSeasons[0].AboutTheContent.AcquisitionDepartment = finalContentResult.AcquisitionDepartment
	contentResult.ContentSeasons[0].AboutTheContent.EnglishSynopsis = finalContentResult.EnglishSynopsis
	contentResult.ContentSeasons[0].AboutTheContent.ArabicSynopsis = finalContentResult.ArabicSynopsis
	contentResult.ContentSeasons[0].AboutTheContent.ProductionYear = finalContentResult.ProductionYear
	contentResult.ContentSeasons[0].AboutTheContent.ProductionHouse = finalContentResult.ProductionHouse
	contentResult.ContentSeasons[0].AboutTheContent.AgeGroup = common.AgeRatings(finalContentResult.AgeGroup, "en")
	/*production countries*/
	var productionCountry []ProductionCountry
	if productionCountryResult := db.Table("production_country ").Select("country_id").Where("about_the_content_info_id=?", finalContentResult.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var countries []int

	for _, prcountries := range productionCountry {
		countries = append(countries, prcountries.CountryId)
	}
	contentResult.ContentSeasons[0].AboutTheContent.ProductionCountries = countries
	if len(tagInfo) < 1 {
		buffer := make([]int, 0)
		contentResult.ContentSeasons[0].AboutTheContent.ProductionCountries = buffer
	}
	/*translation details*/
	contentResult.ContentSeasons[0].Translation.LanguageType = common.LanguageOriginTypes(finalContentResult.MultiTierLanguageType)
	contentResult.ContentSeasons[0].Translation.DubbingLanguage = finalContentResult.DubbingLanguage
	// contentResult.ContentSeasons.Translation.DubbingDialectId = finalContentResult.DubbingDialectId
	contentResult.ContentSeasons[0].Translation.SubtitlingLanguage = finalContentResult.SubtitlingLanguage
	/*episode details*/
	var episodeDetails []FetchEpisodeDetailsMultiTier
	var episodeResult EpisodeDetailsMultiTier
	if err := db.Debug().Table("episode e").Select("e.has_poster_image,e.has_dubbing_script,e.has_subtitling_script,e.number as episode_number,e.third_party_episode_key as episode_key,pi2.duration as length,pi2.video_content_id,e.synopsis_english,e.synopsis_arabic,e.has_poster_image,e.has_dubbing_script,e.has_subtitling_script,e.id as episode_id").
		Joins("join season s on s.id = e.season_id").
		Joins("join playback_item pi2 on pi2.id = e.playback_item_id").
		Where("third_party_season_key = ? ", finalContentResult.SeasonKey).
		Order("e.number asc").
		Find(&episodeDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var allEpisodes []interface{}

	for _, value := range episodeDetails {

		episodeResult.EpisodeNumber = value.EpisodeNumber
		episodeResult.EpisodeKey = value.EpisodeKey
		episodeResult.Length = value.Length
		episodeResult.VideoContentUrl = os.Getenv("IMAGE_URL_GO") + value.VideoContentId
		episodeResult.SynopsisEnglish = value.SynopsisEnglish
		episodeResult.SynopsisArabic = value.SynopsisArabic
		//nontextual data for episodes

		if value.HasPosterImage {
			episodeResult.NonTextualData.PosterImage = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + value.EpisodeId + os.Getenv("POSTER_IMAGE")
		} else {
			episodeResult.NonTextualData.PosterImage = ""
		}
		if value.HasDubbingScript {
			*episodeResult.NonTextualData.DubbingScript = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + value.EpisodeId + os.Getenv("DUBBLING_SCRIPT")
		} else {
			episodeResult.NonTextualData.DubbingScript = nil
		}
		if value.HasSubtitlingScript {
			*episodeResult.NonTextualData.SubtitlingScript = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/" + value.EpisodeId + "/subtitling-script"
		} else {
			episodeResult.NonTextualData.SubtitlingScript = nil
		}
		episodeResult.EpisodeId = value.EpisodeId
		allEpisodes = append(allEpisodes, episodeResult)
	}
	contentResult.ContentSeasons[0].EpisodeResult = allEpisodes
	/*non textual data of content or season*/

	if finalContentResult.HasPosterImage {
		contentResult.ContentSeasons[0].ContentNonTextualData.PosterImage = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/poster-image"
	}
	if finalContentResult.HasOverlayPosterImage {
		contentResult.ContentSeasons[0].ContentNonTextualData.OverlayPosterImage = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/overlay-poster-image"
	}
	if finalContentResult.HasDetailsBackground {
		contentResult.ContentSeasons[0].ContentNonTextualData.DetailsBackground = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/details-background"
	}
	if finalContentResult.HasMobileDetailsBackground {
		contentResult.ContentSeasons[0].ContentNonTextualData.MobileDetailsBackground = IMAGES + finalContentResult.ContentId + "/" + finalContentResult.SeasonId + "/mobile-details-background"
	}

	/*digital rights region season*/
	var digitalRightsRegions []DigitalRightsRegions
	if countryError := db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", finalContentResult.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	var SeasonRights []int
	for _, idarr := range digitalRightsRegions {
		SeasonRights = append(SeasonRights, idarr.CountryId)

	}
	/*CHECKING THE value*/
	if country == "" || country == "all" {
		contentResult.ContentSeasons[0].DigitalRightsRegions = SeasonRights
	} else if country != "all" {
		contentResult.ContentSeasons[0].DigitalRightsRegions = nil
	}
	if len(SeasonRights) < 1 {
		buffer := make([]int, 0)
		contentResult.ContentSeasons[0].DigitalRightsRegions = buffer
	}
	//season id
	contentResult.ContentSeasons[0].SeasonId = finalContentResult.SeasonId
	//SEASON seo details
	contentResult.SeoDetails.EnglishMetaTitle = finalContentResult.EnglishMetaTitle
	contentResult.SeoDetails.ArabicMetaTitle = finalContentResult.ArabicMetaTitle
	contentResult.SeoDetails.EnglishMetaDescription = finalContentResult.EnglishMetaDescription
	contentResult.SeoDetails.ArabicMetaDescription = finalContentResult.ArabicMetaDescription
	contentResult.CreatedAt = finalContentResult.CreatedAt
	contentResult.ModifiedAt = finalContentResult.ModifiedAt
	//content idId
	contentResult.ContentId = finalContentResult.ContentId

	if CountryResult != 0 || (country == "" || country == "all") {
		c.JSON(http.StatusOK, gin.H{"data": contentResult})
	}
}

// GetAllMultiTierDetails- Get All Multi Tier Content Details
// GET /contents/multitier/
// @Description Get All Multi Tier Content Details
// @Tags MultiTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param offset query string false "Zero-based index of the first requested item."
// @Param limit query string false "Maximum number of collection items to return for a single request."
// @Param country path string false "Country code of the user."
// @Success 200  object AllMultiTier
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /contents/multitier [get]
func (hs *HandlerService) GetAllMultiTierDetails(c *gin.Context) {
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
	var contentResult AllMultiTierContent
	var allContents []AllMultiTierContent
	var CountryResult int32
	serverError := common.ServerErrorResponse()
	var finalContentResult []FinalSeasonResult
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if offset == 0 {
		offset = 0
	}
	if limit == 0 {
		limit = 5
	}
	var totalCount int
	/*digital rights*/
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
		fmt.Println(country)
	}
	CountryResult = common.Countrys(country)
	if CountryResult != 0 {
		if err := db.Debug().Table("content c").Select("distinct c.content_key as multi_tier_content_key,s.id,pi2.rights_id, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background").
			Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
			Joins("join content_genre cg on cg.content_id  = c.id").
			Joins("join season s on s.content_id = c.id").
			Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
			Joins("join content_cast cc  on cc.Id  = s.cast_id ").
			Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id").
			Joins("join content_translation ct on ct.id = s. translation_id").
			Joins("join playback_item pi2 on pi2.translation_id = ct.id").
			Joins("join content_rights cr on cr.id = pi2.rights_id").
			Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
			Where("c.status = 1 and c.content_tier =2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null  and crc.country_id = ? ", CountryResult).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		db.Raw("select count(distinct c.id) from content c join season s on s.content_id = c.id join content_translation ct on ct.id = s. translation_id join playback_item pi2 on	pi2.translation_id = ct.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where (c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null  and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null) and crc.country_id = ?", CountryResult).Count(&totalCount)
	} else if country == "" || country == "all" {
		if err := db.Debug().Table("content c").Select("distinct c.content_key as multi_tier_content_key,s.id,pi2.rights_id, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title,cpi.transliterated_title,cpi.notes,c.id as content_id,s.season_key,s.number as season_number,s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background ").
			Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
			Joins("join content_genre cg on cg.content_id  = c.id").
			Joins("join season s on s.content_id = c.id").
			Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
			Joins("join content_cast cc  on cc.Id  = s.cast_id ").
			Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id").
			Joins("join content_translation ct on ct.id = s. translation_id").
			Joins("join playback_item pi2 on pi2.translation_id = ct.id").
			Joins("join content_rights cr on cr.id = pi2.rights_id").
			Where("c.status = 1 and c.content_tier =2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null ").Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		db.Raw("select count(distinct c.id) from content c join season s on s.content_id = c.id join content_translation ct on ct.id = s. translation_id join playback_item pi2 on	pi2.translation_id = ct.id join content_rights cr on cr.id = pi2.rights_id where (c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null  and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null) ").Count(&totalCount)
	} else if country != "all" && CountryResult == 0 {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var contentSeasons []AllContentSeasons
	var contentSeason AllContentSeasons

	for _, eachcontent := range finalContentResult {
		contentResult.ContentKey = eachcontent.MultiTierContentKey
		/*content seasons*/
		contentSeason.ContentId = eachcontent.ContentId
		contentSeason.SeasonKey = eachcontent.SeasonKey
		contentSeason.SeasonNumber = eachcontent.SeasonNumber
		contentSeason.CreatedAt = eachcontent.CreatedAt
		contentSeason.ModifiedAt = eachcontent.ModifiedAt
		/*primary info season*/
		contentSeason.PrimaryInfo.SeasonNumber = eachcontent.SeasonNumber
		contentSeason.PrimaryInfo.OriginalTitle = eachcontent.SeasonOriginalTitle
		contentSeason.PrimaryInfo.AlternativeTitle = eachcontent.SeasonAlternativeTitle
		// contentSeason.PrimaryInfo.ArabicTitle = eachcontent.SeasonArabicTitlecontentResult
		/* Fetch content_cast normal*/
		var contentCast Cast
		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
			Joins("left join actor actor on actor.id =cc.main_actor_id").
			Joins("left join actor actress on actress.id =cc.main_actress_id").
			Where("cc.id=?", eachcontent.CastId).Scan(&contentCast).Error; castResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var varianceTrailers []VarianceTrailers
		if varianceTrailersError := db.Debug().Raw("select * from variance_trailer vt where vt.season_id=? order by vt.order asc ", eachcontent.SeasonId).Find(&varianceTrailers).Error; varianceTrailersError != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var variance VarianceTrailers
		var finalVarianceResult []interface{}

		for _, varianceData := range varianceTrailers {
			variance.Order = varianceData.Order
			variance.VideoTrailerId = varianceData.VideoTrailerId
			variance.EnglishTitle = varianceData.EnglishTitle
			variance.ArabicTitle = varianceData.ArabicTitle
			variance.Duration = varianceData.Duration
			variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
			//const url = "https://z5content-qa.s3.amazonaws.com/"
			variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + eachcontent.ContentId + "/" + eachcontent.SeasonId + "/" + varianceData.Id + "/trailer-poster-image"
			variance.Id = varianceData.Id
			finalVarianceResult = append(finalVarianceResult, variance)
		}

		contentSeason.TrailerInfo = finalVarianceResult

		contentSeason.Cast.CastId = eachcontent.CastId
		contentSeason.Cast.MainActorId = contentCast.MainActorId
		contentSeason.Cast.MainActressId = contentCast.MainActressId
		contentSeason.Cast.MainActorEnglish = contentCast.MainActorEnglish
		contentSeason.Cast.MainActorArabic = contentCast.MainActorArabic
		contentSeason.Cast.MainActressEnglish = contentCast.MainActressEnglish
		contentSeason.Cast.MainActressArabic = contentCast.MainActressArabic
		/*fetching other cast details */
		var contentActor []ContentActor
		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
			Joins("left join actor a on a.id =ca.actor_id").
			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
			Joins("left  join writer w on w.id =cw.writer_id").
			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
			Joins("left join director d on d.id =cd.director_id").
			Where("ca.cast_id=?", eachcontent.CastId).Scan(&contentActor).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var actorEnglish, actorArabic, writerId, writerEnglish, writerArabic, directorId, directorEnglish, directorArabic []string
		for _, actorIds := range contentActor {
			// actorId = append(actorId, actorIds.ActorId)
			actorEnglish = append(actorEnglish, actorIds.ActorEnglish)
			actorArabic = append(actorArabic, actorIds.ActorArabic)
			writerId = append(writerId, actorIds.WriterId)
			writerEnglish = append(writerEnglish, actorIds.WriterEnglish)
			writerArabic = append(writerArabic, actorIds.WriterArabic)
			directorId = append(directorId, actorIds.DirectorId)
			directorEnglish = append(directorEnglish, actorIds.DirectorEnglish)
			directorArabic = append(directorArabic, actorIds.DirectorArabic)
		}
		// contentSeason.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
		contentSeason.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
		contentSeason.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
		// contentSeason.Cast.WriterIds = common.RemoveDuplicateValues(writerId)
		contentSeason.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
		contentSeason.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
		contentSeason.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
		contentSeason.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
		contentSeason.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)
		/* fetching music details */
		var contentMusic []ContentMusic
		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_omposers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
			Joins("left join singer s on s.id =cs.singer_id").
			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
			Where("cs.music_id=?", eachcontent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var singerId, singerEnglish, singerArabic, composerId, composerEnglish, composerArabic, SongWriterId, SongWriterEnglish, SongWriterArabic []string
		for _, musicIds := range contentMusic {
			singerId = append(singerId, musicIds.SingerIds)
			singerEnglish = append(singerEnglish, musicIds.SingersEnglish)
			singerArabic = append(singerArabic, musicIds.SingersArabic)
			composerId = append(composerId, musicIds.MusicComposerIds)
			composerEnglish = append(composerEnglish, musicIds.MusicComposersEnglish)
			composerArabic = append(composerArabic, musicIds.MusicComposersArabic)
			SongWriterId = append(SongWriterId, musicIds.SongWriterIds)
			SongWriterEnglish = append(SongWriterEnglish, musicIds.SongWritersEnglish)
			SongWriterArabic = append(SongWriterArabic, musicIds.SongWritersArabic)
		}
		contentSeason.Music.MusicId = eachcontent.MusicId
		contentSeason.Music.SingerIds = common.RemoveDuplicateValues(singerId)
		contentSeason.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
		contentSeason.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
		contentSeason.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
		contentSeason.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
		contentSeason.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
		contentSeason.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
		contentSeason.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
		contentSeason.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
		/*fetching tag info */
		var contentTagInfo []ContentTag
		db.Table("content_tag ct").Select("tdt.name").
			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
			Where("ct.tag_info_id=?", eachcontent.TagInfoId).Scan(&contentTagInfo)
		var tagInfo []string
		for _, tagInfoIds := range contentTagInfo {
			tagInfo = append(tagInfo, tagInfoIds.Name)
		}
		contentSeason.TagInfo.Tags = tagInfo
		if len(tagInfo) < 1 {
			buffer := make([]string, 0)
			contentSeason.TagInfo.Tags = buffer
		}
		/*about the content*/
		contentSeason.AboutTheContent.OriginalLanguage = eachcontent.OriginalLanguage
		contentSeason.AboutTheContent.Supplier = eachcontent.Supplier
		contentSeason.AboutTheContent.AcquisitionDepartment = eachcontent.AcquisitionDepartment
		contentSeason.AboutTheContent.EnglishSynopsis = eachcontent.EnglishSynopsis
		contentSeason.AboutTheContent.ArabicSynopsis = eachcontent.ArabicSynopsis
		contentSeason.AboutTheContent.ProductionYear = eachcontent.ProductionYear
		contentSeason.AboutTheContent.ProductionHouse = eachcontent.ProductionHouse
		contentSeason.AboutTheContent.AgeGroup = string(eachcontent.AgeGroup)
		/*production countries*/
		var productionCountry []ProductionCountry
		if productionCountryResult := db.Table("production_country ").Select("country_id").Where("about_the_content_info_id=?", eachcontent.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var countries []int

		for _, prcountries := range productionCountry {
			countries = append(countries, prcountries.CountryId)
		}
		contentSeason.AboutTheContent.ProductionCountries = countries
		if len(tagInfo) < 1 {
			buffer := make([]int, 0)
			contentSeason.AboutTheContent.ProductionCountries = buffer
		}
		/*translation details*/
		contentSeason.Translation.LanguageType = common.LanguageOriginTypes(eachcontent.MultiTierLanguageType)
		contentSeason.Translation.DubbingLanguage = eachcontent.DubbingLanguage
		// contentSeason.Translation.DubbingDialectId = eachcontent.DubbingDialectId
		contentSeason.Translation.SubtitlingLanguage = eachcontent.SubtitlingLanguage
		/*non textual data for seasons*/
		if eachcontent.HasPosterImage {
			contentSeason.ContentNonTextualData.PosterImage = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("POSTER_IMAGE")
		}
		if eachcontent.HasOverlayPosterImage {
			contentSeason.ContentNonTextualData.OverlayPosterImage = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("OVERLAY_POSTER_IMAGE")
		}
		if eachcontent.HasDetailsBackground {
			contentSeason.ContentNonTextualData.DetailsBackground = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("DETAILS_BACKGROUND")
		}
		if eachcontent.HasMobileDetailsBackground {
			contentSeason.ContentNonTextualData.MobileDetailsBackground = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("MOBILE_DETAILS_BACKGROUND")
		}
		/*digital rights region season*/
		var digitalRightsRegions []DigitalRightsRegions
		if countryError := db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", eachcontent.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		var SeasonRights []int
		for _, idarr := range digitalRightsRegions {
			SeasonRights = append(SeasonRights, idarr.CountryId)
		}
		/*for digital rights*/
		var IsCheck bool
		for _, value := range SeasonRights {
			if CountryResult == int32(value) {
				IsCheck = true
			}
		}
		if country == "" || country == "all" {
			contentSeason.Rights.DigitalRightsRegions = SeasonRights
		}
		if len(SeasonRights) < 1 {
			buffer := make([]int, 0)
			contentSeason.Rights.DigitalRightsRegions = buffer
		}
		contentSeason.SeasonId = eachcontent.SeasonId
		/*for checking country rights*/

		/*primary info*/
		contentResult.PrimaryInfo.ContentType = eachcontent.ContentType
		contentResult.PrimaryInfo.OriginalTitle = eachcontent.OriginalTitle
		contentResult.PrimaryInfo.AlternativeTitle = eachcontent.AlternativeTitle
		contentResult.PrimaryInfo.ArabicTitle = eachcontent.ArabicTitle
		contentResult.PrimaryInfo.TransliteratedTitle = eachcontent.TransliteratedTitle
		contentResult.PrimaryInfo.Notes = eachcontent.Notes
		/*content genres*/
		var contentGenres []SeasonGenres
		var finalContentGenre []interface{}
		var newContentGenres NewSeasonGenres
		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
			Joins("left join genre g on g.id=cg.genre_id").
			Where("content_id=?", eachcontent.ContentId).Scan(&contentGenres).Error; genreResult != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		for _, tagInfoIds := range contentGenres {
			var contentSubgenre []SeasonSubgenre
			if subgenreResult := db.Table("content_subgenre csg").Select("english_name as sub_gener_english,arabic_name as sub_gener_arabic").
				Joins("left join subgenre sg on sg.id=csg.subgenre_id").
				Where("content_genre_id=?", tagInfoIds.Id).Scan(&contentSubgenre).Error; subgenreResult != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
			var SubgenreEn []string
			var SubgenreAr []string
			for _, data := range contentSubgenre {
				SubgenreEn = append(SubgenreEn, data.SubGenerEnglish)
				SubgenreAr = append(SubgenreAr, data.SubGenerArabic)
				newContentGenres.GenerEnglishName = tagInfoIds.GenerEnglishName
				newContentGenres.GenerArabicName = tagInfoIds.GenerArabicName
				newContentGenres.SubGenerEnglish = SubgenreEn
				newContentGenres.SubGenerArabic = SubgenreAr
				newContentGenres.Id = tagInfoIds.Id
				finalContentGenre = append(finalContentGenre, newContentGenres)
			}
		}
		contentResult.ContentGenres = finalContentGenre
		/*seo details*/
		contentResult.SeoDetails.EnglishMetaTitle = eachcontent.EnglishMetaTitle
		contentResult.SeoDetails.ArabicMetaTitle = eachcontent.ArabicMetaTitle
		contentResult.SeoDetails.EnglishMetaDescription = eachcontent.EnglishMetaDescription
		contentResult.SeoDetails.ArabicMetaDescription = eachcontent.ArabicMetaDescription
		contentResult.CreatedAt = eachcontent.CreatedAt
		contentResult.ModifiedAt = eachcontent.ModifiedAt
		//content id
		contentResult.ContentId = eachcontent.ContentId
		contentSeasons = append(contentSeasons, contentSeason)
		contentResult.ContentSeasons = contentSeasons

		if country == "" || country == "all" {
			allContents = append(allContents, contentResult)
		} else if country != "all" {
			if IsCheck {
				allContents = append(allContents, contentResult)
			}
		}

		// if country != "" {
		// 	if IsCheck {
		// 		allContents = append(allContents, contentResult)
		// 	}
		// } else if country == "" {
		// 	allContents = append(allContents, contentResult)
		// }

	}
	/*Pagination*/
	var pagination Pagination
	pagination.Limit = int(limit)
	pagination.Offset = int(offset)
	pagination.Size = totalCount
	if CountryResult != 0 || country == "" {
		c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})
	}
}
