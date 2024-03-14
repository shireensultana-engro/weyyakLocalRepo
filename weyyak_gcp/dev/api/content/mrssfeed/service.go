package mrssfeed

import (
	"content/common"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	r.GET("/mrss/feed", hs.mrssfeed)

}

//  mrssfeed
// GET /mrss/feed
// @Summary mrss feed
// @Description mrss feed
// @Tags mrss feed
// @AEngGenreept xml
// @Produce  xml
// @SuEngGenreess 200 {array} object c.xml
//@Param offset query string false "Offset"
// @Param limit query string false "Limit"
// @Router /mrss/feed [GET]
//Get all mrssfedd details
// func (hs *HandlerService) mrssfeed(c *gin.Context) {
// 	db := c.MustGet("DB").(*gorm.DB)
// 	var allcontentdetails []AllContentDetails
// 	var Thumbnail thumbnail
// 	var finalresult Rss
// 	var final []Rss
// 	var Channel []channel
// 	var part channel
// 	var Item []item
// 	var It item
// 	var Con content
// 	var href Href
// 	var rel Rel
// 	var english eng
// 	var arabic arb
// 	var limit, offset, current_page int64

// 	if c.Request.URL.Query()["limit"] != nil {
// 		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
// 	}
// 	if c.Request.URL.Query()["page"] != nil {
// 		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
// 	}
// 	if limit == 0 {
// 		limit = 50
// 	}
// 	offset = current_page * limit
// 	if err := db.Debug().Raw("select c.created_at,cpi.original_title as title,pi2.duration,s.number as season_number,e.number as episode_number,c.id as content_ID,sg.season_id as season_id,e.id as episode_id, string_agg(distinct g.english_name::text,', ') as english_genre,string_agg(distinct g.arabic_name::text,', ') as arabic_genre,string_agg(distinct s2.english_name::text, ', ') as english_subgenre,string_agg(distinct s2.arabic_name::text, ', ') as arabic_subgenre, cpi.original_title as english_title, cpi.arabic_title as arabic_title, c.content_type, ct.dubbing_language as language, c.content_key as content_key, e.episode_key as episode_key, c.modified_at from content c left join content_primary_info cpi on cpi.id = c.primary_info_id join season s on s.content_id = c.id join season_genre sg on sg.season_id = s.id join episode e on e.season_id = s.id join playback_item pi2 on pi2.id = e.playback_item_id join genre g on g.id = sg.genre_id join subgenre s2 on s2.genre_id = g.id join content_translation ct on ct.id = s.translation_id where c.content_tier = 2 group by c.created_at,cpi.original_title,pi2.duration,s.number,e.number,c.id,sg.season_id,e.id ,cpi.original_title,cpi.arabic_title,c.content_type, ct.dubbing_language, c.content_key, e.episode_key, c.modified_at union select c.created_at,cpi.original_title,pi2.duration,c.status,c.content_key,c.id,c.cast_id,c.music_id,string_agg( distinct g.english_name ::text,', ')as english_genre,string_agg(distinct g.arabic_name ::text,', ')as arabic_genre,string_agg(distinct s.english_name::text,', ') as english_subgenre,string_agg(distinct s.arabic_name::text,', ') as arabic_subgenre,cpi.original_title,cpi.arabic_title,c.content_type , ct.dubbing_language as language, c.content_key as content_id, c.status as episode_key, c.modified_at from content c join content_genre cg on cg.content_id = c.id join genre g on g.id = cg.genre_id join subgenre s on s.genre_id = cg.genre_id join content_primary_info cpi on cpi.id = c.primary_info_id join content_variance cv on cv.content_id = c.id join playback_item pi2 on pi2.id = cv.playback_item_id left join content_translation ct on ct.id = pi2.translation_id where content_tier = 1 group by c.created_at,cpi.original_title,pi2.duration,cpi.original_title,cpi.arabic_title,c.content_type,ct.dubbing_language,c.content_key,c.modified_at,c.status,c.content_key,c.id,c.cast_id,c.music_id,c.status ,cpi.alternative_title").Limit(limit).Offset(offset).Find(&allcontentdetails).Error; err != nil {
// 		c.XML(http.StatusInternalServerError, err)
// 		return
// 	}
// 	part.Title = "Weyyak Videos"
// 	part.Version = "2"
// 	rel.Rel = "prev"
// 	href.Rel = "next"
// 	rel.Href = "https://apiqafo.engro.in/mrss/feed?page=" + strconv.Itoa(int(current_page)-1)
// 	href.Href = "https://apiqafo.engro.in/mrss/feed?page=" + strconv.Itoa(int(current_page)+1)
// 	part.Href = href
// 	part.Rel = rel
// 	finalresult.Xsd = "http://www.w3.org/2001/XMLSchema"
// 	finalresult.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
// 	finalresult.Media = "http://search.yahoo.com/mrss "
// 	finalresult.Atom = "http://www.w3.org/2005/Atom"
// 	finalresult.OpenSearch = "http://a9.com/-/spec/opensearchrss/1.0"
// 	finalresult.Dfpvideo = "http://api.google.com/dfpvideo"
// 	finalresult.Version = "2.0"
// 	ax := "https://weyyak-content-dev.engro.in/"
// 	ay := "https://weyyak-dev.engro.in/en"
// 	d := "https://weyyak-dev.engro.in/ar"
// 	for _, value := range allcontentdetails {
// 		It.PubDate = value.CreatedAt
// 		It.Title = value.Title
// 		if value.ContentType == "Movie" {
// 			var final []keyvaluee
// 			Thumbnail.URL = ax + value.ContentId + "/poster-image"
// 			english.Duration = value.Duration
// 			arabic.Duration = value.Duration
// 			english.Url = ay + "/player/Movie/" + value.ContentKey
// 			arabic.Url = d + "/player/Movie/" + value.ContentKey
// 			Con.Eng = english
// 			Con.Arb = arabic
// 			someString := value.EnglishGenre
// 			res2 := strings.Split(someString, ",")
// 			var EngGenre keyvaluee
// 			for _, val := range res2 {
// 				EngGenre = keyvaluee{Key: "Weyyak_Genre", Value: val, Type: "string"}
// 				final = append(final, EngGenre)
// 			}
// 			someStrings := value.ArabicGenre
// 			res := strings.Split(someStrings, ",")
// 			var ArbGenre keyvaluee
// 			//var Genrearb []keyvaluee
// 			for _, val := range res {
// 				ArbGenre = keyvaluee{Key: "Weyyak_Genre", Value: val, Type: "string"}
// 				final = append(final, ArbGenre)
// 			}
// 			someStings := value.EnglishSubgenre
// 			res3 := strings.Split(someStings, ",")
// 			var EngSubGen keyvaluee
// 			for _, val := range res3 {
// 				EngSubGen = keyvaluee{Key: "Weyyak_Sub_Genre", Value: val, Type: "string"}
// 				final = append(final, EngSubGen)
// 			}
// 			someStrngs := value.ArabicSubgenre
// 			res1 := strings.Split(someStrngs, ",")
// 			var ArbSubGen keyvaluee
// 			for _, val := range res1 {
// 				ArbSubGen = keyvaluee{Key: "Weyyak_Sub_Genre", Value: val, Type: "string"}
// 				final = append(final, ArbSubGen)
// 			}
// 			EngShow := keyvaluee{Key: "Weyyak_Show_Name", Value: value.EnglishTitle, Type: "string"}
// 			final = append(final, EngShow)
// 			ArbShow := keyvaluee{Key: "Weyyak_Show_Name", Value: value.ArabicTitle, Type: "string"}
// 			final = append(final, ArbShow)
// 			ConType := keyvaluee{Key: "Weyyak_Content_Type", Value: value.ContentType, Type: "string"}
// 			final = append(final, ConType)
// 			Lang := keyvaluee{Key: "Weyyak_Language", Value: value.Language, Type: "string"}
// 			final = append(final, Lang)
// 			nesting := &dfpvideo{}
// 			nesting.Keyvalues = final
// 			out11, _ := xml.MarshalIndent(nesting, " ", "  ")
// 			var data1 data
// 			unmarshalerror := xml.Unmarshal(out11, &data1)
// 			if unmarshalerror != nil {
// 				fmt.Printf("unmarshal error %+v:", unmarshalerror)
// 			}
// 			It.Keyvalues = data1
// 			It.ContentId = value.ContentKey
// 		}
// 		if value.ContentType == "Series" {
// 			var final []keyvaluee
// 			Thumbnail.URL = ax + value.ContentId + "/" + value.SeasonId + "/" + value.EpisodeId + "/poster-image"
// 			english.Duration = value.Duration
// 			arabic.Duration = value.Duration
// 			english.Url = ay + "/player/episode/" + value.EpisodeKey
// 			arabic.Url = d + "/player/episode/" + value.EpisodeKey
// 			Con.Eng = english
// 			Con.Arb = arabic
// 			var seasonnumber string
// 			seasonnumber = value.SeasonNumber
// 			Season := keyvaluee{Key: "Weyyak_Season", Value: seasonnumber, Type: "int"}
// 			final = append(final, Season)
// 			Episode := keyvaluee{Key: "Weyyak_Episode", Value: value.EpisodeNumber, Type: "int"}
// 			final = append(final, Episode)
// 			someString := value.EnglishGenre
// 			res2 := strings.Split(someString, ",")
// 			var EngGenre keyvaluee
// 			for _, val := range res2 {
// 				EngGenre = keyvaluee{Key: "Weyyak_Genre", Value: val, Type: "string"}
// 				final = append(final, EngGenre)
// 			}
// 			someStrings := value.ArabicGenre
// 			res := strings.Split(someStrings, ",")
// 			var ArbGenre keyvaluee
// 			//var Genrearb []keyvaluee
// 			for _, val := range res {
// 				ArbGenre = keyvaluee{Key: "Weyyak_Genre", Value: val, Type: "string"}
// 				final = append(final, ArbGenre)
// 			}
// 			someStings := value.EnglishSubgenre
// 			res3 := strings.Split(someStings, ",")
// 			var EngSubGen keyvaluee
// 			for _, val := range res3 {
// 				EngSubGen = keyvaluee{Key: "Weyyak_Sub_Genre", Value: val, Type: "string"}
// 				final = append(final, EngSubGen)
// 			}
// 			someStrngs := value.ArabicSubgenre
// 			res1 := strings.Split(someStrngs, ",")
// 			var ArbSubGen keyvaluee
// 			for _, val := range res1 {
// 				ArbSubGen = keyvaluee{Key: "Weyyak_Sub_Genre", Value: val, Type: "string"}
// 				final = append(final, ArbSubGen)
// 			}
// 			EngShow := keyvaluee{Key: "Weyyak_Show_Name", Value: value.EnglishTitle, Type: "string"}
// 			final = append(final, EngShow)
// 			ArbShow := keyvaluee{Key: "Weyyak_Show_Name", Value: value.ArabicTitle, Type: "string"}
// 			final = append(final, ArbShow)
// 			ConType := keyvaluee{Key: "Weyyak_Content_Type", Value: value.ContentType, Type: "string"}
// 			final = append(final, ConType)
// 			Lang := keyvaluee{Key: "Weyyak_Language", Value: value.Language, Type: "string"}
// 			final = append(final, Lang)
// 			nesting := &dfpvideo{}
// 			nesting.Keyvalues = final
// 			out11, _ := xml.MarshalIndent(nesting, " ", "  ")
// 			var data1 data
// 			unmarshalerror := xml.Unmarshal(out11, &data1)
// 			if unmarshalerror != nil {
// 				fmt.Printf("unmarshal error %+v:", unmarshalerror)
// 			}
// 			It.Keyvalues = data1
// 			It.ContentId = value.EpisodeKey
// 		}
// 		It.LastModifiedDate = value.ModifiedAt
// 		timeval := value.Duration
// 		var num int = timeval / 600
// 		var a float64 = float64(num)
// 		var b int = int(a)
// 		res := math.Trunc(float64(b))
// 		var x float64 = res
// 		var y int = int(x)
// 		var K []int
// 		for i := 0; i <= y; i++ {
// 			z := i * 600
// 			K = append(K, z)
// 		}
// 		It.Cuepoints = arrayToString(K, ", ")
// 		It.Thumbnail = Thumbnail
// 		It.Content = Con
// 		Item = append(Item, It)
// 	}
// 	part.Item = Item
// 	Channel = append(Channel, part)
// 	finalresult.Channel = Channel
// 	final = append(final, finalresult)
// 	c.XML(http.StatusOK, final)
// }

