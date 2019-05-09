package api

type InviteCode struct {
	ID        uint   `json:"id"`
	Code      string `json:"invite_code"`
	Limit     int    `json:"limit"`
	ExpiredAt Time   `json:"expired_at"`
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

func NewInviteCodeListResponse(list []*InviteCode, total int, pageIndex int) *InviteCodeListResponse {
	if list == nil {
		list = []*InviteCode{}
	}
	return &InviteCodeListResponse{
		List:      list,
		Total:     total,
		PageIndex: pageIndex,
	}
}

type GenerateCodesRequest struct {
	Amount    uint `form:"amount" binding:"required"`
	Limit     uint `form:"limit" binding:"required"`
	ExpiredAt Time `json:"expired_at"`
}
