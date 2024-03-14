package subscription

import (
	"time"
)

type PlanDetails struct {
	Id               int       `gorm:"primaryKey autoIncrement" json:"planId"`
	PlanName         string    `json:"planName"`
	VodType          int       `json:"vodType"`
	PlanType         string    `json:"planType"`
	PlanValidityUnit string    `json:"planValidityUnit"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
	LastModifiedBy   string    `json:"lastModifiedBy"`
	PppId            int       `gorm:"pppId"`
	PppTitle         string    `gorm:"pppTitle"`
	Price            float64   `gorm:"price"`
	Currency         string    `gorm:"currency"`
	SubscriptionFlag bool      `gorm:"subscriptionFlag"`
	NumOfFreeTrials  int       `gorm:"numOfFreeTrials"`
	Status           string    `json:"status"`
	Active           bool      `json:"active,omitempty"`
}

type UserSubscription struct {
	Id                    int       `gorm:"primaryKey autoIncrement" json:"id"`
	UserId                int       `json:"userId"`
	UserEmail             string    `json:"userEmail"`
	UserFirstName         string    `json:"userFirstName"`
	UserLastName          string    `json:"userLastName"`
	PhoneNo               string    `json:"phoneNo"`
	RegistrationDate      time.Time `json:"registrationDate"`
	SubscriptionDate      time.Time `json:"subscriptionDate"`
	SubscriptionStartDate time.Time `json:"subscriptionStartDate"`
	SubscriptionEndDate   time.Time `json:"subscriptionEndDate"`
	PaymentProvider       string    `json:"paymentProvider"`
	SubscriptionStatus    string    `json:"subscriptionStatus"`
	Plan                  string    `json:"plan"`
}

type UserSubscriptionDummy struct {
	Id                    int       `gorm:"primaryKey autoIncrement" json:"id"`
	UserId                int       `json:"userId"`
	UserEmail             string    `json:"userEmail"`
	UserFirstName         string    `json:"userFirstName"`
	UserLastName          string    `json:"userLastName"`
	PhoneNo               string    `json:"phoneNo"`
	RegistrationDate      time.Time `json:"registrationDate"`
	SubscriptionDate      time.Time `json:"subscriptionDate"`
	SubscriptionStartDate string    `json:"subscriptionStartDate"`
	SubscriptionEndDate   string    `json:"subscriptionEndDate"`
	PaymentProvider       string    `json:"paymentProvider"`
	SubscriptionStatus    string    `json:"subscriptionStatus"`
	Plan                  string    `json:"plan"`
}
