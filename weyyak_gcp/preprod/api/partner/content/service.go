package content

import (
	"encoding/json"
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
	srg.POST("/oauth2/token", hs.Login)
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

// GetMenuDetails - Get all menu list details
// GET /v1/get_menu
// @Description Get All menu list details by platform ID
// @Tags Menu
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param device query string true "Device Name"
// @Success 200  object MenuDetails
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/get_menu [get]
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

	rows := db.Debug().Table("page p").Select("p.*").
		Joins("inner join page_target_platform ptp on ptp.page_id=p.id").
		Where("p.is_disabled=false and p.deleted_by_user_id is null and p.third_party_page_key is not null and p.page_type != 16 and ptp.target_platform=?", deviceId).Group("p.id,ptp.page_order_number").Order("ptp.page_order_number asc").Find(&pageDetails).RowsAffected

	if err := db.Table("page p").Select("p.*").
		Joins("inner join page_target_platform ptp on ptp.page_id=p.id").
		Where("p.is_disabled=false and p.deleted_by_user_id is null and p.third_party_page_key is not null and p.page_type != 16 and ptp.target_platform=?", deviceId).Group("p.id,ptp.page_order_number").Order("ptp.page_order_number asc").Limit(limit).Offset(offset).Find(&pageDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, serverError)
		return
	}
	type PageIds struct {
		Id string `json:"id"`
	}
	var pageids []PageIds
	var ids []string
	if err := db.Table("page p").Select("p.id").
		Joins("inner join page_slider ps on ps.page_id=p.id inner join slider s on s.id = ps.slider_id").
		Where("s.deleted_by_user_id  is null and s.is_disabled =false and p.is_disabled=false and p.deleted_by_user_id is null and s.scheduling_start_date <=NOW() and s.scheduling_end_date >=NOW()").Find(&pageids).Error; err != nil {
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
		// if details.EnglishTitle != "" && details.EnglishTitle != "My Playlist" {
		menuResponse = append(menuResponse, menu)
		// }
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

type VideoDurationInfo struct {
	Duration     int           `json:"duration"`
	Thumbnails   []interface{} `json:"thumbnails"`
	UrlTrickplay string        `json:"url_trickplay"`
	UrlVideo     string        `json:"url_video"`
}

// GetVideoDurationDetails -  Get all menu list details
// GET /v1/get_info/:videoId
// @Description Get all menu list details
// @Tags Video
// @Accept  json
// @Security Authorization
// @Produce  json
// @Param videoId path string true "video Id."
// @Param country query string false "Country code of the user."
// @Success 200  object VideoDurationInfo
// @Failure 404 "The object was not found."
// @Failure 500 object ErrorResponse "Internal server error."
// @Router /v1/get_info/{videoId} [get]
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
}
