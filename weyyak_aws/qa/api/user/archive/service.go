package delete

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"user/common"
	l "user/logger"

	// "github.com/robfig/cron"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"github.com/thanhpk/randstr"
	// "gorm.io/gorm"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.POST("/users/archive_initiation", hs.InitiateDelete)
	r.PUT("/users/paycms/inapppurchase", hs.CronUpdateSubscriptionDate)
	r.DELETE("/users/archive_complete", hs.CronDeleteUser)
}
func sendadminmail(userid string, langid int, firstname string, email string, phonenumber string, reasonid string, comments string) {
	templatePath := "DeleteFrontOfficeUserAdminBody.html"
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	langcode := common.LanguageIdToCode(langid)
	//curl call go get paymentdetails
	url := os.Getenv("SUBSCRIPTION_URL") + "/orderDetails?user_id=" + userid + "&include_all=false&language=" + langcode
	response := common.GetCurlCall(url)
	var subscription []Subscription
	json.Unmarshal(response, &subscription)
	content = strings.Replace(content, "@Model.UserFirstName", firstname, 1)
	content = strings.Replace(content, "UserEmail", email, 1)
	content = strings.Replace(content, "@Model.PhoneNumber", phonenumber, 1)
	for _, val := range subscription {
		content = strings.Replace(content, "@Model.Country", val.SubscriptionPlan.CountryName, 1)
		content = strings.Replace(content, "@Model.PlanName", val.SubscriptionPlan.Description, 1)
		content = strings.Replace(content, "@Model.SubscriptionEnd", val.SubscriptionEnd, 1)
		content = strings.Replace(content, "@Model.Price", val.SubscriptionPlan.FinalPrice+" "+val.SubscriptionPlan.Currency, 1)
	}
	if len(subscription) == 0 {
		content = strings.Replace(content, "@Model.Country", "", 1)
		content = strings.Replace(content, "@Model.PlanName", "", 1)
		content = strings.Replace(content, "@Model.SubscriptionEnd", "", 1)
		content = strings.Replace(content, "@Model.Price", "", 1)
	}
	content = strings.Replace(content, "@Model.ReasonId", reasonid, 1)
	content = strings.Replace(content, "@Model.ReasonComments", comments, 1)
	content = strings.Replace(content, "@Model.EmailHeadImageUrl", os.Getenv("EMAILIMAGEBASEURL")+os.Getenv("EMAILHEADIMAGEFILENAME"), 1)
	content = strings.Replace(content, "@Model.EmailContentImageUrl", os.Getenv("EMAILIMAGEBASEURL")+os.Getenv("EMAILCONTENTIMAGEFILENAME"), 1)
	message := template.HTML(content)
	error := common.SendMail(os.Getenv("ADMIN_MAIL"), string(message), "Subject:Weyyak Account Delete Initiate!")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
}

func sendusermail(languageid int, subscribeduser bool, username string, email string, subscriptionenddate string) {
	var templatePath string
	if languageid == 2 && subscribeduser {
		templatePath = "DeleteSubscribeInitiateFrontOfficeUserBodyAR.html"
	} else if languageid == 2 && !subscribeduser {
		templatePath = "DeleteInitiateFrontOfficeUserBodyAR.html"
	} else if languageid != 2 && subscribeduser {
		templatePath = "DeleteSubscribeInitiateFrontOfficeUserBodyEN.html"
	} else if languageid != 2 && !subscribeduser {
		templatePath = "DeleteInitiateFrontOfficeUserBodyEN.html"
	}

	fmt.Println("weyyak Account Delete Initiate", templatePath)
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	content = strings.Replace(content, "@Model.UserFirstName", username, 1)
	if subscribeduser {
		content = strings.Replace(content, "@Model.SubscriptionEndDate", subscriptionenddate, 1)
	}
	content = strings.Replace(content, "@Model.SubscriptionEndDate", subscriptionenddate, 1)
	content = strings.Replace(content, "@Model.Z5HomeUrl", os.Getenv("REDIRECTION_URL"), 1)
	content = strings.Replace(content, "@Model.EmailHeadImageUrl", os.Getenv("EMAILIMAGEBASEURL")+os.Getenv("EMAILHEADIMAGEFILENAME"), 1)
	content = strings.Replace(content, "@Model.EmailContentImageUrl", os.Getenv("EMAILIMAGEBASEURL")+os.Getenv("EMAILCONTENTIMAGEFILENAME"), 1)
	message := template.HTML(content)
	error := common.SendMail(email, string(message), "Subject:Weyyak Account Delete Initiate!")
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
}

