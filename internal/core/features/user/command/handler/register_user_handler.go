package handler

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/requests"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/validators"
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
	"github.com/EngenMe/go-clean-arch-auth/pkg/utils"
)

type RegisterUserHandler struct {
	useCase useCase.UserUseCase
}

func NewRegisterUserHandler(useCase useCase.UserUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{useCase: useCase}
}

func (h *RegisterUserHandler) Handle(request requests.RegisterUserRequest) (
	*userResponse.UserResponse,
	error,
) {
	if err := validators.ValidateRequest(request); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: hashedPassword,
	}

	err = h.useCase.Register(user)
	return &userResponse.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, err
}
