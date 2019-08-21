package api

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminResetAccountPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}
