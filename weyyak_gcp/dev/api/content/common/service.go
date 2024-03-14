package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	// "strconv"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

const BulkInsertLimit int = 3000

var ctx = context.Background()

// Post curl call
func PostCurlCall(method string, url string, data interface{}) ([]byte, error) {
	URL := url
	payloadBytes, _ := json.Marshal(data)
	body := bytes.NewReader(payloadBytes)
	req, _ := http.NewRequest(method, URL, body)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	response, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return response, error
	}
	return response, nil

}

// Get curl call
func GetCurlCall(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	// defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

// LanguageOriginTypes is...
func LanguageOriginTypes(originType string) int {
	OriginTypesArray := map[string]int{"Original": 1, "Dubbed": 2, "Subtitled": 3}
	return OriginTypesArray[originType]
}

// uploadFileToS3 is...
func uploadFileToS3(fileName string) error {
	// fileName := strings.Replace(fileNam, os.Getenv("PDF_URL"), "", -1)
	s := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// open the file for use
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// get the file size and read
	// the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, s3err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("S3_BUCKET")),
		Key:                aws.String(fileName),
		ACL:                aws.String("public-read"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
		StorageClass:       aws.String("Standard"),
	}) //ServerSideEncryption: aws.String("AES256"),
	return s3err
}

// Countrys is...
func OriginalLanguage(lang string) string {
	languageArray := map[string]string{"en": "English", "ar": "Arabic", "fr": "French", "hi": "Hindi", "ir": "Iranian", "ko": "Korean", "pt": "Portuguese", "es": "Spanish", "tr": "Turkish", "ur": "Urdu"}
	return languageArray[lang]
}

// Countrys is...
func Countrys(country string) int32 {
	countryArray := map[string]int{"AF": 4, "AQ": 10, "DZ": 12, "AS": 16, "AD": 20, "AO": 24, "AG": 28, "AZ": 31, "AR": 32, "AU": 36, "AT": 40, "BS": 44, "BH": 48, "BD": 50, "AM": 51, "BB": 52, "BE": 56, "BM": 60, "BT": 64, "BO": 68, "BA": 70, "BW": 72, "BV": 74, "BR": 76, "BZ": 84, "IO": 86, "SB": 90, "VG": 92, "BN": 96, "BG": 100, "MM": 104, "BI": 108, "BY": 112, "KH": 116, "CM": 120, "CA": 124, "CV": 132, "KY": 136, "CF": 140, "LK": 144, "TD": 148, "CL": 152, "CN": 156, "TW": 158, "CX": 162, "CC": 166, "CO": 170, "KM": 174, "YT": 175, "CG": 178, "CD": 180, "CK": 184, "CR": 188, "HR": 191, "CU": 192, "CY": 196, "CZ": 203, "BJ": 204, "DK": 208, "DM": 212, "DO": 214, "EC": 218, "SV": 222, "GQ": 226, "ET": 231, "ER": 232, "EE": 233, "FO": 234, "FK": 238, "GS": 239, "FJ": 242, "FI": 246, "AX": 248, "FR": 250, "GF": 254, "PF": 258, "TF": 260, "DJ": 262, "GA": 266, "GE": 268, "GM": 270, "PS": 275, "DE": 276, "GH": 288, "GI": 292, "KI": 296, "GR": 300, "GL": 304, "GD": 308, "GP": 312, "GU": 316, "GT": 320, "GN": 324, "GY": 328, "HT": 332, "HM": 334, "VA": 336, "HN": 340, "HK": 344, "HU": 348, "IS": 352, "IN": 356, "ID": 360, "IR": 364, "IQ": 368, "IE": 372, "IL": 376, "IT": 380, "CI": 384, "JM": 388, "JP": 392, "KZ": 398, "JO": 400, "KE": 404, "KP": 408, "KR": 410, "KW": 414, "KG": 417, "LA": 418, "LB": 422, "LS": 426, "LV": 428, "LR": 430, "LY": 434, "LI": 438, "LT": 440, "LU": 442, "MO": 446, "MG": 450, "MW": 454, "MY": 458, "MV": 462, "ML": 466, "MT": 470, "MQ": 474, "MR": 478, "MU": 480, "MX": 484, "MC": 492, "MN": 496, "MD": 498, "ME": 499, "MS": 500, "MA": 504, "MZ": 508, "OM": 512, "NA": 516, "NR": 520, "NP": 524, "NL": 528, "CW": 531, "AW": 533, "SX": 534, "BQ": 535, "NC": 540, "VU": 548, "NZ": 554, "NI": 558, "NE": 562, "NG": 566, "NU": 570, "NF": 574, "NO": 578, "MP": 580, "UM": 581, "FM": 583, "MH": 584, "PW": 585, "PK": 586, "PA": 591, "PG": 598, "PY": 600, "PE": 604, "PH": 608, "PN": 612, "PL": 616, "PT": 620, "GW": 624, "TL": 626, "PR": 630, "QA": 634, "RE": 638, "RO": 642, "RU": 643, "RW": 646, "BL": 652, "SH": 654, "KN": 659, "AI": 660, "LC": 662, "PM": 666, "VC": 670, "SM": 674, "ST": 678, "SA": 682, "SN": 686, "RS": 688, "SC": 690, "SL": 694, "SG": 702, "SK": 703, "VN": 704, "SI": 705, "SO": 706, "ZA": 710, "ZW": 716, "ES": 724, "SS": 728, "SD": 729, "EH": 732, "SR": 740, "SJ": 744, "SZ": 748, "SE": 752, "CH": 756, "SY": 760, "TJ": 762, "TH": 764, "TG": 768, "TK": 772, "TO": 776, "TT": 780, "AE": 784, "TN": 788, "TR": 792, "TM": 795, "TC": 796, "TV": 798, "UG": 800, "UA": 804, "MK": 807, "EG": 818, "GB": 826, "GG": 831, "JE": 832, "IM": 833, "TZ": 834, "US": 840, "VI": 850, "BF": 854, "UY": 858, "UZ": 860, "VE": 862, "WF": 876, "WS": 882, "YE": 887, "ZM": 894}
	return int32(countryArray[country])
}

