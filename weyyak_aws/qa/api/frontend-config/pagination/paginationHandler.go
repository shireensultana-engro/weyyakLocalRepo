package pagination

import (
	"fmt"
	"math"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
)

//PaginationService struct
type PaginationService struct {
	Result interface{}
	Error  error
}

//GeneratePaginationRequest func
func GeneratePaginationRequest(c *gin.Context, limitNumber, offsetNumber, pageNumber int64, sort, orderby string) *Pagination {
	// default limit, page & sort parameter
	var limit, page, offset int
	limit = int(limitNumber)
	page = int(pageNumber)
	offset = int(offsetNumber)

	if sort == "" {
		sort = c.DefaultQuery("sort", "id")
	}
	if orderby == "" {
		orderby = c.DefaultQuery("order_by", "asc")
	}

	return &Pagination{Limit: limit, Offset: offset, Page: page, Sort: sort, OrderBy: orderby}
}

// PaginationCall function from url
func PaginationCall(c *gin.Context, pagination *Pagination, fields, where, tName string, tableName interface{}) (PaginationService, int) {
	// var contacts []Content
	db := c.MustGet("DB").(*gorm.DB)
	totalRows, totalPages, fromRow, toRow := 0, 0, 0, 0

	//offset := (pagination.Page - 1) * pagination.Limit
	// count all data
	errCount := db.Table(tName).Where(where).Count(&totalRows).Error

	if errCount != nil {
		return PaginationService{Error: errCount}, totalPages
	}
	// get data with limit, offset & order
	var offset, offset1 int
	if pagination.Offset != 0 {
		offset1 = pagination.Offset
		offset = pagination.Offset
	} else {
		offset1 = pagination.Page * pagination.Limit
		offset = pagination.Page * pagination.Limit
	}
	if totalRows <= offset {
		offset1 = totalRows
	}
	errFind := db.Debug().Table(tName).Select(fields).Where(where).Limit(pagination.Limit).Offset(offset1).Order(pagination.Sort + " " + pagination.OrderBy).Find(tableName).Error

	if errFind != nil {
		return PaginationService{Error: errFind}, totalPages
	}

	pagination.Data = tableName

	pagination.TotalRows = totalRows

	// calculate total pages
	totalPages = int(math.Ceil(float64(totalRows)/float64(pagination.Limit))) - 1

	if pagination.Page == 0 && pagination.Offset == 0 {
		// set from & to row on first page
		fromRow = 0
		toRow = pagination.Limit
	} else {
		// if pagination.Page <= totalPages {
		// calculate from & to row
		fromRow = offset
		toRow = offset + pagination.Limit
		// }
	}

	if toRow > totalRows {
		// set to row with total rows
		toRow = totalRows
	}

	pagination.FromRow = fromRow
	pagination.ToRow = toRow

	return PaginationService{Result: pagination}, totalPages
}

//PaginationServices function
func PaginationServices(context *gin.Context, pagination *Pagination, fields, where, queryParams, tablename string, app interface{}) Response {
	operationResult, totalPages := PaginationCall(context, pagination, fields, where, tablename, app)

	if operationResult.Error != nil {
		return Response{Message: operationResult.Error.Error(), Details: ""}
	}

	var data = operationResult.Result.(*Pagination)

	// get current url path
	urlPath := context.Request.URL.Path
	baseUrl := os.Getenv("BASE_URL")

	// set first & last page pagination response
	//data.FirstPage = baseUrl + fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, 0, pagination.Sort, pagination.OrderBy) + "&" + queryParams
	//data.LastPage = baseUrl + fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, totalPages, pagination.Sort, pagination.OrderBy) + "&" + queryParams
	// pagination.Sort()

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = baseUrl + fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, data.Page-1, pagination.Sort, pagination.OrderBy) + "&" + queryParams
		// pagination.Sort()
	} else if data.Offset > 0 && data.Limit < data.Offset {
		// set previous page pagination response
		data.PreviousPage = baseUrl + fmt.Sprintf("%s?limit=%d&offset=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, data.Offset-data.Limit, pagination.Sort, pagination.OrderBy) + "&" + queryParams
		// pagination.Sort()
	}

	if data.Page < totalPages && data.Offset == 0 {
		// set next page pagination response
		data.NextPage = baseUrl + fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, data.Page+1, pagination.Sort, pagination.OrderBy) + "&" + queryParams
		// paginaton.Sort()
	} else if data.Offset < totalPages && data.Page == 0 {
		// set next page pagination response
		data.NextPage = baseUrl + fmt.Sprintf("%s?limit=%d&offset=%d&sort=%s&order_by=%s", urlPath, pagination.Limit, data.Offset+pagination.Limit, pagination.Sort, pagination.OrderBy) + "&" + queryParams
		// paginaton.Sort()
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return Response{Message: "Success", Details: data}
}
