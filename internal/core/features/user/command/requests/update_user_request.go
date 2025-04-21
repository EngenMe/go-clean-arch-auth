package requests

type UpdateUserRequest struct {
	ID       uint   `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8"`
}
