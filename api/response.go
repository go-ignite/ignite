package api

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrResponse(code int, message string) *ErrResponse {
	return &ErrResponse{
		Code:    code,
		Message: message,
	}
}

type PageData struct {
	Total     int64       `json:"total"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	Data      interface{} `json:"data"`
}

//ServiceResult struct for new created SS service
type ServiceResult struct {
	ID           string `json:"serviceId"`
	Host         string `json:"host"`
	Port         int    `json:"servicePort"`
	Password     string `json:"servicePwd"`
	PackageLimit int    `json:"packageLimit"`
}

type UserInfo struct {
	Id                 int64
	Host               string
	Username           string
	Status             int
	PackageLimit       int
	PackageUsed        string
	PackageLeft        string
	PackageLeftPercent string
	ServicePort        int
	ServicePwd         string
	ServiceMethod      string
	ServiceType        string
	ServiceURL         string
	Expired            string
}
