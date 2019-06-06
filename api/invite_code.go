package api

type InviteCode struct {
	ID        int64  `json:"id"`
	Code      string `json:"invite_code"`
	Limit     int    `json:"limit"`
	ExpiredAt int64  `json:"expired_at"`
}

type InviteCodeListRequest struct {
	PageIndex int `form:"page_index"`
	PageSize  int `form:"page_size"`
}

type InviteCodeListResponse struct {
	List      []*InviteCode `json:"list"`
	Total     int           `json:"total"`
	PageIndex int           `json:"page_index"`
}

type GenerateCodesRequest struct {
	Amount    uint  `form:"amount" binding:"required"`
	Limit     int   `form:"limit" binding:"required"`
	ExpiredAt int64 `form:"expired_at"`
}
