package service

type Rdb interface {
	GetAllUsers() ([]string, error)
	SetAllUsers(users []string)
	InvalidateCache()
}
