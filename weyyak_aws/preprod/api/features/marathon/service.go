package marathon

import (
	"feature/common"
	"strconv"
	"time"

	//	"encoding/json"

	"fmt"
	"net/http"

	// u "pdfGenerator"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// "encoding/csv"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	// sqrg := r.Group("/users")
	qrg := r.Group("/v1")
	qrg.Use(common.ValidateToken())
	qrg.POST("/register/marathon", hs.RegisterUserForMarathon)
	qrg.GET("/login/marathon", hs.loginMarathon)
	r.GET("/v1/getThirtyUsers", hs.GetTopThirty)
	r.POST("/v1/createTenUsers", hs.PostTopTen)
	r.GET("/v1/getTenUsers", hs.GetTopTen)
	r.POST("/v1/createFiveUsers", hs.PostTopFive)
	r.GET("/v1/getFiveUsers", hs.GetTopFive)
	// qrg.GET("/refresh", hs.RefreshButton)
}

// /v1/register/marathon
func (hs *HandlerService) RegisterUserForMarathon(c *gin.Context) {
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "Status": http.StatusUnauthorized})
		return
	}
	db := c.MustGet("DB").(*gorm.DB)
	udb := c.MustGet("UDB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	var request User
	var user Marathon
	userid := c.MustGet("userid")
	if err := udb.Where("id=?", userid).Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Record does not exist.Please provide valid User Id.", "Status": http.StatusBadRequest})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var errorFlag bool
	var Count int
	// var Rank []TopRank
	errorFlag = false
	var emailError ErrorCode
	var watchTime Marathon
	// var Rank []Marathon
	fmt.Println(emailError)
	var finalErrorResponse FinalErrorResponse
	if len(user.NickName) > 60 {
		if request.LanguageId == 2 {
			emailError = ErrorCode{"error_user_nick_name_invalid", "يجب أن يحتوي الاسم المستعار على احرف و ارقام و أن لا يتجاوز ال 60 حرفًا"}
		} else {
			emailError = ErrorCode{"error_user_nick_name_invalid", "Nick name must only contain characters with numbers and up to 60 characters"}
		}
	}
	if !common.RegName(request.NickName) && request.NickName != "" {
		errorFlag = true
		emailError = ErrorCode{"error_user_NickName_invalid", "NickName is invalid."}
	}
	if user.NickName == "" {
		errorFlag = true
		emailError := "Error_user_Nick_name_required"
		finalErrorResponse.Error = emailError
		finalErrorResponse.Description = emailError
	}
	if errorFlag {
		c.JSON(http.StatusBadRequest, finalErrorResponse)
		return
	}

	user.Id = request.Id
	cdb.Debug().Raw("select sum(last_watch_position) as watch_time from view_activity_history where user_id = ? and viewed_at >='2022-03-10 20:32:59.053'", userid).Find(&watchTime)
	db.Debug().Raw("select count(*) from marathon where nick_name=?", user.NickName).Count(&Count)
	if Count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Nick name already exists. Please choose different one. , unique nick_name or nickname not empty"})
		return
	}
	//email validate unique
	(user.WatchTime) = string(watchTime.WatchTime)
	user.UserRegisterAt = time.Now()
	if err := db.Debug().Create(&user).Error; err != nil {
		fmt.Println("errr", err)
		c.JSON(http.StatusBadRequest, gin.H{"Status": "duplicate key value. Please choose unique id"})
		return
	}
	// db.Table("user").Where("id = ?", userid).Update("last_activity_at", time.Now())
	udb.Table("user").Where("id = ?", userid).Update("nick_name", user.NickName)
	if request.LanguageId == 2 {
		c.JSON(http.StatusOK, gin.H{"message": "شكرًا على اشتراكك بماراثون وياك 5"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Thank you for your participation in Marathon 5."})
}

// /v1/login/marathon
func (hs *HandlerService) loginMarathon(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	userid := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "Status": http.StatusUnauthorized})
		return
	}
	var user Marathon
	if data := db.Debug().Raw(`SELECT * FROM "marathon" WHERE id=?`, userid).Scan(&user).Error; data != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": " Please register to marathon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Welcome to Marathon 5."})
}

