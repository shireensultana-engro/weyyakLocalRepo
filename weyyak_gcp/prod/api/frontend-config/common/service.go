package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

const BULK_INSERT_LIMIT = 3000
const BLACK_AREA_PLAYLIST_CONUT = 1
const RED_AREA_PLAYLIST_CONUT = 7
const GREEN_AREA_PLAYLIST_CONUT = 1
const SNP_UNPUBLISHED_CODE = "slider_playlist_unpublished"
const SNP_UNPUBLISHED_MESSAGE = "Slider '#stitle#' contains hidden playlist '#ptitle#'"
const SNP_SCHEDULING_MISMATCH_CODE = "slider_playlist_scheduling_mismatch"
const SNP_SCHEDULING_MISMATCH_MESSAGE = "Slider '#stitle#' scheduling (#sssdate# - #ssedate#) does not match to playlist '#ptitle#' scheduling (#pssdate# - #psedate#)"
const SNP_ITEM_COUNT_MISMATCH_CODE = "slider_playlist_items_count_mismatch"
const SNP_ITEM_COUNT_MISMATCH_MESSAGE = "Slider '#stitle#' contains playlist '#ptitle#' that does not contain minimum of published Movies/Seasons/Series with matching scheduling. Expected count: #ecount#, actual: #acount#"

func GenerateRandomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return strings.ToLower(string(b))
}

// ServerError -- binding struct for error response
type ServerError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
}

