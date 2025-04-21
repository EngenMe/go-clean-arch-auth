package repository

type BaseRepository[T any] interface {
	Create(entity *T) error
	FindByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	List(page, limit int) ([]T, error)
}
