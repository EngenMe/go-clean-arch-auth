package baseRepository

import (
	"errors"
	"gorm.io/gorm"
)

type PostgresRepository[T any] struct {
	DB *gorm.DB
}

func NewPostgresRepository[T any](db *gorm.DB) *PostgresRepository[T] {
	return &PostgresRepository[T]{DB: db}
}

func (r *PostgresRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *PostgresRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	result := r.DB.First(&entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("entity not found")
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *PostgresRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *PostgresRepository[T]) Delete(id uint) error {
	return r.DB.Delete(new(T), id).Error
}

func (r *PostgresRepository[T]) List(page, limit int) ([]T, error) {
	var entities []T
	offset := (page - 1) * limit
	result := r.DB.Offset(offset).Limit(limit).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	return entities, nil
}
