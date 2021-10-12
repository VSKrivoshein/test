package service

type Model struct {
	UserId          string `db:"id"`
	Email           string `validate:"required,email" db:"email"`
	Password        string `validate:"min=6,max=100"`
	PasswordHash    string `db:"password_hash"`
}