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

// GetAllMultiTierDetails- Get All Multi Tier Content Details
// GET /v1/contents/multitier/
// @Description Get All Multi Tier Content Details
// @Tags MultiTier
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param offset query string false "Zero-based index of the first requested item."
// @Param limit query string false "Maximum number of collection items to return for a single request."
// @Param Country query string false "Country code of the user."
// @Success 200  object AllMultiTier
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/contents/multitier [get]
func (hs *HandlerService) GetAllMultiTierDetails(c *gin.Context) {
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

	// var finalContentResult []FinalSeasonResult
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
	var totalCount []FinalSeasonResultContentOneTire
	/*digital rights*/
	// var country string
	// if c.Request.URL.Query()["Country"] != nil {
	// 	country = c.Request.URL.Query()["Country"][0]
	// 	fmt.Println(country)
	// }

	var contentResultFinal []MultiTierContent
	serverError := common.ServerErrorResponse()
	// notFound := common.NotFoundErrorResponse()
	var finalContent []FinalSeasonResultContentOneTire
	// var finalContentResult []FinalSeasonResult
	// var finalContentCount []FinalSeasonResultContentOneTire
	// var finalContentResult FinalSeasonResultContentOneTire
	/*for country rights*/
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	}
	CountryResult := common.Countrys(country)

	if UserId == os.Getenv("WATCH_NOW") {

		if CountryResult != 0 {
			db.Debug().Raw(`
						select
							distinct c.third_party_content_key AS multi_tier_content_key,
							c.id,
							-- atci.supplier,
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
						-- join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.watch_now_supplier = 'true'
							and c.status = 1
							
							and c.content_tier = 2
							
							and atci.supplier != 'Others'
							
							and c.deleted_by_user_id is null
							
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							
							and s.status = 1
						
							and s.deleted_by_user_id is null

							and crc.country_id = ?
							
							and  ((select
										count(e1.*)
									from
										episode e1
									
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
								(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0) 
					
			`, CountryResult).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContent)

			db.Raw(`    select
							distinct c.third_party_content_key
						from
							content c
							
						JOIN content_primary_info cpi on cpi.id = c.primary_info_id
						join season s on s.content_id = c.id
						-- join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.watch_now_supplier = 'true'
							and c.status = 1
							and c.content_tier = 2
							and atci.supplier != 'Others'
							and c.deleted_by_user_id is null
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							and s.status = 1
							and s.deleted_by_user_id is null
							and crc.country_id = ?
							and  ((select
										count(e1.*)
									from
										episode e1
									
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
							(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0)`, CountryResult).Find(&totalCount)

			fmt.Println("--------->", len(totalCount))

		} else if country == "" || country == "all" {

			db.Debug().Raw(`
						select
							distinct c.third_party_content_key AS multi_tier_content_key,
							c.id,
							-- atci.supplier,
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
						-- join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.watch_now_supplier = 'true'
							and c.status = 1
							
							and c.content_tier = 2
							
							and atci.supplier != 'Others'
							
							and c.deleted_by_user_id is null
							
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							
							and s.status = 1
						
							and s.deleted_by_user_id is null

							and  ((select
										count(e1.*)
									from
										episode e1
									
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
								(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0) 
					
			`).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContent)

			db.Raw(`
			select
			distinct c.third_party_content_key
		from
			content c
			
		JOIN content_primary_info cpi on cpi.id = c.primary_info_id
		join season s on s.content_id = c.id
		-- join content_translation ct on ct.id = s.translation_id
		-- join playback_item pi2 on pi2.translation_id = ct.id
		join content_rights cr on cr.id = s.rights_id
		join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
		join content_rights_country crc on crc.content_rights_id = cr.id
			
		where
			c.watch_now_supplier = 'true'
			and c.status = 1
			and c.content_tier = 2
			and atci.supplier != 'Others'
			and c.deleted_by_user_id is null
			-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
			and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
			and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
			and s.status = 1
			and s.deleted_by_user_id is null
			and  ((select
						count(e1.*)
					from
						episode e1
					
					
					join content_rights cr1 on cr1.id = s.rights_id
					join playback_item pi4 on pi4.id = e1.playback_item_id
					
					where
						e1.season_id = s.id
						and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
						and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
						and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
						and e1.status = 1
						and e1.deleted_by_user_id is null
				) > 0 or 
			(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0)
			`).Find(&totalCount)
		}

	} else {
		if CountryResult != 0 {
			db.Debug().Raw(`
						select
							distinct c.third_party_content_key AS multi_tier_content_key,
							c.id,
							-- atci.supplier,
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
						--join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.status = 1
							
							and c.content_tier = 2
							
							and atci.supplier != 'Others'
							
							and c.deleted_by_user_id is null
							
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							
							and s.status = 1
						
							and s.deleted_by_user_id is null

							and crc.country_id = ?
							
							and  ((select
										count(e1.*)
									from
										episode e1
									
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
								(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0) 
					
			`, CountryResult).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContent)

			db.Raw(`
			select
			distinct c.third_party_content_key
		from
			content c
			
		JOIN content_primary_info cpi on cpi.id = c.primary_info_id
		join season s on s.content_id = c.id
		-- join content_translation ct on ct.id = s.translation_id
		-- join playback_item pi2 on pi2.translation_id = ct.id
		join content_rights cr on cr.id = s.rights_id
		join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
		join content_rights_country crc on crc.content_rights_id = cr.id
			
		where
			c.status = 1
			and c.content_tier = 2
			and atci.supplier != 'Others'
			and c.deleted_by_user_id is null
			-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
			and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
			and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
			and s.status = 1
			and s.deleted_by_user_id is null
			and crc.country_id = ?
			and  ((select
						count(e1.*)
					from
						episode e1
					
					
					join content_rights cr1 on cr1.id = s.rights_id
					join playback_item pi4 on pi4.id = e1.playback_item_id
					
					where
						e1.season_id = s.id
						and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
						and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
						and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
						and e1.status = 1
						and e1.deleted_by_user_id is null
				) > 0 or 
			(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0)
			`, CountryResult).Find(&totalCount)
		} else if country == "" || country == "all" {

			db.Debug().Raw(`
						select
							distinct c.third_party_content_key AS multi_tier_content_key,
							c.id,
							-- atci.supplier,
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
						--join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.status = 1
							
							and c.content_tier = 2
							
							and atci.supplier != 'Others'
							
							and c.deleted_by_user_id is null
							
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							
							and s.status = 1
						
							and s.deleted_by_user_id is null

							and  ((select
										count(e1.*)
									from
										episode e1
									
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
								(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0) 
					
			`).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContent)

			db.Raw(`  	select
							distinct c.third_party_content_key
						from
							content c
							
						JOIN content_primary_info cpi on cpi.id = c.primary_info_id
						join season s on s.content_id = c.id
						-- join content_translation ct on ct.id = s.translation_id
						-- join playback_item pi2 on pi2.translation_id = ct.id
						join content_rights cr on cr.id = s.rights_id
						join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
						join content_rights_country crc on crc.content_rights_id = cr.id
							
						where
							c.status = 1
							and c.content_tier = 2
							and atci.supplier != 'Others'
							and c.deleted_by_user_id is null
							-- and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
							and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null)
							and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null)
							and s.status = 1
							and s.deleted_by_user_id is null
							and  ((select
										count(e1.*)
									from
										episode e1
									
									join content_rights cr1 on cr1.id = s.rights_id
									join playback_item pi4 on pi4.id = e1.playback_item_id
									
									where
										e1.season_id = s.id
										and (cr1.digital_rights_start_date <= now() or cr1.digital_rights_start_date is null)
										and (cr1.digital_rights_end_date >= now() or cr1.digital_rights_end_date is null)
										and ( pi4.scheduling_date_time <= NOW() or pi4.scheduling_date_time is null)
										and e1.status = 1
										and e1.deleted_by_user_id is null
								) > 0 or 
							(select count(vt.*) from variance_trailer vt where vt.season_id = s.id) > 0)`).Find(&totalCount)
		}
	}

	// fmt.Println(totalCount)
	// totalCount = len(finalContentCount)

	for _, finalContentResult := range finalContent {

		var contentResult MultiTierContent

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
		// .Error; err != nil {
		// 	serverError.Description = "Genre Id Wrong"
		// 	c.JSON(http.StatusInternalServerError, serverError)
		// 	return
		// }

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

			// if err :=
			db.Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName)
			// .Error; err != nil {
			// 	serverError.Description = "Genre Id Wrong 2"
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			// if err :=
			db.Table("content_subgenre").Where("content_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId)
			// .Error; err != nil {
			// 	serverError.Description = "Genre Id Wrong 3"
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			// if err :=
			db.Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&contentSubGenres)
			// .Error; err != nil {
			// 	serverError.Description = "Genre Id Wrong 4"
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

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

		var contentSeasons []ContentSeasons
		var seasons []Seasons

		// if err :=
		db.Raw(`
			select
				s.*
			from
				season s
				
			join content_rights cr2 on cr2.id = s.rights_id
			join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
			
			where
				s.deleted_by_user_id is null and
				(cr2.digital_rights_start_date <= NOW() or cr2.digital_rights_start_date is null) and 
				(cr2.digital_rights_end_date >= NOW() or cr2.digital_rights_end_date is null) and 
				atci.supplier !='Others' and
				s.content_id = ?
			`, finalContentResult.ContentId).Find(&seasons)

		// .Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, serverError)
		// 	return
		// }

		for _, seasonId := range seasons {

			var seasonPrimaryInfo SeasonPrimaryInfo

			// if err :=
			db.Table("content_primary_info").Where("id =?", seasonId.PrimaryInfoId).Find(&seasonPrimaryInfo)
			// .Error; err != nil {
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

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

			// if err :=
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
						s.status = 1 and 
						s.deleted_by_user_id is null and
						s.id = ?
				`, seasonId.Id).Find(&castList)
			// 	.Error; err != nil {
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			var actorIds, actorEnglishNames, actorArabicNames []string
			var queryActorIds []struct {
				ActorId string `json:"actor_id"`
			}

			// if err :=
			db.Table("content_actor").Where("cast_id =?", castList.CastId).Find(&queryActorIds)
			// .Error; err != nil {
			// 	fmt.Println(888)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, queryActorId := range queryActorIds {
				actorIds = append(actorIds, queryActorId.ActorId)

				var actorName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("actor").Where("id =?", queryActorId.ActorId).Find(&actorName)
				// .Error; err != nil {
				// 	fmt.Println(999)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				actorEnglishNames = append(actorEnglishNames, actorName.EnglishName)
				actorArabicNames = append(actorArabicNames, actorName.ArabicName)
			}

			var writerIds, writerEnglishNames, writerArabicNames []string
			var queryWriterIds []struct {
				WriterId string `json:"writer_id"`
			}

			// if err :=
			db.Table("content_writer").Where("cast_id =?", castList.CastId).Find(&queryWriterIds)
			// .Error; err != nil {
			// 	fmt.Println(10 - 10 - 10)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, queryWriterId := range queryWriterIds {
				writerIds = append(writerIds, queryWriterId.WriterId)

				var writerName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("writer").Where("id =?", queryWriterId.WriterId).Find(&writerName)
				// .Error; err != nil {
				// 	fmt.Println(11 - 11 - 11)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				writerEnglishNames = append(writerEnglishNames, writerName.EnglishName)
				writerArabicNames = append(writerArabicNames, writerName.ArabicName)
			}

			var directorIds, directorEnglishNames, directorArabicNames []string
			var querydirectorIds []struct {
				DirectorId string `json:"director_id"`
			}

			// if err :=
			db.Table("content_director").Where("cast_id =?", castList.CastId).Find(&querydirectorIds)
			// .Error; err != nil {
			// 	fmt.Println(12121212)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, querydirectorId := range querydirectorIds {
				directorIds = append(directorIds, querydirectorId.DirectorId)

				var directorName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("director").Where("id =?", querydirectorId.DirectorId).Find(&directorName)
				// .Error; err != nil {
				// 	fmt.Println(131313)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

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

			// if err :=
			db.Table("content_singer").Where("music_id =?", seasonId.MusicId).Find(&querySingerIds)
			// .Error; err != nil {
			// 	fmt.Println(141414)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, querySingerId := range querySingerIds {
				singerIds = append(singerIds, querySingerId.SingerId)

				var SingerName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("singer").Where("id =?", querySingerId.SingerId).Find(&SingerName)
				// .Error; err != nil {
				// 	fmt.Println(151515)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				singerEnglishNames = append(singerEnglishNames, SingerName.EnglishName)
				singerArabicNames = append(singerArabicNames, SingerName.ArabicName)
			}

			var musicComposerIds, musicComposerEnglishNames, musicComposerArabicNames []string
			var queryMusicComposerIds []struct {
				MusicComposerId string `json:"music_composer_id"`
			}

			// if err :=
			db.Table("content_music_composer").Where("music_id =?", seasonId.MusicId).Find(&queryMusicComposerIds)
			// .Error; err != nil {
			// 	fmt.Println(161616)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, queryMusicComposerId := range queryMusicComposerIds {
				musicComposerIds = append(musicComposerIds, queryMusicComposerId.MusicComposerId)

				var MusicComposerName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("music_composer").Where("id =?", queryMusicComposerId.MusicComposerId).Find(&MusicComposerName)
				// .Error; err != nil {
				// 	fmt.Println(171717)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				musicComposerEnglishNames = append(musicComposerEnglishNames, MusicComposerName.EnglishName)
				musicComposerArabicNames = append(musicComposerArabicNames, MusicComposerName.ArabicName)
			}

			var songWriterIds, songWriterEnglishNames, songWriterArabicNames []string
			var querySongWriterIds []struct {
				SongWriterId string `json:"song_writer_id"`
			}

			// if err :=
			db.Table("content_song_writer").Where("music_id =?", seasonId.MusicId).Find(&querySongWriterIds)
			// .Error; err != nil {
			// 	fmt.Println(181818)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, querySongWriterId := range querySongWriterIds {
				songWriterIds = append(songWriterIds, querySongWriterId.SongWriterId)

				var SongWriterName struct {
					EnglishName string `json:"english_name"`
					ArabicName  string `json:"arabic_name"`
				}

				// if err :=
				db.Table("song_writer").Where("id =?", querySongWriterId.SongWriterId).Find(&SongWriterName)
				// .Error; err != nil {
				// 	fmt.Println(191919)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

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

			// if err :=
			db.Table("content_tag").Where("tag_info_id =?", seasonId.TagInfoId).Find(&queryContentTagIds)
			// .Error; err != nil {
			// 	fmt.Println(171717)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			var ContentTagIds []string
			for _, queryContentTagId := range queryContentTagIds {

				var TextualDataTagName struct {
					Name string `json:"name"`
				}

				// if err :=
				db.Table("textual_data_tag").Where("id =?", queryContentTagId.TextualDataTagId).Find(&TextualDataTagName)
				// .Error; err != nil {
				// 	fmt.Println(181818)
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

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

			// if err :=
			db.Table("season_genre").Where("season_id =?", seasonId.Id).Find(&seasonGenresIds)
			// .Error; err != nil {
			// 	serverError.Description = "Season Genre Id Wrong 1"
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

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

				// if err :=
				db.Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName)
				// .Error; err != nil {
				// 	serverError.Description = "Season Genre Id Wrong 2"
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				// if err :=
				db.Table("season_subgenre").Where("season_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId)
				// .Error; err != nil {
				// 	serverError.Description = "Season Genre Id Wrong 3"
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

				// if err :=
				db.Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&seasonSubGenres)
				// .Error; err != nil {
				// 	serverError.Description = "Season Genre Id Wrong 4"
				// 	c.JSON(http.StatusInternalServerError, serverError)
				// 	return
				// }

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
			// if err :=
			db.Table("variance_trailer").Where("season_id =?", seasonId.Id).Find(&varianceTrailers)
			// .Error; err != nil {
			// 	fmt.Println(191919)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			var finalVarianceResult []interface{}

			for _, varianceData := range varianceTrailers {
				var variance VarianceTrailers
				variance.Order = varianceData.Order
				variance.VideoTrailerId = varianceData.VideoTrailerId
				variance.EnglishTitle = varianceData.EnglishTitle
				variance.ArabicTitle = varianceData.ArabicTitle
				variance.Duration = varianceData.Duration
				variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
				//const url = "https://z5content-qa.s3.amazonaws.com/"
				variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + seasonId.ContentId + "/" + seasonId.Id + "/" + varianceData.Id + "/trailer-poster-image"
				variance.Id = varianceData.Id
				finalVarianceResult = append(finalVarianceResult, variance)
			}

			var aboutTheContent AboutTheContent

			// if err :=
			db.Table("about_the_content_info").Where("id =?", seasonId.AboutTheContentInfoId).Find(&aboutTheContent)
			// .Error; err != nil {
			// 	fmt.Println(120202)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			var ageRatings AgeRatingsCode

			// if err :=
			db.Raw("select split_part(code, '_', 2) as code from age_ratings where id = ?", aboutTheContent.AgeGroup).Find(&ageRatings)
			// .Error; err != nil {
			// 	fmt.Println(212121)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			aboutTheContent.AgeGroup = ageRatings.Code

			var productionCountries []int

			var queryProductionCountries []struct {
				CountryId int `json:"country_id"`
			}

			// if err :=
			db.Table("production_country").Where("about_the_content_info_id =?", seasonId.AboutTheContentInfoId).Find(&queryProductionCountries)
			// .Error; err != nil {
			// 	fmt.Println(22 - 22 - 22)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			for _, queryProductionCountry := range queryProductionCountries {
				productionCountries = append(productionCountries, queryProductionCountry.CountryId)
			}

			aboutTheContent.ProductionCountries = productionCountries

			var translation Translation

			// if err :=
			db.Table("content_translation").Where("id =?", seasonId.TranslationId).Find(&translation)
			// .Error; err != nil {
			// 	fmt.Println(232323)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			LanguageTypeId, _ := strconv.Atoi(translation.LanguageType)
			translation.LanguageType = common.LanguageOriginTypes(LanguageTypeId)
			translation.DubbingDialectName = common.DialectIdname(translation.DubbingDialectId, "en")

			// var episodes []Episode
			var episodeResultFinal []ContentEpisode

			var episodeDetails []FetchEpisodeDetailsMultiTier

			// if err :=
			db.Table("season").Where("id =?", seasonId.Id).Find(&translation)
			// .Error; err != nil {
			// 	fmt.Println(242424)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			// if err :=
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
								join season s on
									s.id = e.season_id
								join playback_item pi2 on
									pi2.id = e.playback_item_id
								join content_rights cr on 
									cr.id = s.rights_id
								where
									e.season_id = ? 
									and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null)
									and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null)
									and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
									and e.status = 1 
									and e.deleted_by_user_id is null

								order by
									e.number asc
			`, seasonId.Id).Find(&episodeDetails)
			// .Error; err != nil {
			// 	fmt.Println(252525)
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

			if len(episodeDetails) == 0 {
				if len(varianceTrailers) == 0 {
					// serverError.Description = "No Episode Found"
					// c.JSON(http.StatusInternalServerError, serverError)
					// return
					episodeDetails = nil
				}
			}

			// if len(episodeDetails) == 0 {
			// 	c.JSON(http.StatusInternalServerError, serverError)
			// 	return
			// }

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

			var flag bool = false

			var digitalRightsRegions []DigitalRightsRegions

			db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", seasonId.RightsId).Scan(&digitalRightsRegions)

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

			if !flag {
				contentSeasons = append(contentSeasons, ContentSeasons{
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
			}

		}

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

		if len(contentSeasons) > 0 {
			contentResultFinal = append(contentResultFinal, contentResult)
		}

	}

	var pagination Pagination
	pagination.Limit = int(limit)
	pagination.Offset = int(offset)
	pagination.Size = len(totalCount)
	// if CountryResult != 0 || country == "" {
	// 	c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})
	// }

	if len(contentResultFinal) > 0 {
		c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": contentResultFinal})
	} else {
		serverError.Description = "No Records Found"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

}

