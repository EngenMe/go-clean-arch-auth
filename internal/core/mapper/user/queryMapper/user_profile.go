package queryMapper

import (
	userResponse "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/query/response"
	userType "github.com/EngenMe/go-clean-arch-auth/internal/core/features/user/types"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/devfeel/mapper"
)

type Mapper interface {
	MapUserToResponse(entity entity.User) (userType.ResponseTypes, error)

	MapRequestToUserEntity(req userType.RequestTypes) (entity.Entity, error)
}

type userProfile struct{}

func (userProfile *userProfile) MapUserToResponse(entity entity.Entity) (
	userType.ResponseTypes,
	error,
) {
	var response userResponse.UserResponse

	if err := mapper.AutoMapper(entity, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (userProfile *userProfile) MapRequestToUserEntity(
	req userType.RequestTypes,
) (
	entity.Entity,
	error,
) {
	var myEntity entity.Entity

	if err := mapper.AutoMapper(req, myEntity); err != nil {
		return nil, err
	}

	return myEntity, nil
}
