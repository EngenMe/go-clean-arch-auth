package baseUseCase

import "github.com/EngenMe/go-clean-arch-auth/internal/Infrastructure/repository"

type BaseUseCase[T any] struct {
	baseRepository repository.BaseRepository[T]
}

func (baseUseCase *BaseUseCase[T]) Add(entity *T) (*T, error) {
	err := baseUseCase.baseRepository.Create(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (baseUseCase *BaseUseCase[T]) GetById(id uint) (*T, error) {
	return baseUseCase.baseRepository.FindByID(id)
}

func (baseUseCase *BaseUseCase[T]) Update(entity *T) (*T, error) {
	err := baseUseCase.baseRepository.Update(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (baseUseCase *BaseUseCase[T]) Delete(id uint) error {
	return baseUseCase.baseRepository.Delete(id)
}

func (baseUseCase *BaseUseCase[T]) List(page, limit int) ([]T, error) {
	return baseUseCase.baseRepository.List(page, limit)
}
