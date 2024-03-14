package playlist

import (
	"context"
	"fmt"
	"frontend_config/common"
	"frontend_config/fragments"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"github.com/thanhpk/randstr"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	//dqrg.Use(ValidateToken())
	brg := r.Group("/api")
	brg.Use(common.ValidateToken())
	brg.GET("/playlists/:playlist_id", hs.GetPlaylistDetailsByPlaylistId)
	brg.POST("/playlists/:playlistId/available", hs.DisablePlaylistByPlaylistId)
	brg.GET("/playlists/summary", hs.GetPlaylistsBasedOnSearchInPage)
	brg.DELETE("/playlists/:Id", hs.DeletePlaylistByPlaylistId)
	brg.GET("/playlists", hs.GetAllPlaylistAndPlaylistsBySearchText)
	brg.POST("/playlists/:playlistId", hs.CreateAndUpdatePlaylist) // Update
	brg.POST("/playlists", hs.CreateAndUpdatePlaylist)             // Insert
	brg.GET("/playlists/:playlist_id/region", hs.GetAllRegionsBasedOnPlaylistId)

	frg := r.Group("/api/playlist")
	frg.Use(common.ValidateToken())
	frg.GET("/sourceitems", hs.GetSourceitemsBySourceFields)
	frg.GET("/sourceitems/searchfilters", hs.GetSearchFiltersForPlaylists)
	frg.GET("/sourceitems/types", hs.getSourceItem)
	frg.GET("/dynamicgroupitems", hs.GetPlaylistItemBasedOnPlaylistItemTypeAndItemTypeID)

	region := r.Group("/api")
	region.Use(common.ValidateToken())

}

// GetAllPlaylistAndPlaylistsBySearchText - Get all playlist and playlists by search text
// GET api/playlists
// @Summary Get all playlist and playlists by search text
// @Description Get all playlist and playlists by search text
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router api/playlists [get]
func (hs *HandlerService) GetAllPlaylistAndPlaylistsBySearchText(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var limit, offset int64
	var searchText string
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.ToLower(c.Request.URL.Query()["searchText"][0])
	}
	if limit == 0 {
		limit = 50
	}

	playlists := []Playlists{}
	var playlist, totalCount []Playlists
	serverError := common.ServerErrorResponse()

	// rawquery := "select p.id,p.english_title,p.arabic_title,p.scheduling_end_date,p.is_disabled,rtrim(substring( string_agg(distinct c.english_name ::varchar, ',' order by c.english_name asc ) from '(:[^,]+,){1,5}'), ',') as region,case when rtrim(substring( string_agg( p2.page_type ::varchar, ',' order by p2.page_type desc ) from '(:[^,]+,){1,1}'), ',') = '1' then 1 else 0 end as is_assigned_to_any_home_page,'' as the_only_playlist_for,string_agg(distinct p2.english_title,',')::varchar as found_in,json_agg(distinct plp.target_platform)::varchar as platforms from playlist p left join page_playlist pp on pp.playlist_id =p.id left join page p2 on p2.id =pp.page_id  left join play_list_platform plp on plp.play_list_id =p.id left join play_list_country plc on p.id = plc.play_list_id left join country c on plc.country_id = c.id where  p.deleted_by_user_id is null "

	rawquery := `
		SELECT
			p.id,
			p.english_title,
			p.arabic_title,
			p.scheduling_end_date,
			p.is_disabled,
			rtrim(substring(string_agg(DISTINCT c.english_name::varchar, ','), '(:[^,]+,){1,5}'), ',') AS region,
			CASE
				WHEN rtrim(substring(string_agg(p2.page_type::varchar, ','), '(:[^,]+,){1,1}'), ',') = '1' THEN 1
				ELSE 0
			END AS is_assigned_to_any_home_page,
			'' AS the_only_playlist_for,
			string_agg(DISTINCT p2.english_title, ',')::varchar AS found_in,
			array_to_json(array_agg(DISTINCT plp.target_platform))::varchar AS platforms
		FROM
			playlist p
		LEFT JOIN page_playlist pp ON
			pp.playlist_id = p.id
		LEFT JOIN page p2 ON
			p2.id = pp.page_id
		LEFT JOIN play_list_platform plp ON
			plp.play_list_id = p.id
		LEFT JOIN play_list_country plc ON
			p.id = plc.play_list_id
		LEFT JOIN country c ON
			plc.country_id = c.id
		WHERE
			p.deleted_by_user_id IS NULL
	`

	if searchText != "" {
		rawquery += " and ( lower(p.english_title) like '%" + searchText + "%' OR  lower(p.arabic_title) like '%" + searchText + "%' )"
	}
	rawquery += " group by p.id order by p.modified_at desc"
	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&playlist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	var singlePlaylist Playlists
	for _, pltid := range playlist {
		singlePlaylist.ArabicTitle = pltid.ArabicTitle
		singlePlaylist.EnglishTitle = pltid.EnglishTitle
		singlePlaylist.FoundIn = pltid.FoundIn
		singlePlaylist.ID = pltid.ID
		singlePlaylist.IsAssignedToAnyHomePage = pltid.IsAssignedToAnyHomePage
		singlePlaylist.IsDisabled = pltid.IsDisabled
		platforms, _ := common.JsonStringToIntSliceOrMap(singlePlaylist.Platforms)
		fmt.Println(platforms)
		singlePlaylist.PublishingPlatforms = platforms
		if len(platforms) < 1 {
			buffer := make([]int, 0)
			singlePlaylist.PublishingPlatforms = buffer
		}
		singlePlaylist.Region = pltid.Region
		singlePlaylist.SchedulingEndDate = pltid.SchedulingEndDate
		singlePlaylist.TheOnlyPlaylistFor = pltid.TheOnlyPlaylistFor
		playlists = append(playlists, singlePlaylist)
	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": playlists})
}

// GetSearchFiltersForPlaylists - Get Search filters for playlists
// GET /playlist/sourceitems/searchfilters
// @Summary  Get Search filters for playlists
// @Description  Get Search filters for playlists
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /playlist/sourceitems/searchfilters [get]
func (hs *HandlerService) GetSearchFiltersForPlaylists(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var sourceitemtypes []SourceItemTypes
	serverError := common.ServerErrorResponse()
	rawquery := "select name,id from search"
	if errCount := db.Debug().Raw(rawquery).Scan(&sourceitemtypes).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sourceitemtypes})
}

