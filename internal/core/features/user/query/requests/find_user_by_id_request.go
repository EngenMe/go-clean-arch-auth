package requests

type FindUserByIdRequest struct {
	ID uint `json:"id" validate:"required"`
}

func (r *FindUserByIdRequest) isRequest() {}