// CountryNames is...
func CountryNames(country int) string {
	countryArray := map[int]string{4: "AF", 10: "AQ", 12: "DZ", 16: "AS", 20: "AD", 24: "AO", 28: "AG", 31: "AZ", 32: "AR", 36: "AU", 40: "AT", 44: "BS", 48: "BH", 50: "BD", 51: "AM", 52: "BB", 56: "BE", 60: "BM", 64: "BT", 68: "BO", 70: "BA", 72: "BW", 74: "BV", 76: "BR", 84: "BZ", 86: "IO", 90: "SB", 92: "VG", 96: "BN", 100: "BG", 104: "MM", 108: "BI", 112: "BY", 116: "KH", 120: "CM", 124: "CA", 132: "CV", 136: "KY", 140: "CF", 144: "LK", 148: "TD", 152: "CL", 156: "CN", 158: "TW", 162: "CX", 166: "CC", 170: "CO", 174: "KM", 175: "YT", 178: "CG", 180: "CD", 184: "CK", 188: "CR", 191: "HR", 192: "CU", 196: "CY", 203: "CZ", 204: "BJ", 208: "DK", 212: "DM", 214: "DO", 218: "EC", 222: "SV", 226: "GQ", 231: "ET", 232: "ER", 233: "EE", 234: "FO", 238: "FK", 239: "GS", 242: "FJ", 246: "FI", 248: "AX", 250: "FR", 254: "GF", 258: "PF", 260: "TF", 262: "DJ", 266: "GA", 268: "GE", 270: "GM", 275: "PS", 276: "DE", 288: "GH", 292: "GI", 296: "KI", 300: "GR", 304: "GL", 308: "GD", 312: "GP", 316: "GU", 320: "GT", 324: "GN", 328: "GY", 332: "HT", 334: "HM", 336: "VA", 340: "HN", 344: "HK", 348: "HU", 352: "IS", 356: "IN", 360: "ID", 364: "IR", 368: "IQ", 372: "IE", 376: "IL", 380: "IT", 384: "CI", 388: "JM", 392: "JP", 398: "KZ", 400: "JO", 404: "KE", 408: "KP", 410: "KR", 414: "KW", 417: "KG", 418: "LA", 422: "LB", 426: "LS", 428: "LV", 430: "LR", 434: "LY", 438: "LI", 440: "LT", 442: "LU", 446: "MO", 450: "MG", 454: "MW", 458: "MY", 462: "MV", 466: "ML", 470: "MT", 474: "MQ", 478: "MR", 480: "MU", 484: "MX", 492: "MC", 496: "MN", 498: "MD", 499: "ME", 500: "MS", 504: "MA", 508: "MZ", 512: "OM", 516: "NA", 520: "NR", 524: "NP", 528: "NL", 531: "CW", 533: "AW", 534: "SX", 535: "BQ", 540: "NC", 548: "VU", 554: "NZ", 558: "NI", 562: "NE", 566: "NG", 570: "NU", 574: "NF", 578: "NO", 580: "MP", 581: "UM", 583: "FM", 584: "MH", 585: "PW", 586: "PK", 591: "PA", 598: "PG", 600: "PY", 604: "PE", 608: "PH", 612: "PN", 616: "PL", 620: "PT", 624: "GW", 626: "TL", 630: "PR", 634: "QA", 638: "RE", 642: "RO", 643: "RU", 646: "RW", 652: "BL", 654: "SH", 659: "KN", 660: "AI", 662: "LC", 666: "PM", 670: "VC", 674: "SM", 678: "ST", 682: "SA", 686: "SN", 688: "RS", 690: "SC", 694: "SL", 702: "SG", 703: "SK", 704: "VN", 705: "SI", 706: "SO", 710: "ZA", 716: "ZW", 724: "ES", 728: "SS", 729: "SD", 732: "EH", 740: "SR", 744: "SJ", 748: "SZ", 752: "SE", 756: "CH", 760: "SY", 762: "TJ", 764: "TH", 768: "TG", 772: "TK", 776: "TO", 780: "TT", 784: "AE", 788: "TN", 792: "TR", 795: "TM", 796: "TC", 798: "TV", 800: "UG", 804: "UA", 807: "MK", 818: "EG", 826: "GB", 831: "GG", 832: "JE", 833: "IM", 834: "TZ", 840: "US", 850: "VI", 854: "BF", 858: "UY", 860: "UZ", 862: "VE", 876: "WF", 882: "WS", 887: "YE", 894: "ZM"}
	return countryArray[country]
}