// DeletePlaylistByPlaylistId - Delete playlist by playlist_id
// DELETE /playlist/{Id}
// @Summary  Delete playlist by playlist_id
// @Description  Delete playlist by playlist_id
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param Playlist Id path string true "playlist Id"
// @Success 200 {array} object c.JSON
// @Router /playlist/{Id} [delete]
func (hs *HandlerService) DeletePlaylistByPlaylistId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var playlist Playlistupdate
	UserId := c.MustGet("userid")
	playlistid := c.Param("Id")
	if err := db.Debug().Table("playlist").Where("id=? and deleted_by_user_id is null ", playlistid).Find(&playlist).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Record does not exist.Please provide valid playlist Id.", "Status": http.StatusBadRequest})
		return
	}
	playlist.DeletedByUserId = UserId.(string)
	playlist.ModifiedAt = time.Now()
	if err := db.Debug().Table("playlist").Where("id=?", playlistid).Update(&playlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Record Deleted Successfully.", "Status": http.StatusOK})
		return

	}
}

// DisablePlaylistByPlaylistId -  Disable Playlist By Playlist Id
// POST /playlists/{Id}/available
// @Summary Disable Playlist By Playlist Id
// @Description Disable Playlist By Playlist Id
// @Tags User
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param Playlist Id path string true "playlist Id"
// @Param body body PlaylistDisable true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /playlists/{Id}/available [post]
func (hs *HandlerService) DisablePlaylistByPlaylistId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var playlist Playlistupdate
	var input PlaylistDisable
	// Check Post data is json formated or Not
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	playlistidRrl := c.Param("playlistId")
	playlistid := input.ID
	if playlistidRrl != playlistid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide valid playlist Id.", "Status": http.StatusBadRequest})
		return
	}
	var totalcount int
	if err := db.Debug().Table("playlist").Where("id=? ", playlistid).Count(&totalcount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	playlist.IsDisabled = input.IsDisabled
	playlist.ModifiedAt = time.Now()
	if totalcount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "error_playlist_not_found", "message": "The specified condition was not met for 'Id'.", "Status": http.StatusBadRequest})
		return
	} else {
		if err := db.Debug().Table("playlist").Where("id=?", playlistid).Update(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		/* update dirty count in playlist_sync table */
		go common.PlaylistSynching(playlistid, c)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
}

// GetPlaylistsBasedOnSearchInPage - Get Playlists Based On Search In Page
// GET /playlists/summary/
// @Summary Get Playlists Based On Search In Page
// @Description Get Playlists Based On Search In Page
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} object c.JSON
// @Router /playlists/summary/ [get]
func (hs *HandlerService) GetPlaylistsBasedOnSearchInPage(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var limit, offset int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if limit == 0 {
		limit = 2
	}
	var searchText string
	playlists := []PlaylistSummary{}
	var playlist, totalCount []PlaylistSummary

	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.ToLower(c.Request.URL.Query()["searchText"][0])
	}

	rawquery := "select id, english_title, arabic_title, is_disabled  from playlist where  deleted_by_user_id is null and is_disabled=false"
	if searchText != "" {
		rawquery += " and ( lower(english_title) like '%" + searchText + "%' OR lower(arabic_title) like '%" + searchText + "%' )"
	}
	rawquery += "  group by id "
	if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&playlist).Error; data != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	var platformDetails []PlanNames
	for _, pltid := range playlist {
		rawpidlist := "select plp.target_platform from playlist p left join play_list_platform plp on p.id = plp.play_list_id where p.id = '" + pltid.ID + "' "
		if dataId := db.Debug().Raw(rawpidlist).Scan(&platformDetails).Error; dataId != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
			return
		}
		var Ids []int
		for _, idarr := range platformDetails {
			Ids = append(Ids, idarr.TargetPlatform)
		}
		pltid.PublishingPlatforms = Ids
		playlists = append(playlists, pltid)
	}
	if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": playlists})
}

