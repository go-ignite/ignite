package api

import "time"

type User struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	PackageLimit     int       `json:"package_limit"`
	MonthTrafficUsed uint64    `json:"month_traffic_used"`
	LastStatsTime    time.Time `json:"last_stats_time"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserRegisterRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UserResisterResponse struct {
	Token string `json:"token"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserInfo struct {
	Name string `json:"name"`
}
