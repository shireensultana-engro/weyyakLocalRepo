package pagination

import (
	_ "github.com/jinzhu/gorm"
)

//Pagination struct
type Pagination struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Page      int    `json:"page"`
	Sort      string `json:"sort"`
	TotalRows int    `json:"total_rows"`
	//FirstPage    string      `json:"first_page"`
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	//LastPage     string      `json:"last_page"`
	FromRow int         `json:"from_row"`
	ToRow   int         `json:"to_row"`
	Data    interface{} `json:"data"`
	OrderBy string      `json:"order_by"`
	Table   string      `json:"table_name"`

	// Searchs      []Search    `json:"searchs"`
}

//Response customized struct
type Response struct {
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