// GetPlaylistDetailsByPlaylistId -  Get Playlist Details By Playlist Id
// GET /playlists/{playlist_id}
// @Summary Get Playlist Details By Playlist Id
// @Description Get Playlist Details By Playlist Id
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router /playlists/{playlist_id} [get]
func (hs *HandlerService) GetPlaylistDetailsByPlaylistId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	playlistid := c.Param("playlist_id")

	var playlist, Result PlaylistDetails

	// rawquery := "select p.id,p.english_title,p.arabic_title,p.is_disabled,p.playlist_type,p.scheduling_start_date,p.scheduling_end_date from playlist p join page_playlist pp on p.id = pp.playlist_id where p.deleted_by_user_id is null and p.id ='" + playlistid + "' "

	// if err := db.Raw(rawquery).Find(&playlist).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
	// 	return
	// }
	/* Hardcoded values There is no meta information for playlist */
	// playlist.EnglishMetaTitle = nil
	// playlist.ArabicMetaTitle = nil
	// playlist.EnglishMetaDescription = nil
	// playlist.ArabicMetaDescription = nil

	//var platformDetails []string
	//var PlaylistRegion []string

	// // For platform start
	// rawpidlist := "select plp.target_platform from playlist p join play_list_platform plp on p.id = plp.play_list_id where p.id = '" + playlist.ID + "' "
	// if dataId := db.Raw(rawpidlist).Scan(&platformDetails).Error; dataId != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
	// 	return
	// }
	// var Ids []int
	// for _, idarr := range platformDetails {
	// 	Ids = append(Ids, idarr.TargetPlatform)
	// }
	// playlist.PublishingPlatforms = Ids
	// if len(Ids) < 1 {
	// 	buffer := make([]int, 0)
	// 	playlist.PublishingPlatforms = buffer
	// }
	// //playlists =append(playlists,pltid)

	// // for regions start
	// regidlist := "select plc.country_id from playlist p  join play_list_country plc on p.id = plc.play_list_id where p.id = '" + playlist.ID + "' "
	// if regs := db.Raw(regidlist).Scan(&PlaylistRegion).Error; regs != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
	// 	return
	// }
	// var Rids []int
	// for _, ridarr := range PlaylistRegion {
	// 	Rids = append(Rids, ridarr.CountryId)
	// }
	// playlist.Regions = Rids
	// if len(Rids) < 1 {
	// 	buffer := make([]int, 0)
	// 	playlist.Regions = buffer
	// }

	/* Platforms */
	// platformDetails = strings.Split(playlist.Platforms, ",")
	// //playlistPlatforms, _ := common.JsonStringToIntSliceOrMap(playlist.Platforms)
	// for _,platforms:=range platformDetails{
	// 	value, _ := strconv.Atoi(platforms)
	// 	playlist.PublishingPlatforms = append(playlist.PublishingPlatforms, value)
	// }
	// if len(playlist.PublishingPlatforms) < 1 {
	// 	buffer := make([]int, 0)
	// 	playlist.PublishingPlatforms = buffer
	// }
	/* Regions */
	// PlaylistRegion = strings.Split(playlist.Country, ",")
	// //playlistRegions, _ := common.JsonStringToIntSliceOrMap(playlist.Country)
	// for _,region:=range PlaylistRegion{
	// 	value, _ := strconv.Atoi(region)
	// 	playlist.Regions = append(playlist.Regions, value)
	// 	fmt.Println(region,"??????????????????")
	// 	fmt.Println(value,"??????????????????")
	// 	fmt.Println(playlist.Regions,"??????????????????")
	// }
	// if len(playlist.Regions) < 1 {
	// 	buffer := make([]int, 0)
	// 	playlist.Regions = buffer
	// }

	if err := db.Debug().Table("playlist p").Select("p.id,p.english_title,p.arabic_title,p.is_disabled,p.playlist_type,p.scheduling_start_date,p.scheduling_end_date,json_agg(distinct plp.target_platform)::varchar as platforms,json_agg(distinct plc.country_id)::varchar as country").Joins("left join play_list_country plc on plc.play_list_id =p.id left join play_list_platform plp on plp.play_list_id =p.id").Where("p.deleted_by_user_id is null and p.id = ?", playlistid).Group("p.id").Find(&Result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	playlist.ID = Result.ID
	playlist.EnglishTitle = Result.EnglishTitle
	playlist.ArabicTitle = Result.ArabicTitle
	playlist.IsDisabled = Result.IsDisabled
	playlist.PlaylistType = Result.PlaylistType
	playlist.SchedulingStartDate = Result.SchedulingStartDate
	playlist.SchedulingEndDate = Result.SchedulingEndDate
	/* Hardcoded values There is no meta information for playlist */
	playlist.EnglishMetaTitle = nil
	playlist.ArabicMetaTitle = nil
	playlist.EnglishMetaDescription = nil
	playlist.ArabicMetaDescription = nil
	/* Platforms */
	if Result.Platforms == "[null]" {
		buffer := make([]int, 0)
		playlist.PublishingPlatforms = buffer
	} else {
		platform, _ := common.JsonStringToIntSliceOrMap(Result.Platforms)
		playlist.PublishingPlatforms = platform
	}
	/* Regions */
	if Result.Country == "[null]" {
		buffer := make([]int, 0)
		playlist.Regions = buffer
	} else {
		region, _ := common.JsonStringToIntSliceOrMap(Result.Country)
		playlist.Regions = region
	}

	//For Pages details
	var pages []Pages
	pagelist := "select p.english_title,p.arabic_title,p.is_disabled,p.page_type as is_home,p.id from page_playlist pp join page p on pp.page_id = p.id where pp.playlist_id = '" + Result.ID + "' "
	if pagedet := db.Debug().Raw(pagelist).Scan(&pages).Error; pagedet != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}
	playlist.Pages = pages

	// For PlaylistItems
	var playlistItem PlaylistItems
	var playlistItems []PlaylistItems
	var contentidata []ContentIds
	//bkp

	contentidlist := " select case when one_tier_content_id is not null then one_tier_content_id when multi_tier_content_id is not null then multi_tier_content_id end as content_id, case when season_id is not null then season_id else null  end as season_id,case when group_by_genre_id is not null then group_by_genre_id else null  end as group_by_genre_id, case when group_by_subgenre_id is not null then group_by_subgenre_id else null  end as group_by_subgenre_id, case when group_by_actor_id is not null then group_by_actor_id	else null  end as group_by_actor_id, case when group_by_writer_id is not null then group_by_writer_id	else null  end as group_by_writer_id, case when group_by_director_id is not null then group_by_director_id	else null  end as group_by_director_id, case when group_by_singer_id is not null then group_by_singer_id	else null  end as group_by_singer_id, case when group_by_music_composer_id is not null then group_by_music_composer_id	else null  end as group_by_music_composer_id, case when group_by_song_writer_id is not null then group_by_song_writer_id	else null  end as group_by_song_writer_id,case when group_by_original_language_code is not null then group_by_original_language_code else null  end as group_by_original_language_code, case when group_by_production_year is not null then group_by_production_year else null  end as group_by_production_year, case when group_by_page_id is not null then group_by_page_id else null  end as group_by_page_id from playlist_item p where playlist_id ='" + Result.ID + "' order by p.order"

	if plcont := db.Debug().Raw(contentidlist).Scan(&contentidata).Error; plcont != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
		return
	}

	for _, contentidarr := range contentidata {

		if contentidarr.ContentId != "" {
			contentdetailsquery := "select cpi.transliterated_title as name,c.content_tier as type, c.id  from content c left join content_primary_info cpi  on c.primary_info_id =cpi.id where c.id ='" + contentidarr.ContentId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.SeasonId != "" {
			contentdetailsquery := "select  distinct cpi.transliterated_title as name,  case WHEN (cpi.transliterated_title != '') THEN '3' end as type,  s.id from  season s  left join content_primary_info cpi on s.primary_info_id = cpi.id where  s.id ='" + contentidarr.SeasonId + "' "

			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByGenreId != "" {
			contentdetailsquery := "select english_name as name,case WHEN (english_name != '') THEN '4' end as type, id  from genre where id ='" + contentidarr.GroupByGenreId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupBySubgenreId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '5' end as type, id from subgenre where id ='" + contentidarr.GroupBySubgenreId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByActorId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '6' end as type, id from actor where id ='" + contentidarr.GroupByActorId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByWriterId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '7' end as type, id from writer where id ='" + contentidarr.GroupByWriterId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByDirectorId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '8' end as type, id from director where id ='" + contentidarr.GroupByDirectorId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupBySingerId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '9' end as type, id from singer where id ='" + contentidarr.GroupBySingerId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByMusicComposerId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '10' end as type, id from music_composer where id ='" + contentidarr.GroupByMusicComposerId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupBySongWriterId != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '11' end as type, id from song_writer where id ='" + contentidarr.GroupBySongWriterId + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByProductionYear != "" && contentidarr.GroupByProductionYear > "0" {
			contentdetailsquery := "select distinct production_year as name,case WHEN ( production_year != 0 ) THEN '12' end as type, production_year as id from about_the_content_info where production_year ='" + contentidarr.GroupByProductionYear + "' "
			if contId := cdb.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByOriginalLanguageCode != "" {
			contentdetailsquery := "select english_name as name,case WHEN ( english_name !='') THEN '13' end as type,code as id from language where  code ilike '%" + contentidarr.GroupByOriginalLanguageCode + "%' "
			if contId := db.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

		if contentidarr.GroupByPageId != "" {
			contentdetailsquery := "select english_title as name,case WHEN ( english_title !='') THEN '14' end as type, id from page where id ='" + contentidarr.GroupByPageId + "' "
			if contId := db.Debug().Raw(contentdetailsquery).Scan(&playlistItem).Error; contId != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "description": "Server error.", "code": "error_server_error", "requestId": randstr.String(32)})
				return
			}
			playlistItems = append(playlistItems, playlistItem)
		}

	}
	playlist.PlaylistItems = playlistItems

	c.JSON(http.StatusOK, gin.H{"data": playlist})
	return
}