// DeviceIds is...
func DeviceIds(device string) int32 {
	deviceArray := map[string]int{"web": 0, "ios": 1, "android": 2, "appletv": 3, "smarttv": 4, "roku": 5, "xbox_one": 6, "playstation": 7, "special": 8, "android_tv": 9, "amazon_fire_tv": 10}
	return int32(deviceArray[device])
}

// DeviceNames is...
func DeviceNames(device int) string {
	deviceArray := map[int]string{0: "web", 1: "ios", 2: "android", 3: "appletv", 4: "smarttv", 5: "roku", 6: "xbox_one", 7: "playstation", 8: "special", 9: "android_tv", 10: "amazon_fire_tv"}
	return deviceArray[device]
}

// ProductNames is...
func ProductNames(product string) int32 {
	ProductArray := map[string]int{"Weyyak": 1, "WeyyakSouthAsian": 2, "Africa": 3, "AfricaSouthAsian": 4, "Europe": 5,
		"Global": 6, "Apac": 7}
	return int32(ProductArray[product])
}

// ContentRightsTypes is...
func ContentRightsTypes(rightType string) int32 {
	ContentRightsTypesArrays := map[string]int{"Avod": 1, "Vod": 2, "Svod": 3, "Tvod": 4}
	return int32(ContentRightsTypesArrays[rightType])
}

func ContentRightsTypesInt(rightType int) string {
	ContentRightsTypesArrays := map[int]string{1: "Avod", 2: "Vod", 3: "Svod", 4: "Tvod"}
	return string(ContentRightsTypesArrays[rightType])
}

func SubscriptionPlansInt(planType int) string {
	ContentRightsTypesArrays := map[int]string{3: "Weyyak Prime"}
	return string(ContentRightsTypesArrays[planType])
}

func ContentPlatformsInt(PublishingPlatforms []int) []string {
	var platforms []string

	for _, i := range PublishingPlatforms {
		platforms = append(platforms, DeviceNames(i))
	}

	return platforms
}

func ContentSubscriptionPlansInt(SubscriptionPlans []int) []string {
	var Plans []string
	for _, i := range SubscriptionPlans {
		Plans = append(Plans, SubscriptionPlansInt(i))
	}
	return Plans
}

func ContentLocationsInt(DigitalRightsRegions []int) []string {
	var location []string

	for _, i := range DigitalRightsRegions {
		location = append(location, CountryNames(i))
	}

	return location
}

// PageTypes is...
func PageTypes(PageTypeID int) string {
	PageTypesArray := map[int]string{0: "VOD", 1: "Home", 8: "Settings", 16: "Favourites"}
	return PageTypesArray[PageTypeID]
}