type PageType struct {
	PageType int `json:"page_type"`
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

// Countrys is...
func Countrys(country string) int32 {
	countryArray := map[string]int{"AF": 4, "AQ": 10, "DZ": 12, "AS": 16, "AD": 20, "AO": 24, "AG": 28, "AZ": 31, "AR": 32, "AU": 36, "AT": 40, "BS": 44, "BH": 48, "BD": 50, "AM": 51, "BB": 52, "BE": 56, "BM": 60, "BT": 64, "BO": 68, "BA": 70, "BW": 72, "BV": 74, "BR": 76, "BZ": 84, "IO": 86, "SB": 90, "VG": 92, "BN": 96, "BG": 100, "MM": 104, "BI": 108, "BY": 112, "KH": 116, "CM": 120, "CA": 124, "CV": 132, "KY": 136, "CF": 140, "LK": 144, "TD": 148, "CL": 152, "CN": 156, "TW": 158, "CX": 162, "CC": 166, "CO": 170, "KM": 174, "YT": 175, "CG": 178, "CD": 180, "CK": 184, "CR": 188, "HR": 191, "CU": 192, "CY": 196, "CZ": 203, "BJ": 204, "DK": 208, "DM": 212, "DO": 214, "EC": 218, "SV": 222, "GQ": 226, "ET": 231, "ER": 232, "EE": 233, "FO": 234, "FK": 238, "GS": 239, "FJ": 242, "FI": 246, "AX": 248, "FR": 250, "GF": 254, "PF": 258, "TF": 260, "DJ": 262, "GA": 266, "GE": 268, "GM": 270, "PS": 275, "DE": 276, "GH": 288, "GI": 292, "KI": 296, "GR": 300, "GL": 304, "GD": 308, "GP": 312, "GU": 316, "GT": 320, "GN": 324, "GY": 328, "HT": 332, "HM": 334, "VA": 336, "HN": 340, "HK": 344, "HU": 348, "IS": 352, "IN": 356, "ID": 360, "IR": 364, "IQ": 368, "IE": 372, "IL": 376, "IT": 380, "CI": 384, "JM": 388, "JP": 392, "KZ": 398, "JO": 400, "KE": 404, "KP": 408, "KR": 410, "KW": 414, "KG": 417, "LA": 418, "LB": 422, "LS": 426, "LV": 428, "LR": 430, "LY": 434, "LI": 438, "LT": 440, "LU": 442, "MO": 446, "MG": 450, "MW": 454, "MY": 458, "MV": 462, "ML": 466, "MT": 470, "MQ": 474, "MR": 478, "MU": 480, "MX": 484, "MC": 492, "MN": 496, "MD": 498, "ME": 499, "MS": 500, "MA": 504, "MZ": 508, "OM": 512, "NA": 516, "NR": 520, "NP": 524, "NL": 528, "CW": 531, "AW": 533, "SX": 534, "BQ": 535, "NC": 540, "VU": 548, "NZ": 554, "NI": 558, "NE": 562, "NG": 566, "NU": 570, "NF": 574, "NO": 578, "MP": 580, "UM": 581, "FM": 583, "MH": 584, "PW": 585, "PK": 586, "PA": 591, "PG": 598, "PY": 600, "PE": 604, "PH": 608, "PN": 612, "PL": 616, "PT": 620, "GW": 624, "TL": 626, "PR": 630, "QA": 634, "RE": 638, "RO": 642, "RU": 643, "RW": 646, "BL": 652, "SH": 654, "KN": 659, "AI": 660, "LC": 662, "PM": 666, "VC": 670, "SM": 674, "ST": 678, "SA": 682, "SN": 686, "RS": 688, "SC": 690, "SL": 694, "SG": 702, "SK": 703, "VN": 704, "SI": 705, "SO": 706, "ZA": 710, "ZW": 716, "ES": 724, "SS": 728, "SD": 729, "EH": 732, "SR": 740, "SJ": 744, "SZ": 748, "SE": 752, "CH": 756, "SY": 760, "TJ": 762, "TH": 764, "TG": 768, "TK": 772, "TO": 776, "TT": 780, "AE": 784, "TN": 788, "TR": 792, "TM": 795, "TC": 796, "TV": 798, "UG": 800, "UA": 804, "MK": 807, "EG": 818, "GB": 826, "GG": 831, "JE": 832, "IM": 833, "TZ": 834, "US": 840, "VI": 850, "BF": 854, "UY": 858, "UZ": 860, "VE": 862, "WF": 876, "WS": 882, "YE": 887, "ZM": 894}
	return int32(countryArray[country])
}

// CountryNames is...
func CountryNames(country int) string {
	countryArray := map[int]string{4: "AF", 8: "AL", 10: "AQ", 12: "DZ", 16: "AS", 20: "AD", 24: "AO", 28: "AG", 31: "AZ", 32: "AR", 36: "AU", 40: "AT", 44: "BS", 48: "BH", 50: "BD", 51: "AM", 52: "BB", 56: "BE", 60: "BM", 64: "BT", 68: "BO", 70: "BA", 72: "BW", 74: "BV", 76: "BR", 84: "BZ", 86: "IO", 90: "SB", 92: "VG", 96: "BN", 100: "BG", 104: "MM", 108: "BI", 112: "BY", 116: "KH", 120: "CM", 124: "CA", 132: "CV", 136: "KY", 140: "CF", 144: "LK", 148: "TD", 152: "CL", 156: "CN", 158: "TW", 162: "CX", 166: "CC", 170: "CO", 174: "KM", 175: "YT", 178: "CG", 180: "CD", 184: "CK", 188: "CR", 191: "HR", 192: "CU", 196: "CY", 203: "CZ", 204: "BJ", 208: "DK", 212: "DM", 214: "DO", 218: "EC", 222: "SV", 226: "GQ", 231: "ET", 232: "ER", 233: "EE", 234: "FO", 238: "FK", 239: "GS", 242: "FJ", 246: "FI", 248: "AX", 250: "FR", 254: "GF", 258: "PF", 260: "TF", 262: "DJ", 266: "GA", 268: "GE", 270: "GM", 275: "PS", 276: "DE", 288: "GH", 292: "GI", 296: "KI", 300: "GR", 304: "GL", 308: "GD", 312: "GP", 316: "GU", 320: "GT", 324: "GN", 328: "GY", 332: "HT", 334: "HM", 336: "VA", 340: "HN", 344: "HK", 348: "HU", 352: "IS", 356: "IN", 360: "ID", 364: "IR", 368: "IQ", 372: "IE", 376: "IL", 380: "IT", 384: "CI", 388: "JM", 392: "JP", 398: "KZ", 400: "JO", 404: "KE", 408: "KP", 410: "KR", 414: "KW", 417: "KG", 418: "LA", 422: "LB", 426: "LS", 428: "LV", 430: "LR", 434: "LY", 438: "LI", 440: "LT", 442: "LU", 446: "MO", 450: "MG", 454: "MW", 458: "MY", 462: "MV", 466: "ML", 470: "MT", 474: "MQ", 478: "MR", 480: "MU", 484: "MX", 492: "MC", 496: "MN", 498: "MD", 499: "ME", 500: "MS", 504: "MA", 508: "MZ", 512: "OM", 516: "NA", 520: "NR", 524: "NP", 528: "NL", 531: "CW", 533: "AW", 534: "SX", 535: "BQ", 540: "NC", 548: "VU", 554: "NZ", 558: "NI", 562: "NE", 566: "NG", 570: "NU", 574: "NF", 578: "NO", 580: "MP", 581: "UM", 583: "FM", 584: "MH", 585: "PW", 586: "PK", 591: "PA", 598: "PG", 600: "PY", 604: "PE", 608: "PH", 612: "PN", 616: "PL", 620: "PT", 624: "GW", 626: "TL", 630: "PR", 634: "QA", 638: "RE", 642: "RO", 643: "RU", 646: "RW", 652: "BL", 654: "SH", 659: "KN", 660: "AI", 662: "LC", 666: "PM", 670: "VC", 674: "SM", 678: "ST", 682: "SA", 686: "SN", 688: "RS", 690: "SC", 694: "SL", 702: "SG", 703: "SK", 704: "VN", 705: "SI", 706: "SO", 710: "ZA", 716: "ZW", 724: "ES", 728: "SS", 729: "SD", 732: "EH", 740: "SR", 744: "SJ", 748: "SZ", 752: "SE", 756: "CH", 760: "SY", 762: "TJ", 764: "TH", 768: "TG", 772: "TK", 776: "TO", 780: "TT", 784: "AE", 788: "TN", 792: "TR", 795: "TM", 796: "TC", 798: "TV", 800: "UG", 804: "UA", 807: "MK", 818: "EG", 826: "GB", 831: "GG", 832: "JE", 833: "IM", 834: "TZ", 840: "US", 850: "VI", 854: "BF", 858: "UY", 860: "UZ", 862: "VE", 876: "WF", 882: "WS", 887: "YE", 894: "ZM"}
	return countryArray[country]
}

// DeviceNames is...
func DeviceNames(device int) string {
	deviceArray := map[int]string{0: "web", 1: "ios", 2: "android", 3: "appletv", 4: "smarttv", 5: "roku", 6: "xbox_one", 7: "playstation", 8: "special", 9: "android_tv", 10: "amazon_fire_tv"}
	return deviceArray[device]
}

// Publishing platforms is...
func PublishingPlatforms(Platform string) int32 {
	PlatformArray := map[string]int{"iOS": 1, "Android": 2, "Smart TV": 4, "Roku": 5, "Xbox One": 6, "PlayStation": 7, "Android TV": 9, "Amazon Fire TV": 10, "AppleTV": 3, "Web": 0}
	return int32(PlatformArray[Platform])
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// requestBody := c.Request.Body
		// jsonData, _ := ioutil.ReadAll(requestBody)
		// fmt.Println(string(jsonData), "+++++++++")
		// c.Set("requestbody",string(jsonData))
		db := c.MustGet("UDB").(*gorm.DB)
		reqToken := c.Request.Header.Get("Authorization")
		if reqToken == "" {
			reqToken = c.Request.Header.Get("authorization")
		}
		replacer := strings.NewReplacer("bearer", "Bearer", "BEARER", "Bearer")
		reqToken = replacer.Replace(reqToken)
		splitToken := strings.Split(reqToken, "Bearer ")
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
			LanguageId       string
			IsBackOfficeUser bool
		}
		var response Details
		// TODO - Redis need to cache
		if reqToken != "" {
			db.Debug().Raw("select data->>'UserID' as userid, data->>'DeviceID' as device_id, data->>'DeviceName' as device_name, data->>'DevicePlatform' as device_platform, data->>'LanguageId' as language_id , data->>'IsBackOfficeUser' as is_back_office_user from oauth2_tokens ot where access = ? and ((data->>'ExpiresAt')::timestamp >= now() at TIME zone 'UTC' or (data->>'ExpiresAt')::timestamp >= now())", reqToken).Scan(&response)
		}
		if response.Userid != "" {
			c.Set("userid", response.Userid)
			c.Set("device_id", response.DeviceId)
			c.Set("device_name", response.DeviceName)
			c.Set("device_platform", response.DevicePlatform)
			c.Set("language_id", response.LanguageId)
			c.Set("is_back_office_user", response.IsBackOfficeUser)
			c.Next()
		} else {
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
		// RequestLogRegister(string(jsonData), c)
	}
}

