package content

import (
	"fmt"
	"log"
	common "masterdata/common"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (hs *HandlerService) GetModifiedContentDetails(c *gin.Context) {
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

	fmt.Println(startdate, enddate)

	enddate = enddate.AddDate(0, 0, 1)
	days := daysBetween(startdate, enddate)

	if days > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide dates within 7 days range"})
	} else {
		var finalOnetierContentResult []FinalSeasonResultOneTire
		var totalCount int

		var contentResultFinal []MultiTierContent2
		var finalContent []FinalSeasonResultContentOneTire
		var totalMultiCount []FinalSeasonResultContentOneTire

		finalData := make(map[string]interface{}, 0)

		var CountryResult int32

		var country string
		if c.Request.URL.Query()["Country"] != nil {
			country = c.Request.URL.Query()["Country"][0]
		}

		CountryResult = common.Countrys(country)
		serverError := common.ServerErrorResponse()

		oneTier := []string{"movie", "livetv", "series"}
		// multiTier := []string{"series", "season", "episode"}

		for _, types := range oneTier {
			var allContents []AllOnetierContent
			var querytype string
			if types == "movie" {
				querytype = "Movie"
			} else {
				querytype = "LiveTV"
			}
			if UserId == os.Getenv("WATCH_NOW") {
				db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
					Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
					Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
					Where(`c.watch_now_supplier = true AND
							c.content_type = ? AND 
							c.deleted_by_user_id is null AND 
							c.modified_at >= ? AND c.modified_at <= ? AND
							c.id IS NOT NULL`, querytype, startdate, enddate).Order("c.modified_at desc").Find(&finalOnetierContentResult)
				if CountryResult != 0 {
					if err := db.Debug().Table("content c").
						Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
						Joins("join content_variance cv on cv.content_id =c.id").
						Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
						Joins("join content_rights cr on cr.id =pi2.rights_id").
						Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
						Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
						Where("c.watch_now_supplier = true and c.content_type = ? and c.deleted_by_user_id is null and crc.country_id = ? and c.modified_at >= ? AND c.modified_at <= ?", querytype, CountryResult, startdate, enddate).Count(&totalCount).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				} else if country == "" || country == "all" {
					if err := db.Debug().Table("content c").
						Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
						Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
						Where("c.watch_now_supplier = true and c.content_type = ? and c.deleted_by_user_id is null and c.modified_at >= ? AND c.modified_at <= ?", querytype, startdate, enddate).Count(&totalCount).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			} else {
				db.Debug().Table("content c").Select(`c.id, c.third_party_content_key as content_key, c.primary_info_id, c.content_type, cpi.original_title, cpi.alternative_title , cpi.arabic_title , cpi.transliterated_title , cpi.notes, c.cast_id, c.music_id, c.tag_info_id, atci.original_language , atci.supplier , atci.acquisition_department , atci.english_synopsis , atci.arabic_synopsis , atci.production_year , atci.production_house , atci.age_group , atci.outro_start as about_outro_start, c.about_the_content_info_id, c.english_meta_title, c.arabic_meta_title, c.english_meta_description, c.arabic_meta_description, c.has_poster_image, c.has_details_background, c.has_mobile_details_background, c.created_at, c.modified_at`).
					Joins("join content_primary_info cpi ON cpi.id = c.primary_info_id").
					Joins("join about_the_content_info atci on atci.id = c.about_the_content_info_id and atci.supplier !='Others'").
					Where(`
				c.content_type = ? AND 
				c.deleted_by_user_id is null AND 
				c.modified_at >= ? AND c.modified_at <= ? AND
				c.id IS NOT NULL`, querytype, startdate, enddate).Order("c.modified_at desc").Find(&finalOnetierContentResult)

				if CountryResult != 0 {
					if err := db.Debug().Table("content c").
						Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
						Joins("join content_variance cv on cv.content_id =c.id").
						Joins("join playback_item pi2 on pi2.id =cv.playback_item_id").
						Joins("join content_rights cr on cr.id =pi2.rights_id").
						Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id").
						Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
						Where("c.content_type = ? and c.deleted_by_user_id is null and crc.country_id = ? and c.modified_at >= ? AND c.modified_at <= ?", querytype, CountryResult, startdate, enddate).Count(&totalCount).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				} else if country == "" || country == "all" {
					if err := db.Debug().Table("content c").
						Joins("join content_primary_info cpi on cpi.id =c.primary_info_id").
						Joins("join about_the_content_info atci on atci.id=c.about_the_content_info_id and atci.supplier !='Others'").
						Where("c.content_type = ? and c.deleted_by_user_id is null and c.modified_at >= ? AND c.modified_at <= ?", querytype, startdate, enddate).Count(&totalCount).Error; err != nil {
						c.JSON(http.StatusInternalServerError, serverError)
						return
					}
				}
			}

			haveVariances := false

			for _, singleContent := range finalOnetierContentResult {
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
							crc.country_id = ? and
							cv.content_id = ?;`, CountryResult, singleContent.Id).Find(&contentVariances)
				} else if country == "" || country == "all" {
					db.Debug().Raw(`select cv.id, pi2.duration as length, pi2.video_content_id, ct.language_type, cv.has_dubbing_script, ct.dubbing_dialect_id,
							cv.has_subtitling_script, ct.dubbing_language, pi2.rights_id, cv.has_overlay_poster_image from content_variance cv
							join playback_item pi2 on pi2.id = cv.playback_item_id
							join content_translation ct on ct.id = pi2.translation_id
							join content_rights cr on cr.id =pi2.rights_id
							where
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
					db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", cv.RightsId).Scan(&digitalRightsRegions)

					var digitalRights []int
					fmt.Println("CountryResult--->", CountryResult)

					if CountryResult == 0 {
						if country != "" {
							if country == "all" || country == "All" {
								if common.CountryCount() == len(digitalRightsRegions) {
									digitalRights = nil
								} else {
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
							flag = true
						}
					}

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
					if cv.LanguageType == 2 {
						if cv.DubbingLanguage == "ar" {
							contentVariance = map[string]interface{}{
								"id":                   cv.Id,
								"length":               cv.Length,
								"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
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
								"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
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
							"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
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
							"videoContentUrl":      os.Getenv("CONTENT_URL") + cv.VideoContentId,
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
				// if len(ContentVariancesRecord) != 0 {
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
			finalData[types] = allContents

			if types == "series" {
				if UserId == os.Getenv("WATCH_NOW") {
					if CountryResult != 0 {
						db.Debug().Raw(`
								select
									c.third_party_content_key AS multi_tier_content_key,
									c.id,
									c.content_type,
									cpi.original_title,
									cpi.alternative_title,
									cpi.arabic_title,
									cpi.transliterated_title,
									cpi.notes,
									c.english_meta_title,
									c.arabic_meta_title,
									c.english_meta_description,
									c.arabic_meta_description,
									c.id AS content_id,
									c.created_at,
									c.modified_at
								from
									content c
									
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
								join season s on s.content_id = c.id
								join content_rights cr on cr.id = s.rights_id
								join content_rights_country crc on crc.content_rights_id = cr.id
									
								where
								c.watch_now_supplier = 'true'
								and c.content_tier = 2
									
									and c.deleted_by_user_id is null
		
									and c.modified_at >= ? AND c.modified_at <= ?
										
									and crc.country_id = ?
							
					`, startdate, enddate, CountryResult).Order("c.modified_at desc").Find(&finalContent)

						db.Raw(`    select c.third_party_content_key
								from
									content c
									
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
								join season s on s.content_id = c.id
								join content_rights cr on cr.id = s.rights_id
								join content_rights_country crc on crc.content_rights_id = cr.id
									
								where
									c.watch_now_supplier = 'true'
									and c.content_tier = 2
									and c.deleted_by_user_id is null
									and c.modified_at >= ? AND c.modified_at <= ?
									and crc.country_id = ?
									`, startdate, enddate, CountryResult).Find(&totalMultiCount)

						fmt.Println("--------->", len(totalMultiCount))

					} else if country == "" || country == "all" {

						db.Debug().Raw(`
								select
									c.third_party_content_key AS multi_tier_content_key,
									c.id,
									c.content_type,
									cpi.original_title,
									cpi.alternative_title,
									cpi.arabic_title,
									cpi.transliterated_title,
									cpi.notes,
									c.english_meta_title,
									c.arabic_meta_title,
									c.english_meta_description,
									c.arabic_meta_description,
									c.id AS content_id,
									c.created_at,
									c.modified_at
								from content c
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
								where
									c.watch_now_supplier = 'true'
		
									and c.content_tier = 2
																		
									and c.deleted_by_user_id is null
		
									and c.modified_at >= ? AND c.modified_at <= ? 
					`, startdate, enddate).Order("c.modified_at desc").Find(&finalContent)

						db.Raw(`
								select
								distinct c.third_party_content_key
							from
								content c
								
							JOIN content_primary_info cpi on cpi.id = c.primary_info_id
							where
								c.watch_now_supplier = 'true'
								and c.content_tier = 2
								and c.deleted_by_user_id is null
								and c.modified_at >= ? AND c.modified_at <= ?
					`, startdate, enddate).Find(&totalMultiCount)
					}
				} else {
					if CountryResult != 0 {
						db.Debug().Raw(`
								select
									c.third_party_content_key AS multi_tier_content_key,
									c.id,
									c.content_type,
									cpi.original_title,
									cpi.alternative_title,
									cpi.arabic_title,
									cpi.transliterated_title,
									cpi.notes,
									c.english_meta_title,
									c.arabic_meta_title,
									c.english_meta_description,
									c.arabic_meta_description,
									c.id AS content_id,
									c.created_at,
									c.modified_at
								from
									content c
									
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
								join season s on s.content_id = c.id
								join content_rights cr on cr.id = s.rights_id
								join content_rights_country crc on crc.content_rights_id = cr.id
									
								where
									c.content_tier = 2
																	
									and c.deleted_by_user_id is null
		
									and c.modified_at >= ? AND c.modified_at <= ?

									and crc.country_id = ?
							
					`, startdate, enddate, CountryResult).Order("c.modified_at desc").Find(&finalContent)

						db.Raw(`
						select
						c.third_party_content_key
					from
						content c
						
					JOIN content_primary_info cpi on cpi.id = c.primary_info_id
					join season s on s.content_id = c.id
					join content_rights cr on cr.id = s.rights_id
					join content_rights_country crc on crc.content_rights_id = cr.id
						
					where
						c.content_tier = 2
							and c.deleted_by_user_id is null
							and c.modified_at >= ? AND c.modified_at <= ?
							and crc.country_id = ?
					`, startdate, enddate, CountryResult).Find(&totalMultiCount)
					} else if country == "" || country == "all" {

						db.Debug().Raw(`
								select
									c.third_party_content_key AS multi_tier_content_key,
									c.id,
									c.content_type,
									cpi.original_title,
									cpi.alternative_title,
									cpi.arabic_title,
									cpi.transliterated_title,
									cpi.notes,
									c.english_meta_title,
									c.arabic_meta_title,
									c.english_meta_description,
									c.arabic_meta_description,
									c.id AS content_id,
									c.created_at,
									c.modified_at
								from
									content c
									
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
									
								where
									c.content_tier = 2
																		
									and c.deleted_by_user_id is null
		
									and c.modified_at >= ? AND c.modified_at <= ?
							
						`, startdate, enddate).Order("c.modified_at desc").Find(&finalContent)

						db.Raw(`  	select
									c.third_party_content_key
								from
									content c
									
								JOIN content_primary_info cpi on cpi.id = c.primary_info_id
									
								where
									c.content_tier = 2
									and c.deleted_by_user_id is null
									and c.modified_at >= ? AND c.modified_at <= ?
							`, startdate, enddate).Find(&totalMultiCount)
					}
				}
				for _, finalContentResult := range finalContent {

					var contentResult MultiTierContent2

					contentResult.ContentKey = finalContentResult.MultiTierContentKey
					contentResult.PrimaryInfo = ContentPrimaryInfo{
						ContentType:         finalContentResult.ContentType,
						OriginalTitle:       finalContentResult.OriginalTitle,
						AlternativeTitle:    finalContentResult.AlternativeTitle,
						ArabicTitle:         finalContentResult.ArabicTitle,
						TransliteratedTitle: finalContentResult.TransliteratedTitle,
						Notes:               finalContentResult.Notes,
					}

					var contentGenres []MultiTierContentGenres

					var contentGenresIds []struct {
						ContentId string `json:"content_id"`
						Order     string `json:"order"`
						Id        string `json:"id"`
						GenreId   string `json:"genre_id"`
					}

					db.Table("content_genre").Where("content_id =?", finalContentResult.ContentId).Find(&contentGenresIds)

					for _, querygenresId := range contentGenresIds {

						var genresName struct {
							EnglishName string `json:"english_name"`
							ArabicName  string `json:"arabic_name"`
						}

						var finalcontentSubGenresEnglish, finalcontentSubGenresArabic []string

						var subgenreId []string

						var contentSubGenres []struct {
							EnglishName string `json:"english_name"`
							ArabicName  string `json:"arabic_name"`
						}

						db.Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName)

						db.Table("content_subgenre").Where("content_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId)

						db.Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&contentSubGenres)

						for _, contentSubGenre := range contentSubGenres {
							finalcontentSubGenresEnglish = append(finalcontentSubGenresEnglish, contentSubGenre.EnglishName)
							finalcontentSubGenresArabic = append(finalcontentSubGenresArabic, contentSubGenre.ArabicName)
						}

						contentGenres = append(contentGenres, MultiTierContentGenres{
							Id:               querygenresId.GenreId,
							GenerEnglishName: genresName.EnglishName,
							GenerArabicName:  genresName.ArabicName,
							SubGenerEnglish:  finalcontentSubGenresEnglish,
							SubGenerArabic:   finalcontentSubGenresArabic,
						})

					}

					contentResult.ContentGenres = contentGenres

					var contentSeasons []ContentSeasons2
					var seasons []Seasons

					db.Raw(`
						select
							s.*
						from
							season s
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						
						where
							s.deleted_by_user_id is null and
							atci.supplier !='Others' and
							s.content_id = ? and s.modified_at >= ? AND s.modified_at <= ?
						`, finalContentResult.ContentId, startdate, enddate).Find(&seasons)

					for _, seasonId := range seasons {
						var seasonPrimaryInfo SeasonPrimaryInfo

						db.Table("content_primary_info").Where("id =?", seasonId.PrimaryInfoId).Find(&seasonPrimaryInfo)

						seasonPrimaryInfo.SeasonNumber = seasonId.Number

						var cast Cast

						var castList struct {
							CastId                 string `json:"cast_id"`
							MainActorId            string `json:"main_actor_id"`
							MainActressId          string `json:"main_actress_id"`
							MainActorEnglishName   string `json:"main_actor_english_name"`
							MainActorArabicName    string `json:"main_actor_arabic_name"`
							MainActressEnglishName string `json:"main_actress_english_name"`
							MainActressArabicName  string `json:"main_actress_arabic_name"`
						}

						db.Raw(`
									select
										s.cast_id as cast_id,
										cc.main_actor_id as main_actor_id,
										cc.main_actress_id as main_actress_id,
										actor.english_name as main_actor_english_name,
										actor.arabic_name as main_actor_arabic_name,
										actress.english_name as main_actress_english_name,
										actress.arabic_name as main_actress_arabic_name
									from
										season s
									join content_cast cc on
										cc.id = s.cast_id
										
									join actor actor on
										actor.id = cc.main_actor_id
									
									join actor actress on
										actress.id = cc.main_actress_id
									
									join content_translation ct on ct.id = s.translation_id
										
									where
										s.deleted_by_user_id is null and
										s.modified_at >= ? AND s.modified_at <= ? and
										s.id = ?
								`, startdate, enddate, seasonId.Id).Find(&castList)

						var actorIds, actorEnglishNames, actorArabicNames []string
						var queryActorIds []struct {
							ActorId string `json:"actor_id"`
						}

						db.Table("content_actor").Where("cast_id =?", castList.CastId).Find(&queryActorIds)

						for _, queryActorId := range queryActorIds {
							actorIds = append(actorIds, queryActorId.ActorId)

							var actorName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("actor").Where("id =?", queryActorId.ActorId).Find(&actorName)

							actorEnglishNames = append(actorEnglishNames, actorName.EnglishName)
							actorArabicNames = append(actorArabicNames, actorName.ArabicName)
						}

						var writerIds, writerEnglishNames, writerArabicNames []string
						var queryWriterIds []struct {
							WriterId string `json:"writer_id"`
						}

						db.Table("content_writer").Where("cast_id =?", castList.CastId).Find(&queryWriterIds)

						for _, queryWriterId := range queryWriterIds {
							writerIds = append(writerIds, queryWriterId.WriterId)

							var writerName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("writer").Where("id =?", queryWriterId.WriterId).Find(&writerName)

							writerEnglishNames = append(writerEnglishNames, writerName.EnglishName)
							writerArabicNames = append(writerArabicNames, writerName.ArabicName)
						}

						var directorIds, directorEnglishNames, directorArabicNames []string
						var querydirectorIds []struct {
							DirectorId string `json:"director_id"`
						}

						db.Table("content_director").Where("cast_id =?", castList.CastId).Find(&querydirectorIds)

						for _, querydirectorId := range querydirectorIds {
							directorIds = append(directorIds, querydirectorId.DirectorId)

							var directorName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("director").Where("id =?", querydirectorId.DirectorId).Find(&directorName)

							directorEnglishNames = append(directorEnglishNames, directorName.EnglishName)
							directorArabicNames = append(directorArabicNames, directorName.ArabicName)
						}

						cast = Cast{
							CastId:             castList.CastId,
							MainActorId:        castList.MainActorId,
							MainActressId:      castList.MainActressId,
							MainActorEnglish:   castList.MainActorEnglishName,
							MainActorArabic:    castList.MainActorArabicName,
							MainActressEnglish: castList.MainActressEnglishName,
							MainActressArabic:  castList.MainActressArabicName,
							ActorIds:           actorIds,
							ActorEnglish:       actorEnglishNames,
							ActorArabic:        actorArabicNames,
							WriterId:           writerIds,
							WriterEnglish:      writerEnglishNames,
							WriterArabic:       writerArabicNames,
							DirectorIds:        directorIds,
							DirectorEnglish:    directorEnglishNames,
							DirectorArabic:     directorArabicNames,
						}

						var music Music

						var singerIds, singerEnglishNames, singerArabicNames []string
						var querySingerIds []struct {
							SingerId string `json:"singer_id"`
						}

						db.Table("content_singer").Where("music_id =?", seasonId.MusicId).Find(&querySingerIds)

						for _, querySingerId := range querySingerIds {
							singerIds = append(singerIds, querySingerId.SingerId)

							var SingerName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("singer").Where("id =?", querySingerId.SingerId).Find(&SingerName)

							singerEnglishNames = append(singerEnglishNames, SingerName.EnglishName)
							singerArabicNames = append(singerArabicNames, SingerName.ArabicName)
						}

						var musicComposerIds, musicComposerEnglishNames, musicComposerArabicNames []string
						var queryMusicComposerIds []struct {
							MusicComposerId string `json:"music_composer_id"`
						}

						db.Table("content_music_composer").Where("music_id =?", seasonId.MusicId).Find(&queryMusicComposerIds)

						for _, queryMusicComposerId := range queryMusicComposerIds {
							musicComposerIds = append(musicComposerIds, queryMusicComposerId.MusicComposerId)

							var MusicComposerName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("music_composer").Where("id =?", queryMusicComposerId.MusicComposerId).Find(&MusicComposerName)

							musicComposerEnglishNames = append(musicComposerEnglishNames, MusicComposerName.EnglishName)
							musicComposerArabicNames = append(musicComposerArabicNames, MusicComposerName.ArabicName)
						}

						var songWriterIds, songWriterEnglishNames, songWriterArabicNames []string
						var querySongWriterIds []struct {
							SongWriterId string `json:"song_writer_id"`
						}

						db.Table("content_song_writer").Where("music_id =?", seasonId.MusicId).Find(&querySongWriterIds)

						for _, querySongWriterId := range querySongWriterIds {
							songWriterIds = append(songWriterIds, querySongWriterId.SongWriterId)

							var SongWriterName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Table("song_writer").Where("id =?", querySongWriterId.SongWriterId).Find(&SongWriterName)

							songWriterEnglishNames = append(songWriterEnglishNames, SongWriterName.EnglishName)
							songWriterArabicNames = append(songWriterArabicNames, SongWriterName.ArabicName)
						}

						music = Music{
							MusicId:               seasonId.MusicId,
							SingerIds:             singerIds,
							SingersEnglish:        singerEnglishNames,
							SingersArabic:         singerArabicNames,
							MusicComposerIds:      musicComposerIds,
							MusicComposersEnglish: musicComposerEnglishNames,
							MusicComposersArabic:  musicComposerArabicNames,
							SongWriterIds:         songWriterIds,
							SongWritersEnglish:    songWriterEnglishNames,
							SongWritersArabic:     songWriterArabicNames,
						}

						var tagInfo TagInfo

						var queryContentTagIds []struct {
							TagInfoId        string `json:"tag_info_id"`
							TextualDataTagId string `json:"textual_data_tag_id"`
						}

						db.Table("content_tag").Where("tag_info_id =?", seasonId.TagInfoId).Find(&queryContentTagIds)

						var ContentTagIds []string
						for _, queryContentTagId := range queryContentTagIds {

							var TextualDataTagName struct {
								Name string `json:"name"`
							}

							db.Table("textual_data_tag").Where("id =?", queryContentTagId.TextualDataTagId).Find(&TextualDataTagName)

							ContentTagIds = append(ContentTagIds, TextualDataTagName.Name)
						}

						tagInfo.Tags = ContentTagIds

						var seasonGenres []MultiTierContentGenres

						var seasonGenresIds []struct {
							ContentId string `json:"content_id"`
							Order     string `json:"order"`
							Id        string `json:"id"`
							GenreId   string `json:"genre_id"`
						}

						db.Table("season_genre").Where("season_id =?", seasonId.Id).Find(&seasonGenresIds)

						for _, querygenresId := range seasonGenresIds {

							var genresName struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							var finalcontentSubGenresEnglish, finalcontentSubGenresArabic []string

							var subgenreId []string

							var seasonSubGenres []struct {
								EnglishName string `json:"english_name"`
								ArabicName  string `json:"arabic_name"`
							}

							db.Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName)

							db.Table("season_subgenre").Where("season_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId)

							db.Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&seasonSubGenres)

							for _, contentSubGenre := range seasonSubGenres {
								finalcontentSubGenresEnglish = append(finalcontentSubGenresEnglish, contentSubGenre.EnglishName)
								finalcontentSubGenresArabic = append(finalcontentSubGenresArabic, contentSubGenre.ArabicName)
							}

							seasonGenres = append(seasonGenres, MultiTierContentGenres{
								Id:               querygenresId.GenreId,
								GenerEnglishName: genresName.EnglishName,
								GenerArabicName:  genresName.ArabicName,
								SubGenerEnglish:  finalcontentSubGenresEnglish,
								SubGenerArabic:   finalcontentSubGenresArabic,
							})
						}

						var varianceTrailers []VarianceTrailers
						db.Table("variance_trailer").Where("season_id =?", seasonId.Id).Find(&varianceTrailers)

						var finalVarianceResult []interface{}

						for _, varianceData := range varianceTrailers {
							var variance VarianceTrailers
							variance.Order = varianceData.Order
							variance.VideoTrailerId = varianceData.VideoTrailerId
							variance.EnglishTitle = varianceData.EnglishTitle
							variance.ArabicTitle = varianceData.ArabicTitle
							variance.Duration = varianceData.Duration
							variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
							variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + seasonId.ContentId + "/" + seasonId.Id + "/" + varianceData.Id + "/trailer-poster-image"
							variance.Id = varianceData.Id
							finalVarianceResult = append(finalVarianceResult, variance)
						}

						var aboutTheContent AboutTheContent

						db.Table("about_the_content_info").Where("id =?", seasonId.AboutTheContentInfoId).Find(&aboutTheContent)

						var ageRatings AgeRatingsCode

						db.Raw("select split_part(code, '_', 2) as code from age_ratings where id = ?", aboutTheContent.AgeGroup).Find(&ageRatings)

						aboutTheContent.AgeGroup = ageRatings.Code

						var productionCountries []int

						var queryProductionCountries []struct {
							CountryId int `json:"country_id"`
						}

						db.Table("production_country").Where("about_the_content_info_id =?", seasonId.AboutTheContentInfoId).Find(&queryProductionCountries)

						for _, queryProductionCountry := range queryProductionCountries {
							productionCountries = append(productionCountries, queryProductionCountry.CountryId)
						}

						var translation Translation2
						aboutTheContent.ProductionCountries = productionCountries

						db.Table("content_translation").Where("id =?", seasonId.TranslationId).Find(&translation)

						LanguageTypeId, _ := strconv.Atoi(translation.LanguageType)
						translation.LanguageType = common.LanguageOriginTypes(LanguageTypeId)
						translation.DubbingDialectName = common.DialectIdname(translation.DubbingDialectId, "en")

						var episodeResultFinal []ContentEpisode

						var episodeDetails []FetchEpisodeDetailsMultiTier

						db.Table("season").Where("id =?", seasonId.Id).Find(&translation)

						db.Raw(`
										select
											e.has_poster_image,
											e.has_dubbing_script,
											e.has_subtitling_script,
											e.number as episode_number,
											e.third_party_episode_key as episode_key,
											pi2.duration as length,
											pi2.video_content_id,
											e.synopsis_english,
											e.synopsis_arabic,
											e.has_poster_image,
											e.has_dubbing_script,
											e.has_subtitling_script,
											e.id as episode_id
										from
											episode e
										join playback_item pi2 on
											pi2.id = e.playback_item_id
										where
											e.season_id = ? 
											and e.deleted_by_user_id is null
											and e.modified_at >= ? AND e.modified_at <= ?
		
										order by
											e.number asc
							`, seasonId.Id, startdate, enddate).Find(&episodeDetails)

						if len(episodeDetails) == 0 {
							if len(varianceTrailers) == 0 {
								episodeDetails = nil
							}
						}

						for _, value := range episodeDetails {
							var episodeResult ContentEpisode

							episodeResult.EpisodeNumber = value.EpisodeNumber
							episodeResult.EpisodeKey = value.EpisodeKey
							episodeResult.Length = value.Length
							episodeResult.VideoContentUrl = os.Getenv("VIDEO_API") + value.VideoContentId
							episodeResult.SynopsisEnglish = value.SynopsisEnglish
							episodeResult.SynopsisArabic = value.SynopsisArabic

							if value.HasPosterImage {
								episodeResult.NonTextualData.PosterImage = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/" + value.EpisodeId + os.Getenv("POSTER_IMAGE")
							} else {
								episodeResult.NonTextualData.PosterImage = ""
							}
							if value.HasDubbingScript {
								*episodeResult.NonTextualData.DubbingScript = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/" + value.EpisodeId + os.Getenv("DUBBLING_SCRIPT")
							} else {
								episodeResult.NonTextualData.DubbingScript = nil
							}
							if value.HasSubtitlingScript {
								*episodeResult.NonTextualData.SubtitlingScript = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/" + value.EpisodeId + "/subtitling-script"
							} else {
								episodeResult.NonTextualData.SubtitlingScript = nil
							}
							episodeResult.EpisodeId = value.EpisodeId

							episodeResultFinal = append(episodeResultFinal, episodeResult)
						}

						var contentNonTextualData ContentNonTextualData

						if seasonId.HasPosterImage {
							contentNonTextualData.PosterImage = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/poster-image"
						}
						if seasonId.HasOverlayPosterImage {
							contentNonTextualData.OverlayPosterImage = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/overlay-poster-image"
						}
						if seasonId.HasDetailsBackground {
							contentNonTextualData.DetailsBackground = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/details-background"
						}
						if seasonId.HasMobileDetailsBackground {
							contentNonTextualData.MobileDetailsBackground = IMAGES + seasonId.ContentId + "/" + seasonId.Id + "/mobile-details-background"
						}

						var digitalRights []int
						var digitalRightsRegions []DigitalRightsRegions

						db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", seasonId.RightsId).Scan(&digitalRightsRegions)

						if CountryResult == 0 {

							if country != "" {
								if country == "all" || country == "All" {

									if common.CountryCount() == len(digitalRightsRegions) {
										digitalRights = nil
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
							for _, region := range digitalRightsRegions {
								if CountryResult == int32(region.CountryId) {
								}
							}
						}

						contentSeasons = append(contentSeasons, ContentSeasons2{
							SeasonId:              seasonId.Id,
							ContentId:             seasonId.ContentId,
							SeasonKey:             seasonId.SeasonKey,
							SeasonNumber:          seasonId.Number,
							CreatedAt:             seasonId.CreatedAt,
							ModifiedAt:            seasonId.ModifiedAt,
							PrimaryInfo:           seasonPrimaryInfo,
							Cast:                  cast,
							Music:                 music,
							TagInfo:               tagInfo,
							SeasonGenres:          seasonGenres,
							TrailerInfo:           finalVarianceResult,
							AboutTheContent:       aboutTheContent,
							Translation:           translation,
							EpisodeResult:         episodeResultFinal,
							ContentNonTextualData: contentNonTextualData,
							DigitalRightsRegions:  digitalRights,
						})

						contentResult.ContentSeasons = contentSeasons

						contentResult.SeoDetails = SeoDetails{
							EnglishMetaTitle:       finalContentResult.EnglishMetaTitle,
							ArabicMetaTitle:        finalContentResult.ArabicMetaTitle,
							EnglishMetaDescription: finalContentResult.EnglishMetaDescription,
							ArabicMetaDescription:  finalContentResult.ArabicMetaDescription,
						}

						contentResult.Id = finalContentResult.ContentId
						contentResult.ModifiedAt = finalContentResult.ModifiedAt
						contentResult.CreatedAt = finalContentResult.CreatedAt

						contentResultFinal = append(contentResultFinal, contentResult)
					}
				}
				finalData[types] = contentResultFinal
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": finalData})
	}
}

func daysBetween(startdate, enddate time.Time) int {
	days := enddate.Sub(startdate).Hours() / 24
	return int(days)
}
