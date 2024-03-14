package marathon

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type LanguageId struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type NickNamecheck struct {
	NickName string `json:"nickName"`
}
type ErrorCode struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PasswordError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type ConfirmationToken struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type DateTimeToken struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type FinalErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
type Marathon struct {
	Id             string    `json:"id"`
	FirstName      string    `json:"firstName"`
	LanguageId     int       `json:"languageId"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email" gorm:"unique"`
	PasswordHash   string    `json:"passwordHash" `
	NickName       string    `json:"nickName"`
	PhoneNumber    string    `json:"phoneNumber" gorm:"unique"`
	CountryName    string    `json:"countryName"`
	Country        int       `json:"country"`
	Alpha2Code     string    `json:"alpha2Code"`
	Rank           int       `json:"rank"`
	WatchTime      string    `json:"watchTime"`
	UserRegisterAt time.Time `json:"user_register_at"`
	// LastActivityAt time.Time `json:"last_activity_at"`
}
type TopRank struct {
	Rank int `json:"rank"`
}
type view_activity_history struct {
	WatchTime int    `json:"watchTime"`
	Id        string `json:"id"`
}
type FinalResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}
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
	NickName             string    `json:"nickName"`
}

// type RequestRegisterUserUsingEmail struct {
// 	Password      string `json:"password"`
// 	LanguageId    int    `json:"languageId"`
// 	PrivacyPolicy bool   `json:"PrivacyPolicy"`
// 	IsAdult       bool   `json:"isAdult"`
// 	IsRecommend   bool   `json:"IsRecommend"`
// 	Email         string `json:"email"`
// 	CountryName   string `json:"countryName"`
// 	Alpha2code    string `json:"Alpha2code"`
// 	Source        string `json:"source"` /* sync usecase */
// 	UserId        string `json:"userId"` /* sync usecase */

// }
// type Emailcheck struct {
// 	Email string `json:"email"`
// }
// type TopThirty struct {
// 	Rank      int    `json:"rank"`
// 	Id        string `json:"id"`
// 	WatchTime string `json:"watchTime"`
// }
type TopThirty struct {
	Rank      int    `json:"rank"`
	Id        string `json:"id"`
	WatchTime string `json:"watchTime"`
	NickName  string `json:"nickName"`
}
type Final struct {
	Lastupdatedtime string `json:"lastUpdatedTime"`
	Top30           []TopThirty
}

type Finals struct {
	Lastupdatedtime string `json:"lastUpdatedTime"`
	TopTen          []topTen
}
type finals struct {
	Lastupdatedtime string `json:"lastUpdatedTime"`
	TopTen          []TopFive
}

type details struct {
	NickName string `json:"nickName"`
}
type topTen struct {
	Rank      int    `json:"rank"`
	Id        string `json:"id"`
	WatchTime string `json:"watchTime"`
	NickName  string `json:"nickName"`
}

type TopFive struct {
	Rank      int    `json:"rank"`
	Id        string `json:"id"`
	WatchTime string `json:"watchTime"`
	NickName  string `json:"nickName"`
}
type FinalFive struct {
	Lastupdatedtime string `json:"lastUpdatedTime"`
	TopTen          []topTen
}
