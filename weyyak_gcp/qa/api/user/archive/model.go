package delete

import (
	"time"
)

type User struct {
	Id                   string    `json:"id" gorm:"primary_key"`
	Country              int       `json:"countryId"`
	FirstName            string    `json:"firstName"`
	IsBackOfficeUser     bool      `json:"isBackOfficeUser"`
	LanguageId           int       `json:"languageId"`
	LastName             string    `json:"lastName"`
	NewslettersEnabled   bool      `json:"newslettersEnabled"`
	PromotionsEnabled    bool      `json:"promotionsEnabled"`
	LastActivityAt       time.Time `json:"lastActivityAt"`
	RegistrationSource   int       `json:"registrationSource"`
	RegisteredAt         time.Time `json:"registeredAt"`
	Email                string    `json:"email" gorm:"unique"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PasswordHash         string    `json:"passwordHash" binding:"required"`
	SecurityStamp        string    `json:"securityStamp"`
	UserName             string    `json:"userName"`
	IsDeleted            bool      `json:"isDeleted"`
	PhoneNumber          string    `json:"phoneNumber" gorm:"unique"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	CallingCode          string    `json:"callingCode"`
	NationalNumber       string    `json:"nationalNumber"`
	CountryName          string    `json:"countryName"`
	SearchableText       string    `json:"searchableText"`
	Paycmsstatus         bool      `json:"payCMSStatus"`
	IsAdult              bool      `json:"isAdult"`
	PrivacyPolicy        bool      `json:"privacyPolicy"`
	IsRecommend          bool      `json:"isRecommend"`
	UserLead             string    `json:"userLead"`
	Performance          bool      `json:"performance"`
	GoogleAnalytics      bool      `json:"googleAnalytics"`
	Firebase             bool      `json:"firebase"`
	AppFlyer             bool      `json:"appFlyer"`
	Advertising          bool      `json:"advertising"`
	Aique                bool      `json:"aique"`
	GoogleAds            bool      `json:"googleAds"`
	FacebookAds          bool      `json:"facebookAds"`
	IsGdprAccepted       bool      `json:"isGdprAccepted"`
	SaltStored           string    `json:"-"`
	Version              int       `json:"-"`
	RoleId               string    `json:"-"`
	DeleteReasonId       int       `json:"deleteReasonId"`
	ReasonDetails        string    `json:"reasonDetails"`
	DeleteInitiatesAt    time.Time `json:"deleteIntiatesAt"`
	IsExport             bool      `json:"isExport"`
	SubscriptionEndDate  time.Time `json:"subscriptionEndDate"`
	OperatorType         string    `json:"operatorType"`
	Recurring            bool      `json:"recurring"`
}
type LanguageDetails struct {
	Code string
}
type UpdateUser struct {
	Id                  string    `json:"id" gorm:"primary_key"`
	IntiateDeleteAt     time.Time `json:"intiateDeleteAt"`
	DeleteReasonId      int       `json:"deleteReasonId"`
	ReasonDetails       string    `json:"reasonDetails"`
	IsExport            bool      `json:"isExport"`
	SubscriptionEndDate time.Time `json:"subscriptionEndDate"`
	OperatorType        string    `json:"operatorType"`
	Recurring           bool      `json:"recurring"`
}

type DeleteIntiate struct {
	DeleteReasonId           int                   `json:"DeleteReasonId"`
	ReasonDetails            string                `json:"reasonDetails"`
	DeleteInitiatesAt        time.Time             `json:"deleteIntiatesAt"`
	SubscriptionEndDate      time.Time             `json:"subscriptionEndDate"`
	SubscriptionPlansEndDate []SubscriptionDetails `json:"subscriptionPlansEndDate"`
	IsExport                 bool                  `json:"isExport"`
	Recurring                bool                  `json:"recurring"`
	OperatorType             string                `json:"operatorType"`
}

type SubscriptionDetails struct {
	Id                  int       `json:"id"`
	PlanName            string    `json:"planName"`
	SubscriptionEndDate time.Time `json:"subscriptionEndDate"`
	OperatorType        string    `json:"operatorType"`
}

