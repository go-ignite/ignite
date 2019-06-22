package api

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserListRequest struct {
	Keyword   string `form:"keyword"`
	PageIndex int    `form:"page_index"`
	PageSize  int    `form:"page_size"`
}

type UserListResponse struct {
	List      []*User `json:"list"`
	Total     int     `json:"total"`
	PageIndex int     `json:"page_index"`
}

type UserLoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserRegisterRequest struct {
	InviteCode string `json:"invite_code" form:"invite_code" binding:"required"`
	Username   string `json:"username" form:"username" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required"`
}

type UserResisterResponse struct {
	Token string `json:"token"`
}
