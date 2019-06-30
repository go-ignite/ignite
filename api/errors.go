package api

var (
	// model
	ErrInviteCodeNotExistOrUnavailable = NewErrResponse(1000, "invite code does not exist or is unavailable")
	ErrInviteCodeExpired               = NewErrResponse(1001, "invite code is expired")
	ErrUserNameExists                  = NewErrResponse(1002, "user name already exists")
	ErrUserDeleted                     = NewErrResponse(1003, "user has been deleted")

	// state
	ErrNodeNotExist                   = NewErrResponse(1100, "node does not exist")
	ErrNodeUnavailable                = NewErrResponse(1101, "node is unavailable")
	ErrNodeHasServicesExceedPortRange = NewErrResponse(1102, "node has services that exceed port range")
	ErrNodeNameExists                 = NewErrResponse(1103, "node name already exists")
	ErrNodeRequestAddressExists       = NewErrResponse(1104, "node request address already exists")
	ErrServiceExists                  = NewErrResponse(1105, "user already has a service on this node")
)

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *ErrResponse) Error() string {
	return r.Message
}

func NewErrResponse(code int, message string) *ErrResponse {
	return &ErrResponse{
		Code:    code,
		Message: message,
	}
}