// GetSourceitemsBySourceFields - Get Source items BySource Fields
// GET api/playlist/sourceitems/
// @Summary Get Source items BySource Fields
// @Description Get Source items BySource Fields
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param searchText path string false "Search Text"
// @Param searchFilter path int false "search Filter"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router api/playlist/sourceitems/ [get]
func (hs *HandlerService) GetSourceitemsBySourceFields(c *gin.Context) {
	/*Authorization*/
	var errorresponse = common.ServerErrorResponse()
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	var limit, offset, searchFilter int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if limit == 0 {
		limit = 10
	}

	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}

	if c.Request.URL.Query()["searchFilter"] != nil {
		searchFilter, _ = strconv.ParseInt(c.Request.URL.Query()["searchFilter"][0], 10, 64)
	}
	var searchText string

	var sourceitemdata, totalCount []SourceItemDatas
	if c.Request.URL.Query()["searchText"] != nil {
		searchText = strings.ToLower(c.Request.URL.Query()["searchText"][0])
	}

	if searchText == "" || searchFilter == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Empty Parameters", "Status": http.StatusBadRequest})
		return
	}

	i := searchFilter
	switch i {

	case 1:
		// content (1-OTC,2-MTC,3-Season)
		rawquery := "select  distinct cpi.transliterated_title as name,  case WHEN (c.content_tier = 1) THEN '1'else '2' end as type,  c.id  from  content c  left join content_primary_info cpi on c.primary_info_id = cpi.id where c.deleted_by_user_id is null  and c.status = 1   and cpi.transliterated_title ilike '%" + searchText + "%' union select CONCAT(cpi.transliterated_title, ' Season ', '№',s.Number)as name,  case WHEN (cpi.transliterated_title != '') THEN '3' end as type,  s.id from  season s  left join content_primary_info cpi on s.primary_info_id = cpi.id where s.deleted_by_user_id is null and cpi.transliterated_title ilike '%" + searchText + "%'"

		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}

		break
	case 2:
		//Genre	-4
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '4' end as type, id from genre where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 3:
		//subGenre-5
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '5' end as type, id from subgenre where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 6:
		//Actor	-6
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '6' end as type, id from actor where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 7:
		//Writer-7
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '7' end as type, id from writer where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 8:
		//Director -8
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '8' end as type, id from director where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 9:
		//singer -9
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '9' end as type, id from singer where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 10:
		//MusicComposer -10
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '10' end as type, id from music_composer where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 11:
		//SongWriter -11
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '11' end as type, id from song_writer where (english_name ilike '%" + searchText + "%' or arabic_name ilike '%" + searchText + "%') "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 12:
		//ProductionYear -12
		rawquery := "select distinct production_year as name,case WHEN ( production_year != 0 ) THEN '12' end as type, production_year as id from about_the_content_info where production_year = " + searchText + "  "
		if data := cdb.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := cdb.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 13:
		//OriginalLanguage -13
		rawquery := "select english_name as name,case WHEN ( english_name !='') THEN '13' end as type,code as id from language where code ilike '%" + searchText + "%' "
		if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 14:
		//page -14
		rawquery := "select english_title as name,case WHEN ( english_title !='') THEN '14' end as type, id from page where (english_title ilike '%" + searchText + "%' or arabic_title ilike '%" + searchText + "%') "
		if data := db.Debug().Raw(rawquery).Limit(limit).Offset(offset).Scan(&sourceitemdata).Error; data != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if errCount := db.Debug().Raw(rawquery).Scan(&totalCount).Error; errCount != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break

	default:
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	pages := map[string]int{
		"size":   len(totalCount),
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": sourceitemdata})
}

// GetAll source item types -  fetches all type of source item types
// GET /api/playlist/sourceitems/types
// @Summary Show a list of all source item types
// @Description get list of all source item types
// @Tags source
// @Accept  json
// @Security Authorization
// @Produce  json
// @Success 200 {array} object c.JSON
// @Router /api/playlist/sourceitems/types [get]
func (hs *HandlerService) getSourceItem(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	var sourceItem []sourceItemTypes
	serverError := common.ServerErrorResponse()
	fields := "name,id"
	if err := db.Debug().Table("source_item_types").Select(fields).Find(&sourceItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sourceItem})
}

