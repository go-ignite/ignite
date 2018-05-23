package models

//Response struct for http response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewErrorResp(message string) *Response {
	return &Response{Message: message}
}

func NewSuccessResp(data interface{}, message ...string) *Response {
	msg := "Success!"
	if len(message) > 0 {
		msg = message[0]
	}
	return &Response{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

//
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

type ServiceConfig struct {
	SSMethods  []string `json:"ssMethods"`
	SSRMethods []string `json:"ssrMethods"`
	Servers    []string `json:"servers"`
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
