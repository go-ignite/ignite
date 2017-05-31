package models

//Response struct for http response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
