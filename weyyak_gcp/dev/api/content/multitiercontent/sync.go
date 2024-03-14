package multitiercontent

import (
	"bytes"
	"content/common"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CreateSyncContent(c, content, request, "en")
// CreateSyncContent(c, content, request, "ar")

func GetCastDetails(c *gin.Context, Cast Cast, lang string) []SyncMovieCast {

	var cast []SyncMovieCast

	// MainActorId := Cast.MainActorId

	// if Cast.MainActorId != nil {
	cast = append(cast, SyncMovieCast{
		Details: "",
		Name:    GetCastName(c, Cast.MainActorId, lang),
		Type:    "mainActor",
	})
	// }

	// if Cast.MainActressId != nil {
	cast = append(cast, SyncMovieCast{
		Details: "",
		Name:    GetCastName(c, Cast.MainActressId, lang),
		Type:    "mainActress",
	})
	// }

	if Cast.Actors != nil {
		for _, ActorID := range Cast.Actors {
			// directors
			cast = append(cast, SyncMovieCast{
				Details: "",
				Name:    GetCastName(c, ActorID, lang),
				Type:    "actors",
			})
		}
	}

	if Cast.Directors != nil {
		for _, DirectorId := range Cast.Directors {
			// directors
			cast = append(cast, SyncMovieCast{
				Details: "",
				Name:    GetDirectorName(c, DirectorId, lang),
				Type:    "director",
			})
		}
	}

	if Cast.Writers != nil {
		for _, WriterID := range Cast.Writers {
			// directors
			cast = append(cast, SyncMovieCast{
				Details: "",
				Name:    GetWriterName(c, WriterID, lang),
				Type:    "writer",
			})
		}
	}

	return cast

}

func GetContentGenreToString(c *gin.Context, genreIds []string, lang string) []string {

	var genreDeatils []string

	db := c.MustGet("DB").(*gorm.DB)

	var genres Genres

	if err := db.Debug().Table("genre").Where("id in (?)", genreIds).Find(&genres).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, serverError)
		return genreDeatils
	}

	if lang == "en" {
		genreDeatils = append(genreDeatils, genres.EnglishName)
	} else if lang == "ar" {
		genreDeatils = append(genreDeatils, genres.ArabicName)
	} else if genres.EnglishName == "" {
		genreDeatils = append(genreDeatils, genres.ArabicName)
	} else if genres.ArabicName == "" {
		genreDeatils = append(genreDeatils, genres.EnglishName)
	}

	return genreDeatils
}

func GetContentSubGenreToString(c *gin.Context, subgenreIds []string, lang string) []string {

	var genreDeatils []string

	db := c.MustGet("DB").(*gorm.DB)

	var genres Genres

	if err := db.Debug().Table("subgenre").Where("id in (?)", subgenreIds).Find(&genres).Error; err != nil {
		return genreDeatils
	}

	if lang == "en" {
		genreDeatils = append(genreDeatils, genres.EnglishName)
	} else if lang == "ar" {
		genreDeatils = append(genreDeatils, genres.ArabicName)
	} else if genres.EnglishName == "" {
		genreDeatils = append(genreDeatils, genres.ArabicName)
	} else if genres.ArabicName == "" {
		genreDeatils = append(genreDeatils, genres.EnglishName)
	}

	return genreDeatils
}

func GetCastName(c *gin.Context, act string, lang string) string {

	db := c.MustGet("DB").(*gorm.DB)

	var cast SyncCast

	if err := db.Debug().Table("actor").Where("id = ?", act).Find(&cast).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, serverError)
		return ""
	}

	if lang == "en" {
		return cast.EnglishName
	} else if lang == "ar" {
		return cast.ArabicName
	}

	return ""
}

func GetDirectorName(c *gin.Context, act string, lang string) string {

	db := c.MustGet("DB").(*gorm.DB)

	var cast SyncCast

	if err := db.Debug().Table("director").Where("id = ?", act).Find(&cast).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, serverError)
		return ""
	}

	if lang == "en" {
		return cast.EnglishName
	} else if lang == "ar" {
		return cast.ArabicName
	}

	return ""
}

func GetWriterName(c *gin.Context, act string, lang string) string {

	db := c.MustGet("DB").(*gorm.DB)

	var cast SyncCast

	if err := db.Debug().Table("writer").Where("id = ?", act).Find(&cast).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, serverError)
		return ""
	}

	if lang == "en" {
		return cast.EnglishName
	} else if lang == "ar" {
		return cast.ArabicName
	}

	return ""
}

func GetTagName(c *gin.Context, tag string, lang string) string {

	db := c.MustGet("DB").(*gorm.DB)

	var tagName Tag

	if err := db.Debug().Table("textual_data_tag").Where("id = ?", tag).Find(&tagName).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, serverError)
		return ""
	}

	return tagName.Name
}

func GetImageForContent(ID string) []SyncMovieImage {

	return []SyncMovieImage{
		{
			ImageCategory: "posterImage",
			ImageURL: []string{
				fmt.Sprintf("https://weyyak-content-dev.engro.in/%s/poster-image", ID),
			},
		},
		{
			ImageCategory: "detailsPageBackground",
			ImageURL: []string{
				fmt.Sprintf("https://weyyak-content-dev.engro.in/%s/details-background", ID),
			},
		},
		{
			ImageURL: []string{
				fmt.Sprintf("https://weyyak-content-dev.engro.in/%s/mobile-details-background", ID),
			},
			ImageCategory: "mobileDetailsPageBackground",
		},
		{
			ImageCategory: "overlayPosterImage",
			ImageURL: []string{
				fmt.Sprintf("https://weyyak-content-dev.engro.in/%s/overlay-poster-image", ID),
			},
		},
	}

}