// func (hs *HandlerService) GetAllMultiTierDetailss(c *gin.Context) {
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
// 	var contentResult AllMultiTierContent
// 	var allContents []AllMultiTierContent
// 	var CountryResult int32
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
// 	/*digital rights*/
// 	var country string
// 	if c.Request.URL.Query()["Country"] != nil {
// 		country = c.Request.URL.Query()["Country"][0]
// 		fmt.Println(country)
// 	}
// 	CountryResult = common.Countrys(country)
// 	if CountryResult != 0 {
// 		if err := db.Debug().Table("content c").Select(`
// 			distinct c.content_key as multi_tier_content_key,
// 			s.id,pi2.rights_id, c.content_type ,cpi.original_title,cpi.alternative_title,cpi.arabic_title as arabic_title,
// 			cpi.transliterated_title as transliterated_title,cpi.notes,c.id as content_id,s.season_key,s.number as season_number,
// 			s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,
// 			cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,s.cast_id,
// 			s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,
// 			atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,
// 			ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,
// 			s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background`).
// 			Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
// 			Joins("join content_genre cg on cg.content_id  = c.id").
// 			Joins("join season s on s.content_id = c.id").
// 			Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
// 			Joins("join content_cast cc  on cc.Id  = s.cast_id ").
// 			Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id").
// 			Joins("join content_translation ct on ct.id = s. translation_id").
// 			Joins("join playback_item pi2 on pi2.translation_id = ct.id").
// 			Joins("join content_rights cr on cr.id = pi2.rights_id").
// 			Joins("join content_rights_country crc on crc.content_rights_id = cr.id").
// 			Where("c.status = 1 and c.content_tier =2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null  and crc.country_id = ? ", CountryResult).Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		db.Raw("select count(distinct c.id) from content c join season s on s.content_id = c.id join content_translation ct on ct.id = s. translation_id join playback_item pi2 on	pi2.translation_id = ct.id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id = cr.id where (c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null  and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null) and crc.country_id = ?", CountryResult).Count(&totalCount)
// 	} else if country == "" || country == "all" {
// 		if err := db.Debug().Table("content c").Select(`
// 			distinct c.content_key as multi_tier_content_key,s.id,pi2.rights_id, c.content_type ,
// 			cpi.original_title,cpi.alternative_title,cpi.arabic_title as arabic_title,
// 			cpi.transliterated_title as transliterated_title,cpi.notes,c.id as content_id,s.season_key,s.number as season_number,
// 			s.created_at as inserted_at,s.modified_at,cpi2.original_title as season_original_title,cpi2.alternative_title as season_alternative_title,
// 			cpi2.arabic_title as season_arabic_title,cpi2.transliterated_title as season_transliterated_title,cpi2.notes as season_notes,
// 			s.cast_id,s.music_id,s.tag_info_id,s.id as season_id,atci.original_language,atci.supplier,atci.acquisition_department,atci.english_synopsis,
// 			atci.arabic_synopsis,atci.production_year,atci.production_house,atci.age_group,s.about_the_content_info_id,ct.language_type as multi_tier_language_type,
// 			ct.dubbing_language,ct.dubbing_dialect_id,ct.subtitling_language,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,s.content_id,s.has_poster_image,s.has_overlay_poster_image,s.has_details_background,s.has_mobile_details_background`).
// 			Joins("join content_primary_info cpi on cpi.id = c.primary_info_id").
// 			Joins("join content_genre cg on cg.content_id  = c.id").
// 			Joins("join season s on s.content_id = c.id").
// 			Joins("join content_primary_info cpi2 on cpi2.id = s.primary_info_id").
// 			Joins("join content_cast cc  on cc.Id  = s.cast_id ").
// 			Joins("join about_the_content_info atci on  atci.Id = s.about_the_content_info_id").
// 			Joins("join content_translation ct on ct.id = s. translation_id").
// 			Joins("join playback_item pi2 on pi2.translation_id = ct.id").
// 			Joins("join content_rights cr on cr.id = pi2.rights_id").
// 			Where("c.status = 1 and c.content_tier =2 and c.deleted_by_user_id is null and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null ").Order("c.modified_at desc").Limit(limit).Offset(offset).Find(&finalContentResult).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		db.Raw("select count(distinct c.id) from content c join season s on s.content_id = c.id join content_translation ct on ct.id = s. translation_id join playback_item pi2 on	pi2.translation_id = ct.id join content_rights cr on cr.id = pi2.rights_id where (c.status = 1 and c.content_tier = 2 and c.deleted_by_user_id is null  and (pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and (cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and s.status = 1 and s.deleted_by_user_id is null) ").Count(&totalCount)
// 	} else if country != "all" && CountryResult == 0 {
// 		c.JSON(http.StatusInternalServerError, serverError)
// 		return
// 	}
// 	var contentSeasons []AllContentSeasons
// 	var contentSeason AllContentSeasons