// /v1/getThirtyUsers
func (hs *HandlerService) GetTopThirty(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)

	var erroresponse = common.ServerErrorResponse()
	var totalUser TopThirty
	var finalusers Final
	var Top30users, finalResult []TopThirty
	var userdetails details
	var usersDetails []Marathon
	var usersDetail []topTen
	var req []string

	if err := db.Find(&usersDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	if err := db.Find(&usersDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	db.Debug().Raw("truncate table top_thirty").Scan(&usersDetail)
	for _, user := range usersDetails {
		req = append(req, user.Id)
	}
	if err := cdb.Debug().Raw("select id,watch_time from (select user_id as id,sum(last_watch_position) as watch_time from view_activity_history vah  where user_id in (?) group by user_id order by watch_time desc limit 30) as foo where watch_time !=0", req).Find(&Top30users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	dt := time.Now()
	updated := fmt.Sprintf(dt.Format("01-02-2006 15:04:05"))
	finalusers.Lastupdatedtime = updated
	for key, val := range Top30users {
		totalUser.Rank = key + 1
		totalUser.Id = val.Id
		db.Debug().Raw("select nick_name from marathon where id=?", val.Id).Find(&userdetails)
		totalUser.NickName = userdetails.NickName
		timeval, _ := strconv.Atoi(val.WatchTime)
		hours := timeval / 3600
		minutes := (timeval % 3600) / 60
		seconds := timeval % 60
		timee := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		totalUser.WatchTime = timee
		finalResult = append(finalResult, totalUser)
		if inserterr := db.Debug().Table("top_thirty").Create(&totalUser).Error; inserterr != nil {
			fmt.Println("error", inserterr.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "500"})
			return
		}
		if updateWatchTime := db.Debug().Model(&usersDetails).Where("id=?", val.Id).Update("watch_time", totalUser.WatchTime).Error; updateWatchTime != nil {
			fmt.Println("error", updateWatchTime.Error)
			c.JSON(http.StatusBadRequest, gin.H{"message": "500"})
			return
		}

		const BULK_INSERT_LIMIT = 3000
	}
	finalusers.Top30 = finalResult

	c.JSON(http.StatusOK, &finalusers)
}

// /v1/createTenUsers
func (hs *HandlerService) PostTopTen(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var erroresponse = common.ServerErrorResponse()
	var usersDetails []topTen
	var top []topTen

	if err := db.Find(&usersDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	db.Debug().Raw("truncate table top_ten").Scan(&usersDetails)
	if err := c.ShouldBindJSON(&top); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}

	for key, userdetails := range top {
		fmt.Println("user details", userdetails)
		userdetails.Rank = key + 1

		fmt.Println(userdetails.Id, "ooooo")
		fmt.Println(userdetails.NickName, "ooo")

		if inserterr := db.Debug().Table("top_ten").Create(&userdetails).Error; inserterr != nil {
			fmt.Println("error", inserterr.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "500"})
			return
		}
		const BULK_INSERT_LIMIT = 3000
	}
	fmt.Println(top)
	c.JSON(http.StatusOK, &top)
}

// /v1/getTenUsers
func (hs *HandlerService) GetTopTen(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var finalResult []topTen
	db.Find(&finalResult)
	c.JSON(http.StatusOK, gin.H{"Data": finalResult})
}

// /v1/createFiveUsers
func (hs *HandlerService) PostTopFive(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var erroresponse = common.ServerErrorResponse()
	var usersDetails []TopFive
	var top []TopFive
	if err := db.Find(&usersDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, erroresponse)
		return
	}
	db.Debug().Raw("truncate table top_five").Scan(&usersDetails)
	if err := c.ShouldBindJSON(&top); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}
	fmt.Println("body", top)
	// db.Table("user").Where("id = ?", userdetails.Id).Update(&Top10, null)
	for key, userdetails := range top {
		fmt.Println("user details", userdetails)
		userdetails.Rank = key + 1
		if inserterr := db.Debug().Table("top_five").Create(&userdetails).Error; inserterr != nil {
			fmt.Println("error", inserterr.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "500"})
			return
		}
		const BULK_INSERT_LIMIT = 3000
	}
	fmt.Println(top)
	c.JSON(http.StatusOK, gin.H{"status": "successfully got top five"})
}

// /v1/getFiveUsers
func (hs *HandlerService) GetTopFive(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var finalResult []TopFive
	db.Find(&finalResult)
	c.JSON(http.StatusOK, gin.H{"Data": finalResult})

}

// /api/v1/refresh
// func (hs *HandlerService) RefreshButton(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	cdb := c.MustGet("CDB").(*gorm.DB)
// 	var erroresponse = common.ServerErrorResponse()
// 	var viewActivity []view_activity_history
// 	var usersDetails []Marathon
// 	var req []string
// 	var top_thirty []interface{}
// 	if err := db.Find(&usersDetails).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, erroresponse)
// 		return
// 	}
// 	for _, user := range usersDetails {
// 		req = append(req, user.Id)
// 	}
// 	db.Debug().Raw("truncate table top_thirty").Scan(&top_thirty)
// 	cdb.Debug().Raw("select id,watch_time from (select user_id as id,sum(last_watch_position) as watch_time from view_activity_history vah  where user_id in (?) group by user_id order by watch_time desc limit 30 ) as foo where watch_time !=0", req).Scan(&viewActivity)
// 	for _, v := range viewActivity {
// 		top_thirty = append(top_thirty, v)
// 	}
// 	err1 := gormbulk.BulkInsert(db.Debug().Table("top_thirty"), top_thirty, 3000)
// 	if err1 != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err1.Error(), "status": http.StatusBadRequest})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"Data": "users"})

// }
