package common

type RequestPagination struct {
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=9223372036854775807" swaggertype:"integer"`
	Page     int `form:"page" binding:"omitempty,min=-1,max=9223372036854775807" swaggertype:"integer"`
}

type ResponsePagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size,omitempty"`
}
