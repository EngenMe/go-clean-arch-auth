package handler

import (
	"errors"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/requests"
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"github.com/EngenMe/go-clean-arch-auth/pkg/utils"
)

type LoginHandler struct {
	useCase useCase.UserUseCase
}

func (h *LoginHandler) Handle(request requests.LoginRequest) (
	*userResponse.UserResponse,
	error,
) {
	user, err := h.useCase.FindUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if err := utils.VerifyPassword(
		user.Password,
		request.Password,
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &userResponse.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
