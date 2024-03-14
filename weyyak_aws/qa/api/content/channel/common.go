package channel

import (
	// "content/common"

	"bytes"
	l "content/logger"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/* channel Thumbnail upload */
func (hs *HandlerService) UploadChannelThumbnailImage(c *gin.Context) {
	/*Authorization*/
	// if c.MustGet("AuthorizationRequired") == 1 || c.MustGet("is_back_office_user") == false {
	// 	l.JSON(c, http.StatusUnauthorized, gin.H{"message": "Authorization has been denied for this request."})
	// 	return
	// }
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Println("Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYOGUWMUMGAQLMW3U",                     // id
			"Jb0NV2eHwXAJg6UADb5vs3BgAyuUsvhgREi/hWRj", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Println("from s3 session", err)
		fmt.Println("Could not upload file -- session")
	}
	fileName, errr := UploadFileToS3(s, file, fileHeader, "programlogo")
	if errr != nil {
		fmt.Println("from s3 upload", errr)
		fmt.Println("Could not upload file")
	}
	fmt.Printf("Image uploaded successfully: %v", fileName)
	filetrim := strings.Split(fileName, "/")
	l.JSON(c, http.StatusOK, gin.H{"data": filetrim[1]})
	return
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, imagetype string) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	tempFileName := "temp/" + imagetype + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET")),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		StorageClass:         aws.String("STANDARD"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("Unable to upload %q, %v", tempFileName, err)
	}
	fmt.Printf("Successfully uploaded %q", tempFileName)
	return tempFileName, err
}

/*Uploade image Based on Page Id*/
func programLogoUpload(logopath string, channelName string, programName string) {
	bucketName := os.Getenv("S3_BUCKET")

	item := logopath
	filetrim := strings.Split(item, "_")
	Destination := channelName + "/" + programName + "/" + filetrim[0]
	source := bucketName + "/" + "temp/" + item
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("S3_ID"),     // id
			os.Getenv("S3_SECRET"), // secret
			""),                    // token can be left blank for now
	})
	/* Copy object from one directory to another*/
	svc := s3.New(s)
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(source),
		Key:        aws.String(Destination),
	}
	result, err := svc.CopyObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeObjectNotInActiveTierError:
				fmt.Println(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
	url := "https://" + bucketName + ".s3.ap-south-1.amazonaws.com/" + Destination
	// don't worry about errors
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filetrim[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	errorr := SizeUploadFileToS3(s, filetrim[0], channelName, programName)
	if errorr != nil {
		fmt.Println("error in uploading size upload", errorr)
	}
	fmt.Println("Success!")
}

// SizeUploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func SizeUploadFileToS3(s *session.Session, fileName string, chaneelname string, progranName string) error {
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
	sizeValue := [17]string{
		"100x100/",
		"150x150/",
		"200x200/",
		"250x250/",
		"270x270/",
		"300x300/",
		"420x420/",
		"450x450/",
		"570x570/",
		"600x600/",
		"620x620/",
		"800x384/",
		"800x800/",
		"811x811/",
		"900x900/",
		"2048x670/",
		"1125x240/",
	}
	var er error
	for i := 0; i < len(sizeValue); i++ {
		s3file := sizeValue[i] + chaneelname + "/" + progranName + "/" + fileName
		_, er = s3.New(s).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(os.Getenv("S3_BUCKET")),
			Key:                  aws.String(s3file),
			ACL:                  aws.String("public-read"),
			Body:                 bytes.NewReader(buffer),
			ContentLength:        aws.Int64(size),
			ContentType:          aws.String(http.DetectContentType(buffer)),
			ContentDisposition:   aws.String("attachment"),
			StorageClass:         aws.String("STANDARD"),
			ServerSideEncryption: aws.String("AES256"),
		})
		if er != nil {
			fmt.Println("Unable to upload", er)
		}
		fmt.Printf("Successfully uploaded %q", fileName)
	}
	return er
}