// TODO:images upload functionality is not complete
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
		Bucket:               aws.String(os.Getenv("S3_BUCKET")),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	return s3err
}
func ContentsByPlansQuery(countryId, PlanId string) string {
	query := "SELECT c.content_key,c.id,c.content_tier,c.created_at from content c join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=c.about_the_content_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and crc.country_id=" + countryId + " and crp.subscription_plan_id =" + PlanId + " group by c.content_key,c.id,c.content_tier,c.created_at UNION SELECT c.content_key,c.id,c.content_tier,c.created_at from content c join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null and crc.country_id=" + countryId + " and crp.subscription_plan_id =" + PlanId + " group by c.content_key,c.id,c.content_tier,c.created_at"
	return query
}
func ContentRatingQuery(contentkey string) string {
	query := "SELECT cpi.transliterated_title,c.content_type,pi1.duration as length,c.id,c.content_key from content c join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=c.about_the_content_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and c.content_key = " + contentkey + " group by  cpi.transliterated_title,c.content_type,pi1.duration,c.id,c.content_key UNION SELECT cpi.transliterated_title,c.content_type,(select sum(pi1.duration) from playback_item pi1 join episode e on e.playback_item_id =pi1.id join season s on s.id=e.season_id join content c on c.id =s.content_id where c.content_key =" + contentkey + ") as length,c.id,c.content_key from content c join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null)  and c.content_key = " + contentkey + " and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null group by cpi.transliterated_title,c.content_type,pi1.duration,c.id,c.content_key"
	return query
}
func MultitierContentQuery(contentId string, language string) (string, string, string, string) {
	groupBy := "c.content_key,c.content_type,s.english_meta_title,s.arabic_meta_title,s.english_meta_description,s.arabic_meta_description,c.created_at,c.modified_at,cpi.transliterated_title,cpi.arabic_title,atci.age_group,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,cpi.original_title,pi1.duration,pi1.video_content_id,cr.digital_rights_type,e.number,s.number,s.cast_id,pi1.scheduling_date_time,cr.digital_rights_start_date,cr.digital_rights_end_date,s.has_poster_image,s.has_details_background,s.has_mobile_details_background"
	fields := "s.cast_id,c.content_key as id,c.content_type,c.created_at as inserted_at,c.modified_at as modified_at,cpi.transliterated_title as friendly_url,atci.age_group as age_rating,atci.production_year,(select sum(pi1.duration) from playback_item pi1 join episode e on e.playback_item_id =pi1.id join season s on s.id =e.season_id where s.content_id ='" + contentId + "') as length,pi1.video_content_id as video_id,cr.digital_rights_type,pi1.scheduling_date_time,cr.digital_rights_start_date,cr.digital_rights_end_date,s.has_poster_image,s.has_details_background,s.has_mobile_details_background"
	if language == "en" {
		fields += ", s.english_meta_title as seo_title,s.english_meta_description as seo_description,cpi.transliterated_title as title,cpi.arabic_title as translated_title,atci.english_synopsis as synopsis"
	} else {
		fields += ", s.arabic_meta_title as seo_title,s.arabic_meta_description as seo_description,cpi.arabic_title as title,cpi.transliterated_title as translated_title,atci.arabic_synopsis as synopsis"
	}
	join := "join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.id=? and c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null and crc.country_id=?"
	return fields, join, Where, groupBy
}
func SeasonDetailsQuery(language string) (string, string, string, string) {
	fields := "s.id,s.season_key,s.number as season_number,cr.digital_rights_type,ct.language_type"
	if language == "en" {
		fields += " ,s.english_meta_description as seo_description,s.english_meta_title as seo_title,cpi.transliterated_title as title,atci.english_synopsis as synopsis"
	} else {
		fields += " ,s.arabic_meta_description as seo_description,s.arabic_meta_title as seo_title,cpi.arabic_title as title,atci.arabic_synopsis as synopsis"
	}
	join := "join content c on c.id=s.content_id join content_primary_info cpi on cpi.id =s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi2 on pi2.rights_id = s.rights_id join content_translation ct on ct.id =pi2.translation_id join content_rights cr on cr.id =s.rights_id join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.id=? and s.status =1 and s.deleted_by_user_id is null and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and (pi2.scheduling_date_time <=NOW() or pi2.scheduling_date_time is null) and crc.country_id =?"
	groupBy := "s.id,s.season_key,s.number,s.english_meta_description,s.english_meta_title,cpi.transliterated_title,cr.digital_rights_type,ct.language_type,s.arabic_meta_description,s.arabic_meta_title,cpi.arabic_title,atci.arabic_synopsis,atci.english_synopsis"
	return fields, join, Where, groupBy
}

