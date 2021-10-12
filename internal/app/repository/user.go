package repository

import (
	"fmt"
	e "github.com/VSKrivoshein/test/internal/app/custom_err"
	"github.com/VSKrivoshein/test/internal/app/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"net/http"
	"os"
)

type db struct {
	cache  service.Rdb
	client *sqlx.DB
}

func New(cache service.Rdb) service.Repository {
	var dbUrl = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_SSL_MODE"),
	)

	logrus.Infof("db url: %v", dbUrl)

	client, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		logrus.Fatalf("Postgres connection failed %s", err.Error())
	}

	logrus.Info("DB is connected")

	return &db{
		cache:  cache,
		client: client,
	}
}

func (d *db) Create(user *service.Model) error {
	row, err := d.client.NamedQuery(`
		INSERT INTO users
		VALUES (gen_random_uuid(), :email, :password_hash)
		RETURNING id, email, password_hash;
	`, user)

	if err != nil {
		return e.New(err, e.ErrCreatingUser, codes.InvalidArgument)
	}

	if row.Next() {
		if err := row.StructScan(user); err != nil {
			return e.New(err, e.ErrUnexpected, codes.Internal)
		}
	}

	go d.cache.InvalidateCache()

	return nil
}

func (d *db) Get(user *service.Model) error {
	err := d.client.Get(user, `
		SELECT  id, password_hash, password_hash
		FROM users
		WHERE email = $1;
		`, user.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return e.New(err, e.ErrGetUser, codes.InvalidArgument)
		}
		return e.New(err, e.ErrUnexpected, codes.Internal)
	}

	return nil
}

func (d *db) Delete(user *service.Model) error {
	res, err := d.client.Exec(`
		DELETE FROM users
		WHERE email=$1;
	`, user.Email)

	if err != nil {
		return e.New(err, e.ErrDeletingUser, codes.Internal)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return e.New(err, e.ErrDeletingUser, http.StatusInternalServerError)
	}

	if count == 0 {
		return e.New(err, e.ErrDeletingUserNotFound, http.StatusInternalServerError)
	}

	go d.cache.InvalidateCache()

	return nil
}

func (d *db) GetAll() ([]string, error) {
	users, err := d.cache.GetAllUsers()
	if err == nil && len(users) > 0 {
		return users, nil
	}

	var data []struct {
		Email string `db:"email"`
	}

	if err := d.client.Select(&data, `SELECT email FROM users;`); err != nil {
		return nil, e.New(err, e.ErrUnexpected, codes.Internal)
	}

	emails := make([]string, 0, len(data))
	for _, str := range data {
		emails = append(emails, str.Email)
	}

	go d.cache.SetAllUsers(emails)

	return emails, nil
}
