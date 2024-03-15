package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const BulkInsertLimit int = 3000
const BlackPlaylistCount int = 1
const RedPlaylistCount int = 7
const GreenPlaylistCount int = 1

//Post curl call
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

//Get curl call
func GetCurlCall(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	// defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

//TODO:images upload functionality is not complete
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

//Countrys is...
func Countrys(country string) int32 {
	countryArray := map[string]int{"AF": 4, "AQ": 10, "DZ": 12, "AS": 16, "AD": 20, "AO": 24, "AG": 28, "AZ": 31, "AR": 32, "AU": 36, "AT": 40, "BS": 44, "BH": 48, "BD": 50, "AM": 51, "BB": 52, "BE": 56, "BM": 60, "BT": 64, "BO": 68, "BA": 70, "BW": 72, "BV": 74, "BR": 76, "BZ": 84, "IO": 86, "SB": 90, "VG": 92, "BN": 96, "BG": 100, "MM": 104, "BI": 108, "BY": 112, "KH": 116, "CM": 120, "CA": 124, "CV": 132, "KY": 136, "CF": 140, "LK": 144, "TD": 148, "CL": 152, "CN": 156, "TW": 158, "CX": 162, "CC": 166, "CO": 170, "KM": 174, "YT": 175, "CG": 178, "CD": 180, "CK": 184, "CR": 188, "HR": 191, "CU": 192, "CY": 196, "CZ": 203, "BJ": 204, "DK": 208, "DM": 212, "DO": 214, "EC": 218, "SV": 222, "GQ": 226, "ET": 231, "ER": 232, "EE": 233, "FO": 234, "FK": 238, "GS": 239, "FJ": 242, "FI": 246, "AX": 248, "FR": 250, "GF": 254, "PF": 258, "TF": 260, "DJ": 262, "GA": 266, "GE": 268, "GM": 270, "PS": 275, "DE": 276, "GH": 288, "GI": 292, "KI": 296, "GR": 300, "GL": 304, "GD": 308, "GP": 312, "GU": 316, "GT": 320, "GN": 324, "GY": 328, "HT": 332, "HM": 334, "VA": 336, "HN": 340, "HK": 344, "HU": 348, "IS": 352, "IN": 356, "ID": 360, "IR": 364, "IQ": 368, "IE": 372, "IL": 376, "IT": 380, "CI": 384, "JM": 388, "JP": 392, "KZ": 398, "JO": 400, "KE": 404, "KP": 408, "KR": 410, "KW": 414, "KG": 417, "LA": 418, "LB": 422, "LS": 426, "LV": 428, "LR": 430, "LY": 434, "LI": 438, "LT": 440, "LU": 442, "MO": 446, "MG": 450, "MW": 454, "MY": 458, "MV": 462, "ML": 466, "MT": 470, "MQ": 474, "MR": 478, "MU": 480, "MX": 484, "MC": 492, "MN": 496, "MD": 498, "ME": 499, "MS": 500, "MA": 504, "MZ": 508, "OM": 512, "NA": 516, "NR": 520, "NP": 524, "NL": 528, "CW": 531, "AW": 533, "SX": 534, "BQ": 535, "NC": 540, "VU": 548, "NZ": 554, "NI": 558, "NE": 562, "NG": 566, "NU": 570, "NF": 574, "NO": 578, "MP": 580, "UM": 581, "FM": 583, "MH": 584, "PW": 585, "PK": 586, "PA": 591, "PG": 598, "PY": 600, "PE": 604, "PH": 608, "PN": 612, "PL": 616, "PT": 620, "GW": 624, "TL": 626, "PR": 630, "QA": 634, "RE": 638, "RO": 642, "RU": 643, "RW": 646, "BL": 652, "SH": 654, "KN": 659, "AI": 660, "LC": 662, "PM": 666, "VC": 670, "SM": 674, "ST": 678, "SA": 682, "SN": 686, "RS": 688, "SC": 690, "SL": 694, "SG": 702, "SK": 703, "VN": 704, "SI": 705, "SO": 706, "ZA": 710, "ZW": 716, "ES": 724, "SS": 728, "SD": 729, "EH": 732, "SR": 740, "SJ": 744, "SZ": 748, "SE": 752, "CH": 756, "SY": 760, "TJ": 762, "TH": 764, "TG": 768, "TK": 772, "TO": 776, "TT": 780, "AE": 784, "TN": 788, "TR": 792, "TM": 795, "TC": 796, "TV": 798, "UG": 800, "UA": 804, "MK": 807, "EG": 818, "GB": 826, "GG": 831, "JE": 832, "IM": 833, "TZ": 834, "US": 840, "VI": 850, "BF": 854, "UY": 858, "UZ": 860, "VE": 862, "WF": 876, "WS": 882, "YE": 887, "ZM": 894}
	return int32(countryArray[country])
}

// Country count is
func CountryCount() int {
	// countryArray := map[string]int{"AF": 4, "AQ": 10, "DZ": 12, "AS": 16, "AD": 20, "AO": 24, "AG": 28, "AZ": 31, "AR": 32, "AU": 36, "AT": 40, "BS": 44, "BH": 48, "BD": 50, "AM": 51, "BB": 52, "BE": 56, "BM": 60, "BT": 64, "BO": 68, "BA": 70, "BW": 72, "BV": 74, "BR": 76, "BZ": 84, "IO": 86, "SB": 90, "VG": 92, "BN": 96, "BG": 100, "MM": 104, "BI": 108, "BY": 112, "KH": 116, "CM": 120, "CA": 124, "CV": 132, "KY": 136, "CF": 140, "LK": 144, "TD": 148, "CL": 152, "CN": 156, "TW": 158, "CX": 162, "CC": 166, "CO": 170, "KM": 174, "YT": 175, "CG": 178, "CD": 180, "CK": 184, "CR": 188, "HR": 191, "CU": 192, "CY": 196, "CZ": 203, "BJ": 204, "DK": 208, "DM": 212, "DO": 214, "EC": 218, "SV": 222, "GQ": 226, "ET": 231, "ER": 232, "EE": 233, "FO": 234, "FK": 238, "GS": 239, "FJ": 242, "FI": 246, "AX": 248, "FR": 250, "GF": 254, "PF": 258, "TF": 260, "DJ": 262, "GA": 266, "GE": 268, "GM": 270, "PS": 275, "DE": 276, "GH": 288, "GI": 292, "KI": 296, "GR": 300, "GL": 304, "GD": 308, "GP": 312, "GU": 316, "GT": 320, "GN": 324, "GY": 328, "HT": 332, "HM": 334, "VA": 336, "HN": 340, "HK": 344, "HU": 348, "IS": 352, "IN": 356, "ID": 360, "IR": 364, "IQ": 368, "IE": 372, "IL": 376, "IT": 380, "CI": 384, "JM": 388, "JP": 392, "KZ": 398, "JO": 400, "KE": 404, "KP": 408, "KR": 410, "KW": 414, "KG": 417, "LA": 418, "LB": 422, "LS": 426, "LV": 428, "LR": 430, "LY": 434, "LI": 438, "LT": 440, "LU": 442, "MO": 446, "MG": 450, "MW": 454, "MY": 458, "MV": 462, "ML": 466, "MT": 470, "MQ": 474, "MR": 478, "MU": 480, "MX": 484, "MC": 492, "MN": 496, "MD": 498, "ME": 499, "MS": 500, "MA": 504, "MZ": 508, "OM": 512, "NA": 516, "NR": 520, "NP": 524, "NL": 528, "CW": 531, "AW": 533, "SX": 534, "BQ": 535, "NC": 540, "VU": 548, "NZ": 554, "NI": 558, "NE": 562, "NG": 566, "NU": 570, "NF": 574, "NO": 578, "MP": 580, "UM": 581, "FM": 583, "MH": 584, "PW": 585, "PK": 586, "PA": 591, "PG": 598, "PY": 600, "PE": 604, "PH": 608, "PN": 612, "PL": 616, "PT": 620, "GW": 624, "TL": 626, "PR": 630, "QA": 634, "RE": 638, "RO": 642, "RU": 643, "RW": 646, "BL": 652, "SH": 654, "KN": 659, "AI": 660, "LC": 662, "PM": 666, "VC": 670, "SM": 674, "ST": 678, "SA": 682, "SN": 686, "RS": 688, "SC": 690, "SL": 694, "SG": 702, "SK": 703, "VN": 704, "SI": 705, "SO": 706, "ZA": 710, "ZW": 716, "ES": 724, "SS": 728, "SD": 729, "EH": 732, "SR": 740, "SJ": 744, "SZ": 748, "SE": 752, "CH": 756, "SY": 760, "TJ": 762, "TH": 764, "TG": 768, "TK": 772, "TO": 776, "TT": 780, "AE": 784, "TN": 788, "TR": 792, "TM": 795, "TC": 796, "TV": 798, "UG": 800, "UA": 804, "MK": 807, "EG": 818, "GB": 826, "GG": 831, "JE": 832, "IM": 833, "TZ": 834, "US": 840, "VI": 850, "BF": 854, "UY": 858, "UZ": 860, "VE": 862, "WF": 876, "WS": 882, "YE": 887, "ZM": 894}
	return 241
}

//CountryNames is...
func CountryNames(country int) string {
	countryArray := map[int]string{4: "AF", 10: "AQ", 12: "DZ", 16: "AS", 20: "AD", 24: "AO", 28: "AG", 31: "AZ", 32: "AR", 36: "AU", 40: "AT", 44: "BS", 48: "BH", 50: "BD", 51: "AM", 52: "BB", 56: "BE", 60: "BM", 64: "BT", 68: "BO", 70: "BA", 72: "BW", 74: "BV", 76: "BR", 84: "BZ", 86: "IO", 90: "SB", 92: "VG", 96: "BN", 100: "BG", 104: "MM", 108: "BI", 112: "BY", 116: "KH", 120: "CM", 124: "CA", 132: "CV", 136: "KY", 140: "CF", 144: "LK", 148: "TD", 152: "CL", 156: "CN", 158: "TW", 162: "CX", 166: "CC", 170: "CO", 174: "KM", 175: "YT", 178: "CG", 180: "CD", 184: "CK", 188: "CR", 191: "HR", 192: "CU", 196: "CY", 203: "CZ", 204: "BJ", 208: "DK", 212: "DM", 214: "DO", 218: "EC", 222: "SV", 226: "GQ", 231: "ET", 232: "ER", 233: "EE", 234: "FO", 238: "FK", 239: "GS", 242: "FJ", 246: "FI", 248: "AX", 250: "FR", 254: "GF", 258: "PF", 260: "TF", 262: "DJ", 266: "GA", 268: "GE", 270: "GM", 275: "PS", 276: "DE", 288: "GH", 292: "GI", 296: "KI", 300: "GR", 304: "GL", 308: "GD", 312: "GP", 316: "GU", 320: "GT", 324: "GN", 328: "GY", 332: "HT", 334: "HM", 336: "VA", 340: "HN", 344: "HK", 348: "HU", 352: "IS", 356: "IN", 360: "ID", 364: "IR", 368: "IQ", 372: "IE", 376: "IL", 380: "IT", 384: "CI", 388: "JM", 392: "JP", 398: "KZ", 400: "JO", 404: "KE", 408: "KP", 410: "KR", 414: "KW", 417: "KG", 418: "LA", 422: "LB", 426: "LS", 428: "LV", 430: "LR", 434: "LY", 438: "LI", 440: "LT", 442: "LU", 446: "MO", 450: "MG", 454: "MW", 458: "MY", 462: "MV", 466: "ML", 470: "MT", 474: "MQ", 478: "MR", 480: "MU", 484: "MX", 492: "MC", 496: "MN", 498: "MD", 499: "ME", 500: "MS", 504: "MA", 508: "MZ", 512: "OM", 516: "NA", 520: "NR", 524: "NP", 528: "NL", 531: "CW", 533: "AW", 534: "SX", 535: "BQ", 540: "NC", 548: "VU", 554: "NZ", 558: "NI", 562: "NE", 566: "NG", 570: "NU", 574: "NF", 578: "NO", 580: "MP", 581: "UM", 583: "FM", 584: "MH", 585: "PW", 586: "PK", 591: "PA", 598: "PG", 600: "PY", 604: "PE", 608: "PH", 612: "PN", 616: "PL", 620: "PT", 624: "GW", 626: "TL", 630: "PR", 634: "QA", 638: "RE", 642: "RO", 643: "RU", 646: "RW", 652: "BL", 654: "SH", 659: "KN", 660: "AI", 662: "LC", 666: "PM", 670: "VC", 674: "SM", 678: "ST", 682: "SA", 686: "SN", 688: "RS", 690: "SC", 694: "SL", 702: "SG", 703: "SK", 704: "VN", 705: "SI", 706: "SO", 710: "ZA", 716: "ZW", 724: "ES", 728: "SS", 729: "SD", 732: "EH", 740: "SR", 744: "SJ", 748: "SZ", 752: "SE", 756: "CH", 760: "SY", 762: "TJ", 764: "TH", 768: "TG", 772: "TK", 776: "TO", 780: "TT", 784: "AE", 788: "TN", 792: "TR", 795: "TM", 796: "TC", 798: "TV", 800: "UG", 804: "UA", 807: "MK", 818: "EG", 826: "GB", 831: "GG", 832: "JE", 833: "IM", 834: "TZ", 840: "US", 850: "VI", 854: "BF", 858: "UY", 860: "UZ", 862: "VE", 876: "WF", 882: "WS", 887: "YE", 894: "ZM"}
	return countryArray[country]
}

//DeviceIds is...
func DeviceIds(device string) int32 {
	deviceArray := map[string]int{"web": 0, "ios": 1, "android": 2, "appletv": 3, "smarttv": 4, "roku": 5, "xbox_one": 6, "playstation": 7, "special": 8, "android_tv": 9, "amazon_fire_tv": 10}
	return int32(deviceArray[device])
}

//DeviceNames is...
func DeviceNames(device int) string {
	deviceArray := map[int]string{0: "web", 1: "ios", 2: "android", 3: "appletv", 4: "smarttv", 5: "roku", 6: "xbox_one", 7: "playstation", 8: "special", 9: "android_tv", 10: "amazon_fire_tv"}
	return deviceArray[device]
}

//ProductNames is...
func ProductNames(product string) int32 {
	ProductArray := map[string]int{"Weyyak": 1, "WeyyakSouthAsian": 2, "Africa": 3, "AfricaSouthAsian": 4, "Europe": 5,
		"Global": 6, "Apac": 7}
	return int32(ProductArray[product])
}

//ContentRightsTypes is...
func ContentRightsTypes(rightType string) int32 {
	ContentRightsTypesArrays := map[string]int{"Avod": 1, "Vod": 2, "Svod": 3, "Tvod": 4}
	return int32(ContentRightsTypesArrays[rightType])
}

//PageTypes is...
func PageTypes(PageTypeID int) string {
	PageTypesArray := map[int]string{0: "VOD", 1: "Home", 8: "Settings", 16: "Favourites"}
	return PageTypesArray[PageTypeID]
}

//AgeRatings is...
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
			// c.JSON(http.StatusUnauthorized, gin.H{"Error": "Authentication Required."})
			// return
			reqToken = splitToken[1]
		} else {
			c.Set("AuthorizationRequired", 1)
		}

		// fmt.Println("kkkkkkkkkkk", reqToken)
		type Details struct {
			Userid         string
			DeviceId       string
			DeviceName     string
			DevicePlatform string
			LanguageId     int
		}
		var response Details
		// TODO - Redis need to cache
		if reqToken != "" {
			db.Raw("select data->>'UserID' as userid, data->>'DeviceID' as device_id, data->>'DeviceName' as device_name, data->>'DevicePlatform' as device_platform, data->>'LanguageId' as language_id  from oauth2_tokens ot where access = ? and (data->>'ExpiresAt')::timestamp >= now()", reqToken).Scan(&response)
		}
		if response.Userid != "" {
			c.Set("userid", response.Userid)
			c.Set("device_id", response.DeviceId)
			c.Set("device_name", response.DeviceName)
			c.Set("device_platform", response.DevicePlatform)
			c.Set("language_id", response.LanguageId)
			c.Next()
		} else {
			// TODO - API call to .NET api to confirm whether user token is valid or not
			c.Set("userid", "")
			c.Next()
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

//ServerError -- binding struct for error response
type ServerError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
	Code        string `json:"code"`
	RequestId   string `json:"requestId"`
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

//ContentLanguageOriginTypes is...
func LanguageOriginTypes(originType int) string {
	OriginTypesArray := map[int]string{1: "Original", 2: "Dubbed", 3: "Subtitled"}
	return OriginTypesArray[originType]
}

func RemoveDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//SliderTypes is...
func SliderTypes(SliderTypeID int) string {
	SliderTypesArray := map[int]string{1: "Layout A – Smart TV", 2: "Layout B - STV / Website / Apple TV", 3: "layout C - STV - Website - Apple TV"}
	return SliderTypesArray[SliderTypeID]
}

func DialectIdname(DialectId int, language string) string {
	DubbingDialectId := int(DialectId)
	AgeRatingArray := map[int]map[string]string{1: {"EnglishName": "Egyptian", "ArabicName": "مصري"}, 2: {"EnglishName": "Syrian", "ArabicName": "سوري"}, 3: {"EnglishName": "Khaliji", "ArabicName": "خليجي"}, 4: {"EnglishName": "Kuwaiti", "ArabicName": "كويتي"}, 5: {"EnglishName": "Emarati", "ArabicName": "إماراتي"}, 6: {"EnglishName": "Saudi", "ArabicName": "سعودي"}, 7: {"EnglishName": "Formal Arabic", "ArabicName": "فصحة"}, 8: {"EnglishName": "PAN ARAB", "ArabicName": "عربي متعدد اللهجات"}, 9: {"EnglishName": "Beduin", "ArabicName": "بدوي"}, 10: {"EnglishName": "Sa’edi", "ArabicName": "صعيدي"}, 11: {"EnglishName": "Iraqi", "ArabicName": "عراقي"}, 12: {"EnglishName": "Lebanese", "ArabicName": "لبناني"}, 13: {"EnglishName": "Jordanian", "ArabicName": "أردني"}, 14: {"EnglishName": "Maghrebi", "ArabicName": "مغربي"}}
	var DubbingDialectName string
	if language == "en" && DubbingDialectId > 0 {
		DubbingDialectName = AgeRatingArray[DubbingDialectId]["EnglishName"]
	} else if language == "ar" && DubbingDialectId > 0 {
		DubbingDialectName = AgeRatingArray[DubbingDialectId]["ArabicName"]
	}
	return DubbingDialectName
}

type RedisCacheResponse struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

type RedisErrorResponse struct {
	Error string `json:"error`
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
