package service

type Repository interface {
	Create(user *Model) error
	Delete(user *Model) error
	Get(user *Model) error
	GetAll() ([]string, error)
}
