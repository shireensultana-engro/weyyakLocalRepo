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
	filename := "configstagemigration.json"
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

	filename := "configstagemigration.json"

	// Storing to GCP Cloud Storage
	if err := UploadFileToGcp(jsonString, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	//storing in DB
	db := c.MustGet("FCDB").(*gorm.DB)
	var applicationSetting ApplicationSetting
	applicationSetting.Value = string(jsonString)
	if err := db.Model(&applicationSetting).Where("name='ConfigEndpointResponseBody'").Update(&applicationSetting).Error; err != nil {
		c.JSON(http.StatusCreated, gin.H{"message": "config updated"})
		return
	}

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
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDGsE6CDHUcaOzC\nwQDlukKEPNKttBSXK0BgbQBNTMaO4sq9AfF9C71OyV1BuPaBKoXlQqSjrT8pEgy+\npXfNFHopGDIjloacFe2BjlV1e5w0FoanFc7F6wMoXtVRx8nkkprXYH913rryNoXW\nKjFa9D/8Uta4vZbkkTVEbZDyyDvLMYakgRf+uBqniG3WwyKTuyrrZHCLR/kpIXOq\nmTx/lip1WrCGckopmeRyeJK7yPLyb219A/L1gaPTZ6RtyA1t+L20xZcznAU2Heoa\nRg61tB5rdVR9V5GzP+5e1z77aGJ8D/YjsF33feXofWVG9DK1ur7eHydXqbI21s5n\nMdoe3hinAgMBAAECggEACPPexbtb94i+ylO5/5/x6VV4RL1thBY901p1/gWsmLmd\nWwDgabsCr68hFZoI+W76E4d6NanTw7z9eAWqrUHz8zAU6keZDyVpy0GnliCYvDwb\nmyG/bGmXfdUwFXxEG5mDgprlg2Ei2VEdnLXP/YIt8+ejpzVbvDbSyJ+LPXVKRHCF\nfPdpwCnybcWkeT3hiN8ovkomZiYr8bnt7sjuGOeQMjvJvq6sd0J+sXofgL43yfgD\n2GkqyXTR5kURonXMrX2uV6mMlLayaA+QvGRuAcbMxTMH5g03JsESKg6vh+PvEE6e\n/ki7G6+9xbPAOuRW2oxQmcR52QASrjwdEfpdpvMJSQKBgQDvn81m4voyySlVnzqZ\nbVo9xsNjtSpjcl7LwObJFQe2vpiUn3rj39M9biu/xLfXPLluFsQib7sZwRbBx6Lu\neln7UvLNzYZYSTUzgNeeV+VUp2OOIrWJ51oteamM+9B9W41lx7nYkvXLT5xKWszq\ntAfdSr8UaHWiF/WjuLaf2uRxdQKBgQDURFeQM8gVfiNBlC0CR4V08yMzeZD2jqIW\n5xu0XM1La4LpqsQm7CcOs8eUu3KA3NNm9XZpvTSbULyQgZw5G3eW1U5ND3St0Ddy\nYHXyTJKp3Q3HTO0AQcdU9AQISw+059dSeD+RaBlb7JP0osNZGDlTp5FZHPu9vK0J\n53yPXoeiKwKBgDTs28Ysvcw3yAxkReIbWAIrA37jRcB/Q1bHfXHOVkzTngm9i7wG\n9LYtvjX18hD1FZOuLZXZjb6reiZEvMTlezhaYsx354NacAi3HWiYy0s+SWvcWLJj\nyfQfWgaMm8kETp+7VF30X5uPMtrtYTM5nj8PQlL0m364wgVuR8/Y3fn9AoGACFfn\nWTOv2ahrmlhIrJ5DEKW97HgKyqYwmNXcsOo055ICQ00DCMSfhGRso9v6VDZZ2OIt\nFVrqhnBV+RgfG9+Ig9U+jqjc3Tgh9cz01eFMooCd2gecCTaMrzooLmtE4sd6HzO6\ny+xbktFpv2PmacoZ9r/PZsFM49hWtNz0eG4uxqECgYBvHsKOzUjj3E5SstP9JHMa\n/S2bpSebhhHtCKsODW/D6LdhMMMjJvVOilAsABnmX/ewFlcUi6xTxeJ9DYIWURID\n3o86noVLpigGzMCYihbAlxM1Od7UQJAleiIbVcu7JevGCvscFn5E0xz/+p6pMPZG\nHC44cZG1B99QIIpJUt9cqg==\n-----END PRIVATE KEY-----\n",
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
