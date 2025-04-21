package useCase

import (
	"errors"
	"github.com/EngenMe/go-clean-arch-auth/internal/Infrastructure/repository"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/EngenMe/go-clean-arch-auth/internal/useCase/baseUseCase"
	"github.com/EngenMe/go-clean-arch-auth/pkg/utils"
)

type UserUseCase interface {
	Register(user *entity.User) error
	Login(email, password string) (*entity.User, error)
	baseUseCase.DefaultUseCase[entity.User]
}

type userUseCase struct {
	userRepo repository.UserRepository
	baseUseCase.DefaultUseCase[entity.User]
}

func NewUserUseCase(
	userRepo repository.UserRepository,
) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Register(user *entity.User) error {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}
	return u.userRepo.Create(user)
}

func (u *userUseCase) Login(email, password string) (*entity.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if utils.VerifyPassword(user.Password, password) != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
