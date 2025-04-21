package handler

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/requests"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/validators"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
)

type DeleteUserHandler struct {
	useCase useCase.UserUseCase
}

func (h *DeleteUserHandler) Handle(request requests.DeleteUserRequest) error {
	if err := validators.ValidateRequest(request); err != nil {
		return err
	}

	return h.useCase.DeleteUser(request.ID)
}
