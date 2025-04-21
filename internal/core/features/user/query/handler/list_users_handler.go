package handler

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/command/validators"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/requests"
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
)

type ListUsersHandler struct {
	useCase useCase.UserUseCase
}

func (h *ListUsersHandler) Handle(request *requests.ListUsersRequest) (
	[]*userResponse.UserResponse,
	error,
) {
	if err := validators.ValidateRequest(request); err != nil {
		return nil, err
	}

	users, err := h.useCase.ListUsers(request.Page, request.Limit)
	if err != nil {
		return nil, err
	}

	var userResponses []*userResponse.UserResponse
	for _, user := range users {
		userResponses = append(
			userResponses, &userResponse.UserResponse{
				ID:        user.ID,
				Email:     user.Email,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		)
	}

	return userResponses, nil
}
