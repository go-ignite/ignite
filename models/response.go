package models

//Response struct for http response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
