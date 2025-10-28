package postgres

import (
	"database/sql"
	"errors"
	"gophermart/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

var errEmptyDSN = errors.New("DSN is empty")
var errEmptyPassword = errors.New("password is empty")

type Repo struct {
	db *sql.DB
}

func New(dsn string) (*Repo, error) {
	if dsn == "" {
		return nil, errEmptyDSN
	}

	// Запускаем миграции
	if err := runMigration(dsn); err != nil {
		return nil, err
	}

	// Подключаемся к БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	repo := &Repo{
		db: db,
	}

	return repo, nil
}

func runMigration(dsn string) error {
	m, err := migrate.New(
		"file://./migrations",
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r *Repo) AddUser(login, password string) error {
	pwdHash, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		"INSERT INTO users (login, password_hash) VALUES ($1, $2)",
		login,
		pwdHash,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrLoginAlreadyTaken
		}
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	if password == "" {
		return "", errEmptyPassword
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (r *Repo) AuthenticateUser(login, password string) error {
	var hashedPwdFromDB string
	err := r.db.QueryRow(
		"SELECT password_hash FROM users WHERE login = $1",
		login,
	).Scan(&hashedPwdFromDB)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrUserLoginNotFound
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPwdFromDB), []byte(password))
	if err != nil {
		return repository.ErrWrongPassword
	}

	return nil
}

func (r *Repo) GetBalance(login string) (current, withdrawn float64, err error) {
	err = r.db.QueryRow(
		"SELECT balance, withdrawn FROM users WHERE login = $1",
		login,
	).Scan(&current, &withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, repository.ErrUserLoginNotFound
		}
		return 0, 0, err
	}

	return current, withdrawn, err
}