// API for initiating delete process
func (hs *HandlerService) InitiateDelete(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var user User
	userid := c.MustGet("userid")
	if c.MustGet("AuthorizationRequired") == 1 {
		l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request.", "Status": http.StatusUnauthorized})
		return
	}

	if err := db.Debug().Where("id=?", userid).Find(&user).Error; err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": "Record does not exist.Please provide valid User Id.", "Status": http.StatusBadRequest})
		return
	}
	var input DeleteIntiate
	if err := c.ShouldBindJSON(&input); err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "Status": http.StatusBadRequest})
		return
	}

	var WeyyakUser string
	if user.FirstName != "" {
		WeyyakUser = user.FirstName
	} else {
		WeyyakUser = "Weyyak User"
	}

	input.DeleteInitiatesAt = time.Now()
	subscriptionPlansEndDate, err := json.Marshal(input.SubscriptionPlansEndDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	var maxTime time.Time

	for _, t := range input.SubscriptionPlansEndDate {
		if t.SubscriptionEndDate.After(maxTime) {
			maxTime = t.SubscriptionEndDate
		}
	}

	if err := db.Debug().Table("user").Where("id=?", userid).Updates(map[string]interface{}{
		"delete_initiates_at":   input.DeleteInitiatesAt,
		"delete_reason_id":      input.DeleteReasonId,
		"operator_type":         input.OperatorType,
		"reason_details":        input.ReasonDetails,
		"recurring":             input.Recurring,
		"subscription_end_date": input.SubscriptionEndDate,
		// "subscription_end_date":       maxTime, // BLU TV change
		"subscription_plans_end_date": subscriptionPlansEndDate,
	}).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}

	if input.SubscriptionEndDate.IsZero() {
		// user.SubscriptionEndDate.IsZero() gives true for unsubscribed user
		// for unsubscribed users
		if user.LanguageId == 2 {
			go sendusermail(2, false, WeyyakUser, user.Email, input.SubscriptionEndDate.Format("2006-01-02 15:04:05"))
		} else {
			go sendusermail(1, false, WeyyakUser, user.Email, input.SubscriptionEndDate.Format("2006-01-02 15:04:05"))
		}
	} else {
		// for subscribed users
		if user.LanguageId == 2 {
			go sendusermail(2, true, WeyyakUser, user.Email, input.SubscriptionEndDate.Format("2006-01-02 15:04:05"))
		} else {
			go sendusermail(1, true, WeyyakUser, user.Email, input.SubscriptionEndDate.Format("2006-01-02 15:04:05"))
		}
	}

	var deleteReason DeletionReason

	if err := db.Table("delete_reasons").Where("id = ?", input.DeleteReasonId).Find(&deleteReason).Error; err != nil {
		// l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		// return
		fmt.Println("Delete User Reason Not Found")
	}

	go sendadminmail(user.Id, user.LanguageId, WeyyakUser, user.Email, user.PhoneNumber, deleteReason.Reasons, input.ReasonDetails)
	l.JSON(c, http.StatusOK, gin.H{"message": "delete process initiated"})
}

// Cron for updating subscription end date
func (hs *HandlerService) CronUpdateSubscriptionDate(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var user []User
	// var subscriptionend, recurring, operatortype string
	// var userids []string
	if err := db.Raw(`select * from public.user where subscription_end_date is not null`).Scan(&user).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}

	for _, v := range user {
		langcode := common.LanguageIdToCode(v.LanguageId)
		url1 := os.Getenv("SUBSCRIPTION_URL") + "/orderDetails?user_id=" + v.Id + "&include_all=false&language=" + langcode
		response := common.GetCurlCall(url1)
		var subscription []Subscription
		json.Unmarshal(response, &subscription)

		var maxTime time.Time
		var OperatorType string
		var subscriptionDetails []SubscriptionDetails
		var recur bool

		for _, val := range subscription {
			layout := "2006-01-02 15:04:05"
			SubscriptionEnd, err := time.Parse(layout, val.SubscriptionEnd)
			if err != nil {
				fmt.Println(err)
				return
			}

			if SubscriptionEnd.After(maxTime) {
				maxTime = SubscriptionEnd
				OperatorType = val.PaymentProvider

				if val.Recurring {
					recur = true
				} else {
					recur = false
				}
			}

			SubscriptionPlanId, err := strconv.Atoi(val.SubscriptionPlan.Id)
			if err != nil {
				fmt.Println("Wrong Plan id")
				return
			}

			subscriptionDetails = append(subscriptionDetails, SubscriptionDetails{
				Id:                  SubscriptionPlanId,
				PlanName:            val.SubscriptionPlan.Title,
				SubscriptionEndDate: SubscriptionEnd,
				OperatorType:        val.PaymentProvider,
			})

			// 	operatortype = operatortype + " when id = '" + val.UserId + "' then '" + val.PaymentProvider + "'"
			// 	recurring = recurring + " when id = '" + val.UserId + "' then " + recur
			// 	subscriptionend = subscriptionend + " when id = '" + val.UserId + "' then TO_TIMESTAMP('" + val.SubscriptionEnd + "','YYYY-MM-DD HH24:MI:SS')"
			// 	userids = append(userids, val.UserId)
		}

		subscriptionPlansEndDate, err := json.Marshal(subscriptionDetails)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := db.Debug().Table("user").Where("id=?", v.Id).Updates(map[string]interface{}{
			"operator_type":               OperatorType,
			"subscription_end_date":       maxTime,
			"subscription_plans_end_date": subscriptionPlansEndDate,
			"recurring":                   recur,
		}).Error; err != nil {
			l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
			return
		}

	}
	// query := "update public.user  set operator_type = (case " + operatortype + " end),recurring = (case " + recurring + " end),subscription_end_date  = (case " + subscriptionend + " end) "
	// if err := db.Debug().Exec(query+"where id in (?)", userids).Error; err != nil {
	// 	l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
	// 	return
	// }

	l.JSON(c, http.StatusOK, gin.H{"message": "updated subscription date"})
}

