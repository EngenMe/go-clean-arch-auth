package handler

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/requests"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/validators"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"github.com/EngenMe/go-clean-arch-auth/pkg/utils"
)

type UpdateUserHandler struct {
	useCase useCase.UserUseCase
}

func (h *UpdateUserHandler) Handle(request requests.UpdateUserRequest) error {
	if err := validators.ValidateRequest(request); err != nil {
		return err
	}

	user := &entity.User{
		ID:    request.ID,
		Email: request.Email,
		Name:  request.Name,
	}

	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return h.useCase.UpdateUser(user)
}