// AgeRatings is...
func AgeRatings(ageRating int, language string) string {
	Rating := int(ageRating)
	AgeRatingArray := map[int]map[string]string{1: {"EnglishName": "G – General Audiences", "ArabicName": "لجميع الأعمار"}, 2: {"EnglishName": "PG – Parental Guidance Suggested", "ArabicName": "بإشراف الوالدين"}, 3: {"EnglishName": "PG 13 – Parents Strongly Cautioned", "ArabicName": "غير مناسبة للأطفال تحت سن الــ 13"}, 4: {"EnglishName": "PG 15 - Parents Strongly Cautioned", "ArabicName": "غير مناسبة للأطفال تحت سن الــ 15"}, 5: {"EnglishName": "15+", "ArabicName": "15+"}, 6: {"EnglishName": "18+", "ArabicName": "18+"}, 7: {"EnglishName": "NR - Not Rated by MPAA", "ArabicName": "لم يتم تقييمه من قبل الرقابة"}}
	var AgeRatingString string
	if language == "en" && Rating > 0 {
		AgeRatingString = AgeRatingArray[Rating]["EnglishName"]
	} else if language == "ar" && Rating > 0 {
		AgeRatingString = AgeRatingArray[Rating]["ArabicName"]
	}
	return AgeRatingString
}
func FindString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Println("Authorization token: ", c.Request.Header.Get("Authorization"))
		db := c.MustGet("UDB").(*gorm.DB)
		// requestBody := c.Request.Body
		// jsonData, _ := ioutil.ReadAll(requestBody)
		// fmt.Println(string(jsonData), "+++++++++")
		reqToken := c.Request.Header.Get("Authorization")
		if reqToken == "" {
			reqToken = c.Request.Header.Get("authorization")
		}
		// RequestLogRegister(string(jsonData), c)
		replacer := strings.NewReplacer("bearer", "Bearer", "BEARER", "Bearer")
		reqToken = replacer.Replace(reqToken)
		splitToken := strings.Split(reqToken, "Bearer ")
		//fmt.Println("split token:", splitToken)
		if len(splitToken) == 2 {
			c.Set("AuthorizationRequired", 0)
			// c.Abort()
			// l.JSON(c, http.StatusUnauthorized, gin.H{"Error": "Authentication Required."})
			// return
			reqToken = splitToken[1]
		} else {
			c.Set("AuthorizationRequired", 1)
		}

		// fmt.Println("kkkkkkkkkkk", reqToken)
		type Details struct {
			Userid           string
			DeviceId         string
			DeviceName       string
			DevicePlatform   string
			LanguageId       int
			IsBackOfficeUser bool
		}
		var response Details
		// TODO - Redis need to cache
		if reqToken != "" {
			db.Debug().Raw("select data->>'UserID' as userid, data->>'DeviceID' as device_id, data->>'DeviceName' as device_name, data->>'DevicePlatform' as device_platform, data->>'LanguageId' as language_id, data->>'IsBackOfficeUser' as is_back_office_user  from oauth2_tokens ot where access = ? and (data->>'ExpiresAt')::timestamp >= now()", reqToken).Scan(&response)
		}
		if response.Userid != "" {
			//fmt.Println("Weyyak2 token storing")
			c.Set("userid", response.Userid)
			c.Set("device_id", response.DeviceId)
			c.Set("device_name", response.DeviceName)
			c.Set("device_platform", response.DevicePlatform)
			c.Set("language_id", response.LanguageId)
			c.Set("is_back_office_user", response.IsBackOfficeUser)
			c.Next()
		} else {
			//fmt.Println("Weyyak1 token storing")
			url := os.Getenv("DOTNET_URL") + reqToken
			method := "GET"
			client := &http.Client{}
			req, err := http.NewRequest(method, url, nil)

			if err != nil {
				fmt.Println(err)
				return
			}
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
			// fmt.Println(string(body))
			type DotNetResponse struct {
				UserId        string `json:"user_id"`
				UserName      string `json:"user_name"`
				LanguageId    string `json:"language_id"`
				Role          string `json:"role"`
				DeviceId      string `json:"device_id"`
				ExpiresAt     string `json:"expire_at"`
				IsPersistence bool   `json:"is_persistence"`
				IssueedAt     string `json:"issueed_at"`
			}
			var dotNetToken DotNetResponse
			json.Unmarshal(body, &dotNetToken)
			if dotNetToken.UserId != "" {
				c.Set("userid", dotNetToken.UserId)
				c.Set("device_id", dotNetToken.DeviceId)
				c.Set("language_id", dotNetToken.LanguageId)
				c.Set("is_back_office_user", false)
			} else {
				c.Set("AuthorizationRequired", 1)
				c.Set("userid", "")
				c.Next()
			}
		}
	}
}
func DeleteEmpty(s []int) []int {
	var r []int
	for _, str := range s {
		if str != 0 {
			r = append(r, str)
		}
	}
	return r
}
func GenerateRandomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return strings.ToLower(string(b))
}

func ServerErrorResponse() ServerError {
	var serverError ServerError
	serverError.Error = SERVER_ERROR
	serverError.Description = EN_SERVER_ERROR_DESCRIPTION
	serverError.Code = SERVER_ERROR_CODE
	serverError.RequestId = GenerateRandomString(32)
	return serverError
}
func NotFoundErrorResponse() ServerError {
	var serverError ServerError
	serverError.Error = NOT_FOUND_ERROR
	serverError.Description = NOT_FOUND_ERROR_DESCRIPTION
	serverError.Code = NOT_FOUND_ERROR_CODE
	serverError.RequestId = GenerateRandomString(32)
	return serverError
}

const BULK_INSERT_LIMIT = 3000
const USERID = "2f634603-ce5b-eb11-831d-020666e39080"
const EN_SERVER_ERROR_DESCRIPTION = "Server error"
const AR_SERVER_ERROR_DESCRIPTION = "خطأ في الخادم"
const SERVER_ERROR = "server_error"
const SERVER_ERROR_CODE = "error_server_error"
const NOT_FOUND_ERROR = "not_found"
const NOT_FOUND_ERROR_CODE = ""
const NOT_FOUND_ERROR_DESCRIPTION = "Not found"
const BAD_REQUEST_ERROR = "invalid_request"
const BAD_REQUEST_ERROR_CODE = "error_validation_failed"
const BAD_REQUEST_ERROR_DESCRIPTION = "Validation failed."

// ContentLanguageOriginTypes is...
func ContentLanguageOriginTypes(languageType string) int {
	ContentLanguageOriginTypesArrays := map[string]int{"Original": 1, "Dubbed": 2, "Subtitled": 3}
	return ContentLanguageOriginTypesArrays[languageType]
}

// ContentLanguageOriginTypes is...
func ContentLanguageOriginTypesName(languageType int) string {
	ContentLanguageOriginTypesArrays := map[int]string{1: "Original", 2: "Dubbed", 3: "Subtitled"}
	return ContentLanguageOriginTypesArrays[languageType]
}

