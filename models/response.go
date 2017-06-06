package models

//Response struct for http response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

//ServiceResponse struct for creating new SS service
type ServiceResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ServicePort  int    `json:"servicePort"`
	ServicePwd   string `json:"servicePwd"`
	PackageLimit int    `json:"packageLimit"`
}
