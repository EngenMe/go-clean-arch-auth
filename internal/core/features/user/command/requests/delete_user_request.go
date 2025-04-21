package requests

type DeleteUserRequest struct {
	ID uint `json:"id" validate:"required"`
}