func GetVideoDuration(VideoId string) (string, int, int) {
	type VideoDurationInfo struct {
		Duration     int           `json:"duration"`
		Thumbnails   []interface{} `json:"thumbnails"`
		UrlTrickplay string        `json:"url_trickplay"`
		UrlVideo     string        `json:"url_video"`
	}
	response, err := GetCurlCall(os.Getenv("VIDEO_API") + VideoId)
	if err != nil {
		return err.Error(), http.StatusBadRequest, 0
	}
	var DataResponse VideoDurationInfo
	json.Unmarshal(response, &DataResponse)
	return "", http.StatusOK, DataResponse.Duration
}

/*Monitering if content get modified or updated dirty count will increase in content_sync table */
type ContentSync struct {
	ContentId  string `json:"contentId"`
	DirtyCount int    `json:"dirtyCount"`
}

func ContentSynching(contentId string, c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	var contentSync ContentSync
	if countError := fdb.Debug().Table("content_sync").Where("content_id=?", contentId).Find(&contentSync).Error; countError != nil {
		fmt.Println("err", countError)
	}
	if contentSync.DirtyCount > 0 {
		contentSync.ContentId = contentId
		contentSync.DirtyCount = contentSync.DirtyCount + 1
		if updateError := fdb.Debug().Table("content_sync").Where("content_id=?", contentId).Update(&contentSync).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	} else {
		contentSync.ContentId = contentId
		contentSync.DirtyCount = contentSync.DirtyCount + 1
		if updateError := fdb.Debug().Table("content_sync").Create(&contentSync).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	}
}

type PageSync struct {
	PageId     string `json:"pageId"`
	DirtyCount int    `json:"dirtyCount"`
	PageKey    int    `json:"pageKey"`
}

type ContentIds struct {
	contentId string
}

/* sync tables updation */
func PageSyncWithContentId(c *gin.Context) {
	fcdb := c.MustGet("FCDB").(*gorm.DB)
	fdb := c.MustGet("FDB").(*gorm.DB)
	//fmt.Println("AAAAAAAAAAAAAAAAA")
	var pageSync []PageSync
	var contentIds []ContentIds
	if fetchResult := fdb.Debug().Table("content_sync cs").Select("cs.content_id").Where("cs.dirty_count>?", 0).Find(&contentIds).Error; fetchResult != nil {
		fmt.Println(fetchResult)
		return
	}
	for _, contentIds := range contentIds {
		/*Fetch page details with content Id through playlist and slider */
		if sliderError := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").
			Joins("join playlist p on p.id = pi2.playlist_id").
			Joins("join page_playlist pp on pp.playlist_id = p.id").
			Joins("join page p2 on p2.id = pp.page_id").
			Joins("join slider s on s.black_area_playlist_id = p.id or s.red_area_playlist_id =p.id or s.green_area_playlist_id =p.id").
			Where("pi2.one_tier_content_id =? or pi2.multi_tier_content_id =?", contentIds.contentId, contentIds.contentId).
			Find(&pageSync).Error; sliderError != nil {
			fmt.Println(sliderError)
			return
		}
		if len(pageSync) < 1 {
			/*Fetch page details with content Id through playlist */
			if playError := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").
				Joins("join playlist p on p.id = pi2.playlist_id").
				Joins("join page_playlist pp on pp.playlist_id = p.id").
				Joins("join page p2 on p2.id = pp.page_id").
				Where("pi2.one_tier_content_id =? or pi2.multi_tier_content_id =?", contentIds.contentId, contentIds.contentId).
				Find(&pageSync).Error; playError != nil {
				fmt.Println(playError)
				return
			}
		}
		for _, pkeys := range pageSync {
			var pageSyncDetails PageSync
			if countError := fdb.Debug().Table("page_sync ps").Select("ps.page_id,ps.page_key,ps.dirty_count").Where("ps.page_id=?", pkeys.PageId).Find(&pageSyncDetails).Error; countError != nil {
				fmt.Println("err", countError)
				return
			}
			var insertPageSync PageSync
			if pageSyncDetails.DirtyCount > 0 {
				insertPageSync.PageId = pageSyncDetails.PageId
				insertPageSync.DirtyCount = pageSyncDetails.DirtyCount + 1
				insertPageSync.PageKey = pageSyncDetails.PageKey
				if updateError := fdb.Debug().Table("page_sync").Where("page_id=?", pkeys.PageId).Update(&insertPageSync).Error; updateError != nil {
					fmt.Println(updateError)
					return
				}
			} else {
				insertPageSync.PageId = pageSyncDetails.PageId
				insertPageSync.DirtyCount = 1
				insertPageSync.PageKey = pageSyncDetails.PageKey
				if updateError := fdb.Debug().Table("page_sync").Create(&insertPageSync).Error; updateError != nil {
					fmt.Println(updateError)
					return
				}
			}
		}
		var deletecontentSync ContentSync
		if deleteError := fdb.Debug().Where("content_id=?", contentIds.contentId).Delete(&deletecontentSync).Error; deleteError != nil {
			fmt.Println(deleteError)
			return
		}
	}

}