func OnetierContentQuery(contentId string, language string) (string, string, string, string) {
	groupBy := "c.content_key,c.content_type,c.english_meta_title,c.arabic_meta_title,c.english_meta_description,c.arabic_meta_description,c.created_at,c.modified_at,cpi.transliterated_title,cpi.arabic_title,atci.age_group,atci.english_synopsis,atci.arabic_synopsis,atci.production_year,cpi.original_title,cv.has_all_rights,pi1.duration,cv.id,pi1.video_content_id,cr.digital_rights_type,c.cast_id,pi1.scheduling_date_time,cr.digital_rights_start_date,cr.digital_rights_end_date,c.has_poster_image,c.has_details_background,c.has_mobile_details_background"
	fields := "c.cast_id,c.content_key as id,c.content_type,c.created_at as inserted_at,c.modified_at as modifiedAt,cpi.transliterated_title as friendly_url,atci.age_group as age_rating,atci.production_year,cv.has_all_rights as geoblock,cv.id as content_version_id,pi1.duration as length,pi1.video_content_id as video_id,cr.digital_rights_type,c.english_meta_title as seo_title,c.english_meta_description as seo_description,cpi.transliterated_title as title,cpi.arabic_title as translated_title,atci.english_synopsis as synopsis,pi1.scheduling_date_time,cr.digital_rights_start_date,cr.digital_rights_end_date,c.has_poster_image,c.has_details_background,c.has_mobile_details_background"
	if language == "en" {
		fields += ", c.english_meta_title as seo_title,c.english_meta_description as seo_description,cpi.transliterated_title as title,cpi.arabic_title as translated_title,atci.english_synopsis as synopsis"
	} else {
		fields += ", c.arabic_meta_title as seo_title,c.arabic_meta_description as seo_description,cpi.arabic_title as title,cpi.transliterated_title as translated_title,atci.arabic_synopsis as synopsis"
	}
	join := "join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=c.about_the_content_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.id=? and c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null)  and crc.country_id=?"
	return fields, join, Where, groupBy
}
func MovieDetailsQuery(language string) (string, string, string, string) {
	groupBy := "crc.content_rights_id,cr.digital_rights_type,cv.order,crp.subscription_plan_id,cpi.transliterated_title,c.created_at,cpi.arabic_title,cv.id"
	fields := "cv.id,c.created_at as insertedAt,cr.digital_rights_type"
	if language == "en" {
		fields += ", cpi.transliterated_title as title"
	} else {
		fields += ", cpi.arabic_title as title"
	}
	join := "join content_primary_info cpi on cpi.id =c.primary_info_id join content_variance cv on cv.content_id=c.id join playback_item pi on pi.id =cv.playback_item_id join content_rights cr on cr.id =pi.rights_id join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id"
	Where := "content_id =? and cv.deleted_by_user_id is null and cv.status=1 and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and crc.country_id=?"
	return fields, join, Where, groupBy
}
func ContentActorsQuery(language string) (string, string, string, string) {
	fields := "cc.main_actor_id,cc.main_actress_id"
	if language == "en" {
		fields += ",string_agg(CAST(a.english_name AS varchar), ',') as actors"
	} else {
		fields += ",string_agg(CAST(a.arabic_name AS varchar), ',') as actors"
	}
	join := "full outer join content_actor ca on ca.cast_id=cc.id full outer join actor a on a.id=ca.actor_id"
	groupBy := "cc.main_actor_id,cc.main_actress_id"
	Where := "cc.id =?"
	return fields, join, Where, groupBy
}
func ContentGenresQuery(language string) (string, string, string, string) {
	var fields string
	if language == "en" {
		fields += "g.english_name as name"
	} else {
		fields += "g.arabic_name as name"
	}
	join := "join content_genre cg  on cg.genre_id=g.id join content c on c.id=cg.content_id"
	groupBy := "g.english_name,g.id"
	Where := "c.id =?"
	return fields, join, Where, groupBy
}
func ContentPlansQuery(content_type int) (string, string, string) {
	fields := "crp.subscription_plan_id"
	var join, Where string
	if content_type == 1 {
		join = "join playback_item pi2 on pi2.id =cv.playback_item_id join content_rights_plan crp on crp.rights_id = pi2.rights_id"
		Where = "cv.id =?"
	} else {
		join = "join season s on s.rights_id = crp.rights_id join content c on c.id =s.content_id"
		Where = "s.id =?"
	}

	return fields, join, Where
}
func SeasonEpisodesQuery(language string) (string, string, string, string) {
	fields := "e.id as episode_id,e.episode_key as id,c.content_key as series_id,e.number as episode_number,pi1.video_content_id as video_id,pi1.duration as length,e.created_at as insertedAt,cr.digital_rights_type as digitalRighttype"
	if language == "en" {
		fields += ",e.synopsis_english as synopsis,cpi.transliterated_title as title"
	} else {
		fields += ",e.synopsis_arabic as synopsis,cpi.arabic_title as title"
	}
	join := "join season s on s.id =e.season_id join content c on c.id = s.content_id join content_primary_info cpi on cpi.id = e.primary_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id =pi1.rights_id"
	Where := "e.season_id=? and e.status =1 and e.deleted_by_user_id is null"
	groupBy := "e.id,c.content_key,pi1.duration,cr.digital_rights_type,e.synopsis_english,cpi.transliterated_title,cpi.arabic_title,e.synopsis_arabic"
	return fields, join, Where, groupBy
}
func GetSeriesQuery(language string) (string, string, string, string) {
	fields := "c.id as content_id,s.id as season_id,c.content_key as id,c.content_tier,min(pi1.video_content_id) as video_id,Replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url"
	if language == "en" {
		fields += ",cpi.transliterated_title as title"
	} else {
		fields += ",cpi.arabic_title as title"
	}
	join := "join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id =e.primary_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=s.rights_id join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and c.status = 1 and c.deleted_by_user_id is null and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null and crc.country_id=? and lower(c.content_type)= ?"
	groupBy := "c.id,s.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at"
	return fields, join, Where, groupBy
}
func GetMoviesQuery(language string) (string, string, string, string) {
	fields := "c.id as content_id,cv.id as season_id,c.content_key as id,c.content_tier,min(pi1.video_content_id) as video_id,Replace(lower(cpi.transliterated_title), ' ', '_') as friendly_url"
	if language == "en" {
		fields += ",cpi.transliterated_title as title"
	} else {
		fields += ",cpi.arabic_title as title"
	}
	join := "join content_primary_info cpi on cpi.id=c.primary_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and crc.country_id=? and lower(c.content_type)= ?"
	groupBy := "c.id,cv.id,c.content_key,cpi.transliterated_title,cpi.arabic_title,c.created_at"
	return fields, join, Where, groupBy
}
func GetMovieTrailerQuery(language string) (string, string, string) {
	fields := "vt.id,vt.content_variance_id as variance_id,c.id as content_id,c.content_tier,vt.video_trailer_id as video_id,vt.duration as length"
	if language == "en" {
		fields += ",vt.english_title as title"
	} else {
		fields += ",vt.arabic_title as title"
	}
	join := "join content_variance cv on cv.id = vt.content_variance_id join content c on c.id = cv.content_id join playback_item pi2 on pi2.id =cv.playback_item_id join content_rights cr on cr.id = pi2.rights_id join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.content_key =? and crc.country_id =?"
	return fields, join, Where
}
func GetSeasonTrailerQuery(language string) (string, string, string) {
	fields := "vt.id,vt.season_id as variance_id,c.id as content_id,c.content_tier,vt.video_trailer_id as video_id,vt.duration as length"
	if language == "en" {
		fields += ",vt.english_title as title"
	} else {
		fields += ",vt.arabic_title as title"
	}
	join := "join season s on s.id = vt.season_id join content c on c.id = s.content_id join content_rights cr on cr.id = s.rights_id join content_rights_country crc on crc.content_rights_id =cr.id"
	Where := "c.content_key =? and crc.country_id =?"
	return fields, join, Where
}
func ContentTagsQuery() (string, string, string, string) {
	fields := "tdt.name"
	join := "join content_tag ct on ct.textual_data_tag_id =tdt.id join content c on c.tag_info_id = ct.tag_info_id"
	Where := "c.id =?"
	groupBy := "tdt.name"
	return fields, join, Where, groupBy
}
func MeadiaObjectQuery(language string) string {
	fields := "atci.age_group::text as age_rating,lower(c.content_type) as content_type,REPLACE (lower(cpi.transliterated_title), ' ', '_') as friendly_url,c.content_key as id,c.created_at as insertedAt,min(pi1.duration) as length,min(pi1.video_content_id) as video_id,c.id as content_id,c.content_tier"
	if language == "en" {
		fields += ",cpi.transliterated_title as title"
	} else {
		fields += ",cpi.arabic_title as title"
	}
	query := "SELECT " + fields + ",min(cv.id::text) as variance_id from content c join content_primary_info cpi on cpi.id=c.primary_info_id join about_the_content_info atci on atci.id=c.about_the_content_info_id join content_variance cv on cv.content_id =c.id join playback_item pi1 on pi1.id =cv.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and c.content_key in(?) group by c.content_type,c.content_key,c.id,c.content_tier,c.created_at,cpi.transliterated_title,atci.age_group UNION SELECT " + fields + ",min(s.id::text) as variance_id from content c join content_genre cg on cg.content_id =c.id join content_subgenre cs on cs.content_genre_id =cg.id join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null and c.content_key in(?) group by c.content_type,c.content_key,c.id,c.content_tier,c.created_at,cpi.transliterated_title,atci.age_group"
	return query
}

