// pkg/infrastructure/postgres/user.go
package postgres

import (
	"context"
	"errors"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ user.UserRepository = new(UserRepoPostgres)

type UserRepoPostgres struct {
	Pool *pgxpool.Pool
}

func NewUserRepoPostgres(pool *pgxpool.Pool) *UserRepoPostgres {
	return &UserRepoPostgres{Pool: pool}
}

func (r *UserRepoPostgres) GetUser(ctx context.Context, User *user.User) (*user.User, error) {
	// yakovlev: по идее тут перед логином я должен проверять, есть ли сессия (?) в таблице sessions в пг
	// также я должен проверять не протухла ли она (поле expires_at)

	// ну кстати не в репозитории это делать. Наверно это надо уносить куда-то логику выше
	var u user.User
	err := r.Pool.QueryRow(ctx, "SELECT id, login, password FROM users WHERE login = $1", User.Username).
		Scan(&u.UserID, &u.Username)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepoPostgres) Create(ctx context.Context, User *user.User, Password string) (*user.User, error) {
	err := r.Pool.QueryRow(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		User.Username, Password,
	).Scan(&User.UserID)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (r *UserRepoPostgres) GetUserPassword(ctx context.Context, user *user.User) (string, error) {
	var Password string
	err := r.Pool.QueryRow(ctx, "select password from users where username = $1", user.Username).Scan(Password)
	if err == pgx.ErrNoRows {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", err
	}
	return Password, nil
}
