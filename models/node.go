package models

type NodeResp struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	Address   string `json:"address"`
	Available bool   `json:"available"`
}

type AddNodeReq struct {
	Name    string `json:"name" binding:"required"`
	Comment string `json:"comment"`
	Address string `json:"address" binding:"required"`
}

type UpdateNodeReq struct {
	Name    string `json:"name" binding:"required"`
	Comment string `json:"comment"`
}