// 	for _, eachcontent := range finalContentResult {

// 		contentResult.ContentKey = eachcontent.MultiTierContentKey
// 		/*content seasons*/
// 		contentSeason.ContentId = eachcontent.ContentId
// 		contentSeason.SeasonKey = eachcontent.SeasonKey
// 		contentSeason.SeasonNumber = eachcontent.SeasonNumber
// 		contentSeason.CreatedAt = eachcontent.CreatedAt
// 		contentSeason.ModifiedAt = eachcontent.ModifiedAt
// 		/*primary info season*/
// 		contentSeason.PrimaryInfo.SeasonNumber = eachcontent.SeasonNumber
// 		contentSeason.PrimaryInfo.OriginalTitle = eachcontent.SeasonOriginalTitle
// 		contentSeason.PrimaryInfo.AlternativeTitle = eachcontent.SeasonAlternativeTitle
// 		contentSeason.PrimaryInfo.ArabicTitle = eachcontent.ArabicTitle
// 		contentSeason.PrimaryInfo.TransliteratedTitle = eachcontent.TransliteratedTitle
// 		contentSeason.PrimaryInfo.Notes = eachcontent.Notes
// 		/* Fetch content_cast normal*/
// 		var contentCast Cast
// 		if castResult := db.Table("content_cast cc").Select("cc.main_actor_id,cc.main_actress_id,actor.english_name as main_actor_english,actor.arabic_name as main_actor_arabic,actress.english_name as main_actress_english,actress.arabic_name as main_actress_arabic").
// 			Joins("left join actor actor on actor.id =cc.main_actor_id").
// 			Joins("left join actor actress on actress.id =cc.main_actress_id").
// 			Where("cc.id=?", eachcontent.CastId).Scan(&contentCast).Error; castResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}

