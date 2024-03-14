package content

import (
	common "masterdata/common"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// const IMAGES = "https://content.weyyak.com/"

// GetPageDetails - Get  page details
// GET /v1/get_page/:pageId
// @Description Get All menu list details by platform ID
// @Tags Menu
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param pageId path string true "Page Id"
// @Param Country query string false "Country code of the user."
// @Success 200  object PageDetails
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/get_page/{pageId} [get]
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
	_ = countryId
	serverError := common.ServerErrorResponse()
	db := c.MustGet("FCDB").(*gorm.DB)

	var details PageDetails
	var menu MenuDatas
	if err := db.Table("page p").Select("p.*").Where("p.is_disabled=false and p.deleted_by_user_id is null and p.page_key=?", pagekey).Find(&details).Error; err != nil {
		serverError.Description = "Page Query Failed"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	type PageIds struct {
		Id string `json:"id"`
	}
	var pageids []PageIds
	var ids []string
	if err := db.Table("page p").Select("p.id").Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null  and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW() and p.page_key=?", pagekey).Find(&pageids).Error; err != nil {
		serverError.Description = "Page Query Failed2"
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

	// var blackPlaylistCount, redPlaylistCount, greenPlaylistCount int
	var sliders []Slider
	var featuredDetails FeaturedDetails

	if err := db.Select("s.*").Table("slider s").
		Joins("inner join page_slider ps on ps.slider_id=s.id").
		Where("s.deleted_by_user_id  is null and s.is_disabled = false and ps.page_id=? and (s.scheduling_start_date <=NOW() or ps.order =0) and (s.scheduling_end_date >=NOW()  or ps.order =0)", details.Id).Order("ps.order desc").Find(&sliders).Error; err != nil && err.Error() != "record not found" {
		serverError.Description = "Slider Query Failed"
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	for _, slider := range sliders {

		var featuredPlaylists []FeaturedPlaylists

		if slider.BlackAreaPlaylistId != "" || slider.RedAreaPlaylistId != "" || slider.GreenAreaPlaylistId != "" {

			playlists, _ := SliderPlaylists(slider.BlackAreaPlaylistId, slider.RedAreaPlaylistId, slider.GreenAreaPlaylistId, c)

			for _, playlist := range playlists {

				contentIds, err := PlaylistItemContents(playlist.ID, c)
				if err != nil {
					serverError.Description = "No Contents Found for " + playlist.ID
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

				playlistContents, err := GetPlayListDetailss(c, contentIds)
				if err != nil {
					serverError.Description = "No Playlist Found"
					c.JSON(http.StatusInternalServerError, serverError)
					return
				}

				featuredPlaylists = append(featuredPlaylists, FeaturedPlaylists{
					ID:           int32(playlist.PlaylistKey),
					PlaylistType: playlist.PlaylistType,
					Content:      playlistContents,
				})

			}
		}

		featuredDetails = FeaturedDetails{
			ID:        int64(slider.SliderKey),
			Type:      common.SliderTypes(slider.Type),
			Playlists: featuredPlaylists,
		}

	}

	if len(featuredDetails.Playlists) > 0 {
		menu.Featured = &[]FeaturedDetails{featuredDetails}[0]
	} else {
		menu.Featured = nil
	}

	// var playlists []MenuPlaylists

	// _ = playlists

	var playlists []Playlist
	if err := db.Select("p.id,english_title,arabic_title,p.scheduling_start_date,p.scheduling_end_date,p.deleted_by_user_id,p.is_disabled,p.created_at,p.playlist_key,p.modified_at,p.playlist_type").Table("page_playlist pp").Joins("join playlist p on p.id =pp.playlist_id").Where("p.is_disabled =false and p.deleted_by_user_id is null and pp.page_id =? and (p.scheduling_start_date <=now() or p.scheduling_start_date is null) and (p.scheduling_end_date >=now() or p.scheduling_end_date is null)", details.Id).Order("pp.order asc").Find(&playlists).Error; err != nil {
		serverError.Description = "No Playlist Found2"
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
			serverError.Description = "No Playlist Content Not Found"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}
		// var Ids []string
		// for _, content := range contentIds {
		// 	Ids = append(Ids, content.ContentId)
		// }

		var playlistContent []PlaylistContent

		playlistContents, err := GetPlayListDetailss(c, contentIds)
		if err != nil {
			serverError.Description = "No Playlist Details Not Found"
			c.JSON(http.StatusInternalServerError, serverError)
			return
		}

		for _, playlistContentId := range playlistContents {

			if playlistContentId.ContentId != "" {
				playlistContent = append(playlistContent, playlistContentId)
			}

		}

		pagePlaylist.Content = playlistContent

		pagePlaylists = append(pagePlaylists, pagePlaylist)

	}

	menu.Playlists = pagePlaylists

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func GetPlayListDetailss(c *gin.Context, contentIds []PlaylistContentIds) ([]PlaylistContent, error) {
	var playlistContents []PlaylistContent
	cdb := c.MustGet("DB").(*gorm.DB)
	serverError := common.ServerErrorResponse()

	for _, contentId := range contentIds {

		var content ContentDetails
		var queryContents ContentDetails

		if err := cdb.Table("content").Where("id = ?", contentId.ContentId).Find(&queryContents).Limit(1).Error; err != nil {
			serverError.Description = "Content Missing " + contentId.ContentId
			c.JSON(http.StatusInternalServerError, serverError)
			return nil, err
		}

		var contentImageryDetails ContentImageryDetails

		if queryContents.ContentType == "Series" || queryContents.ContentType == "Program" {

			cdb.Debug().Raw(`
					select
						c.id,
						c.content_key as content_key,
						ar.english_name as age_rating,
						--	pi.video_content_id as video_id,
						(
							select
									pi2.video_content_id
										from
											playback_item pi2
										where
											pi2.id in (select e1.playback_item_id from episode e1 where season_id = s.id ORDER by e1.number asc limit 1) 
						) as video_id,
						replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,
						lower(c.content_type) as content_type,
						atci.english_synopsis as synopsis_english,
						atci.arabic_synopsis as synopsis_arabic,
						s.english_meta_title as seo_title_english,
						s.arabic_meta_title as seo_title_arabic,
						s.english_meta_description as seo_description_english,
						s.arabic_meta_description as seo_description_arabic,
						pi.duration as length,
						cpi.transliterated_title as title_english,
						cpi.arabic_title as title_arabic,
						s.english_meta_title as seo_title,
						c.created_at as inserted_at,
						c.modified_at,
						s.id as varience_id
					from
						content c
						
					join season s on s.content_id = c.id
					
					join episode e on e.season_id = s.id
					
					join about_the_content_info atci on atci.id = s.about_the_content_info_id
					
					join age_ratings ar on ar.id = atci.age_group 
					
					join playback_item pi on pi.translation_id = s.translation_id
					
					join content_primary_info cpi on cpi.id = c.primary_info_id
					
					join content_rights cr on cr.id = pi.rights_id
					
					--join content_rights_country crc on crc.content_rights_id = cr.id
						
					where
						c.id = ?
						and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null)
						and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null)
						and ( pi.scheduling_date_time <= NOW() or pi.scheduling_date_time is null)
						and c.status = 1
						and c.deleted_by_user_id is null
						and s.status = 1
						and s.deleted_by_user_id is null
						and e.status = 1
						and e.deleted_by_user_id is null
					
					GROUP BY
						c.id,
						c.content_key,
						pi.video_content_id,
						c.average_rating,
						ar.english_name,
						cpi.transliterated_title,
						c.content_type,
						atci.english_synopsis,
						atci.arabic_synopsis,
						s.english_meta_title,
						s.arabic_meta_title,
						s.english_meta_description,
						s.arabic_meta_description,
						cpi.transliterated_title,
						cpi.arabic_title,
						s.english_meta_title,
						c.created_at,
						c.modified_at,
						pi.duration,
						e.number,
						s.id
					
					limit 1;				
			`, contentId.ContentId).Find(&content)

			contentImageryDetails.Thumbnail = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("POSTER_IMAGE")
			contentImageryDetails.Backdrop = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("DETAILS_BACKGROUND")
			contentImageryDetails.MobileImg = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("MOBILE_DETAILS_BACKGROUND")
			contentImageryDetails.FeaturedImg = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("POSTER_IMAGE")
			contentImageryDetails.OverlayPoster = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("OVERLAY_POSTER_IMAGE")

		} else if queryContents.ContentType == "Movie" || queryContents.ContentType == "LiveTV" || queryContents.ContentType == "Play" {
			cdb.Debug().Raw(`
					select
						c.id,
						c.content_key,
						round(cast(c.average_rating as numeric), 1) as average_rating,
						ar.english_name as age_rating,
						pi2.video_content_id as video_id,
						replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,
						lower(c.content_type) as content_type,
						atci.english_synopsis as synopsis_english,
						atci.arabic_synopsis as synopsis_arabic,
						c.english_meta_title as seo_title_english,
						c.arabic_meta_title as seo_title_arabic,
						c.english_meta_description as seo_description_english,
						c.arabic_meta_description as seo_description_arabic,
						min(pi2.duration) as length,
						cpi.transliterated_title as title_english,
						cpi.arabic_title as title_arabic,
						c.english_meta_description as seo_title,
						c.created_at as inserted_at,
						c.modified_at,
						cv.id as varience_id
					from
						content c
					join content_primary_info cpi on
						cpi.id = c.primary_info_id
					join about_the_content_info atci on
						atci.id = c.about_the_content_info_id
					join content_variance cv on
						cv.content_id = c.id
					join playback_item pi2 on
						pi2.id = cv.playback_item_id
					join content_rights cr on
						cr.id = pi2.rights_id
					join content_rights_country crc on
						crc.content_rights_id = cr.id
						
					join age_ratings ar on ar.id = atci.age_group 
					
					where
						c.id = ?
						and (cr.digital_rights_start_date <= now() or cr.digital_rights_start_date is null)
						and (cr.digital_rights_end_date >= now() or cr.digital_rights_end_date is null)
						and c.status = 1
						and c.deleted_by_user_id is null
						and cv.status = 1
						and cv.deleted_by_user_id is null
					
					GROUP BY
						c.id,
						c.content_key,
						ar.english_name,
						pi2.video_content_id,
						c.average_rating,
						cpi.transliterated_title,
						c.content_type,
						atci.english_synopsis,
						atci.arabic_synopsis,
						c.english_meta_title,
						c.arabic_meta_title,
						c.english_meta_description,
						c.arabic_meta_description,
						cpi.transliterated_title,
						cpi.arabic_title,
						c.english_meta_title,
						c.created_at,
						c.modified_at,
						pi2.duration,
						cv.id
					
					ORDER BY
						pi2.duration asc
									
					limit 1;
			`, contentId.ContentId).Find(&content)

			contentImageryDetails.Thumbnail = os.Getenv("IMAGE_URL") + content.Id + os.Getenv("POSTER_IMAGE")
			contentImageryDetails.Backdrop = os.Getenv("IMAGE_URL") + content.Id + os.Getenv("DETAILS_BACKGROUND")
			contentImageryDetails.MobileImg = os.Getenv("IMAGE_URL") + content.Id + os.Getenv("MOBILE_DETAILS_BACKGROUND")
			contentImageryDetails.FeaturedImg = os.Getenv("IMAGE_URL") + content.Id + os.Getenv("POSTER_IMAGE")
			contentImageryDetails.OverlayPoster = os.Getenv("IMAGE_URL") + content.Id + "/" + content.VarienceId + os.Getenv("OVERLAY_POSTER_IMAGE")
		}

		if content.ContentKey != 0 {
			playlistContents = append(playlistContents, PlaylistContent{
				ContentId:             content.Id,
				ContentKey:            content.ContentKey,
				AgeRating:             content.AgeRating,
				VideoId:               content.VideoId,
				FriendlyUrl:           content.FriendlyUrl,
				ContentType:           content.ContentType,
				SynopsisEnglish:       content.SynopsisEnglish,
				SynopsisArabic:        content.SynopsisArabic,
				SeoTitleEnglish:       content.SeoTitleEnglish,
				SeoTitleArabic:        content.SeoTitleArabic,
				SeoDescriptionEnglish: content.SeoDescriptionEnglish,
				SeoDescriptionArabic:  content.SeoDescriptionArabic,
				Length:                content.Length,
				TitleEnglish:          content.TitleEnglish,
				TitleArabic:           content.TitleArabic,
				SeoTitle:              content.SeoTitle,
				Imagery:               contentImageryDetails,
				InsertedAt:            content.InsertedAt,
				ModifiedAt:            content.ModifiedAt,
			})
		}
	}

	return playlistContents, nil
}

func (hs *HandlerService) GetPageDetailss(c *gin.Context) {
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
	cdb := c.MustGet("DB").(*gorm.DB)
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
	if err := db.Select("s.*").Table("slider s").
		Joins("inner join page_slider ps on ps.slider_id=s.id").
		Where("s.deleted_by_user_id  is null and s.is_disabled =false and ps.page_id=? and (s.scheduling_start_date <=NOW() or ps.order =0) and (s.scheduling_end_date >=NOW()  or ps.order =0)", details.Id).
		Limit(1).Order("ps.order desc").Find(&slider).Error; err != nil && err.Error() != "record not found" {
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

							var contentUpdates AboutTheContentInfo

							cdb.Raw(`
								select atci.age_group from season s 
								join about_the_content_info atci on atci.id = s.about_the_content_info_id
								where content_id = ?
							`, content.ContentId).Find(&contentUpdates)

							AgeGroup, _ := strconv.Atoi(contentUpdates.AgeGroup)
							content.AgeRating = common.AgeRatings(AgeGroup, "en")
							content.SeoTitle = content.SeoDescriptionEnglish
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
		// if len(featuredDetails.Playlists) > 0 {
		// 	menu.Featured = &featuredDetails
		// }
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
					var contentUpdates AboutTheContentInfo
					cdb.Raw(`
								select atci.age_group from season s 
								join about_the_content_info atci on atci.id = s.about_the_content_info_id
								where content_id = ?
							`, content.ContentId).Find(&contentUpdates)

					db.Table("page p").Select("p.*").Where("p.is_disabled=false and p.deleted_by_user_id is null and p.page_key=?", pagekey).Find(&details)

					AgeGroup, _ := strconv.Atoi(contentUpdates.AgeGroup)
					content.AgeRating = common.AgeRatings(AgeGroup, "en")
					content.SeoTitle = content.SeoDescriptionEnglish
					contents = append(contents, content)
				}
			}
		}
		pagePlaylist.Content = contents
		pagePlaylists = append(pagePlaylists, pagePlaylist)
	}
	menu.Playlists = pagePlaylists
	c.JSON(http.StatusOK, gin.H{"data": menu})
	// return
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
	if err := db.Raw(`
				select
					c.id,
					c.content_key,
					round(cast(c.average_rating as numeric), 1) as average_rating,
					min(pi2.video_content_id) as video_id,
					replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,
					lower(c.content_type) as content_type,
					atci.english_synopsis as synopsis_english,
					atci.arabic_synopsis as synopsis_arabic,
					c.english_meta_title as seo_title_english,
					c.arabic_meta_title as seo_title_arabic,
					c.english_meta_description as seo_description_english,
					c.arabic_meta_description as seo_description_arabic,
					min(pi2.duration) as length,
					cpi.transliterated_title as title_english,
					cpi.arabic_title as title_arabic,
					c.english_meta_description as seo_title,
					c.created_at as inserted_at,
					c.modified_at,
					cv.id as varience_id
				from
					content c
				join content_primary_info cpi on
					cpi.id = c.primary_info_id
				join about_the_content_info atci on
					atci.id = c.about_the_content_info_id
				join content_variance cv on
					cv.content_id = c.id
				join playback_item pi2 on
					pi2.id = cv.playback_item_id
				join content_rights cr on
					cr.id = pi2.rights_id
				join content_rights_country crc on
					crc.content_rights_id = cr.id
				where
					crc.country_id = ?
					and (cr.digital_rights_start_date <= now()
						or cr.digital_rights_start_date is null)
					and (cr.digital_rights_end_date >= now()
						or cr.digital_rights_end_date is null)
					and c.status = 1
					and c.deleted_by_user_id is null
					and c.id in(?)
					and cv.status = 1
					and cv.deleted_by_user_id is null
				group by
					c.id,
					c.content_key,
					c.average_rating,
					cpi.transliterated_title,
					c.content_type,
					atci.english_synopsis,
					atci.arabic_synopsis,
					c.english_meta_title,
					c.arabic_meta_title,
					c.english_meta_description,
					c.arabic_meta_description,
					cpi.transliterated_title,
					cpi.arabic_title,
					c.english_meta_title,
					c.created_at,
					c.modified_at,
					cv.id
				union
				select
					c.id,
					c.content_key,
					round(cast(c.average_rating as numeric), 1) as average_rating,
					min(pi2.video_content_id) as video_id,
					replace(lower(cpi.transliterated_title), ' ', '-') as friendly_url,
					lower(c.content_type) as content_type,
					atci.english_synopsis as synopsis_english,
					atci.arabic_synopsis as synopsis_arabic,
					s.english_meta_title as seo_title_english,
					s.arabic_meta_title as seo_title_arabic,
					s.english_meta_description as seo_description_english,
					s.arabic_meta_description as seo_description_arabic,
					min(pi2.duration) as length,
					cpi.transliterated_title as title_english,
					cpi.arabic_title as title_arabic,
					s.english_meta_title as seo_title,
					c.created_at as inserted_at,
					c.modified_at,
					s.id as varience_id
				from
					content c
				join season s on
					s.content_id = c.id
				join episode e on
					e.season_id = s.id
				join content_primary_info cpi on
					cpi.id = s.primary_info_id
				join about_the_content_info atci on
					atci.id = s.about_the_content_info_id
				join playback_item pi2 on
					pi2.id = (select pi1.id from playback_item pi1 where pi1.id = e.playback_item_id order by pi1.duration asc limit 1) 
				join content_rights cr on
					cr.id = pi2.rights_id
				join content_rights_country crc on
					crc.content_rights_id = cr.id
				where
					crc.country_id = ?
					and (cr.digital_rights_start_date <= now()
						or cr.digital_rights_start_date is null)
					and (cr.digital_rights_end_date >= now()
						or cr.digital_rights_end_date is null)
					and c.status = 1
					and c.deleted_by_user_id is null
					and c.id in(?)
					and s.status = 1
					and s.deleted_by_user_id is null
					and e.status = 1
					and e.deleted_by_user_id is null
				group by
					c.id,
					c.content_key,
					c.average_rating,
					cpi.transliterated_title,
					c.content_type,
					atci.english_synopsis,
					atci.arabic_synopsis,
					s.english_meta_title,
					s.arabic_meta_title,
					s.english_meta_description,
					s.arabic_meta_description,
					cpi.transliterated_title,
					cpi.arabic_title,
					s.english_meta_title,
					c.created_at,
					c.modified_at,
					s.id
	`, country, contentIds, country, contentIds).Find(&cDetails).Error; err != nil {
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
