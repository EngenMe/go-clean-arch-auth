package handler

import (
	"errors"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/requests"
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	"github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/types"
	mapper "github.com/EngenMe/go-clean-arch-auth/internal/core/mapper/user/queryMapper"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase"
)

type FindUserByIdHandler struct {
	useCase useCase.UserUseCase
	mapper  mapper.Mapper
}

func (h *FindUserByIdHandler) Handle(request *requests.FindUserByIdRequest) (
	types.ResponseTypes,
	error,
) {

	user, err := h.useCase.GetById(request.ID)
	if err != nil {
		return &userResponse.UserResponse{}, errors.New("invalid user id")
	}

	userRes, err := h.mapper.MapUserToResponse(user)
	if err != nil {
		return &userResponse.UserResponse{}, errors.New("mapper map user to response failed")
	}

	return userRes, nil
}