// 		var varianceTrailers []VarianceTrailers
// 		if varianceTrailersError := db.Debug().Raw("select * from variance_trailer vt where vt.season_id=? order by vt.order asc ", eachcontent.SeasonId).Find(&varianceTrailers).Error; varianceTrailersError != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var variance VarianceTrailers
// 		var finalVarianceResult []interface{}

// 		for _, varianceData := range varianceTrailers {
// 			variance.Order = varianceData.Order
// 			variance.VideoTrailerId = varianceData.VideoTrailerId
// 			variance.EnglishTitle = varianceData.EnglishTitle
// 			variance.ArabicTitle = varianceData.ArabicTitle
// 			variance.Duration = varianceData.Duration
// 			variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
// 			//const url = "https://z5content-qa.s3.amazonaws.com/"
// 			variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + eachcontent.ContentId + "/" + eachcontent.SeasonId + "/" + varianceData.Id + "/trailer-poster-image"
// 			variance.Id = varianceData.Id
// 			finalVarianceResult = append(finalVarianceResult, variance)
// 		}

// 		contentSeason.TrailerInfo = finalVarianceResult

// 		contentSeason.Cast.CastId = eachcontent.CastId
// 		contentSeason.Cast.MainActorId = contentCast.MainActorId
// 		contentSeason.Cast.MainActressId = contentCast.MainActressId
// 		contentSeason.Cast.MainActorEnglish = contentCast.MainActorEnglish
// 		contentSeason.Cast.MainActorArabic = contentCast.MainActorArabic
// 		contentSeason.Cast.MainActressEnglish = contentCast.MainActressEnglish
// 		contentSeason.Cast.MainActressArabic = contentCast.MainActressArabic
// 		/*fetching other cast details */
// 		var contentActor []ContentActor
// 		if actorResult := db.Table("content_actor ca").Select("a.english_name as actor_english,a.arabic_name as actor_arabic,a.id as actor_id,w.id as writer_id,w.english_name as writer_english,w.arabic_name as writer_arabic,d.id as director_id,d.english_name as director_english,d.arabic_name as director_arabic").
// 			Joins("left join actor a on a.id =ca.actor_id").
// 			Joins("left join content_writer cw on cw.cast_id =ca.cast_id").
// 			Joins("left  join writer w on w.id =cw.writer_id").
// 			Joins("left join content_director cd on cd.cast_id =ca.cast_id").
// 			Joins("left join director d on d.id =cd.director_id").
// 			Where("ca.cast_id=?", eachcontent.CastId).Scan(&contentActor).Error; actorResult != nil {
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
// 		contentSeason.Cast.ActorIds = common.RemoveDuplicateValues(actorId)
// 		contentSeason.Cast.ActorEnglish = common.RemoveDuplicateValues(actorEnglish)
// 		contentSeason.Cast.ActorArabic = common.RemoveDuplicateValues(actorArabic)
// 		contentSeason.Cast.WriterId = common.RemoveDuplicateValues(writerId)
// 		contentSeason.Cast.WriterEnglish = common.RemoveDuplicateValues(writerEnglish)
// 		contentSeason.Cast.WriterArabic = common.RemoveDuplicateValues(writerArabic)
// 		contentSeason.Cast.DirectorIds = common.RemoveDuplicateValues(directorId)
// 		contentSeason.Cast.DirectorEnglish = common.RemoveDuplicateValues(directorEnglish)
// 		contentSeason.Cast.DirectorArabic = common.RemoveDuplicateValues(directorArabic)
// 		/* fetching music details */
// 		var contentMusic []ContentMusic
// 		if actorResult := db.Table("content_singer cs").Select("s.id as singer_ids,s.english_name as singers_english,s.arabic_name as singers_arabic,mc.id as music_composer_ids,mc.english_name as music_composers_english ,mc.arabic_name as music_omposers_arabic,sw.id as song_writer_ids,sw.english_name as song_writers_english,sw.arabic_name as song_writers_arabic").
// 			Joins("left join singer s on s.id =cs.singer_id").
// 			Joins("left join content_music_composer cmc on cmc.music_id =cs.music_id").
// 			Joins("left join music_composer mc on mc.id =cmc.music_composer_id").
// 			Joins("left join content_song_writer csw on csw.music_id =cs.music_id ").
// 			Joins("left join song_writer sw on sw.id =csw.song_writer_id").
// 			Where("cs.music_id=?", eachcontent.MusicId).Scan(&contentMusic).Error; actorResult != nil {
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
// 		contentSeason.Music.MusicId = eachcontent.MusicId
// 		contentSeason.Music.SingerIds = common.RemoveDuplicateValues(singerId)
// 		contentSeason.Music.SingersEnglish = common.RemoveDuplicateValues(singerEnglish)
// 		contentSeason.Music.SingersArabic = common.RemoveDuplicateValues(singerArabic)
// 		contentSeason.Music.MusicComposerIds = common.RemoveDuplicateValues(composerId)
// 		contentSeason.Music.MusicComposersEnglish = common.RemoveDuplicateValues(composerEnglish)
// 		contentSeason.Music.MusicComposersArabic = common.RemoveDuplicateValues(composerArabic)
// 		contentSeason.Music.SongWriterIds = common.RemoveDuplicateValues(SongWriterId)
// 		contentSeason.Music.SongWritersEnglish = common.RemoveDuplicateValues(SongWriterEnglish)
// 		contentSeason.Music.SongWritersArabic = common.RemoveDuplicateValues(SongWriterArabic)
// 		/*fetching tag info */
// 		var contentTagInfo []ContentTag
// 		db.Table("content_tag ct").Select("tdt.name").
// 			Joins("left join textual_data_tag tdt on tdt.id =ct.textual_data_tag_id").
// 			Where("ct.tag_info_id=?", eachcontent.TagInfoId).Scan(&contentTagInfo)
// 		var tagInfo []string
// 		for _, tagInfoIds := range contentTagInfo {
// 			tagInfo = append(tagInfo, tagInfoIds.Name)
// 		}
// 		contentSeason.TagInfo.Tags = tagInfo
// 		if len(tagInfo) < 1 {
// 			buffer := make([]string, 0)
// 			contentSeason.TagInfo.Tags = buffer
// 		}
// 		/*about the content*/
// 		contentSeason.AboutTheContent.OriginalLanguage = eachcontent.OriginalLanguage
// 		contentSeason.AboutTheContent.Supplier = eachcontent.Supplier
// 		contentSeason.AboutTheContent.AcquisitionDepartment = eachcontent.AcquisitionDepartment
// 		contentSeason.AboutTheContent.EnglishSynopsis = eachcontent.EnglishSynopsis
// 		contentSeason.AboutTheContent.ArabicSynopsis = eachcontent.ArabicSynopsis
// 		contentSeason.AboutTheContent.ProductionYear = eachcontent.ProductionYear
// 		contentSeason.AboutTheContent.ProductionHouse = eachcontent.ProductionHouse
// 		contentSeason.AboutTheContent.AgeGroup = common.AgeRatings(eachcontent.AgeGroup, "en")
// 		/*production countries*/
// 		var productionCountry []ProductionCountry
// 		if productionCountryResult := db.Table("production_country ").Select("country_id").Where("about_the_content_info_id=?", eachcontent.AboutTheContentInfoId).Scan(&productionCountry).Error; productionCountryResult != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var countries []int

