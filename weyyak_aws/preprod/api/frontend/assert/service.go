package assert

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	r.GET(".well-known/assertlinks.json", hs.GetAssertJSON)
}

func (hs *HandlerService) GetAssertJSON(c *gin.Context) {
	// Initialize an AWS session with your access keys and region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")), // Replace with your region
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("S3_ID"),     // id
			os.Getenv("S3_SECRET"), // secret
			"",
		),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Initialize an S3 service client
	svc := s3.New(sess)

	// Specify the S3 bucket and key for the JSON file you want to fetch
	bucket := os.Getenv("ASSERT_BUCKET")
	key := os.Getenv("ASSERT_FILE")

	// Retrieve the JSON file from the S3 bucket
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Convert the JSON file to bytes
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bytes := buf.Bytes()

	// Return the JSON data as a response with the appropriate content type header
	c.Data(http.StatusOK, "application/json", bytes)
}
