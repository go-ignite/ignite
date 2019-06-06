package api

type Node struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Comment           string `json:"comment"`
	RequestAddress    string `json:"request_address"`
	ConnectionAddress string `json:"connection_address"`
	PortFrom          int    `json:"port_from"`
	PortTo            int    `json:"port_to"`
}

type AddNodeRequest struct {
	RequestAddress string `json:"request_address" binding:"required"`
	UpdateNodeRequest
}

type UpdateNodeRequest struct {
	Name              string `json:"name" binding:"required"`
	Comment           string `json:"comment"`
	ConnectionAddress string `json:"connection_address" binding:"required"`
	PortFrom          int    `json:"port_from" binding:"required"`
	PortTo            int    `json:"port_to" binding:"required,gtfield=PortFrom"`
}