// 		for _, prcountries := range productionCountry {
// 			countries = append(countries, prcountries.CountryId)
// 		}
// 		contentSeason.AboutTheContent.ProductionCountries = countries
// 		if len(tagInfo) < 1 {
// 			buffer := make([]int, 0)
// 			contentSeason.AboutTheContent.ProductionCountries = buffer
// 		}
// 		/*translation details*/
// 		contentSeason.Translation.LanguageType = common.LanguageOriginTypes(eachcontent.MultiTierLanguageType)
// 		contentSeason.Translation.DubbingLanguage = eachcontent.DubbingLanguage
// 		contentSeason.Translation.DubbingDialectId = eachcontent.DubbingDialectId
// 		contentSeason.Translation.DubbingDialectName = common.DialectIdname(eachcontent.DubbingDialectId, "en")
// 		contentSeason.Translation.SubtitlingLanguage = eachcontent.SubtitlingLanguage
// 		/*non textual data for seasons*/
// 		if eachcontent.HasPosterImage {
// 			contentSeason.ContentNonTextualData.PosterImage = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("POSTER_IMAGE")
// 		}
// 		if eachcontent.HasOverlayPosterImage {
// 			contentSeason.ContentNonTextualData.OverlayPosterImage = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("OVERLAY_POSTER_IMAGE")
// 		}
// 		if eachcontent.HasDetailsBackground {
// 			contentSeason.ContentNonTextualData.DetailsBackground = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("DETAILS_BACKGROUND")
// 		}
// 		if eachcontent.HasMobileDetailsBackground {
// 			contentSeason.ContentNonTextualData.MobileDetailsBackground = IMAGES + eachcontent.ContentId + "/" + eachcontent.SeasonId + os.Getenv("MOBILE_DETAILS_BACKGROUND")
// 		}
// 		/*digital rights region season*/
// 		var digitalRightsRegions []DigitalRightsRegions
// 		if countryError := db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", eachcontent.RightsId).Scan(&digitalRightsRegions).Error; countryError != nil {
// 			c.JSON(http.StatusInternalServerError, serverError)
// 			return
// 		}
// 		var SeasonRights []int
// 		for _, idarr := range digitalRightsRegions {
// 			SeasonRights = append(SeasonRights, idarr.CountryId)
// 		}
// 		/*for digital rights*/
// 		var IsCheck bool
// 		for _, value := range SeasonRights {
// 			if CountryResult == int32(value) {
// 				IsCheck = true
// 			}
// 		}
// 		if country == "" || country == "all" {
// 			contentSeason.Rights.DigitalRightsRegions = SeasonRights
// 		}
// 		if len(SeasonRights) < 1 {
// 			buffer := make([]int, 0)
// 			contentSeason.Rights.DigitalRightsRegions = buffer
// 		}
// 		contentSeason.SeasonId = eachcontent.SeasonId
// 		/*for checking country rights*/