func CreateSyncContent(c *gin.Context, content Content, request OnetierContentRequest, lang string) {

	var (
		rights             []SyncMovieRights
		products           []int
		varianceTrailer    []VarianceTrailers
		tags               []string
		contentGenreIds    []string
		contentsubGenreIds []string
		MetaTitle          string
		MetaDescription    string
		Synopsis           string
	)

	// request.TextualData.ContentVariances

	for _, variance := range *request.TextualData.ContentVariances {

		products = variance.Products

		rights = append(rights, SyncMovieRights{
			Platform:           common.ContentPlatformsInt(variance.PublishingPlatforms),
			Rights:             common.ContentRightsTypesInt(variance.DigitalRightsType),
			SchedulingDateTime: variance.SchedulingDateTime,
			StartDate:          variance.DigitalRightsStartDate,
			EndDate:            variance.DigitalRightsEndDate,
			Location:           common.ContentLocationsInt(variance.DigitalRightsRegions),
			Plan:               common.ContentSubscriptionPlansInt(variance.SubscriptionPlans),
		})

		varianceTrailer = append(varianceTrailer, variance.VarianceTrailers...)

		// for _, a := range variance.VarianceTrailers {
		// 	varianceTrailer = append(varianceTrailer, a)
		// }

		// variance.VarianceTrailers

	}

	for _, i := range *request.TextualData.ContentGenres {
		// for _, i := range request.TextualData.ContentGenres {
		contentGenreIds = append(contentGenreIds, i.GenreId)
		contentsubGenreIds = append(contentsubGenreIds, i.SubgenresId...)
	}

	fmt.Println("contentGenreIds--->", contentGenreIds)
	fmt.Println("contentGenreIds--->", contentsubGenreIds)

	for _, tagID := range request.TextualData.TagInfo.Tags {
		// for _, tagID := range request.TextualData.TagInfo.Tags {
		tags = append(tags, GetTagName(c, tagID, lang))
	}

	if lang == "en" {
		MetaTitle = request.TextualData.SeoDetails.EnglishMetaTitle
		MetaDescription = request.TextualData.SeoDetails.EnglishMetaDescription
		Synopsis = request.TextualData.AboutTheContent.EnglishSynopsis
	} else if lang == "ar" {
		MetaTitle = request.TextualData.SeoDetails.ArabicMetaTitle
		MetaDescription = request.TextualData.SeoDetails.ArabicMetaDescription
		Synopsis = request.TextualData.AboutTheContent.ArabicSynopsis
	}

	Date := SyncMovieDate{
		Date: time.Now(),
	}

	payload := SyncMoviePayload{
		ID:                    content.Id,
		ContentType:           request.TextualData.PrimaryInfo.ContentType,
		Notes:                 request.TextualData.PrimaryInfo.Notes,
		Rights:                rights,
		AcquisitionDepartment: request.TextualData.AboutTheContent.AcquisitionDepartment,
		Status:                "1",
		Products:              products,
		Genre:                 GetContentGenreToString(c, contentGenreIds, lang),
		Name:                  request.TextualData.PrimaryInfo.OriginalTitle,
		SiteID:                "1",
		Trailers:              varianceTrailer,
		TransliteratedTitle:   request.TextualData.PrimaryInfo.TransliteratedTitle,
		Supplier:              request.TextualData.AboutTheContent.Supplier,
		Synopsis:              Synopsis,
		Cast:                  GetCastDetails(c, *request.TextualData.Cast, lang),
		IntroStart:            request.TextualData.PrimaryInfo.IntroStart,
		Key:                   content.ContentKey,
		Language:              common.OriginalLanguage(request.TextualData.AboutTheContent.OriginalLanguage),
		MetaTitle:             MetaTitle,
		ProductionHouse:       request.TextualData.AboutTheContent.ProductionHouse,
		Tags:                  tags,
		AgeGroup:              common.AgeRatings(request.TextualData.AboutTheContent.AgeGroup, lang),
		Images:                GetImageForContent(content.Id),
		OriginalLanguage:      lang,
		OutroStart:            request.TextualData.PrimaryInfo.OutroStart,
		ProductionYear:        request.TextualData.AboutTheContent.ProductionYear,
		SubGenre:              GetContentSubGenreToString(c, contentsubGenreIds, lang),
		CreatedAt:             Date,
		MetaDescription:       MetaDescription,
	}

	// Define the URL for the POST request
	url := "https://xms-dev.weyyak.com/api/config/v2/1/movie/" + lang

	_ = url
	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal JSON payload:", err)
		// return
	}

	fmt.Println("-------------------->")
	fmt.Println("-------------------->")
	fmt.Println("-------------------->")
	fmt.Println("-------------------->")
	fmt.Println("-------------------->", string(jsonPayload))
	fmt.Println("-------------------->")
	fmt.Println("-------------------->")
	fmt.Println("-------------------->")
	fmt.Println("-------------------->")

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Failed to send POST request:", err)
		// return
	}
	defer resp.Body.Close()

	// Check the response status code
	fmt.Println("Response Status Code:", resp.StatusCode)
}
