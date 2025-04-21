package repository

import (
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
)

type UserRepository interface {
	BaseRepository[entity.User]
	FindByEmail(email string) (*entity.User, error)
}