// /* sync tables updation */
// func PageSyncWithContentId(contentId string, c *gin.Context) {
// 	fcdb := c.MustGet("FCDB").(*gorm.DB)
// 	fdb := c.MustGet("FDB").(*gorm.DB)
// 	fmt.Println("AAAAAAAAAAAAAAAAA")
// 	var pageSync []PageSync
// 	/*Fetch page details with content Id through playlist and slider */
// 	if sliderError := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").
// 		Joins("join playlist p on p.id = pi2.playlist_id").
// 		Joins("join page_playlist pp on pp.playlist_id = p.id").
// 		Joins("join page p2 on p2.id = pp.page_id").
// 		Joins("join slider s on s.black_area_playlist_id = p.id or s.red_area_playlist_id =p.id or s.green_area_playlist_id =p.id").
// 		Where("pi2.one_tier_content_id =? or pi2.multi_tier_content_id =?", contentId, contentId).
// 		Find(&pageSync).Error; sliderError != nil {
// 		fmt.Println(sliderError)
// 		return
// 	}
// 	if len(pageSync) < 1 {
// 		/*Fetch page details with content Id through playlist */
// 		if playError := fcdb.Debug().Table("playlist_item pi2").Select("p2.page_key,p2.id as page_id").
// 			Joins("join playlist p on p.id = pi2.playlist_id").
// 			Joins("join page_playlist pp on pp.playlist_id = p.id").
// 			Joins("join page p2 on p2.id = pp.page_id").
// 			Where("pi2.one_tier_content_id =? or pi2.multi_tier_content_id =?", contentId, contentId).
// 			Find(&pageSync).Error; playError != nil {
// 			fmt.Println(playError)
// 			return
// 		}
// 	}
// 	for _, pkeys := range pageSync {
// 		var pageSyncDetails PageSync
// 		if countError := fdb.Debug().Table("page_sync ps").Select("ps.page_id,ps.page_key,ps.dirty_count").Where("ps.page_id=?", pkeys.PageId).Find(&pageSyncDetails).Error; countError != nil {
// 			fmt.Println("err", countError)
// 			return
// 		}
// 		var insertPageSync PageSync
// 		if pageSyncDetails.DirtyCount > 0 {
// 			insertPageSync.PageId = pageSyncDetails.PageId
// 			insertPageSync.DirtyCount = pageSyncDetails.DirtyCount + 1
// 			insertPageSync.PageKey = pageSyncDetails.PageKey
// 			if updateError := fdb.Debug().Table("page_sync").Where("page_id=?", pkeys.PageId).Update(&insertPageSync).Error; updateError != nil {
// 				fmt.Println(updateError)
// 				return
// 			}
// 		} else {
// 			insertPageSync.PageId = pageSyncDetails.PageId
// 			insertPageSync.DirtyCount = 1
// 			insertPageSync.PageKey = pageSyncDetails.PageKey
// 			if updateError := fdb.Debug().Table("page_sync").Create(&insertPageSync).Error; updateError != nil {
// 				fmt.Println(updateError)
// 				return
// 			}
// 		}
// 	}
// 	var deletecontentSync ContentSync
// 	if deleteError := fdb.Where("content_id=?", contentId).Delete(&deletecontentSync).Error; deleteError != nil {
// 		fmt.Println(deleteError)
// 		return
// 	}
// }

/* create redis key for get all contents based on content type */

