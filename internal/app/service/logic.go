package service

import (
	"fmt"
	e "github.com/VSKrivoshein/test/internal/app/custom_err"
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
)

type service struct {
	repo      Repository
	broker    Broker
	validator *validator.Validate
}

func New(repo Repository, broker Broker) User {
	return &service{
		repo:      repo,
		broker:    broker,
		validator: validator.New(),
	}
}

func (s *service) Create(mail string, password string) error {

	user := Model{
		Email:        mail,
		Password:     password,
		PasswordHash: "",
	}

	if err := s.validator.Struct(user); err != nil {
		return e.New(err, err, codes.InvalidArgument)
	}

	passwordHash, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	user.PasswordHash = passwordHash

	if err := s.repo.Create(&user); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	s.broker.Log(user.Email)

	return nil
}

func (s *service) Delete(mail string, password string) error {
	user := Model{
		Email:    mail,
		Password: password,
	}

	if err := s.validator.Struct(user); err != nil {
		return e.New(err, err, codes.InvalidArgument)
	}

	if err := s.repo.Delete(&user); err != nil {
		return fmt.Errorf(e.GetInfo(), err)
	}

	return nil
}

func (s *service) GetAll() ([]string, error) {
	usersList, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf(e.GetInfo(), err)
	}
	return usersList, nil
}
