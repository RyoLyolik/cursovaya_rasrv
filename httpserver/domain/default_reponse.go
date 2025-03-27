package domain

type DefaultResponse struct {
	Status string         `json:"status"`
	Error  *errorResponse `json:"error,omitempty"`
	Data   any            `json:"data,omitempty"`
}

var (
	ErrResp = DefaultResponse{
		Status: "Error",
		Error:  &errorResponse{},
	}
	OkResp = DefaultResponse{
		Status: "Ok",
	}
)