func CreateRedisKeyForContentType(contentType string, c *gin.Context) {
	/*delete Redis keys */
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	fmt.Println("inside redis")
	contentType = strings.ToLower(contentType)
	searchPattern := os.Getenv("REDIS_CONTENT_KEY") + contentType + "*"
	if len(os.Args) > 1 {
		searchPattern = os.Args[1]
	}
	iter := rdb.Keys(searchPattern)
	for _, Rkey := range iter.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

func CreateRedisKeyForContentTypeMTC(c *gin.Context) {
	/*delete Redis keys */
	contentType1 := "Series"
	contentType2 := "Program"
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	fmt.Println("inside redis MTC")
	contentType1 = strings.ToLower(contentType1)
	searchPattern1 := os.Getenv("REDIS_CONTENT_KEY") + contentType1 + "*"
	if len(os.Args) > 1 {
		searchPattern1 = os.Args[1]
	}
	iter1 := rdb.Keys(searchPattern1)
	for _, Rkey := range iter1.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
	contentType2 = strings.ToLower(contentType2)
	searchPattern2 := os.Getenv("REDIS_CONTENT_KEY") + contentType2 + "*"
	if len(os.Args) > 1 {
		searchPattern2 = os.Args[1]
	}
	iter2 := rdb.Keys(searchPattern2)
	for _, Rkey := range iter2.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

func CreateRedisKeyForContentTypeOTC(c *gin.Context) {
	/*delete Redis keys */
	contentType1 := "Movie"
	contentType2 := "LiveTV"
	contentType3 := "Play"
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	fmt.Println("inside redis OTC")
	contentType1 = strings.ToLower(contentType1)
	searchPattern1 := os.Getenv("REDIS_CONTENT_KEY") + contentType1 + "*"
	if len(os.Args) > 1 {
		searchPattern1 = os.Args[1]
	}
	iter1 := rdb.Keys(searchPattern1)
	for _, Rkey := range iter1.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
	contentType2 = strings.ToLower(contentType2)
	searchPattern2 := os.Getenv("REDIS_CONTENT_KEY") + contentType2 + "*"
	if len(os.Args) > 1 {
		searchPattern2 = os.Args[1]
	}
	iter2 := rdb.Keys(searchPattern2)
	for _, Rkey := range iter2.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
	contentType3 = strings.ToLower(contentType3)
	searchPattern3 := os.Getenv("REDIS_CONTENT_KEY") + contentType3 + "*"
	if len(os.Args) > 1 {
		searchPattern2 = os.Args[1]
	}
	iter3 := rdb.Keys(searchPattern3)
	for _, Rkey := range iter3.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

// need to implement in future for redis keys generation of contentType
/*Prepare  Redis Cache for all contents based on contentType */
// fcdb := c.MustGet("FCDB").(*gorm.DB)
// languages := [2]string{"en", "ar"}
// var regions Regions
// if regionfetch := fcdb.Raw("select string_agg(distinct c.id::text,',')as region from country c").Find(&regions).Error; regionfetch != nil {
// 	l.JSON(c, http.StatusInternalServerError, regionfetch)
// 	return
// }
// var country []string
// country = strings.Split(regions.Region, ",")
// for _, language := range languages {
// 	for _, countryId := range country {
// 		key := os.Getenv("REDIS_CONTENT_KEY") + contentType + language + countryId
// 		key = strings.ReplaceAll(key, "/", "_")
// 		key = strings.ReplaceAll(key, "?", "_")
// 		cId, _ := strconv.Atoi(countryId)
// 		ccode := CountryNames(cId)
// 		/*use any one url
// 		url := "https://ynk2yz6oak.execute-api.ap-south-1.amazonaws.com/weyyak-fo-ms-api-qa/v1/en/contents/contentType?contentType=" + contentType + "&Country=" /+ ccode + "&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0"*/
// 		url := "https://apiqafo.engro.in/v1/" + language + "/contents/contentType?contentType=" + contentType + "&Country=" + ccode + "&pageNo=1&OrderBy=desc&RowCountPerPage=50&IsPaging=0"
// 		method := "GET"
// 		client := &http.Client{}
// 		req, err := http.NewRequest(method, url, nil)

// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		res, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		defer res.Body.Close()
// 		body, err := ioutil.ReadAll(res.Body)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		Result := []AllAvailableSeasons{}
// 		json.Unmarshal(body, &Result)
// 		/* Create Redis Key for content Type */
// 		jsonData, _ := json.Marshal(Result)
// 		var request RedisCacheRequest
// 		Redisurl := os.Getenv("REDIS_CACHE_URL")
// 		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
// 		request.Value = string(jsonData)
// 		_, posterr := PostCurlCall("POST", Redisurl, request)
// 		if err != nil {
// 			fmt.Println(posterr)
// 			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
// 			return
// 		}
// 	}
// }

// func RequestLogRegister(reqBody string, c *gin.Context) {
// 	// func() {
// 	db := c.MustGet("SDB").(*gorm.DB)
// 	requestMethod := c.Request.Method
// 	requestURL := c.Request.Host + c.Request.URL.Path
// 	// requestBody := c.Request.Body
// 	// jsonData, err := ioutil.ReadAll(requestBody)
// 	// fmt.Println(jsonData)
// 	// if err != nil {
// 	// 	fmt.Println("error for request body conversion:", err)
// 	// }
// 	reqToken := c.Request.Header.Get("Authorization")
// 	if reqToken == "" {
// 		reqToken = c.Request.Header.Get("authorization")
// 	}
// 	type RequestLog struct {
// 		Id        int
// 		Method    string
// 		Url       string
// 		RawBody   string
// 		JsnBody   string
// 		UserToken string
// 	}
// 	var requestLog RequestLog
// 	requestLog.Method = requestMethod
// 	requestLog.Url = requestURL
// 	// requestLog.RawBody = requestBody
// 	fmt.Println(reqBody, "requestBody")
// 	requestLog.JsnBody = reqBody
// 	requestLog.UserToken = reqToken
// 	if c.Request.URL.Path != "/health" {
// 		if err := db.Table("request_log").Create(&requestLog).Error; err != nil {
// 			// l.JSON(c, http.StatusInternalServerError, err)
// 			// return
// 			fmt.Println("error in logging:", err)
// 		}
// 		fmt.Println("Register Log Id: ", requestLog.Id)
// 	}
// 	// }
// }

func CreateRedisKeyForContent(contentKey string, c *gin.Context) {
	/*delete Redis keys */
	// use contenttype for specific content
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	fmt.Println("inside redis 1")
	contentKey = strings.ToLower(contentKey)
	searchPattern := os.Getenv("REDIS_CONTENT_KEY") + contentKey + "*"
	if len(os.Args) > 1 {
		searchPattern = os.Args[1]
	}
	iter := rdb.Keys(searchPattern)
	for _, Rkey := range iter.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

//need to implement in future redis for single content
/*fcdb := c.MustGet("FCDB").(*gorm.DB)
languages := [2]string{"en", "ar"}
var regions Regions
if regionfetch := fcdb.Raw("select string_agg(distinct c.id::text,',')as region from country c").Find(&regions).Error; regionfetch != nil {
	l.JSON(c, http.StatusInternalServerError, regionfetch)
	return
}
var country []string
country = strings.Split(regions.Region, ",")
for _, language := range languages {
	for _, countryId := range country {
		key := os.Getenv("REDIS_CONTENT_KEY") + contentKey + language + countryId
		key = strings.ReplaceAll(key, "/", "_")
		key = strings.ReplaceAll(key, "?", "_")
		cId, _ := strconv.Atoi(countryId)
		ccode := CountryNames(cId)

		var url string
		// if contentType == "Movie" {
		// 	url = "https://apiqafo.engro.in/v1/en/contents/moviedetails?Country=" + ccode + "&contentkey=" + contentKey
		// } else if contentType == "Series" {
		// 	url = "https://apiqafo.engro.in/v1/en/series?contentkey=" + contentKey + "&country=" + ccode
		// }
		url = "https://apiqafo.engro.in/v1/" + language + "/contents/moviedetailsresponse?Country=" + ccode + "&contentkey=" + contentKey

		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
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
		Result := PlaylistContent{}
		json.Unmarshal(body, &Result)
		//fmt.Println(Result, "?+?+?+?+?+?+?+?+")
		//  Create Redis Key for content Type
		jsonData, _ := json.Marshal(Result)
		var request RedisCacheRequest
		Redisurl := os.Getenv("REDIS_CACHE_URL")
		request.Key = key //pageKey + language + strconv.Itoa(country) + strconv.Itoa(platform)
		request.Value = string(jsonData)
		_, posterr := PostCurlCall("POST", Redisurl, request)
		if err != nil {
			fmt.Println(posterr)
			l.JSON(c, http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}
}*/

func PlaylistSynching(PlaylistId string, c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	var playlistSync PlaylistSync
	playlist := fdb.Debug().Table("playlist_sync").Where("playlist_id=?", PlaylistId).Find(&playlistSync)
	totalcount := int(playlist.RowsAffected)
	if totalcount < 1 {
		playlistSync.PlaylistId = PlaylistId
		playlistSync.DirtyCount = 1
		if updateError := fdb.Debug().Table("playlist_sync").Create(&playlistSync).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	} else {
		playlistSync.PlaylistId = PlaylistId
		playlistSync.DirtyCount = playlistSync.DirtyCount + 1
		if updateError := fdb.Debug().Table("playlist_sync").Where("playlist_id=?", PlaylistId).Update(&playlistSync).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	}
}

func SliderSynching(SliderId string, c *gin.Context) {
	fdb := c.MustGet("FDB").(*gorm.DB)
	var sliderSync SliderSync
	slider := fdb.Debug().Table("slider_sync").Where("slider_id=?", SliderId).Find(&sliderSync)
	totalcount := int(slider.RowsAffected)
	if totalcount < 1 {
		sliderSync.SliderId = SliderId
		sliderSync.DirtyCount = 1
		if updateError := fdb.Debug().Table("slider_sync").Create(&sliderSync).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	} else {
		var updateslider SliderSync
		updateslider.SliderId = sliderSync.SliderId
		updateslider.DirtyCount = sliderSync.DirtyCount + 1
		if updateError := fdb.Debug().Table("slider_sync").Where("slider_id=?", SliderId).Update(&updateslider).Error; updateError != nil {
			fmt.Println(updateError)
			return
		}
	}
}

func GetRedisDataWithKey(key string) (string, error) {
	resp, err := http.Get(os.Getenv("REDIS_CACHE_URL") + "/" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func PostRedisDataWithKey(key string, data []byte) error {

	type RedisPostData struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	}

	finalResult := RedisPostData{
		Key:   key,
		Value: string(data),
	}

	jsonString, _ := json.Marshal(&finalResult)

	body := bytes.NewReader(jsonString)

	req, err := http.NewRequest("POST", os.Getenv("REDIS_CACHE_URL"), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func ClearRedisKeyFollowKeys(c *gin.Context, pageKey string) {
	/*delete Redis keys for pages */
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	searchPattern := "*" + pageKey + "*"

	if len(os.Args) > 1 {
		searchPattern = os.Args[1]
	}

	iter := rdb.Keys(searchPattern)
	for _, Rkey := range iter.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

func ClearRedisKeyKeys(c *gin.Context, pageKey string) {
	/*delete Redis keys for pages */
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	searchPattern := pageKey

	if len(os.Args) > 1 {
		searchPattern = os.Args[1]
	}

	iter := rdb.Keys(searchPattern)
	for _, Rkey := range iter.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}
