package channel

import (
	// "content/common"

	"bytes"
	l "content/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
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

func programLogoUploadGcp(logopath string, channelName string, programName string) {
	bucketName := os.Getenv("BUCKET_NAME")

	item := logopath
	filetrim := strings.Split(item, "_")
	Destination := channelName + "/" + programName + "/" + filetrim[0]
	source := "temp/" + item // Assuming temp is a local directory.

	ctx := context.Background()

	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", gcperr)
	}
	defer client.Close()

	// Copy object from one bucket to another.
	src := client.Bucket(bucketName).Object(source)
	dst := client.Bucket(bucketName).Object(Destination)
	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		log.Fatalf("CopyObject failed: %v", err)
	}

	// Generate the public URL for the uploaded file.
	// url := "https://storage.googleapis.com/" + bucketName + "/" + Destination
	url := os.Getenv("IMAGERY_URL") + "/" + Destination

	// Don't worry about errors.
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	// Open a file for writing.
	file, err := os.Create(filetrim[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to dump the response body to the file.
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = SizeUploadFileToGcp(ctx, client, filetrim[0], channelName, programName, bucketName)
	if err != nil {
		fmt.Println("error in uploading size upload", err)
	}
	fmt.Println("Success!")
}

func SizeUploadFileToGcp(ctx context.Context, client *storage.Client, fileName string, channelName string, programName string, bucketName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
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

	for i := 0; i < len(sizeValue); i++ {
		gcsfile := sizeValue[i] + channelName + "/" + programName + "/" + fileName

		wc := client.Bucket(bucketName).Object(gcsfile).NewWriter(ctx)
		wc.ContentType = http.DetectContentType(buffer)
		wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

		_, err = wc.Write(buffer)
		if err != nil {
			return fmt.Errorf("Unable to upload %s: %v", fileName, err)
		}

		if err := wc.Close(); err != nil {
			return fmt.Errorf("Unable to close writer for %s: %v", fileName, err)
		}

		fmt.Printf("Successfully uploaded %q\n", fileName)
	}

	return nil
}

func getGCPClient() (*storage.Client, error) {
	data := map[string]interface{}{
		// "client_id":       "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com",
		// "client_secret":    "d-FL95Q19q7MQmFpd7hHD0Ty",
		// "quota_project_id": "engro-project-392708",
		// "refresh_token":    "1//0gCu2SwEAITTxCgYIARAAGBASNwF-L9IrXoW2jiRehyvfOj0yt3jnt5FXmYdlmkXXNIDjKzt5O1a3USJtclNE6sMSlr_W_Mw4xes",
		// "type":             "authorized_user",

		"type":                        os.Getenv("TYPE"),
		"project_id":                  os.Getenv("PROJECT_ID"),
		"private_key_id":              os.Getenv("PRIVATE_KEY_ID"),
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCrrOzb27NS3vj2\n/mspxBm/22giPXnCND0DYKLRcQuzsYv0elnbqmcGSJfCqs9C59PGJzRS6RV+ie/N\nhgflwUx3I3zMtd1fgVsjZFYYwmIDsaShkaK0B3eYWGXqYw8swG67qctk3GubNYNl\n04p6qIsUZgTpmH8Jjvn7abwntDvfxYbK9QVC2rjXzzA0pABfrgo3CVZJiXO8t4Ly\na8P5jkNY8Uze6Us+L0XxYL8T08vD09Cqde2kEEMBM4R6okdtgj3Vp8G287Dj0OUe\nMBGoj6YgsH32ZOXxxeQp003tMETaJbwQx6HTKr4N0CDwQnQkMbVUN2N8wOd4e0+2\nG+U36cGVAgMBAAECggEAFn+JDZ8TNyRleD5gs46G2VqFoRxxXSlqEuE9NTlyu8/k\nHtv8nrRhirSaFDbnsUWfE/QwqpTv7i9hhTZayUS1zVSR7GSrvZ0UNo/Vq1T+HKx7\n03i52+IGov54DL7X+ZjBFPLsPCxEJd5eI/Vpy9KpYg5PTSsLqv2udmulmYZzOktP\nYeV/qAaV/h/uQa+yTkxz9q0lixganx+ZSiC/3iTLwQLTI+Em8ayjVcIGQ/A9j6X1\nVCOxHBvy3bcIgZe+ZImwoWvko8ryaHWrdCKz7zVgXPZ9aT6B+VW0qqJGsHS0F5m7\nK0EC8fkdMlRufEiw6DChWUmspg7FYNW3fL7boAXemQKBgQDoOzEwr9khlO32ZXSs\nqIKRGNoL5pZPekVHPc/LI6713Vwg4g52xmtT0ZwgjkUB9QF4CimYVGLHytZG6P0G\nSBAdf4JMeeuBkJtmkXnYdJAlbwRNTiHWz409yAJ9hIyPafLZFKYxMLAqj95wnBxc\nMGq9accLaLIUtG8WGfSUrs5fwwKBgQC9PxTrAMl+ewm2O+a86du+BGsiy8fUvQZX\nJ9xayx9ARjEJXv1cgD4z59mQDn6gzBLrDcH+KY2ZSZUmvPof5LkXUlXXplJWh1Qj\nYvpMx2IOdu2OFFfydtyvq/JbXaEMrvUGU3+pvCF1e7Wxf+jlCTZM4yKwg1Ba9FyT\nCUaPlJFbxwKBgC0wv4y622TWh0voSEEE9Ytoq52fPGaw42ROme3svrInZjMb6jag\nu+fupRQMu077L1L9n0R+P06joPjhg8NCKKik1GUvYG2xBxx5eJ1vaVFvfgXRC3Ky\npsh78Egej/+kXVZy1zhBQja2ElIVfstNvKepOst0jxrKVceWO2rnbU9jAoGAdDH6\nNvxpuyXyZZjL6GwyRq5R1bCHRqC09uh7jKewzXcLfrR7HcOD7bzKQYAU0cfbScVN\nui9rSJX8ZSec794woxgjqt/tKEG5MG0CQAgftb/hxd3Jzg6bG6WYje6kBrSZr0Ov\nW9kuNgM6IPznU1FfrL+9OeG2gdIN0R3d3CSdR1sCgYBxpi1DXXBCXVeU1wIb9XBA\nwewiFSAabF/UtiF7CHkGSMN1lMe/R1AFKM8Irrqbbm0jl00BZ5fgVYV/wVaYZtDw\nPZQmGeO3yGi6FanLnBaxE/bKjk+RkaORM8QoaYGghX59TNoFzHNE1rF0w1lMdrlN\nnsFelOtLls3xNrtNNMxHeg==\n-----END PRIVATE KEY-----\n",
		"client_email":                os.Getenv("CLIENT_EMAIL"),
		"client_id":                   os.Getenv("CLIENT_ID"),
		"auth_uri":                    os.Getenv("AUTH_URI"),
		"token_uri":                   os.Getenv("TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("CLIENT_X509_CERT_URL"),
		"universe_domain":             os.Getenv("UNIVERSE_DOMAIN"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, jsonData, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		fmt.Println("Error creating credentials:", err)
	}
	client, err := storage.NewClient(ctx, option.WithCredentials(creds))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client, err
}
