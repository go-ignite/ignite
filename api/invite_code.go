package api

import "time"

type InviteCode struct {
	ID        int64     `json:"id"`
	Code      string    `json:"invite_code"`
	Limit     int       `json:"limit"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type GenerateCodesRequest struct {
	Amount    uint      `json:"amount" form:"amount" binding:"required"`
	Limit     int       `json:"limit" form:"limit" binding:"required"`
	ExpiredAt time.Time `json:"expired_at" form:"expired_at" binding:"required"`
}