// GetPlaylistItemBasedOnPlaylistItemTypeAndItemTypeID -Get Playlist item based on playlist item type and item type ID
// GET api/playlist/dynamicgroupitems
// @Summary Get Playlist item based on playlist item type and item type ID
// @Description Get Playlist item based on playlist item type and item type ID
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param Id path string false "Id"
// @Param PlaylistItemType path int false "PlaylistItem Type"
// @Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Success 200 {array} object c.JSON
// @Router api/playlist/dynamicgroupitems [get]
func (hs *HandlerService) GetPlaylistItemBasedOnPlaylistItemTypeAndItemTypeID(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	var limit, offset, PlaylistItemType int64
	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if limit == 0 {
		limit = 2
	}
	if c.Request.URL.Query()["offset"] != nil {
		offset, _ = strconv.ParseInt(c.Request.URL.Query()["offset"][0], 10, 64)
	}
	if c.Request.URL.Query()["playlistItemType"] != nil {
		PlaylistItemType, _ = strconv.ParseInt(c.Request.URL.Query()["playlistItemType"][0], 10, 64)
	}

	var errorresponse = common.ServerErrorResponse()
	var Id string
	var playListItem []PlaylistItems
	var totalcount int
	if c.Request.URL.Query()["id"] != nil {
		Id = strings.ToLower(c.Request.URL.Query()["id"][0])
	}
	if Id == "" || PlaylistItemType == 0 {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	i := PlaylistItemType

	switch i {
	case 4:
		Query := "select cpi.transliterated_title as name,c.content_tier as type,c.id from content c join content_genre  cg on cg.content_id =c.id join content_primary_info cpi on cpi.id =c.primary_info_id where cg.genre_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join season_genre sg on sg.season_id =s.id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where sg.genre_id  ='" + Id + "' and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 5:
		Query := "select cpi.transliterated_title as name,c.content_tier as type,c.id from content c join content_genre cg on cg.content_id = c.id join content_subgenre cs on cs.content_genre_id =cg.id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.subgenre_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join season_genre sg2 on sg2.season_id =s.id join season_subgenre ss on ss.season_genre_id =sg2.id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where ss.subgenre_id  ='" + Id + "' and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 6:
		Query := "select cpi.transliterated_title as name,c.content_tier as type,c.id from content c join content_cast cc on cc.id=c.cast_id join actor a on a.id =cc.main_actor_id or a.id=cc.main_actress_id join content_primary_info cpi on cpi.id =c.primary_info_id where a.id='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join content_cast cc2 on cc2.id =s.cast_id join actor a2 on a2.id =cc2.main_actor_id or a2.id=cc2.main_actress_id join content_primary_info cpi2 on cpi2.id=s.primary_info_id where a2.id='" + Id + "' and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 7:
		Query := "select cpi.transliterated_title as name ,c.content_tier as type,c.id from content c join content_writer cw on cw.cast_id =c.cast_id join content_primary_info cpi on cpi.id =c.primary_info_id where cw.writer_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join content_writer cw2 on cw2.cast_id =s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where cw2.writer_id ='" + Id + "' and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 8:
		Query := "select cpi.transliterated_title as name ,c.content_tier as type,c.id from content c join content_director cd on cd.cast_id =c.cast_id join content_primary_info cpi on cpi.id=c.primary_info_id where cd.director_id ='" + Id + "'  and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join content_director cd1 on cd1.cast_id =s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where cd1.director_id ='" + Id + "' and s.deleted_by_user_id is null "
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 9:
		Query := "select cpi.transliterated_title as name ,c.content_tier as type,c.id from  content c join content_singer cs on cs.music_id =c.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.singer_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join content_singer cs1 on cs1.music_id =s.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where cs1.singer_id ='" + Id + "' and s.deleted_by_user_id is null "
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 10:
		Query := "select cpi.transliterated_title as name ,c.content_tier as type,c.id from  content_music_composer cmc join content c on c.music_id  =cmc.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cmc.music_composer_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from content_music_composer cmc1 join season s on s.music_id =cmc1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where cmc1.music_composer_id  ='" + Id + "' and s.deleted_by_user_id is null "
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 11:
		Query := "select cpi.transliterated_title as name ,c.content_tier as type,c.id from  content_song_writer csw join content c on c.music_id  =csw.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where csw.song_writer_id ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from content_song_writer csw1 join season s on s.music_id =csw1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where csw1.song_writer_id  ='" + Id + "' and s.deleted_by_user_id is null "
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 12:
		Query := "select cpi.transliterated_title as name,c.content_tier as type,c.id from content c join about_the_content_info  atci on atci.id =c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.production_year =" + Id + " and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join about_the_content_info atci2 on atci2.id=s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where atci2.production_year =" + Id + " and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 13:
		Query := "select cpi.transliterated_title as name,c.content_tier as type,c.id from content c join about_the_content_info  atci on atci.id =c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.original_language ='" + Id + "' and c.deleted_by_user_id is null union select CONCAT(cpi2.transliterated_title,' season ',',№',s.Number)as name,'3',s.id from season s join about_the_content_info atci2 on atci2.id=s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id where atci2.original_language ='" + Id + "' and s.deleted_by_user_id is null"
		rows := cdb.Debug().Raw(Query).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		if actorResult := cdb.Raw(Query).Limit(limit).Offset(offset).Scan(&playListItem).Error; actorResult != nil {
			fmt.Println(actorResult)
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		break
	case 14:
		var groupByPageId []GroupByPageId
		pageQuery := "select pic.content_id from page_playlist pp join playlist_item pi2 on pi2.playlist_id =pp.playlist_id join playlist_item_content pic on pic.playlist_item_id =pi2.id where pp.page_id ='" + Id + "'"
		if pageerror := db.Debug().Raw(pageQuery).Scan(&groupByPageId).Error; pageerror != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		var contentIds []string
		for _, val := range groupByPageId {
			contentIds = append(contentIds, val.ContentId)
		}
		if pageResult := cdb.Debug().Table("content as c").Select("cpi.transliterated_title as name,c.content_tier as type,c.id").
			Joins("join content_primary_info cpi on cpi.id=c.primary_info_id").
			Where("c.id in (?) and c.deleted_by_user_id is null", contentIds).Limit(limit).Offset(offset).Scan(&playListItem).Error; pageResult != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		rows := cdb.Debug().Table("content as c").Select("cpi.transliterated_title as name,c.content_tier as type,c.id").
			Joins("join content_primary_info cpi on cpi.id=c.primary_info_id").
			Where("c.id in (?) and c.deleted_by_user_id is null", contentIds).Find(&playListItem)
		totalcount = int(rows.RowsAffected)
		break
	default:
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}

	pages := map[string]int{
		"size":   totalcount,
		"offset": int(offset),
		"limit":  int(limit),
	}
	c.JSON(http.StatusOK, gin.H{"pagination": pages, "data": playListItem})
}

// CreateAndUpdatePlaylist -  Create And Update Playlist
// POST /playlists/{playlistId}
// @Summary Create And Update Playlist By Playlist Id
// @Description Create And Update Playlist By Playlist Id
// @Tags Playlist
// @Accept  json
// @Produce  json
// @Param playlistId path string false "playlist Id"
// @Param body body PlaylistInputs true "Raw JSON string"
// @Success 200 {array} object c.JSON
// @Router /playlists/{playlistId} [post]
func (hs *HandlerService) CreateAndUpdatePlaylist(c *gin.Context) {
	var errorresponse = common.ServerErrorResponse()
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}

	fcdb := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	ctx := context.Background()
	fctx := fcdb.BeginTx(ctx, nil)
	//cdtx := cdb.BeginTx(ctx, nil)
	playlistid := c.Param("playlistId")
	var errorFlag bool
	errorFlag = false

	var input PlaylistInputs
	/* Check post data is json formated or Not */
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
		return
	}
	/* Validations */
	var PlaylistItemError common.PlaylistItems
	/*for _, contents := range input.PlaylistItems {
		if contents.Type == 1 {
			var count int
			cdb.Table("content c").Joins("join content_variance cv on cv.content_id =c.id join playback_item pi2 on pi2.id =cv.playback_item_id join content_rights cr on cr.id =pi2.rights_id").Where("c.id =? and cr.digital_rights_start_date <=? and cr.digital_rights_end_date >=?", contents.Id, input.SchedulingStartDate, input.SchedulingEndDate).Count(&count)
			if count < 1 {
				errorFlag = true
				PlaylistItemError = common.PlaylistItems{Code: "error_playlist_not_all_playlist_items_match_playlist_scheduling", Description: "Please, make sure the scheduling settings match for the next items: " + contents.Name + "."}
			}
		}
		if contents.Type == 2 {
			var tcount int
			cdb.Table("content c").Joins("join season s on s.content_id =c.id join content_rights cr on cr.id =s.rights_id").Where("s.content_id=? and cr.digital_rights_start_date <= ? and cr.digital_rights_end_date >= ?", contents.Id, input.SchedulingStartDate, input.SchedulingEndDate).Count(&tcount)
			if tcount < 1 {
				errorFlag = true
				PlaylistItemError = common.PlaylistItems{Code: "error_playlist_not_all_playlist_items_match_playlist_scheduling", Description: "Please, make sure the scheduling settings match for the next items: " + contents.Name + "."}
			}
		}
		if contents.Type == 3 {
			var tcount int
			cdb.Table("season s").Joins("join content_rights cr on cr.id = s.rights_id").Where("s.id = ? and cr.digital_rights_start_date <= ? and cr.digital_rights_end_date >= ?", contents.Id, input.SchedulingStartDate, input.SchedulingEndDate).Count(&tcount)
			if tcount < 1 {
				errorFlag = true
				PlaylistItemError = common.PlaylistItems{Code: "error_playlist_not_all_playlist_items_match_playlist_scheduling", Description: "Please, make sure the scheduling settings match for the next items: " + contents.Name + "."}
			}
		}
	}*/
	var invalid common.Invalidsslider
	if PlaylistItemError.Code != "" {
		invalid.PlaylistItems = PlaylistItemError
	}
	var finalErrorResponse common.FinalErrorResponseslider
	finalErrorResponse = common.FinalErrorResponseslider{Error: "invalid_request", Description: "Validation failed.", Code: "error_validation_failed", RequestId: randstr.String(32), Invalid: invalid}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}
	/* End of Validations */
	var playlist Playlist
	/*update play_list_platform */
	var platform PlayListPlatform
	var platforms []interface{}
	/*update Page IDs Start */
	var pageplaylist PagePlaylist
	var pageplaylists []interface{}
	/*update play_list_country as regions */
	var region PlayListCountry
	var regions []interface{}
	/* for playlist items */
	var playlistItem PlaylistItem
	//var  playlistItems []interface{}

	if playlistid != "" {
		fmt.Println("----------------------------------------------------FOR  Update----------------------------------------------------------------------------")
		//playlistid := input.ID
		if err := fcdb.Debug().Where("id=? ", playlistid).Find(&playlist).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request", "Status": http.StatusBadRequest})
			return
		}

		/* Update playlist */
		playlist.EnglishTitle = input.EnglishTitle
		playlist.ArabicTitle = input.ArabicTitle
		playlist.ModifiedAt = time.Now()
		playlist.ModifiedAt = time.Now()
		playlist.SchedulingStartDate = input.SchedulingStartDate
		playlist.SchedulingEndDate = input.SchedulingEndDate
		if input.Playlisttype == "" {
			playlist.PlaylistType = "content"
		} else {
			playlist.PlaylistType = input.Playlisttype
		}
		if err := fctx.Debug().Table("playlist").Where("id=?", playlistid).Update(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}

		/*update play_list_platform START */
		if err := fctx.Debug().Table("play_list_platform").Where("play_list_id=?", playlistid).Delete(&platform).Error; err != nil {
			fctx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}

		if input.PublishingPlatforms != nil {
			for _, platformId := range input.PublishingPlatforms {
				platform.PlayListId = playlistid
				platform.PlayListKey = "0"
				platform.TargetPlatform = platformId
				platforms = append(platforms, platform)
			}

			err := gormbulk.BulkInsert(fctx, platforms, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/*update play_list_platform END */

		/*update Page IDs Start */
		if input.PagesIds != nil {
			// if err := fctx.Debug().Table("page_playlist").Where("playlist_id=?", playlistid).Delete(&pageplaylist).Error; err != nil {
			// 	fctx.Rollback()
			// 	c.JSON(http.StatusInternalServerError, errorresponse)
			// 	return
			// }
			for _, pageId := range input.PagesIds {
				var pageOrder PageOrder
				if err := fcdb.Debug().Table("page_playlist pp").Select("max(pp.order) as order").Where("page_id=?", pageId).Find(&pageOrder).Error; err != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
				var check int
				if err := fcdb.Debug().Raw("select count(*) from page_playlist where playlist_id = ? and page_id = ? ", playlistid, pageId).Count(&check).Error; err != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
				fmt.Println("PRINTITNG ", check)
				if check < 1 {
					pageplaylist.PlaylistId = playlistid
					pageplaylist.PageId = pageId
					pageplaylist.Order = pageOrder.Order + 1
					pageplaylists = append(pageplaylists, pageplaylist)
				}
			}
			err := gormbulk.BulkInsert(fctx, pageplaylists, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/*update Page EDs END */

		/*update play_list_country as regions END */

		if err := fctx.Debug().Table("play_list_country").Where("play_list_id=?", playlistid).Delete(&region).Error; err != nil {
			fctx.Rollback()
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if input.Regions != nil {
			for _, regionId := range input.Regions {
				region.PlayListId = playlistid
				region.CountryId = regionId
				regions = append(regions, region)
			}

			err := gormbulk.BulkInsert(fctx, regions, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/*update play_list_country as  regions END */

		/*update  For PlaylistItems Start */
		if input.PlaylistItems != nil {

			if err := fctx.Debug().Table("playlist_item").Where("playlist_id=?", playlistid).Delete(&playlistItem).Error; err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}

			for val, items := range input.PlaylistItems {
				var playlistItem PlaylistItem
				var playlistItemContent PlaylistItemContent
				playlistItem.PlaylistId = playlistid
				playlistItem.Order = val + 1
				itemtypes := items.Type
				switch itemtypes {
				case 1:
					playlistItem.OneTierContentId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.ContentId = &items.Id
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 2:
					playlistItem.MultiTierContentId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.ContentId = &items.Id
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 3:
					playlistItem.SeasonId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.SeasonId = &items.Id
					type SeasonDetails struct {
						ContentId *string
					}
					var seasonContentId SeasonDetails
					cdb.Debug().Raw("select content_id from season where id=?", &items.Id).Find(&seasonContentId)
					playlistItemContent.ContentId = seasonContentId.ContentId
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 4:
					playlistItem.GroupByGenreId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select distinct c.id as content_id,s.id as season_id from content c join season s on s.content_id =c.id join content_genre cg on cg.content_id =c.id join season_genre sg on sg.season_id = s.id join genre g on g.id = cg.genre_id or g.id=sg.genre_id where g.id=? and c.deleted_by_user_id is null", &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 5:
					playlistItem.GroupBySubgenreId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_genre cg on cg.content_id = c.id join content_subgenre cs on cs.content_genre_id =cg.id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.subgenre_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join season_genre sg2 on sg2.season_id =s.id join season_subgenre ss on ss.season_genre_id =sg2.id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where ss.subgenre_id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 6:
					playlistItem.GroupByActorId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					// cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_cast cc on cc.id = c.cast_id join actor a on a.id = cc.main_actor_id or a.id = cc.main_actress_id join content_primary_info cpi on cpi.id =c.primary_info_id where a.id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_cast cc2 on cc2.id = s.cast_id join actor a2 on a2.id = cc2.main_actor_id or a2.id = cc2.main_actress_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where a2.id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_cast cc on cc.id = c.cast_id join content_actor ca on ca.cast_id = cc.id join actor a on a.id = ca.actor_id join content_primary_info cpi on cpi.id =c.primary_info_id where a.id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_cast cc2 on cc2.id = s.cast_id join content_actor ca1 on ca1.cast_id = cc2.id join actor a2 on a2.id = ca1.actor_id  join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where a2.id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 7:
					playlistItem.GroupByWriterId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_writer cw on cw.cast_id = c.cast_id join content_primary_info cpi on cpi.id =c.primary_info_id where cw.writer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_writer cw2 on cw2.cast_id = s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cw2.writer_id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 8:
					playlistItem.GroupByDirectorId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_director cd on cd.cast_id = c.cast_id join content_primary_info cpi on cpi.id =c.primary_info_id where cd.director_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_director cd1 on cd1.cast_id = s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cd1.director_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 9:
					playlistItem.GroupBySingerId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_singer cs on cs.music_id = c.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.singer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_singer cs1 on cs1.music_id = s.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cs1.singer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 10:
					playlistItem.GroupByMusicComposerId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content_music_composer cmc join content c on c.music_id = cmc.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cmc.music_composer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from content_music_composer cmc1 join season s on s.music_id = cmc1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cmc1.music_composer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 11:
					playlistItem.GroupBySongWriterId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content_song_writer csw join content c on c.music_id = csw.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where csw.song_writer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from content_song_writer csw1 join season s on s.music_id = csw1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where csw1.song_writer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 12:
					years := items.Id
					yearsd, _ := strconv.Atoi(years)
					playlistItem.GroupByProductionYear = yearsd
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.production_year = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join about_the_content_info atci2 on atci2.id = s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where atci2.production_year = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 13:
					playlistItem.GroupByOriginalLanguageCode = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.original_language = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join about_the_content_info atci2 on atci2.id = s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where atci2.original_language = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 14:
					playlistItem.GroupByPageId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				}

			}
		}
		/*update  For PlaylistItems End */

		/*commit changes*/
		if err := fctx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		if len(input.PagesIds) != 0 {
			response := make(chan fragments.FragmentUpdate)
			go fragments.UpdatePlaylistFragment(playlistid, "", c, response, 0, 0, 0)
			outPut := <-response
			if outPut.Err != nil {
				errorresponse.Description = outPut.Err.Error()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		} else {
			fdb := c.MustGet("FDB").(*gorm.DB)
			var playlistFragment fragments.PlaylistFragment
			fdb.Debug().Where("playlist_id=?", playlistid).Delete(&playlistFragment)
		}
		// Status Message
		/* update dirty count in playlist_sync table */
		go common.PlaylistSynching(playlistid, c)
		finPlaylistId := map[string]string{"id": playlistid}
		c.JSON(http.StatusOK, gin.H{"data": finPlaylistId})
		return
	} else {
		fmt.Println("----------------------------------------------------FOR  Create----------------------------------------------------------------------------")

		type PlaylistKey struct {
			Key                   int `json:"key"`
			ThirdPartyPlaylistKey int `json:"third_party_playlist_key"`
		}
		var playlistKey PlaylistKey
		if err := fcdb.Debug().Table("playlist").Select("max(playlist_key)+1 as key,max(third_party_playlist_key) as third_party_playlist_key").Find(&playlistKey).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		// for creating third party playlist key
		playlist.ThirdPartyPlaylistKey = playlistKey.ThirdPartyPlaylistKey + 1
		// for removing sync below line is commented
		playlist.Id = input.PlaylistId // playlist id for creating old playlists with .net
		playlist.EnglishTitle = input.EnglishTitle
		playlist.ArabicTitle = input.ArabicTitle
		playlist.CreatedAt = time.Now()
		playlist.ModifiedAt = time.Now()
		// for removing sync below line is commented and replaced with playlist.PlaylistKey = playlistKey.Key
		//playlist.PlaylistKey = input.PlaylistKey
		playlist.PlaylistKey = playlistKey.Key
		playlist.SchedulingEndDate = input.SchedulingEndDate
		playlist.SchedulingStartDate = input.SchedulingStartDate
		if input.Playlisttype == "" {
			playlist.PlaylistType = "content"
		} else {
			playlist.PlaylistType = input.Playlisttype
		}
		if err := fctx.Debug().Create(&playlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, errorresponse)
			return
		}
		/* for playlist_platform Start */

		playlistid := playlist.Id

		if input.PublishingPlatforms != nil {
			for _, platformId := range input.PublishingPlatforms {
				platform.PlayListId = playlistid
				platform.PlayListKey = "0"
				platform.TargetPlatform = platformId
				platforms = append(platforms, platform)
			}

			err := gormbulk.BulkInsert(fctx, platforms, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/* for playlist_platform End */

		/* For Page Id Add Start */
		if len(input.PagesIds) != 0 {
			for _, pageId := range input.PagesIds {
				var pageOrder PageOrder
				if err := fcdb.Debug().Table("page_playlist pp").Select("max(pp.order) as order").Where("pp.page_id=?", pageId).Find(&pageOrder).Error; err != nil {
					c.JSON(http.StatusInternalServerError, errorresponse)
					return
				}
				pageplaylist.PlaylistId = playlistid
				pageplaylist.PageId = pageId
				pageplaylist.Order = pageOrder.Order + 1
				pageplaylists = append(pageplaylists, pageplaylist)
			}
			err := gormbulk.BulkInsert(fctx, pageplaylists, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/* For Page Id Add End */

		/* For play_list_country as regions Start */
		if input.Regions != nil {
			for _, regionId := range input.Regions {
				region.PlayListId = playlistid
				region.CountryId = regionId
				regions = append(regions, region)
			}
			err := gormbulk.BulkInsert(fctx, regions, common.BULK_INSERT_LIMIT)
			if err != nil {
				fctx.Rollback()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/* For play_list_country as regions End */

		/* create  For PlaylistItems Start */
		if input.PlaylistItems != nil {
			for val, items := range input.PlaylistItems {
				var playlistItem PlaylistItem
				var playlistItemContent PlaylistItemContent
				playlistItem.PlaylistId = playlistid
				playlistItem.Order = val + 1
				itemtypes := items.Type
				switch itemtypes {
				case 1:
					playlistItem.OneTierContentId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.ContentId = &items.Id
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 2:
					playlistItem.MultiTierContentId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.ContentId = &items.Id
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 3:
					playlistItem.SeasonId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					playlistItemContent.PlaylistItemId = playlistItem.Id
					playlistItemContent.SeasonId = &items.Id
					type SeasonDetails struct {
						ContentId *string
					}
					var seasonContentId SeasonDetails
					cdb.Debug().Raw("select content_id from season where id=?", &items.Id).Find(&seasonContentId)
					playlistItemContent.ContentId = seasonContentId.ContentId
					if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				case 4:
					playlistItem.GroupByGenreId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_genre cg on cg.content_id = c.id join content_primary_info cpi on cpi.id =c.primary_info_id where cg.genre_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join season_genre sg on sg.season_id = s.id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where sg.genre_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					// cdb.Debug().Raw("select distinct c.id as content_id,s.id as season_id from content c join season s on s.content_id =c.id join content_genre cg on cg.content_id =c.id join season_genre sg on sg.season_id = s.id join genre g on g.id = cg.genre_id or g.id=sg.genre_id where g.id=? and c.deleted_by_user_id is null", &items.Id).Find(&seasonContentId)
					fmt.Println("Playlist Genre Id: ", seasonContentId)
					for _, val := range seasonContentId {
						fmt.Println("Playlist Genre Id Values: ", val.ContentId, val.SeasonId)
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 5:
					playlistItem.GroupBySubgenreId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_genre cg on cg.content_id = c.id join content_subgenre cs on cs.content_genre_id =cg.id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.subgenre_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join season_genre sg2 on sg2.season_id =s.id join season_subgenre ss on ss.season_genre_id =sg2.id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where ss.subgenre_id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					fmt.Println("Playlist Sub-Genre Id: ", seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 6:
					playlistItem.GroupByActorId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_cast cc on cc.id = c.cast_id join content_actor ca on ca.cast_id = cc.id join actor a on a.id = ca.actor_id join content_primary_info cpi on cpi.id =c.primary_info_id where a.id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_cast cc2 on cc2.id = s.cast_id join content_actor ca1 on ca1.cast_id = cc2.id join actor a2 on a2.id = ca1.actor_id  join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where a2.id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					fmt.Println("Playlist Actor: ", seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 7:
					playlistItem.GroupByWriterId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_writer cw on cw.cast_id = c.cast_id join content_primary_info cpi on cpi.id =c.primary_info_id where cw.writer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_writer cw2 on cw2.cast_id = s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cw2.writer_id  = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 8:
					playlistItem.GroupByDirectorId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_director cd on cd.cast_id = c.cast_id join content_primary_info cpi on cpi.id =c.primary_info_id where cd.director_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_director cd1 on cd1.cast_id = s.cast_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cd1.director_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 9:
					playlistItem.GroupBySingerId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join content_singer cs on cs.music_id = c.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cs.singer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join content_singer cs1 on cs1.music_id = s.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cs1.singer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 10:
					playlistItem.GroupByMusicComposerId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content_music_composer cmc join content c on c.music_id = cmc.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where cmc.music_composer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from content_music_composer cmc1 join season s on s.music_id = cmc1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where cmc1.music_composer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 11:
					playlistItem.GroupBySongWriterId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content_song_writer csw join content c on c.music_id = csw.music_id join content_primary_info cpi on cpi.id =c.primary_info_id where csw.song_writer_id = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from content_song_writer csw1 join season s on s.music_id = csw1.music_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where csw1.song_writer_id = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 12:
					years := items.Id
					yearsd, _ := strconv.Atoi(years)
					playlistItem.GroupByProductionYear = yearsd
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.production_year = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join about_the_content_info atci2 on atci2.id = s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where atci2.production_year = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 13:
					playlistItem.GroupByOriginalLanguageCode = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					type SeasonDetails struct {
						ContentId *string
						SeasonId  *string
					}
					var seasonContentId []SeasonDetails
					cdb.Debug().Raw("select c.id as content_id,null as season_id from content c join about_the_content_info atci on atci.id = c.about_the_content_info_id join content_primary_info cpi on cpi.id =c.primary_info_id where atci.original_language = ? and c.deleted_by_user_id is null union select c1.id as content_id,s.id as season_id from season s join about_the_content_info atci2 on atci2.id = s.about_the_content_info_id join content_primary_info cpi2 on cpi2.id =s.primary_info_id join content c1 on c1.id =s.content_id where atci2.original_language = ? and s.deleted_by_user_id is null", &items.Id, &items.Id).Find(&seasonContentId)
					for _, val := range seasonContentId {
						playlistItemContent.PlaylistItemId = playlistItem.Id
						playlistItemContent.SeasonId = val.SeasonId
						playlistItemContent.ContentId = val.ContentId
						if err1 := fcdb.Debug().Create(&playlistItemContent).Error; err1 != nil {
							c.JSON(http.StatusInternalServerError, errorresponse)
							return
						}
					}
					break
				case 14:
					playlistItem.GroupByPageId = &items.Id
					if err := fcdb.Debug().Create(&playlistItem).Error; err != nil {
						c.JSON(http.StatusInternalServerError, errorresponse)
						return
					}
					break
				}
			}
		}
		/*create  For PlaylistItems End */

		/*commit changes*/
		if err := fctx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
			return
		}
		if len(input.PagesIds) != 0 {
			response := make(chan fragments.FragmentUpdate)
			go fragments.UpdatePlaylistFragment(playlistid, "", c, response, 0, 0, 0)
			outPut := <-response
			if outPut.Err != nil {
				errorresponse.Description = outPut.Err.Error()
				c.JSON(http.StatusInternalServerError, errorresponse)
				return
			}
		}
		/* Status Message */
		/* update dirty count in playlist_sync table */
		go common.PlaylistSynching(playlistid, c)
		finPlaylistId := map[string]string{"id": playlistid}
		c.JSON(http.StatusOK, gin.H{"data": finPlaylistId})
		return

	}

}

// GetAllRegionsBasedOnPlaylistId - Get All Regions Based On PlaylistId
// GET /api/playlist_id/:playlistId/region
// @Summary Get All Regions Based On PlaylistId
// @Description Get All Regions Based On PlaylistId
// @Tags Playlist
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param playlist_id path string true "Playlist Id"
// @Success 200 {array} object c.JSON
// @Router /api/playlist_id/{playlistId}/region [get]
func (hs *HandlerService) GetAllRegionsBasedOnPlaylistId(c *gin.Context) {
	/*Authorization*/
	if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Param("playlist_id")
	var errorresponse = common.ServerErrorResponse()
	var sourceItemTypes []SourceItemTypes
	/*Fetch playlists regions*/
	if resultError := db.Debug().Table("play_list_country plc").Select("c.english_name as name,c.id").
		Joins("left join country c on c.id =plc.country_id").
		Where("plc.play_list_id=?", id).Find(&sourceItemTypes).Error; resultError != nil {
		c.JSON(http.StatusInternalServerError, errorresponse)
	}
	c.JSON(http.StatusOK, gin.H{"data": sourceItemTypes})
}
