package service

type User interface {
	Create(mail string, password string) error
	Delete(mail string, password string) error
	GetAll() ([]string, error)
}
