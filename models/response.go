package models

//Response struct for http response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//ServiceResponse struct for creating new SS service
type ServiceResult struct {
	ID       string `json:"serviceId"`
	Port     int    `json:"servicePort"`
	Password string `json:"servicePwd"`
}
