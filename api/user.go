package api

type User struct {
	ID   uint   `json:"id"`
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

func NewUserListResponse(list []*User, total int, pageIndex int) *UserListResponse {
	if list == nil {
		list = []*User{}
	}
	return &UserListResponse{
		List:      list,
		Total:     total,
		PageIndex: pageIndex,
	}
}

type UserLoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

func NewUserLoginResponse(token string) *UserLoginResponse {
	return &UserLoginResponse{
		Token: token,
	}
}

type UserRegisterRequest struct {
	InviteCode      string `form:"invite_code" binding:"required"`
	Username        string `form:"username" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required"`
}

type UserResisterResponse struct {
	Token string `json:"token"`
}

func NewUserRegisterResponse(token string) *UserResisterResponse {
	return &UserResisterResponse{
		Token: token,
	}
}
