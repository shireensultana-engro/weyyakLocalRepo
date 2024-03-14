package config

import (
	"context"
	"encoding/json"
	"fmt"
	"frontend_service/common"
	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type HandlerService struct{}

// All the services should be protected by auth token
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// Setup Routes
	r.GET("/allconfig", hs.GetAllConfig)
	r.POST("/setconfig", hs.SetConfigGcp)
	r.GET("/getconfig", hs.GetConfig)
	r.GET("/config", hs.GetConfig)
}

// GetAllConfig -  fetches all config
// GET /allconfig
// @Summary Show a list of all country's
// @Description get list of all country's
// @Tags Config
// @Accept  json
// @Produce  json
// @Success 200 {array} object ConfigurationDetails
// @Router /allconfig [get]
func (hs *HandlerService) GetAllConfig(c *gin.Context) {
	db := c.MustGet("FCDB").(*gorm.DB)
	var config ApplicationSetting
	if err := db.Where("name='ConfigEndpointResponseBody'").Find(&config).Error; err != nil {
		serverError := common.ServerErrorResponse("ar")
		c.JSON(http.StatusInternalServerError, gin.H{"Message": serverError})
		return
	}

	var formated ConfigurationDetails
	if config.Value != "" {
		data := fmt.Sprintf("%v", config.Value)
		if err := json.Unmarshal([]byte(data), &formated); err != nil {
			serverError := common.ServerErrorResponse("ar")
			c.JSON(http.StatusInternalServerError, gin.H{"Message": serverError})
			return
		}
	}
	c.JSON(http.StatusOK, formated)
	return
}

func (hs *HandlerService) SetConfig(c *gin.Context) {
	var request interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "err in binding")
		return
	}
	jsonString, _ := json.Marshal(&request)
	filename := "configprodmigration.json"
	ioutil.WriteFile(filename, jsonString, os.ModePerm)
	common.UploadFileToS3(jsonString, filename)
}

