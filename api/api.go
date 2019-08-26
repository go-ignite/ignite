package api

type PagingRequest struct {
	Keyword   string `form:"keyword" json:"keyword,omitempty"`
	PageIndex int    `form:"page_index" binding:"gt=0" json:"page_index"`
	PageSize  int    `form:"page_size" binding:"gt=0" json:"page_size"`
}

type PagingResponse struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
	PagingRequest
}