func EpisodeMeadiaObjectQuery(language string) string {
	fields := "e.episode_key as id,c.content_key as series_id,s.number as season_number,atci.age_group as age_rating,s.season_key as season_id,e.number as episode_number,'episode' as content_type,REPLACE (lower(cpi.transliterated_title), ' ', '_') as friendly_url,c.content_key as id,e.episode_key,c.created_at as insertedAt,min(pi1.duration) as length,min(pi1.video_content_id) as video_id,c.id as content_id,c.content_tier,min(s.id::text) as variance_id,e.id as episode_id"
	if language == "en" {
		fields += ",cpi.transliterated_title as series_title,cpi.transliterated_title as title"
	} else {
		fields += ",cpi.arabic_title as series_title,cpi.arabic_title as title"
	}
	query := "SELECT " + fields + " from content c join content_genre cg on cg.content_id =c.id join content_subgenre cs on cs.content_genre_id =cg.id join season s on s.content_id =c.id join episode e on e.season_id =s.id join content_primary_info cpi on cpi.id=s.primary_info_id join about_the_content_info atci on atci.id=s.about_the_content_info_id join playback_item pi1 on pi1.id =e.playback_item_id join content_rights cr on cr.id=pi1.rights_id join content_translation ct on ct.id=pi1.translation_id full outer join content_rights_country crc on crc.content_rights_id =cr.id full outer join content_rights_plan crp on crp.rights_id =cr.id where c.status = 1 and c.deleted_by_user_id is null and ( pi1.scheduling_date_time  <= NOW() or pi1.scheduling_date_time is null) and (cr.digital_rights_start_date <=NOW() or cr.digital_rights_start_date is null) and (cr.digital_rights_end_date >=NOW() or cr.digital_rights_end_date is null) and s.status =1 and s.deleted_by_user_id is null and e.status =1 and e.deleted_by_user_id is null and e.episode_key in(?) group by c.content_type,c.content_key,c.id,c.content_tier,c.created_at,cpi.transliterated_title,atci.age_group,e.episode_key,s.number,s.season_key,e.number,e.id"
	return query
}