func (hs *HandlerService) mrssfeed(c *gin.Context) {
	var result Rss
	var contentdetails []ContentDetails
	//var content []Content
	var limit, offset, current_page int64

	if c.Request.URL.Query()["limit"] != nil {
		limit, _ = strconv.ParseInt(c.Request.URL.Query()["limit"][0], 10, 64)
	}
	if c.Request.URL.Query()["page"] != nil {
		current_page, _ = strconv.ParseInt(c.Request.URL.Query()["page"][0], 10, 64)
	}
	if limit == 0 {
		limit = 100
	}
	if current_page <= 0 {
		current_page = 1
	}
	//below are static if need change
	result.Xsd = "http://www.w3.org/2001/XMLSchema"
	result.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	result.Media = "http://search.yahoo.com/mrss/"
	result.Atom = "http://www.w3.org/2005/Atom"
	result.OpenSearch = "http://a9.com/-/spec/opensearchrss/1.0/"
	result.Dfpvideo = "http://api.google.com/dfpvideo"
	result.Channel.Title = "Weyyak Videos"
	result.Version = "2.0"
	var links Link
	if current_page > 2 {
		for i := 0; i < 2; i++ {
			if i == 0 {
				links.Rel = "prev"
				links.Href = "https://apiqabo.engro.in/mrss/feed?page=" + strconv.Itoa(int(current_page-1))
				result.Channel.Link = append(result.Channel.Link, links)
			} else if i == 1 {
				links.Rel = "next"
				links.Href = "https://apiqabo.engro.in/mrss/feed?page=" + strconv.Itoa(int(current_page+1))
				result.Channel.Link = append(result.Channel.Link, links)
			}
		}
	} else {
		links.Rel = "next"
		links.Href = "https://apiqabo.engro.in/mrss/feed?page=" + strconv.Itoa(int(current_page+1))
		result.Channel.Link = append(result.Channel.Link, links)
	}
	result.Channel.Version = "2"
	offset = current_page * limit
	db := c.MustGet("DB").(*gorm.DB)
	db.Debug().Raw("select modified_at,created_at,content_key,content_type,transliterated_title,id,season_id,episode_id,duration,season_number,episode_number,arabic_title,original_language from (select c.modified_at,c.created_at,c.content_key,c.content_type,cpi.transliterated_title,c.id,null as season_id,null as episode_id,pi2.duration,0 as season_number,0 as episode_number,cpi.arabic_title,atci.original_language from content c join content_variance cv on cv.content_id = c.id  join content_primary_info cpi on cpi.id = c.primary_info_id join playback_item pi2 on pi2.id = cv.playback_item_id join about_the_content_info atci on atci.id = c.about_the_content_info_id where c.deleted_by_user_id is null and c.status = 1 and cv.status = 1 and cv.deleted_by_user_id is null union select c.modified_at,c.created_at,e.episode_key as content_key,'episode' as content_type,cpi.transliterated_title,c.id,s.id as season_id,e.id as episode_id,pi2.duration,s.number as season_number,e.number as episode_number,cpi.arabic_title,atci.original_language from content c join season s on s.content_id  =  c.id join episode e on e.season_id = s.id join playback_item pi2 on pi2.id  = e.playback_item_id join content_primary_info cpi on cpi.id = e.primary_info_id join about_the_content_info atci on atci.id = s.about_the_content_info_id) as foo order by modified_at desc").Limit(int(limit)).Offset(int(offset)).Find(&contentdetails)
	//var makestring string
	var contentids []string
	for _, val := range contentdetails {
		// j := i + 1
		// value := strconv.Itoa(j)
		// samplestring := " when cg.content_id ='" + val.Id + "' then " + value
		// makestring = makestring + samplestring
		contentids = append(contentids, val.Id)
	}
	var genredetails []GenreDetails
	var subgenredetails []SubGenreDetails
	db.Debug().Raw("select string_agg(g.english_name , ',' order by cg.order) as genres_english,string_agg(g.arabic_name , ',' order by cg.order) as genres_arabic,cg.content_id from content_genre cg join genre g on g.id = cg.genre_id where cg.content_id in (?) group by cg.content_id", contentids).Find(&genredetails)
	//order by case "+makestring+" end"
	db.Debug().Raw("select string_agg(s.english_name , ',' order by cg.order) as subgenres_english,string_agg(s.arabic_name , ',' order by cg.order) as subgenres_arabic,cg.content_id from content_genre cg join genre g on g.id = cg.genre_id join content_subgenre cs on cs.content_genre_id = cg.id join subgenre s on s.id = cs.subgenre_id where cg.content_id in (?) group by cg.content_id", contentids).Find(&subgenredetails)
	for _, val := range contentdetails {
		var temp Item
		temp.Title = val.TransliteratedTitle
		temp.ContentId = strconv.Itoa(val.ContentKey)
		temp.LastModifiedDate = val.ModifiedAt.Format(time.RFC1123Z)
		temp.PubDate = val.CreatedAt.Format(time.RFC1123Z)
		res := val.Duration / 600
		var cuestring string
		for i := 0; i <= res; i++ {
			if cuestring != "" {
				cuestring = cuestring + "," + strconv.Itoa(i*600)
			} else if cuestring == "" {
				cuestring = "0"
			}
		}
		temp.Cuepoints = cuestring
		// imagery
		if val.SeasonId != "" {
			temp.Thumbnail.URL = "https://weyyak-content-dev.engro.in/" + val.Id + "/" + val.SeasonId + "/" + val.EpisodeId + "/poster-image"

		} else {
			temp.Thumbnail.URL = "https://weyyak-content-dev.engro.in/" + val.Id + "/poster-image"
		}
		for i := 0; i < 2; i++ {
			var temps Content
			if i == 0 {
				temps.Duration = strconv.Itoa(val.Duration)
				temps.URL = "https://weyyakqa.z5.com/en/player/" + val.ContentType + "/" + strconv.Itoa(val.ContentKey) + "/"
			} else {
				temps.Duration = strconv.Itoa(val.Duration)
				temps.URL = "https://weyyakqa.z5.com/ar/player/" + val.ContentType + "/" + strconv.Itoa(val.ContentKey) + "/"
			}
			temp.Content = append(temp.Content, temps)
		}
		if val.SeasonId != "" {
			for i := 0; i < 2; i++ {
				var tempkeyvalues Keyvalues
				if i == 0 {
					tempkeyvalues.Key = "Weyyak_Season"
					tempkeyvalues.Value = strconv.Itoa(val.SeasonNumber)
					tempkeyvalues.Type = "int"
				} else {
					tempkeyvalues.Key = "Weyyak_Episode"
					tempkeyvalues.Value = strconv.Itoa(val.EpisodeNumber)
					tempkeyvalues.Type = "int"
				}
				temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
			}
		}
		for _, value := range genredetails {
			var tempkeyvalues Keyvalues
			if val.Id == value.ContentId {
				genreenglishtrim := strings.Split(value.GenresEnglish, ",")
				for _, val := range genreenglishtrim {
					tempkeyvalues.Key = "Weyyak_Genre"
					tempkeyvalues.Value = val
					tempkeyvalues.Type = "string"
					temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
				}
				genrearabictrim := strings.Split(value.GenresArabic, ",")
				for _, val := range genrearabictrim {
					tempkeyvalues.Key = "Weyyak_Genre"
					tempkeyvalues.Value = val
					tempkeyvalues.Type = "string"
					temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
				}
			}
		}
		for _, value := range subgenredetails {
			var tempkeyvalues Keyvalues
			if val.Id == value.ContentId {
				subgenreenglishtrim := strings.Split(value.SubgenresEnglish, ",")
				for _, val := range subgenreenglishtrim {
					tempkeyvalues.Key = "Weyyak_Sub_Genre"
					tempkeyvalues.Value = val
					tempkeyvalues.Type = "string"
					temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
				}
				subgenrearabictrim := strings.Split(value.SubgenresArabic, ",")
				for _, val := range subgenrearabictrim {
					tempkeyvalues.Key = "Weyyak_Sub_Genre"
					tempkeyvalues.Value = val
					tempkeyvalues.Type = "string"
					temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
				}
			}
		}
		for i := 0; i < 4; i++ {
			var tempkeyvalues Keyvalues
			if i == 0 {
				tempkeyvalues.Key = "Weyyak_Show_Name"
				tempkeyvalues.Value = val.TransliteratedTitle
				tempkeyvalues.Type = "string"
			} else if i == 1 {
				tempkeyvalues.Key = "Weyyak_Show_Name"
				tempkeyvalues.Value = val.ArabicTitle
				tempkeyvalues.Type = "string"
			} else if i == 2 {
				tempkeyvalues.Key = "Weyyak_Content_Type"
				tempkeyvalues.Value = val.ContentType
				tempkeyvalues.Type = "string"
			} else {
				tempkeyvalues.Key = "Weyyak_Language"
				tempkeyvalues.Value = common.OriginalLanguage(val.OriginalLanguage)
				tempkeyvalues.Type = "string"
			}
			temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
		}

		result.Channel.Item = append(result.Channel.Item, temp)
	}
	c.XML(http.StatusOK, result)
}

// func arrayToString(a []int, delim string) string {
// 	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
// }

// func (p keyvaluee) String() string {
// 	return fmt.Sprintf("keyvalue Key=%v, Value=%v, Type=%v",
// 		p.Key, p.Value, p.Type)
// }

// var genresarray []string
// length := len(genrearabictrim) + len(genreenglishtrim)
// for i := 1; i <= length; i++ {
// 	j := 0
// 	k := 0
// 	if i/2 != 0 {
// 		genresarray = append(genresarray, genreenglishtrim[j])
// 		j++
// 	} else {
// 		genresarray = append(genresarray, genrearabictrim[k])
// 		k++
// 	}
// }
// for _, val := range genresarray {
// 	tempkeyvalues.Key = "Weyyak_Genre"
// 	tempkeyvalues.Value = val
// 	tempkeyvalues.Type = "string"
// 	temp.Keyvalues = append(temp.Keyvalues, tempkeyvalues)
// }
