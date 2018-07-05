package models

type NodeResp struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	Address   string `json:"address"`
	ConnectIP string `json:"connect_ip"`
	PortFrom  int    `json:"port_from"`
	PortTo    int    `json:"port_to"`
	Available bool   `json:"available"`
}

type AddNodeReq struct {
	Name       string `json:"name" binding:"required"`
	Comment    string `json:"comment"`
	Address    string `json:"address" binding:"required"`
	Connection string `json:"connection" binding:"required"`
	PortFrom   int    `json:"port_from" binding:"required"`
	PortTo     int    `json:"port_to" binding:"required"`
}

type UpdateNodeReq struct {
	Name       string `json:"name" binding:"required"`
	Comment    string `json:"comment"`
	Connection string `json:"connection" binding:"required"`
	PortFrom   int    `json:"port_from" binding:"required"`
	PortTo     int    `json:"port_to" binding:"required"`
}
