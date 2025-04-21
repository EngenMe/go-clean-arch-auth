package respositories

import (
	"errors"
	"github.com/EngenMe/go-clean-arch-auth/internal/Infrastructure/repository/baseRepository"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	*baseRepository.PostgresRepository[entity.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	baseRepo := baseRepository.NewPostgresRepository[entity.User](db)
	return &UserRepository{PostgresRepository: baseRepo}
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}