// Cron for permanently deleting the user Account
func (hs *HandlerService) CronDeleteUser(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	cdb := c.MustGet("CDB").(*gorm.DB)
	var user []User
	var exceededusers []User
	var intitatedusers []string
	if err := db.Raw(`select * from public.user where delete_initiates_at is not null`).Scan(&user).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	for _, val := range user {
		intitatedusers = append(intitatedusers, val.Id)
	}
	if err := db.Debug().Raw("select * from public.user where ((subscription_end_date + interval '30 DAY'<Now() and recurring = false) or (delete_initiates_at + interval '30 DAY'<Now() and subscription_end_date is null)) and id in (?)", intitatedusers).Scan(&exceededusers).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	var deleteduserids []string
	var deletedUser []interface{}
	for _, v := range exceededusers {
		deleteduserids = append(deleteduserids, v.Id)
		var WeyyakUser string
		if v.FirstName != "" {
			WeyyakUser = v.FirstName
		} else {
			WeyyakUser = "Weyyak User"
		}
		fmt.Println(v.Email)
		go sendmailfordeleteduser(v.LanguageId, WeyyakUser, v.Email)
		go UpdatePaymentsDB(v.Id)
		deletedUser = append(deletedUser, v)
	}
	// deleting in user table and adding in deleted_user table
	err := gormbulk.BulkInsert(db.Debug().Table("deleted_user"), deletedUser, 3000)
	if err != nil {
		l.JSON(c, http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	if err := cdb.Debug().Exec("DELETE FROM public.playlisted_content  WHERE user_id in (?)", deleteduserids).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	if err := cdb.Debug().Exec("DELETE FROM view_activity WHERE user_id in (?)", deleteduserids).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}
	if err := db.Debug().Exec("DELETE FROM public.user WHERE id in (?)", deleteduserids).Error; err != nil {
		l.JSON(c, http.StatusInternalServerError, FinalResponse{"server_error", "Server error", "error_server_error", randstr.String(32)})
		return
	}

	// go sendmailfordeleteduser(1, "WeyyakUser", "usermanagement001@mailnesia.com")
	l.JSON(c, http.StatusOK, gin.H{"View": "deleted"})
}

func sendmailfordeleteduser(languageid int, username string, email string) {
	var templatePath string
	var subject string
	if languageid == 2 {
		templatePath = "DeleteCompleteFrontOfficeUserBodyAR.html"
		subject = "Subject:تم حذف حساب وياك بنجاح"
	} else {
		templatePath = "DeleteCompleteFrontOfficeUserBodyEN.html"
		subject = "Subject:Weyyak Account Deleted Successfully"
	}
	con, _ := ioutil.ReadFile(templatePath)
	content := string(con)
	content = strings.Replace(content, "@Model.UserFirstName", username, 1)
	message := template.HTML(content)
	error := common.SendMail(email, string(message), subject)
	if error != nil {
		fmt.Println("Email has not sent- ", error)
	}
}

func UpdatePaymentsDB(userid string) {
	url := os.Getenv("USER_DELETE_URL") + userid
	common.GetCurlCall(url)
}
