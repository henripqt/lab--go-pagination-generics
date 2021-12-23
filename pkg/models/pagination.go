package models

// PaginationRequest is a struct that represents a pagination request
type PaginationRequest struct {
	Page     int      `json:"page" query:"page"`
	PerPage  int      `json:"perPage" query:"per_page"`
	OrderBy  []string `json:"orderBy" query:"order_by"`
	OrderDir string   `json:"orderDir" query:"order_dir"`
}

// PagingResponse is a struct that represents the response of a paginated request
type PagingResponse struct {
	Page       int         `json:"page"`
	Items      interface{} `json:"items"`
	PerPage    int         `json:"perPage"`
	PrevPage   int         `json:"prevPage"`
	NextPage   int         `json:"nextPage"`
	TotalPage  int         `json:"totalPage"`
	TotalItems int64       `json:"totalItems"`
}
