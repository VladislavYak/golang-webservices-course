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
	var u user.User
	err := r.Pool.QueryRow(ctx, "SELECT id, login, password FROM users WHERE login = $1", User.Username).
		Scan(&u.UserID, &u.Username, &u.Password)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepoPostgres) Create(ctx context.Context, user *user.User) (*user.User, error) {
	err := r.Pool.QueryRow(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Username, user.Password,
	).Scan(&user.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
