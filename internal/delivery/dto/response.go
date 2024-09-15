package dto

type ErrResponse struct {
	Reason string `json:"reason"`
}

func NewErrResponse(err error) *ErrResponse {
	return &ErrResponse{
		Reason: err.Error(),
	}
}