// 		/*primary info*/
// 		contentResult.PrimaryInfo.ContentType = eachcontent.ContentType
// 		contentResult.PrimaryInfo.OriginalTitle = eachcontent.OriginalTitle
// 		contentResult.PrimaryInfo.AlternativeTitle = eachcontent.AlternativeTitle
// 		contentResult.PrimaryInfo.ArabicTitle = eachcontent.ArabicTitle
// 		contentResult.PrimaryInfo.TransliteratedTitle = eachcontent.TransliteratedTitle
// 		contentResult.PrimaryInfo.Notes = eachcontent.Notes
// 		/*content genres*/
// 		var contentGenres []SeasonGenres
// 		var finalContentGenre []interface{}
// 		var newContentGenres NewSeasonGenres
// 		if genreResult := db.Table("content_genre cg").Select("cg.id,g.english_name as gener_english_name,g.arabic_name as gener_arabic_name").
// 			Joins("left join genre g on g.id=cg.genre_id").
// 			Where("content_id=?", eachcontent.ContentId).Scan(&contentGenres).Error; genreResult != nil {
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
// 				finalContentGenre = append(finalContentGenre, newContentGenres)
// 			}
// 		}
// 		contentResult.ContentGenres = finalContentGenre
// 		/*seo details*/
// 		contentResult.SeoDetails.EnglishMetaTitle = eachcontent.EnglishMetaTitle
// 		contentResult.SeoDetails.ArabicMetaTitle = eachcontent.ArabicMetaTitle
// 		contentResult.SeoDetails.EnglishMetaDescription = eachcontent.EnglishMetaDescription
// 		contentResult.SeoDetails.ArabicMetaDescription = eachcontent.ArabicMetaDescription
// 		contentResult.CreatedAt = eachcontent.CreatedAt
// 		contentResult.ModifiedAt = eachcontent.ModifiedAt
// 		//content id
// 		contentResult.ContentId = eachcontent.ContentId
// 		contentSeasons = append(contentSeasons, contentSeason)
// 		contentResult.ContentSeasons = contentSeasons

// 		if country == "" || country == "all" {
// 			allContents = append(allContents, contentResult)
// 		} else if country != "all" {
// 			if IsCheck {
// 				allContents = append(allContents, contentResult)
// 			}
// 		}

// 		// if country != "" {
// 		// 	if IsCheck {
// 		// 		allContents = append(allContents, contentResult)
// 		// 	}
// 		// } else if country == "" {
// 		// 	allContents = append(allContents, contentResult)
// 		// }

