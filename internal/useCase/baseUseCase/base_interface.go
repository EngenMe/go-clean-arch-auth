package baseUseCase

type DefaultUseCase[T any] interface {
	GetById(id uint) (T, error)
	Add(entity *T) error
	Update(entity *T) error
	Delete(id uint) error
	List(page int, limit int) ([]T, error)
}