// LanguageOriginTypes is...
func LanguageOriginTypes(originType string) int {
	OriginTypesArray := map[string]int{"Original": 1, "Dubbed": 2, "Subtitled": 3}
	return OriginTypesArray[originType]
}

// DeviceIds is...
func DeviceIds(device string) int32 {
	deviceArray := map[string]int{"web": 0, "ios": 1, "android": 2, "appletv": 3, "smarttv": 4, "roku": 5, "xbox_one": 6, "playstation": 7, "special": 8, "android_tv": 9, "amazon_fire_tv": 10}
	return int32(deviceArray[device])
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

// PageTypes is...
func PageTypes(PageTypeID int) string {
	PageTypesArray := map[int]string{0: "VOD", 1: "Home", 8: "Settings", 16: "Favourites"}
	return PageTypesArray[PageTypeID]
}

// SliderTypes is...
func SliderTypes(SliderTypeID int) string {
	SliderTypesArray := map[int]string{1: "Layout A – Smart TV", 2: "Layout B - STV / Website / Apple TV", 3: "layout C - STV - Website - Apple TV"}
	return SliderTypesArray[SliderTypeID]
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

/* Alphanumaric RegEx */
func AlphaNumaricRegex(phonenumber string) bool {
	Re := regexp.MustCompile(`^[A-Z]{1}[a-zA-z\s]+[0-9\s]`)
	return Re.MatchString(phonenumber)
}

func NumberRegex(SliderName string) bool {
	Re := regexp.MustCompile(`[0-9]`)
	return Re.MatchString(SliderName)
}

/* Split String to slice(int) */
func JsonStringToIntSliceOrMap(data string) ([]int, error) {
	output := make([]int, 1000)
	err := json.Unmarshal([]byte(data), &output)
	if err != nil {
		return nil, err
	}
	sort.Ints(output)
	return output, nil
}

// validations for slider
type Invalidsslider struct {
	NameError              interface{} `json:"name,omitempty"`
	RedAreaPlaylistError   interface{} `json:"redAreaPlaylistId,omitempty"`
	GreenAreaPlaylistError interface{} `json:"greenAreaPlaylistId,omitempty"`
	BlackAreaPlaylistError interface{} `json:"blackAreaPlaylistId,omitempty"`
	PlaylistItems          interface{} `json:"playlistItems,omitempty"`
	PagesIds               interface{} `json:"pagesIds,omitempty"`
}
type FinalErrorResponseslider struct {
	Error       string         `json:"error"`
	Description string         `json:"description"`
	Code        string         `json:"code"`
	RequestId   string         `json:"requestId"`
	Invalid     Invalidsslider `json:"invalid,omitempty"`
}

type NameError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type RedAreaPlaylistError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type GreenAreaPlaylistError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type BlackAreaPlaylistError struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PlaylistItems struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
type PagesIds struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

/* Monitering if page,playlist or slider get modified or updated dirty count will increase in respected table */
type PageSync struct {
	PageId     string `json:"pageId"`
	DirtyCount int    `json:"dirtyCount"`
	PageKey    int    `json:"pageKey"`
}
type PlaylistSync struct {
	PlaylistId string `json:"playlistId"`
	DirtyCount int    `json:"dirtyCount"`
}
type SliderSync struct {
	SliderId   string `json:"sliderId"`
	DirtyCount int    `json:"dirtyCount"`
}

func PageSynching(PageId string, pageKey int, c *gin.Context) string {
	fdb := c.MustGet("FDB").(*gorm.DB)
	var pageSync PageSync
	page := fdb.Debug().Table("page_sync").Where("page_id=?", PageId).Find(&pageSync)
	totalcount := int(page.RowsAffected)
	if totalcount < 1 {
		pageSync.PageId = PageId
		pageSync.DirtyCount = 1
		pageSync.PageKey = pageKey
		if updateError := fdb.Debug().Table("page_sync").Create(&pageSync).Error; updateError != nil {
			fmt.Println(updateError)
			return "failure"
		}
	} else {
		pageSync.PageId = PageId
		pageSync.DirtyCount = pageSync.DirtyCount + 1
		pageSync.PageKey = pageKey
		if updateError := fdb.Debug().Table("page_sync").Where("page_id=?", PageId).Update(&pageSync).Error; updateError != nil {
			fmt.Println(updateError)
			return "failure"
		}
	}
	return "success"
}

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

func RequestLogRegister(reqBody string, c *gin.Context) {
	// func() {
	db := c.MustGet("SDB").(*gorm.DB)
	requestMethod := c.Request.Method
	requestURL := c.Request.Host + c.Request.URL.Path
	// requestBody := c.Request.Body
	// jsonData, err := ioutil.ReadAll(requestBody)
	// fmt.Println(jsonData)
	// if err != nil {
	// 	fmt.Println("error for request body conversion:", err)
	// }
	reqToken := c.Request.Header.Get("Authorization")
	if reqToken == "" {
		reqToken = c.Request.Header.Get("authorization")
	}
	type RequestLog struct {
		Id        int
		Method    string
		Url       string
		RawBody   string
		JsnBody   string
		UserToken string
	}
	var requestLog RequestLog
	requestLog.Method = requestMethod
	requestLog.Url = requestURL
	// requestLog.RawBody = requestBody
	fmt.Println(reqBody, "requestBody")
	requestLog.JsnBody = reqBody
	requestLog.UserToken = reqToken
	if c.Request.URL.Path != "/health" {
		if err := db.Debug().Table("request_log").Create(&requestLog).Error; err != nil {
			// l.JSON(c, http.StatusInternalServerError, err)
			// return
			fmt.Println("error in logging:", err)
		}
		fmt.Println("Register Log Id: ", requestLog.Id)
	}
	// }
}

func ClearRedisKeyForPages(pageKey string, c *gin.Context) {
	/*delete Redis keys for pages */
	rdb := c.MustGet("REDIS_CLIENT").(*redis.Client)
	fmt.Println("inside redis clear page level")
	searchPattern := c.Request.Host + pageKey + "*"
	if len(os.Args) > 1 {
		searchPattern = os.Args[1]
	}
	iter := rdb.Keys(searchPattern)
	for _, Rkey := range iter.Val() {
		rdb.Del(Rkey)
		fmt.Println("Redis key " + Rkey + "is deleted")
	}
}

func UploadFileToS3(buffer []byte, filename string) (string, error) {
	s, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})

	tempFileName := filename
	// config settingshere you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET")),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		StorageClass:         aws.String("STANDARD"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("Unable to upload %q, %v", tempFileName, err)
	} else {
		fmt.Printf("Successfully uploaded %q", tempFileName)
	}
	return tempFileName, err

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

func DeleteRedisDataWithKey(key string) error {
	req, err := http.NewRequest("DELETE", os.Getenv("REDIS_CACHE_URL")+"/"+key, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
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
