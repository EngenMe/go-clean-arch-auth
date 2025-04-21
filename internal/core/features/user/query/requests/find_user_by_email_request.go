package requests

type FindUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r FindUserByIdRequest) isReqeust() {}
