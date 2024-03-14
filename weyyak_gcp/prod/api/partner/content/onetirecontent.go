package content

import (
	"fmt"
	common "masterdata/common"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllOneTierContentDetails - Get All One Tier Contents Details
// GET /v1/contents/onetier/all/
// @Description Get All One Tier Contents Details
// @Tags OneTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param offset query string false "Zero-based index of the first requested item." default(0)
// @Param limit query string false "Maximum number of collection items to return for a single request." default(5)
// @Param Country query string false "Country code of the user."
// @Success 200  object Response
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/contents/onetier/all [get]
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
	err := common.CheckAuthorization(c, UserId)
	fmt.Println("@@@@@@@@@@@@@@",UserId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization", "status": http.StatusUnauthorized})
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	// udb := c.MustGet("UDB").(*gorm.DB)

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
			Where(`c.watch_now_supplier = true AND
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.third_party_content_key != 0 AND 
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.id IS NOT NULL`).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult)

		// db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
		// 	Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
		// 	Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
		// 	Where(`c.watch_now_supplier = true AND
		// 	c.status = 1 AND
		// 	c.content_tier = 1 AND
		// 	c.third_party_content_key != 0 AND
		// 	c.deleted_by_user_id is null AND
		// 	c.id IS NOT NULL`).Count(&totalCount)

		if CountryResult != 0 {
			if err := db.Debug().Table("content c").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
				Joins("join content_variance cv on cv.content_id =c.id").
				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
				Joins("join content_translation ct on ct.id =pi2.translation_id").
				Joins("join content_rights cr on cr.id =pi2.rights_id").
				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
				Where("c.watch_now_supplier = true and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Count(&totalCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country == "" || country == "all" {
			if err := db.Debug().Table("content c").
				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
				Joins("join content_variance cv on cv.content_id =c.id").
				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
				Joins("join content_translation ct on ct.id =pi2.translation_id").
				Joins("join content_rights cr on cr.id =pi2.rights_id").
				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
				Where("c.watch_now_supplier = true and c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Count(&totalCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		}

	} else {
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.third_party_content_key != 0 AND
			c.deleted_by_user_id is null AND 
			-- c.id = '0c5647c7-1676-443a-a4a9-6d18a4230d5a' AND
			c.id IS NOT NULL`).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult)

		// db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
		// 	Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
		// 	Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
		// 	Where(`
		// 	c.status = 1 AND
		// 	c.content_tier = 1 AND
		// 	c.third_party_content_key != 0 AND
		// 	c.deleted_by_user_id is null AND
		// 	c.id IS NOT NULL`).Count(&totalCount)

		if CountryResult != 0 {
			if err := db.Debug().Table("content c").Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
				Joins("join content_variance cv on cv.content_id =c.id").
				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
				Joins("join content_translation ct on ct.id =pi2.translation_id").
				Joins("join content_rights cr on cr.id =pi2.rights_id").
				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
				Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
				Where(" c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null) and crc.country_id = ?", CountryResult).Count(&totalCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country == "" || country == "all" {
			if err := db.Debug().Table("content c").
				Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
				Joins("join content_variance cv on cv.content_id =c.id").
				Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
				Joins("join content_translation ct on ct.id =pi2.translation_id").
				Joins("join content_rights cr on cr.id =pi2.rights_id").
				Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
				Where("c.status = 1 and c.content_tier = 1 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW()	or cr.digital_rights_end_date is null)").Count(&totalCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		}
		// else {
		// 	continue
		// }
	}

	haveVariances := false

	for _, singleContent := range finalContentResult {

		var contentVariances []ContentVariancesSource

		var ContentVariancesRecord []map[string]interface{}

		flag := false

		if CountryResult != 0 {
			db.Debug().Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
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
			db.Debug().Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
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

		if len(contentVariances) > 0 {
			haveVariances = true
		}

		for _, cv := range contentVariances {
			var dubbingScript string
			var subtitlingScript, OverlayPosterImage string

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
			// var trailerInfo []VarianceTrailers
			// var digitalRights []int

			db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

			// if CountryResult == 0 {

			// 	db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

			// 	for _, region := range digitalRightsRegions {
			// 		var digitalRightsCountry Country
			// 		if country != "" && country != "all" && country != "All" {
			// 			udb.Table("country").Where("id=?", region.CountryId).Find(&digitalRightsCountry)
			// 			if country == digitalRightsCountry.Alpha2code {
			// 				digitalRights = append(digitalRights, region.CountryId)
			// 			}
			// 		} else {
			// 			// fmt.Println(common.CountryCount(), "----------------", len(digitalRightsRegions))
			// 			if common.CountryCount() == len(digitalRightsRegions) {
			// 				digitalRights = append(digitalRights, region.CountryId)
			// 			} else {
			// 				haveVariances = false
			// 			}

			// 		}
			// 	}
			// }

			var digitalRights []int

			fmt.Println("CountryResult--->", CountryResult)

			if CountryResult == 0 {

				if country != "" {
					if country == "all" || country == "All" {

						if common.CountryCount() == len(digitalRightsRegions) {
							digitalRights = nil
						} else {
							// c.JSON(http.StatusInternalServerError, serverError)
							// return
							flag = true
						}

					}
				} else if country == "" {
					for _, region := range digitalRightsRegions {
						digitalRights = append(digitalRights, region.CountryId)
					}
				}

			} else if CountryResult == 1 {

				for _, region := range digitalRightsRegions {
					digitalRights = append(digitalRights, region.CountryId)
				}

			} else {
				var isCheck bool = false
				for _, region := range digitalRightsRegions {
					if CountryResult == int32(region.CountryId) {
						isCheck = true
					}
				}

				if !isCheck {
					// c.JSON(http.StatusInternalServerError, serverError)
					// return
					flag = true
				}
			}

			// db.Table("variance_trailer").Where("content_variance_id=?", cv.Id).Scan(&trailerInfo)

			var trailerInfo, finaltrailerInfo []VarianceTrailers

			db.Table("variance_trailer").Where("content_variance_id=?", cv.Id).Scan(&trailerInfo)
			var url = "https://" + os.Getenv("S3_BUCKET") + ".s3.amazonaws.com"
			for _, trailer := range trailerInfo {
				var posterTrailerImage string
				if trailer.HasTrailerPosterImage {
					posterTrailerImage = url + "/" + singleContent.Id + "/" + cv.Id + "/" + trailer.Id + "/trailer-poster-image"
				} else {
					posterTrailerImage = ""
				}

				finaltrailerInfo = append(finaltrailerInfo, VarianceTrailers{
					Order:                 trailer.Order,
					VideoTrailerId:        trailer.VideoTrailerId,
					EnglishTitle:          trailer.EnglishTitle,
					ArabicTitle:           trailer.ArabicTitle,
					Duration:              trailer.Duration,
					HasTrailerPosterImage: trailer.HasTrailerPosterImage,
					TrailerPosterImage:    posterTrailerImage,
					Id:                    trailer.Id,
					SeasonId:              trailer.SeasonId,
				})

			}

			if cv.HasOverlayPosterImage {
				OverlayPosterImage = os.Getenv("IMAGERY_URL") + singleContent.Id + "/" + cv.Id + "/overlay-poster-image"
			}

			var contentVariance map[string]interface{}
			// {1: "Original", 2: "Dubbed", 3: "Subtitled"}
			if cv.LanguageType == 2 {

				if cv.DubbingLanguage == "ar" {
					contentVariance = map[string]interface{}{
						"id":                   cv.Id,
						"length":               cv.Length,
						"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
						"languageType":         common.LanguageOriginTypes(cv.LanguageType),
						"dubbingScript":        dubbingScript,
						"subtitlingScript":     subtitlingScript,
						"dubbingLanguage":      cv.DubbingLanguage,
						"dubbingDialectName":   common.DialectIdname(cv.DubbingDialectId, "en"),
						"digitalRightsRegions": digitalRights,
						"trailerInfo":          finaltrailerInfo,
						"overlayPosterImage":   OverlayPosterImage,
					}
				} else {
					contentVariance = map[string]interface{}{
						"id":                   cv.Id,
						"length":               cv.Length,
						"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
						"languageType":         common.LanguageOriginTypes(cv.LanguageType),
						"dubbingScript":        dubbingScript,
						"subtitlingScript":     subtitlingScript,
						"dubbingLanguage":      cv.DubbingLanguage,
						"digitalRightsRegions": digitalRights,
						"trailerInfo":          finaltrailerInfo,
						"overlayPosterImage":   OverlayPosterImage,
					}
				}

			} else if cv.LanguageType == 3 {
				contentVariance = map[string]interface{}{
					"id":                   cv.Id,
					"length":               cv.Length,
					"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
					"languageType":         common.LanguageOriginTypes(cv.LanguageType),
					"dubbingScript":        dubbingScript,
					"subtitlingScript":     subtitlingScript,
					"digitalRightsRegions": digitalRights,
					"trailerInfo":          finaltrailerInfo,
					"overlayPosterImage":   OverlayPosterImage,
				}
			} else {
				contentVariance = map[string]interface{}{
					"id":                   cv.Id,
					"length":               cv.Length,
					"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
					"languageType":         common.LanguageOriginTypes(cv.LanguageType),
					"dubbingScript":        dubbingScript,
					"subtitlingScript":     subtitlingScript,
					"digitalRightsRegions": digitalRights,
					"trailerInfo":          finaltrailerInfo,
					"overlayPosterImage":   OverlayPosterImage,
				}
			}

			if !flag {
				fmt.Println("frggggggggggggggggggg")
				ContentVariancesRecord = append(ContentVariancesRecord, contentVariance)
			}

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

		_ = haveVariances
		if len(ContentVariancesRecord) != 0 {
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

	}

	pagination := Pagination{
		Size:   totalCount,
		Limit:  int(limit),
		Offset: int(offset),
	}

	// m, _ := json.Marshal(RedisResponse{
	// 	Pagination: pagination,
	// 	Data:       allContents,
	// })

	// err = common.PostRedisDataWithKey(redisKey, m)
	// if err != nil {
	// 	fmt.Println("Redis Value Not Updated")
	// }

	c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})
}

// GetOneTierContentDetailsBasedonContentID
// GET /v1/contents/onetier/:contentId
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
// @Router /v1/contents/onetier/{contentId} [get]
func (hs *HandlerService) GetOneTierContentDetailsBasedonContentID(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	UserId := c.MustGet("userid")
	if UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization1", "status": http.StatusUnauthorized})
		return
	}
	err := common.CheckAuthorization(c, UserId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization2", "status": http.StatusUnauthorized})
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	udb := c.MustGet("UDB").(*gorm.DB)

	serverError := common.ServerErrorResponse()
	notFound := common.NotFoundErrorResponse()

	var finalContentResult []FinalSeasonResultOneTire
	var allContents AllOnetierContent

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
	} else if c.Request.URL.Query()["country"] != nil {
		country = c.Request.URL.Query()["country"][0]
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
			Where(`c.watch_now_supplier = 'true' AND
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			c.third_party_content_key = ? `, content_key).Find(&finalContentResult)

	} else {
		fmt.Println("inside else ..........")		
		db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
			Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
			Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
			Where(`
			c.status = 1 AND 
			c.content_tier = 1 AND 
			c.deleted_by_user_id is null AND 
			c.third_party_content_key = ? `, content_key).Find(&finalContentResult)
	}

	if len(finalContentResult) == 0 {
		c.JSON(http.StatusInternalServerError, serverError)
		return
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
						cv.status = 1 and 
						cv.deleted_by_user_id IS NULL and 
						cv.content_id = ?;`, CountryResult, singleContent.Id).Find(&contentVariances)

			if len(contentVariances) == 0 {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
		} else if country == "" || country == "all" || country == "All" {

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
						cv.status = 1 and 
						cv.deleted_by_user_id IS NULL and 
						cv.content_id = ?;`, singleContent.Id).Find(&contentVariances)

			if len(contentVariances) == 0 {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}
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

			var trailerInfo, finaltrailerInfo []VarianceTrailers
			var digitalRights []int
			var digitalRightsRegions []DigitalRightsRegions
			if CountryResult == 0 {

				db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

				for _, region := range digitalRightsRegions {
					var digitalRightsCountry Country
					if country != "" && country != "all" && country != "All" {
						udb.Table("country").Where("id=?", region.CountryId).Find(&digitalRightsCountry)
						if country == digitalRightsCountry.Alpha2code {
							digitalRights = append(digitalRights, region.CountryId)
						}
					} else {

						if common.CountryCount() == len(digitalRightsRegions) && country == "" {
							digitalRights = append(digitalRights, region.CountryId)
						}

						// else {
						// 	c.JSON(http.StatusInternalServerError, serverError)
						// 	return
						// }

					}
				}
			}

			db.Table("variance_trailer").Where("content_variance_id=?", cv.Id).Scan(&trailerInfo)
			var url = "https://" + os.Getenv("S3_BUCKET") + ".s3.amazonaws.com"
			for _, trailer := range trailerInfo {
				var posterTrailerImage string
				if trailer.HasTrailerPosterImage {
					posterTrailerImage = url + "/" + singleContent.Id + "/" + cv.Id + "/" + trailer.Id + "/trailer-poster-image"
				} else {
					posterTrailerImage = ""
				}

				finaltrailerInfo = append(finaltrailerInfo, VarianceTrailers{
					Order:                 trailer.Order,
					VideoTrailerId:        trailer.VideoTrailerId,
					EnglishTitle:          trailer.EnglishTitle,
					ArabicTitle:           trailer.ArabicTitle,
					Duration:              trailer.Duration,
					HasTrailerPosterImage: trailer.HasTrailerPosterImage,
					TrailerPosterImage:    posterTrailerImage,
					Id:                    trailer.Id,
					SeasonId:              trailer.SeasonId,
				})

			}

			if cv.HasOverlayPosterImage {
				OverlayPosterImage = os.Getenv("IMAGERY_URL") + singleContent.Id + "/" + cv.Id + "/overlay-poster-image"
			}

			var contentVariance map[string]interface{}
			// {1: "Original", 2: "Dubbed", 3: "Subtitled"}
			if cv.LanguageType == 2 {

				if cv.DubbingLanguage == "ar" {
					contentVariance = map[string]interface{}{
						"id":                   cv.Id,
						"length":               cv.Length,
						"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
						"languageType":         common.LanguageOriginTypes(cv.LanguageType),
						"dubbingScript":        dubbingScript,
						"subtitlingScript":     subtitlingScript,
						"dubbingLanguage":      cv.DubbingLanguage,
						"dubbingDialectName":   common.DialectIdname(cv.DubbingDialectId, "en"),
						"digitalRightsRegions": digitalRights,
						"trailerInfo":          finaltrailerInfo,
						"overlayPosterImage":   OverlayPosterImage,
					}
				} else {
					contentVariance = map[string]interface{}{
						"id":                   cv.Id,
						"length":               cv.Length,
						"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
						"languageType":         common.LanguageOriginTypes(cv.LanguageType),
						"dubbingScript":        dubbingScript,
						"subtitlingScript":     subtitlingScript,
						"dubbingLanguage":      cv.DubbingLanguage,
						"digitalRightsRegions": digitalRights,
						"trailerInfo":          finaltrailerInfo,
						"overlayPosterImage":   OverlayPosterImage,
					}
				}

			} else if cv.LanguageType == 3 {
				contentVariance = map[string]interface{}{
					"id":                   cv.Id,
					"length":               cv.Length,
					"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
					"languageType":         common.LanguageOriginTypes(cv.LanguageType),
					"dubbingScript":        dubbingScript,
					"subtitlingScript":     subtitlingScript,
					"digitalRightsRegions": digitalRights,
					"trailerInfo":          finaltrailerInfo,
					"overlayPosterImage":   OverlayPosterImage,
				}
			} else {
				contentVariance = map[string]interface{}{
					"id":                   cv.Id,
					"length":               cv.Length,
					"videoContentUrl":      os.Getenv("VIDEO_API") + cv.VideoContentId,
					"languageType":         common.LanguageOriginTypes(cv.LanguageType),
					"dubbingScript":        dubbingScript,
					"subtitlingScript":     subtitlingScript,
					"digitalRightsRegions": digitalRights,
					"trailerInfo":          finaltrailerInfo,
					"overlayPosterImage":   OverlayPosterImage,
				}
			}

			ContentVariancesRecord = append(ContentVariancesRecord, contentVariance)

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
			Where("cs.music_id=?", singleContent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
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

		allContents = AllOnetierContent{
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
		}

	}

	c.JSON(http.StatusOK, gin.H{"data": allContents})

}