// 	}
// 	/*Pagination*/
// 	var pagination Pagination
// 	pagination.Limit = int(limit)
// 	pagination.Offset = int(offset)
// 	pagination.Size = totalCount
// 	if CountryResult != 0 || country == "" {
// 		c.JSON(http.StatusOK, gin.H{"pagination": pagination, "data": allContents})
// 	}
// }

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
// @Router /v1/contents/multitier/{contentId} [get]
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
	var finalContentResult FinalSeasonResultContentOneTire
	/*for country rights*/
	var country string
	if c.Request.URL.Query()["Country"] != nil {
		country = c.Request.URL.Query()["Country"][0]
	}
	CountryResult := common.Countrys(country)
	_ = CountryResult
	ContentKey, _ := strconv.Atoi(c.Param("contentId"))
	var count int
	if err := db.Table("content").Where("third_party_content_key=?", ContentKey).Count(&count).Error; err != nil {
		serverError.Description = "Content Key Wrong"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	if count < 1 {
		c.JSON(http.StatusNotFound, notFound)
		return
	}

	if UserId == os.Getenv("WATCH_NOW") {

		db.Debug().Raw(`
					select
						c.third_party_content_key as multi_tier_content_key,
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
						c.id as content_id,
						c.created_at,
						c.modified_at
					from
						content c
					join content_primary_info cpi on
						cpi.id = c.primary_info_id
					join content_genre cg on
						cg.content_id = c.id
					where
							c.watch_now_supplier = 'true'
							and c.third_party_content_key = ?
							and c.status = 1
							and c.content_tier = 2
							and c.deleted_by_user_id is null;
		`, ContentKey).Find(&finalContentResult)

	} else {

		db.Debug().Raw(`
					select
						c.third_party_content_key as multi_tier_content_key,
						c.content_type ,
						cpi.original_title,
						cpi.alternative_title,
						cpi.arabic_title,
						cpi.transliterated_title,
						cpi.notes,
						c.english_meta_title,
						c.arabic_meta_title,
						c.english_meta_description,
						c.arabic_meta_description,
						c.id as content_id,
						c.created_at,
						c.modified_at
					from
						content c
					join content_primary_info cpi on
						cpi.id = c.primary_info_id
					join content_genre cg on
						cg.content_id = c.id
					where
							c.third_party_content_key = ?
							and c.status = 1
							and c.content_tier = 2
							and c.deleted_by_user_id is null
					limit 1;
		`, ContentKey).Find(&finalContentResult)

	}

	contentResult.ContentKey = ContentKey
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

	if err := db.Debug().Table("content_genre").Where("content_id =?", finalContentResult.ContentId).Find(&contentGenresIds).Error; err != nil {
		serverError.Description = "Genre Id Wrong"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

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

		if err := db.Debug().Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName).Error; err != nil {
			serverError.Description = "Genre Id Wrong 2"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		if err := db.Debug().Table("content_subgenre").Where("content_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId).Error; err != nil {
			serverError.Description = "Genre Id Wrong 3"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		if err := db.Debug().Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&contentSubGenres).Error; err != nil {
			serverError.Description = "Genre Id Wrong 4"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

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

	var contentSeasons []ContentSeasons
	var seasons []Seasons

	if err :=
		db.Debug().Raw(`
			select
				s.*
			from
				season s
				
			join content_rights cr2 on cr2.id = s.rights_id
			join about_the_content_info atci on  atci.Id = s.about_the_content_info_id
			
			where
				s.deleted_by_user_id is null and
				(cr2.digital_rights_start_date <= NOW() or cr2.digital_rights_start_date is null) and 
				(cr2.digital_rights_end_date >= NOW() or cr2.digital_rights_end_date is null) and 
				atci.supplier !='Others' and
				s.content_id = ?
		`, finalContentResult.ContentId).Find(&seasons).Error; err != nil {
		serverError.Description = "No Season Found"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	for _, seasonId := range seasons {

		var seasonPrimaryInfo SeasonPrimaryInfo

		if err := db.Table("content_primary_info").Where("id =?", seasonId.PrimaryInfoId).Find(&seasonPrimaryInfo).Error; err != nil {
			serverError.Description = "Season Primary Information wrong"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

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

		if err := db.Debug().Raw(`
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
					--join playback_item pi2 on pi2.translation_id = ct.id
					--join content_rights cr on cr.id = pi2.rights_id
						
					where
					--	( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null) and 
					--	(cr.digital_rights_start_date <= NOW() or cr.digital_rights_start_date is null) and 
					--	(cr.digital_rights_end_date >= NOW() or cr.digital_rights_end_date is null) and 
						s.status = 1 and 
					--	s.deleted_by_user_id is null and
						s.id = ?
			`, seasonId.Id).Find(&castList).Error; err != nil {
			serverError.Description = "Season Cast Details getting wrong"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var actorIds, actorEnglishNames, actorArabicNames []string
		var queryActorIds []struct {
			ActorId string `json:"actor_id"`
		}

		if err := db.Table("content_actor").Where("cast_id =?", castList.CastId).Find(&queryActorIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, queryActorId := range queryActorIds {
			actorIds = append(actorIds, queryActorId.ActorId)

			var actorName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("actor").Where("id =?", queryActorId.ActorId).Find(&actorName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			actorEnglishNames = append(actorEnglishNames, actorName.EnglishName)
			actorArabicNames = append(actorArabicNames, actorName.ArabicName)
		}

		var writerIds, writerEnglishNames, writerArabicNames []string
		var queryWriterIds []struct {
			WriterId string `json:"writer_id"`
		}

		if err := db.Table("content_writer").Where("cast_id =?", castList.CastId).Find(&queryWriterIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, queryWriterId := range queryWriterIds {
			writerIds = append(writerIds, queryWriterId.WriterId)

			var writerName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("writer").Where("id =?", queryWriterId.WriterId).Find(&writerName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			writerEnglishNames = append(writerEnglishNames, writerName.EnglishName)
			writerArabicNames = append(writerArabicNames, writerName.ArabicName)
		}

		var directorIds, directorEnglishNames, directorArabicNames []string
		var querydirectorIds []struct {
			DirectorId string `json:"director_id"`
		}

		if err := db.Table("content_director").Where("cast_id =?", castList.CastId).Find(&querydirectorIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, querydirectorId := range querydirectorIds {
			directorIds = append(directorIds, querydirectorId.DirectorId)

			var directorName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("director").Where("id =?", querydirectorId.DirectorId).Find(&directorName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

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

		if err := db.Table("content_singer").Where("music_id =?", seasonId.MusicId).Find(&querySingerIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, querySingerId := range querySingerIds {
			singerIds = append(singerIds, querySingerId.SingerId)

			var SingerName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("singer").Where("id =?", querySingerId.SingerId).Find(&SingerName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			singerEnglishNames = append(singerEnglishNames, SingerName.EnglishName)
			singerArabicNames = append(singerArabicNames, SingerName.ArabicName)
		}

		var musicComposerIds, musicComposerEnglishNames, musicComposerArabicNames []string
		var queryMusicComposerIds []struct {
			MusicComposerId string `json:"music_composer_id"`
		}

		if err := db.Table("content_music_composer").Where("music_id =?", seasonId.MusicId).Find(&queryMusicComposerIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, queryMusicComposerId := range queryMusicComposerIds {
			musicComposerIds = append(musicComposerIds, queryMusicComposerId.MusicComposerId)

			var MusicComposerName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("music_composer").Where("id =?", queryMusicComposerId.MusicComposerId).Find(&MusicComposerName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			musicComposerEnglishNames = append(musicComposerEnglishNames, MusicComposerName.EnglishName)
			musicComposerArabicNames = append(musicComposerArabicNames, MusicComposerName.ArabicName)
		}

		var songWriterIds, songWriterEnglishNames, songWriterArabicNames []string
		var querySongWriterIds []struct {
			SongWriterId string `json:"song_writer_id"`
		}

		if err := db.Table("content_song_writer").Where("music_id =?", seasonId.MusicId).Find(&querySongWriterIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, querySongWriterId := range querySongWriterIds {
			songWriterIds = append(songWriterIds, querySongWriterId.SongWriterId)

			var SongWriterName struct {
				EnglishName string `json:"english_name"`
				ArabicName  string `json:"arabic_name"`
			}

			if err := db.Table("song_writer").Where("id =?", querySongWriterId.SongWriterId).Find(&SongWriterName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

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

		if err := db.Table("content_tag").Where("tag_info_id =?", seasonId.TagInfoId).Find(&queryContentTagIds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var ContentTagIds []string
		for _, queryContentTagId := range queryContentTagIds {

			var TextualDataTagName struct {
				Name string `json:"name"`
			}

			if err := db.Table("textual_data_tag").Where("id =?", queryContentTagId.TextualDataTagId).Find(&TextualDataTagName).Error; err != nil {
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

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

		if err := db.Debug().Table("season_genre").Where("season_id =?", seasonId.Id).Find(&seasonGenresIds).Error; err != nil {
			serverError.Description = "Season Genre Id Wrong 1"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

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

			if err := db.Debug().Select("english_name, arabic_name").Table("genre").Where("id = ?", querygenresId.GenreId).Find(&genresName).Error; err != nil {
				serverError.Description = "Season Genre Id Wrong 2"
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Debug().Table("season_subgenre").Where("season_genre_id =?", querygenresId.Id).Pluck("subgenre_id", &subgenreId).Error; err != nil {
				serverError.Description = "Season Genre Id Wrong 3"
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

			if err := db.Debug().Select("english_name, arabic_name").Table("subgenre").Where("id in (?)", subgenreId).Find(&seasonSubGenres).Error; err != nil {
				serverError.Description = "Season Genre Id Wrong 4"
				c.JSON(http.StatusInternalServerError, serverError)
				return
			}

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
		if err := db.Table("variance_trailer").Where("season_id =?", seasonId.Id).Find(&varianceTrailers).Error; err != nil {
			serverError.Description = "Variance trailer for season not found"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var finalVarianceResult []interface{}

		for _, varianceData := range varianceTrailers {
			var variance VarianceTrailers
			variance.Order = varianceData.Order
			variance.VideoTrailerId = varianceData.VideoTrailerId
			variance.EnglishTitle = varianceData.EnglishTitle
			variance.ArabicTitle = varianceData.ArabicTitle
			variance.Duration = varianceData.Duration
			variance.HasTrailerPosterImage = varianceData.HasTrailerPosterImage
			//const url = "https://z5content-qa.s3.amazonaws.com/"
			variance.TrailerPosterImage = os.Getenv("IMAGERY_URL") + seasonId.ContentId + "/" + seasonId.Id + "/" + varianceData.Id + "/trailer-poster-image"
			variance.Id = varianceData.Id
			finalVarianceResult = append(finalVarianceResult, variance)
		}

		var aboutTheContent AboutTheContent

		if err := db.Debug().Table("about_the_content_info").Where("id =?", seasonId.AboutTheContentInfoId).Find(&aboutTheContent).Error; err != nil {
			serverError.Description = "Season About the content info not found"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		var ageRatings AgeRatingsCode

		if err := db.Debug().Raw("select split_part(code, '_', 2) as code from age_ratings where id = ?", aboutTheContent.AgeGroup).Find(&ageRatings).Error; err != nil {
			serverError.Description = "Season Age group split issue"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		aboutTheContent.AgeGroup = ageRatings.Code

		var productionCountries []int

		var queryProductionCountries []struct {
			CountryId int `json:"country_id"`
		}

		if err := db.Table("production_country").Where("about_the_content_info_id =?", seasonId.AboutTheContentInfoId).Find(&queryProductionCountries).Error; err != nil {
			serverError.Description = "Season About the content info not found 2"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, queryProductionCountry := range queryProductionCountries {
			productionCountries = append(productionCountries, queryProductionCountry.CountryId)
		}

		aboutTheContent.ProductionCountries = productionCountries

		var translation Translation

		if err := db.Table("content_translation").Where("id =?", seasonId.TranslationId).Find(&translation).Error; err != nil {
			serverError.Description = "Season translation details is not getting fetched"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		LanguageTypeId, _ := strconv.Atoi(translation.LanguageType)
		translation.LanguageType = common.LanguageOriginTypes(LanguageTypeId)
		translation.DubbingDialectName = common.DialectIdname(translation.DubbingDialectId, "en")

		// var episodes []Episode
		var episodeResultFinal []ContentEpisode

		var episodeDetails []FetchEpisodeDetailsMultiTier

		if err := db.Debug().Table("season").Where("id =?", seasonId.Id).Find(&translation).Error; err != nil {
			serverError.Description = "Season details not getting fetched"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		if err := db.Debug().Raw(`
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
							join season s on
								s.id = e.season_id
							join playback_item pi2 on
								pi2.id = e.playback_item_id
							join content_rights cr on 
								cr.id = s.rights_id
							where
								e.season_id = ? 
								and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null)
								and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null)
								and ( pi2.scheduling_date_time <= NOW() or pi2.scheduling_date_time is null)
								and e.status = 1 
								and e.deleted_by_user_id is null

							order by
								e.number asc
		`, seasonId.Id).Find(&episodeDetails).Error; err != nil {
			serverError.Description = "Season episode details is not getting fetched"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		if len(episodeDetails) == 0 {
			if len(varianceTrailers) == 0 {
				serverError.Description = "No Episode Found"
				c.JSON(http.StatusInternalServerError, serverError)
				return
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

		var flag bool = false

		var digitalRightsRegions []DigitalRightsRegions

		db.Table("content_rights_country").Select("country_id").Where("content_rights_id=?", seasonId.RightsId).Scan(&digitalRightsRegions)

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

		if !flag {
			contentSeasons = append(contentSeasons, ContentSeasons{
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
		}

	}

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

	if len(contentSeasons) > 0 {
		c.JSON(http.StatusOK, gin.H{"data": contentResult})
	} else {
		serverError.Description = "No Record Found"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

}
