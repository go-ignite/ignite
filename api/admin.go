package api

type AdminLoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}