func (hs *HandlerService) GetConfig(c *gin.Context) {
	var result interface{}
	// url changed to env variable
	// url, err := common.GetCurlCall("https://z5content-uat.s3.ap-south-1.amazonaws.com/configuat.json")
	url, err := common.GetCurlCall(os.Getenv("GCP_URLFORCONFIG"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println("err11", err)
	fmt.Println("yelll-", string(url))
	err = json.Unmarshal(url, &result)
	fmt.Println("err2", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, result)
}

func (hs *HandlerService) SetConfigGcp(c *gin.Context) {
	var request interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	// Marshal the request into JSON
	jsonString, err := json.Marshal(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	filename := "configprodmigration.json"

	// Storing to GCP Cloud Storage
	if err := UploadFileToGcp(jsonString, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	//storing in DB
	// db := c.MustGet("FCDB").(*gorm.DB)
	// var applicationSetting ApplicationSetting
	// applicationSetting.Value = string(jsonString)
	// if err := db.Model(&applicationSetting).Where("name='ConfigEndpointResponseBody'").Update(&applicationSetting).Error; err != nil {
	// 	c.JSON(http.StatusCreated, gin.H{"message": "config updated"})
	// 	return
	// }

}

func UploadFileToGcp(buffer []byte, filename string) error {
	ctx := context.Background()

	// Set up the GCP Cloud Storage client
	client, gcperr := getGCPClient()
	if gcperr != nil {
		fmt.Println("from gcp Connection", gcperr)
		// return gcperr
	}
	defer client.Close()
	bucketName := os.Getenv("BUCKET_NAME")

	// Open a GCP Cloud Storage handle to upload to the bucket
	wc := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	wc.ContentType = "application/json"

	// Write the buffer to GCP Cloud Storage
	if _, err := wc.Write(buffer); err != nil {
		return err
	}

	// Close the GCP Cloud Storage handle
	err := wc.Close()
	if err != nil {
		fmt.Println("Error closing GCP Cloud Storage handle:", err)
		return err
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
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC1nukPNB8e8iXM\nMhHr1iuvSGwienHPd4avprk0yUIXAlBZzwvbcK9i8V5yfpZzz6RUQcwshPOs5k9r\n3hBMy7zTGiZeyh1tPHdSumn3c4o7vL90RivGKff0VFvbPk4GdcuUFrEOJEH5gMS3\nYBTJtKKhOxKK3hqG/e0WQVjPJidfZwIKFthq9+z+d/4GMpldJAI3CPRpw9l+xzoC\n+vOueZ0aUCaSMvLKqgVzsKp5+YxGAEZbdxeYPMZGffJlZVedBwFNnyELBL8uKmVi\nmnCABwMjCTRwL3bPSgJ9mHLm2FiIK3heJ6Tg5HFjjIHIrxcdbVG57lKoXOt2wKed\n23l24T9tAgMBAAECggEAAvM8+unWbG6qjzmvLPtn1kzLpXEoEEEd8ssxMqJIqCOM\nLHCGOubJnZXZ4evNMbH3BcjHirUcWTvluUW2Rh4GiA/KIdEKIdoXL1bzORTMvG7d\nhoOI/69agNtAgwIp6ZTO+K24QODQnBrNPtccJ7cXaamqoFI4XgsHc7Q2jsfNC2bp\nAaIi4ZZHLhhQf94KFmfqOsVMhX7nmaBjaVZrpIfSM+5g0ESKbYbaLgdg/yVwpbdQ\nrLjolOvZw5r8e8ZVstdfU/GwihHuNsbgTbU511IeUYd+YmxoCZ1fkJq0Xf0uJ6Cy\nz1byOXfFOps8RurZhR1hkUknfeBaTBGVrujlHcoV6wKBgQDgBPKgOlzEWh/I5Epz\nviJZKa3TTJS2kinIGDYjbiJNp5NucQZfkJu0xBn0vztwmFnIdYnIOV7kiWuWulVM\nzjC3KoSBiC+GVGjABukuU/dlcWpbSttRuKg90gJ/gOtF4FuYLWhZkJGM2iCndanv\nkFmylCMoq6aiPnC73VGX2mfvpwKBgQDPjHTAzU1RMWiymF9yO25RZ6jzcyB6hPXf\n2NG2YJ0luM41pZMx5DlRFi8ky6YGK2gFwNsyBBXhRh5AdciGD96FwbqtLEv6eKCn\nC1BxZceYAdA/P6Aa6h/4Wv7J67THKEYbwzGgPYE5Js+jmJXG3tKhtPqnlWa7Zx6H\nLS5uNM7aywKBgF4z1m9oe3AaUflhfqlzT/BcpXsQXgz0I9u/yqxVeNlc2ZN8tehj\n4AZA3IVeETnE5yRzwM/QyEWkP/jvPEWDA1tS5sutoAaF4lK11UKlDoi7C7V+IgIY\ne68ba++AH++PbBTvK01WjM5FP6wLv70832tH/gzxOa5KQY/OfqwzrLdLAoGBAMoa\nAoLQJ+7dRw9KEv9AYf9BCqLtw32qxWYRUrzePYhC+gIBVmEp1KpiCMwyxluRnvyj\nPI7qrYes6L5qMzZgc5YZ/LauwNmI5x9ihBW4P3CEq407XqN2wmTr7tke/e1FCWf1\nXfiki5XkdiLe7VI3HjI68i2H7P6lvnNxCppkL92bAoGAERQ76Setyiz6VOiy+AJ6\nRVcQg2PDpoP0woWj7AOstCwbP91AT1h0Wq/aXRm1lk10Yvq0zm8RNzAUkPrfqLwx\n1pC5SQPrut0h+RZaQPRzUJERxfzMzej/WStGh51E9gRyFdSKfL/iOQ6ZT8PysEwy\ntcDR1sRoS24TmmtlgwyP91o=\n-----END PRIVATE KEY-----\n",
		"client_email":                os.Getenv("CLIENT_EMAIL"),
		"client_id":                   os.Getenv("CLIENT_ID"),
		"auth_uri":                    os.Getenv("AUTH_URI"),
		"token_uri":                   os.Getenv("TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("CLIENT_X509_CERT_URL"),
		"universe_domain":             os.Getenv("UNIVERSE_DOMAIN"),
		// "access_key":    "GOOGYK5RO4TBNDZF3EAVZRPN",
		// "secret": "WvHkEry+SLLLoeRu/25aJTi6Zj5Ii68Mi5UDQ3rS",
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
