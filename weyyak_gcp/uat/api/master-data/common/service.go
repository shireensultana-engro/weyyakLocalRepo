package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
)

const BulkInsertLimit int = 3000

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

//LanguageOriginTypes is...
func LanguageOriginTypes(originType string) int {
	OriginTypesArray := map[string]int{"Original": 1, "Dubbed": 2, "Subtitled": 3}
	return OriginTypesArray[originType]
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
