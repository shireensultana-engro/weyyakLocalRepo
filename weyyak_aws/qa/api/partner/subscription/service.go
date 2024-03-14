package subscription

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	srg := r.Group("api/")
	subscription := srg.Group("subscriptions")
	subscription.GET("/", hs.GetUserSubscriptionDetails)
	subscription.POST("/", hs.PostSubscription)
	plan := srg.Group("plans")
	plan.GET("/", hs.GetAllPlansDetails)
	plan.DELETE("/:id", hs.DeletePlanDetail)
	plan.POST("", hs.PostPlanDetails)
	plan.PUT("/:id", hs.UpdatePlan)

}

func (hs *HandlerService) UpdatePlan(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var plan PlanDetails
	if err := db.Debug().Table("plan_details").Where("id = ?", c.Param("id")).First(&plan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	// Validate input
	var updatePlan PlanDetails
	if err := c.ShouldBindJSON(&updatePlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("updateplan", updatePlan)
	db.Debug().Model(&plan).Updates(updatePlan)
	url := "https://zpapi.wyk.z5.com/index.php?c=PremiumPlans&m=updatePlansData"
	method := "POST"
	var price string = strconv.FormatFloat(plan.Price, 'E', -1, 32)
	var pppId string = strconv.FormatInt(int64(plan.PppId), 10)
	// pppId := strconv.Itoa(plan.PppId)
	nubmerOfFreeTrials := strconv.Itoa(plan.NumOfFreeTrials)
	payload := strings.NewReader(`{
		"pppid":"` + pppId + `",
		"title":"` + plan.PppTitle + `",
	  "price":"` + price + `",
	 "currency":"` + plan.Currency + `",
	  "free_trail_days":"` + nubmerOfFreeTrials + `"
  }`)
	fmt.Println("payloadddd", payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	fmt.Println("request,,,,,,", req)
	fmt.Println("errorrrrrr", err)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

	fmt.Println("plannnnnnnnn", plan)
	c.JSON(http.StatusOK, gin.H{"data": plan})
}

func (hs *HandlerService) GetAllPlansDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var plans []PlanDetails

	offset, _ := strconv.Atoi(c.Query("offset"))
	perPage, _ := strconv.Atoi(c.Query("limit"))

	if (offset == 0) && (perPage == 0) {
		offset = 0
		perPage = 25
	} else if offset == 0 {
		offset = 0
	}

	if perPage == 0 {
		perPage = 25
	}
	rows := db.Debug().Raw("select * from plan_details where active=true").Offset(offset).Limit(perPage).Find(&plans)
	if rows.RowsAffected == 0 {
		c.JSON(200, gin.H{
			"data": "No data available",
		})
		return
	} else if rows.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"data": plans,
		})
		return
	}

}

func (hs *HandlerService) DeletePlanDetail(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	plan_detail_id := c.Param("id")

	if plan_detail_id == "" {
		c.JSON(500, gin.H{
			"message": "ID cannot be empty",
		})
		return
	}

	var plan PlanDetails
	db.First(&plan, plan_detail_id)

	if plan.Id != 0 {
		db.Delete(&plan)
		c.JSON(200, gin.H{"success": "Deleted Successfully"})
	} else {
		c.JSON(404, gin.H{"error": "Not Found"})
	}

}

func (hs *HandlerService) PostPlanDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var plan PlanDetails
	plan.Active = true
	plan.LastModifiedDate = time.Now()
	body := c.Request.Body
	datas, _ := ioutil.ReadAll(body)
	json.Unmarshal(datas, &plan)
	err := db.Debug().Table("plan_details").Create(&plan)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": err.Error})
	} else {
		c.JSON(200, gin.H{"data": plan})
	}

}
func (hs *HandlerService) GetUserSubscriptionDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var UserSubscription []UserSubscription
	var result []UserSubscriptionDummy

	offset, _ := strconv.Atoi(c.Query("offset"))
	perPage, _ := strconv.Atoi(c.Query("limit"))

	if (offset == 0) && (perPage == 0) {
		offset = 0
		perPage = 25
	} else if offset == 0 {
		offset = 0
	}

	if perPage == 0 {
		perPage = 25
	}
	rows := db.Debug().Raw("select * from user_subscription").Offset(offset).Limit(perPage).Find(&UserSubscription)

	for _, v := range UserSubscription {
		var dummy UserSubscriptionDummy
		dummy.Id = v.Id
		dummy.UserId = v.UserId
		dummy.UserEmail = v.UserEmail
		dummy.UserFirstName = v.UserFirstName
		dummy.UserLastName = v.UserLastName
		dummy.PhoneNo = v.PhoneNo
		dummy.RegistrationDate = v.RegistrationDate
		dummy.SubscriptionDate = v.SubscriptionDate

		subscription_start_date := v.SubscriptionStartDate
		new_subscription_start_date := subscription_start_date.Format("2006-01-02")
		dummy.SubscriptionStartDate = new_subscription_start_date
		subscription_end_date := v.SubscriptionEndDate
		new_subscription_end_date := subscription_end_date.Format("2006-01-02")
		dummy.SubscriptionEndDate = new_subscription_end_date
		dummy.PaymentProvider = v.PaymentProvider
		dummy.SubscriptionStatus = v.SubscriptionStatus
		dummy.Plan = v.Plan
		result = append(result, dummy)
	}
	if rows.RowsAffected == 0 {
		c.JSON(200, gin.H{
			"data": "No data available",
		})
		return
	} else if rows.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"data": result,
		})
		return
	}

}

func (hs *HandlerService) PostSubscription(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var subscription UserSubscriptionDummy
	subscription.RegistrationDate = time.Now()
	subscription.SubscriptionDate = time.Now()
	body := c.Request.Body
	datas, _ := ioutil.ReadAll(body)
	json.Unmarshal(datas, &subscription)
	err := db.Debug().Table("user_subscription").Create(&subscription)

	if err.Error != nil {
		c.JSON(500, gin.H{"error": err.Error})
	} else {
		c.JSON(200, gin.H{"data": subscription})
	}

}
