// pkg/infrastructure/postgres/user.go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/auth"
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
	fmt.Println("before insertion Create")
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

func (r *UserRepoPostgres) AddJWT(ctx context.Context, Token string, Claims *auth.JwtCustomClaims) error {
	fmt.Println("before Add JWT")

	fmt.Println("")
	_, err := r.Pool.Exec(ctx,
		"INSERT INTO sessions (user_id, token, issued_at, expires_at) VALUES ($1, $2, $3, $4)",
		Claims.UserID, Token, Claims.IssuedAt.Time, Claims.ExpiresAt.Time,
	)

	fmt.Println("errrrrr", err)

	if err != nil {
		return err
	}

	return nil
}

// INSERT INTO sessions (user_id, token, issued_at, expires_at) VALUES
// (1, 'sample_jwt_token_1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour'),
// (2, 'sample_jwt_token_2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour');