type Subscription struct {
	OrderId          string `json:"order_id"`
	UserId           string `json:"user_id"`
	Identifier       string `json:"identifier"`
	SubscriptionPlan struct {
		Id                         string      `json:"id"`
		AssetType                  string      `json:"asset_type"`
		SubscriptionPlanType       string      `json:"subscription_plan_type"`
		Title                      string      `json:"title"`
		PromoCode                  interface{} `json:"promo_code"`
		CouponType                 string      `json:"coupon_type"`
		OriginalTitle              string      `json:"original_title"`
		System                     string      `json:"system"`
		Description                string      `json:"description"`
		BillingCycleType           string      `json:"billing_cycle_type"`
		BillingFrequency           int         `json:"billing_frequency"`
		Price                      string      `json:"price"`
		DiscountPrice              string      `json:"discount_price"`
		FinalPrice                 string      `json:"final_price"`
		Currency                   string      `json:"currency"`
		CountryCode                string      `json:"country_code"`
		CountryName                string      `json:"country_name"`
		NoOfFreeTrialDays          string      `json:"no_of_free_trial_days"`
		Start                      string      `json:"start"`
		End                        string      `json:"end"`
		OnlyAvailableWithPromotion bool        `json:"only_available_with_promotion"`
		Recurring                  bool        `json:"recurring"`
		PaymentProviders           string      `json:"payment_providers"`
		AssetTypes                 []string    `json:"asset_types"`
		NumberOfSupportedDevices   int         `json:"number_of_supported_devices"`
	} `json:"subscription_plan"`
	SubscriptionStart string `json:"subscription_start"`
	SubscriptionEnd   string `json:"subscription_end"`
	Recurring         bool   `json:"recurring"`
	RecurringEnabled  bool   `json:"recurring_enabled"`
	PaymentProvider   string `json:"payment_provider"`
	FreeTrialDays     string `json:"free_trial_days"`
	OrderDate         string `json:"order_date"`
}
type SubscriptionEn struct {
	SubscriptionEnd string `json:"subscription_end"`
}

type DeletedUser struct {
	Id                   string    `json:"id"`
	Country              int       `json:"countryId"`
	FirstName            string    `json:"firstName"`
	IsBackOfficeUser     bool      `json:"isBackOfficeUser"`
	LanguageId           int       `json:"languageId"`
	LastName             string    `json:"lastName"`
	NewslettersEnabled   bool      `json:"newslettersEnabled"`
	PromotionsEnabled    bool      `json:"promotionsEnabled"`
	LastActivityAt       time.Time `json:"lastActivityAt"`
	RegistrationSource   int       `json:"registrationSource"`
	RegisteredAt         time.Time `json:"registeredAt"`
	Email                string    `json:"email" gorm:"unique"`
	EmailConfirmed       bool      `json:"emailConfirmed"`
	PasswordHash         string    `json:"passwordHash" binding:"required"`
	SecurityStamp        string    `json:"securityStamp"`
	UserName             string    `json:"userName"`
	IsDeleted            bool      `json:"isDeleted"`
	PhoneNumber          string    `json:"phoneNumber" gorm:"unique"`
	PhoneNumberConfirmed bool      `json:"phoneNumberConfirmed"`
	CallingCode          string    `json:"callingCode"`
	NationalNumber       string    `json:"nationalNumber"`
	CountryName          string    `json:"countryName"`
	SearchableText       string    `json:"searchableText"`
	Paycmsstatus         bool      `json:"payCMSStatus"`
	IsAdult              bool      `json:"isAdult"`
	PrivacyPolicy        bool      `json:"privacyPolicy"`
	IsRecommend          bool      `json:"isRecommend"`
	UserLead             string    `json:"userLead"`
	Performance          bool      `json:"performance"`
	GoogleAnalytics      bool      `json:"googleAnalytics"`
	Firebase             bool      `json:"firebase"`
	AppFlyer             bool      `json:"appFlyer"`
	Advertising          bool      `json:"advertising"`
	Aique                bool      `json:"aique"`
	GoogleAds            bool      `json:"googleAds"`
	FacebookAds          bool      `json:"facebookAds"`
	IsGdprAccepted       bool      `json:"isGdprAccepted"`
	SaltStored           string    `json:"-"`
	Version              int       `json:"-"`
	RoleId               string    `json:"-"`
	DeleteReasonId       int       `json:"deleteReasonId"`
	ReasonDetails        string    `json:"reasonDetails"`
	IntiateDeleteAT      time.Time `json:"intiateDeleteAt"`
	IsExport             bool      `json:"isExport"`
	SubscriptionEndDate  time.Time `json:"subscriptionEndDate"`
	OperatorType         string    `json:"operatorType"`
	Recurring            bool      `json:"recurring"`
}
type FinalResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}
type Deletion struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DeletionReason struct {
	Reasons string `json:"reasons"`
}
