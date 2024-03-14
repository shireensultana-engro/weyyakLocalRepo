package common

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ses"
)

const USERID = "2f634603-ce5b-eb11-831d-020666e39080"

var (
	host       = "smtp.pepipost.com"
	username   = "weyyakpepi"
	password   = "90082f@f33a8b"
	portNumber = "587"
	from       = "weyyak@weyyak.com"
)

func GetUserFromToken(token string) string {
	//TODO - Decode of token and return user details(ID)
	//Return UserID
	// reqToken := r.Header.Get("Authorization")
	// splitToken := strings.Split(reqToken, "Bearer ")
	// reqToken = splitToken[1]

	userId := "2f634603-ce5b-eb11-831d-020666e39080"
	return userId
}

func PostCurlCall(method string, url string, data interface{}) ([]byte, error) {
	URL := url
	payloadBytes, _ := json.Marshal(data)
	body := bytes.NewReader(payloadBytes)
	req, _ := http.NewRequest(method, URL, body)
	req.Header.Add("content-type", "application/json")
	// req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	fmt.Println("errrrrrrr", err)
	defer res.Body.Close()
	response, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return response, error

	}
	return response, nil

}

func GetCurlCall(url string) []byte {

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	// defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body
}

func VerifyHashPassword(hashPassword string, password string, version int, salt_stored string) bool {
	if version == 1 {
		hashedPassword, err := base64.StdEncoding.DecodeString(hashPassword)
		if err != nil {
			fmt.Printf("Error decoding string: %s ", err.Error())
			return false
		}
		salt_dst := make([]byte, 0x10)
		hashedKey := make([]byte, 0x20)
		copy(salt_dst, hashedPassword[1:17])
		copy(hashedKey, hashedPassword[17:49])
		password_hash := pbkdf2.Key([]byte(password), salt_dst, 1000, 32, sha1.New)
		return bytes.Equal(password_hash, hashedKey)
	} else {
		salt_stored, _ := base64.StdEncoding.DecodeString(salt_stored)
		hashPassword, _ := base64.StdEncoding.DecodeString(hashPassword)
		password_hash := pbkdf2.Key([]byte(password), []byte(salt_stored), 1000, 49, sha1.New)
		return bytes.Equal(password_hash, hashPassword)
	}

}
func HashPassword(password string) (string, string) {
	salt, _ := GetRandomBytes(0x10)
	hashedPassword := pbkdf2.Key([]byte(password), salt, 1000, 49, sha1.New)
	return base64.StdEncoding.EncodeToString(hashedPassword), base64.StdEncoding.EncodeToString(salt)
}

// GetRandomBytes returns len random looking bytes
func GetRandomBytes(len int) ([]byte, error) {
	key := make([]byte, len)
	// TODO: rand could fill less bytes then len
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.Wrap(err, "error getting random bytes")
	}
	return key, nil
}

func SendMail(to string, message string, subject string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// subject := "Subject: Welcome to Weyyak!"
	from := "weyyak@weyyak.com"
	password := "90082f@f33a8b"
	toList := []string{to}
	host := "smtp.pepipost.com"
	port := "587"
	// msg := "Hello email from golang api"
	body := []byte(subject + "\n" + mime + "\n" + message)
	auth := smtp.PlainAuth("", "weyyakpepi", password, host)
	err := smtp.SendMail(host+":"+port, auth, from, toList, body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully sent mail to all user in toList")
	return nil
}

func AwsSendMail(to string, message string, subject string) string {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("SES_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("SES_ID"),     // id
			os.Getenv("SES_SECRET"), // secret
			""),                     // token can be left blank for now
	})
	fmt.Println(sess)

	svc := ses.New(sess)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			// CcAddresses: []*string{
			// 	aws.String("recipient3@example.com"),
			// },
			ToAddresses: []*string{
				aws.String(to),
				// aws.String("recipient2@example.com"),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message),
				},
				// Text: &ses.Content{
				// 	Charset: aws.String("UTF-8"),
				// 	Data:    aws.String("This is the message body in text format."),
				// },
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		// ReturnPath:    aws.String(""),
		// ReturnPathArn: aws.String(""),
		Source: aws.String("support@weyyak.com"),
		// SourceArn:     aws.String(""),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		// 	if aerr, ok := err.(awserr.Error); ok {
		// 		switch aerr.Code() {
		// 		case ses.ErrCodeMessageRejected:
		// 			fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		// 		case ses.ErrCodeMailFromDomainNotVerifiedException:
		// 			fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		// 		case ses.ErrCodeConfigurationSetDoesNotExistException:
		// 			fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		// 		case ses.ErrCodeConfigurationSetSendingPausedException:
		// 			fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
		// 		case ses.ErrCodeAccountSendingPausedException:
		// 			fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
		// 		default:
		// 			fmt.Println(aerr.Error())
		// 		}
		// } else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error(), "herrrr")
		// }
		return err.Error()
	}

	fmt.Println(result)
	return ""
}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func RegEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
func RegName(nickName string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
	return Re.MatchString(nickName)
}
func RegMobile(phonenumber string) bool {
	Re := regexp.MustCompile(`^\+[1-9]{1}[0-9]{3,14}$`)
	return Re.MatchString(phonenumber)
}

// func IsUUID(uuid string) bool {
// 	Re := regexp.MustCompile(`/^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$/i`)
// 	return Re.MatchString(uuid)
// }

func ValidTime(time string) bool {
	Re := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2}.\d{9}) (\+\d{4}).(\w+)`)
	return Re.MatchString(time)
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("UDB").(*gorm.DB)
		reqToken := c.Request.Header.Get("Authorization")
		replacer := strings.NewReplacer("bearer", "Bearer", "BEARER", "Bearer")
		reqToken = replacer.Replace(reqToken)
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) == 2 {
			c.Set("AuthorizationRequired", 0)
			reqToken = splitToken[1]
		}
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
		fmt.Println(reqToken, response, "ggggggggggg")
		if response.Userid != "" {
			c.Set("userid", response.Userid)
			c.Set("device_id", response.DeviceId)
			c.Set("device_name", response.DeviceName)
			c.Set("device_platform", response.DevicePlatform)
			c.Set("language_id", response.LanguageId)
			c.Set("is_back_office_user", response.IsBackOfficeUser)
			c.Next()
			fmt.Println(reqToken, response.DeviceId, "ggggggggggg")
		} else {
			url := "https://uat-api.weyyak.z5.com/v1/ar/oauth2/tokendata?access_token=" + reqToken
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

//Countrys is...
func Countrys(country string) int32 {
	countryArray := map[string]int{"AF": 4, "AQ": 10, "DZ": 12, "AS": 16, "AD": 20, "AO": 24, "AG": 28, "AZ": 31, "AR": 32, "AU": 36, "AT": 40, "BS": 44, "BH": 48, "BD": 50, "AM": 51, "BB": 52, "BE": 56, "BM": 60, "BT": 64, "BO": 68, "BA": 70, "BW": 72, "BV": 74, "BR": 76, "BZ": 84, "IO": 86, "SB": 90, "VG": 92, "BN": 96, "BG": 100, "MM": 104, "BI": 108, "BY": 112, "KH": 116, "CM": 120, "CA": 124, "CV": 132, "KY": 136, "CF": 140, "LK": 144, "TD": 148, "CL": 152, "CN": 156, "TW": 158, "CX": 162, "CC": 166, "CO": 170, "KM": 174, "YT": 175, "CG": 178, "CD": 180, "CK": 184, "CR": 188, "HR": 191, "CU": 192, "CY": 196, "CZ": 203, "BJ": 204, "DK": 208, "DM": 212, "DO": 214, "EC": 218, "SV": 222, "GQ": 226, "ET": 231, "ER": 232, "EE": 233, "FO": 234, "FK": 238, "GS": 239, "FJ": 242, "FI": 246, "AX": 248, "FR": 250, "GF": 254, "PF": 258, "TF": 260, "DJ": 262, "GA": 266, "GE": 268, "GM": 270, "PS": 275, "DE": 276, "GH": 288, "GI": 292, "KI": 296, "GR": 300, "GL": 304, "GD": 308, "GP": 312, "GU": 316, "GT": 320, "GN": 324, "GY": 328, "HT": 332, "HM": 334, "VA": 336, "HN": 340, "HK": 344, "HU": 348, "IS": 352, "IN": 356, "ID": 360, "IR": 364, "IQ": 368, "IE": 372, "IL": 376, "IT": 380, "CI": 384, "JM": 388, "JP": 392, "KZ": 398, "JO": 400, "KE": 404, "KP": 408, "KR": 410, "KW": 414, "KG": 417, "LA": 418, "LB": 422, "LS": 426, "LV": 428, "LR": 430, "LY": 434, "LI": 438, "LT": 440, "LU": 442, "MO": 446, "MG": 450, "MW": 454, "MY": 458, "MV": 462, "ML": 466, "MT": 470, "MQ": 474, "MR": 478, "MU": 480, "MX": 484, "MC": 492, "MN": 496, "MD": 498, "ME": 499, "MS": 500, "MA": 504, "MZ": 508, "OM": 512, "NA": 516, "NR": 520, "NP": 524, "NL": 528, "CW": 531, "AW": 533, "SX": 534, "BQ": 535, "NC": 540, "VU": 548, "NZ": 554, "NI": 558, "NE": 562, "NG": 566, "NU": 570, "NF": 574, "NO": 578, "MP": 580, "UM": 581, "FM": 583, "MH": 584, "PW": 585, "PK": 586, "PA": 591, "PG": 598, "PY": 600, "PE": 604, "PH": 608, "PN": 612, "PL": 616, "PT": 620, "GW": 624, "TL": 626, "PR": 630, "QA": 634, "RE": 638, "RO": 642, "RU": 643, "RW": 646, "BL": 652, "SH": 654, "KN": 659, "AI": 660, "LC": 662, "PM": 666, "VC": 670, "SM": 674, "ST": 678, "SA": 682, "SN": 686, "RS": 688, "SC": 690, "SL": 694, "SG": 702, "SK": 703, "VN": 704, "SI": 705, "SO": 706, "ZA": 710, "ZW": 716, "ES": 724, "SS": 728, "SD": 729, "EH": 732, "SR": 740, "SJ": 744, "SZ": 748, "SE": 752, "CH": 756, "SY": 760, "TJ": 762, "TH": 764, "TG": 768, "TK": 772, "TO": 776, "TT": 780, "AE": 784, "TN": 788, "TR": 792, "TM": 795, "TC": 796, "TV": 798, "UG": 800, "UA": 804, "MK": 807, "EG": 818, "GB": 826, "GG": 831, "JE": 832, "IM": 833, "TZ": 834, "US": 840, "VI": 850, "BF": 854, "UY": 858, "UZ": 860, "VE": 862, "WF": 876, "WS": 882, "YE": 887, "ZM": 894}
	return int32(countryArray[country])
}

// func GenerateRandomString(length int) string {
// 	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// 	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
// 	b := make([]rune, 32)
// 	for i := range b {
// 		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
// 	}
// 	return strings.ToLower(string(b))
// }

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
	serverError.RequestId = EncodeToString(32)
	return serverError
}
func NotFoundErrorResponse() ServerError {
	var serverError ServerError
	serverError.Error = NOT_FOUND_ERROR
	serverError.Description = NOT_FOUND_ERROR_DESCRIPTION
	serverError.Code = NOT_FOUND_ERROR_CODE
	serverError.RequestId = EncodeToString(32)
	return serverError
}

const BULK_INSERT_LIMIT = 3000

// const USERID = "2f634603-ce5b-eb11-831d-020666e39080"
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

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func New() *Sender {
	auth := smtp.PlainAuth("", username, password, host)
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.auth, from, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	fmt.Println(withAttachments, "ddddddd.......")
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=us-ascii\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			fmt.Println(k, http.DetectContentType(v), "iiiiiiiiiiii")
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			// buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", "text/csv"))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func JsonStringToIntSliceOrMap(data string) ([]int, error) {

	output := make([]int, 1000)

	err := json.Unmarshal([]byte(data), &output)

	if err != nil {
		return nil, err
	}
	sort.Ints(output)
	return output, nil
}

//DeviceNames is...
func DeviceNames(device int) string {
	deviceArray := map[int]string{0: "web", 1: "ios", 2: "android", 3: "appletv", 4: "smarttv", 5: "roku", 6: "xbox_one", 7: "playstation", 8: "special", 9: "android_tv", 10: "amazon_fire_tv"}
	return deviceArray[device]
}
func DeviceName(device string) string {
	deviceArray := map[string]string{"0": "web", "1": "ios", "2": "android", "3": "appletv", "4": "smarttv", "5": "roku", "6": "xbox_one", "7": "playstation", "8": "special", "9": "android_tv", "10": "amazon_fire_tv"}
	return deviceArray[device]
}

func RegistrationSource(source int) string {
	registrationarr := map[int]string{1: "Email", 2: "Twitter", 3: "Facebook", 4: "Mobile", 5: "Apple"}
	return registrationarr[source]
}

func DupCount(list []string) map[string]int {

	duplicate_frequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}
