// pkg/infrastructure/postgres/user.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"

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
	const op = "GetUser"
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
		return nil, errors.Wrap(err, op)
	}
	return &u, nil
}

func (r *UserRepoPostgres) Create(ctx context.Context, User *user.User, Password string) (*user.User, error) {
	const op = "Create"

	fmt.Println("before insertion Create")

	err := r.Pool.QueryRow(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		User.Username, Password,
	).Scan(&User.UserID)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return User, nil
}

func (r *UserRepoPostgres) GetUserPassword(ctx context.Context, user *user.User) (string, error) {
	const op = "GetUserPassword"

	var Password string
	err := r.Pool.QueryRow(ctx, "select password from users where username = $1", user.Username).Scan(Password)
	if err == pgx.ErrNoRows {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", errors.Wrap(err, op)
	}
	return Password, nil
}

func (r *UserRepoPostgres) AddJWT(ctx context.Context, Token string, UserID string, IssuedAt time.Time, ExpiresAt time.Time) error {
	const op = "AddJWT"
	fmt.Println("before Add JWT")

	fmt.Println("")
	_, err := r.Pool.Exec(ctx,
		"INSERT INTO sessions (user_id, token, issued_at, expires_at) VALUES ($1, $2, $3, $4)",
		UserID, Token, IssuedAt, ExpiresAt,
	)

	fmt.Println("errrrrr", err)

	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// INSERT INTO sessions (user_id, token, issued_at, expires_at) VALUES
// (1, 'sample_jwt_token_1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour'),
// (2, 'sample_jwt_token_2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 hour');

func (r *UserRepoPostgres) ValidateJWT(ctx context.Context, Token string, ExpiresAt time.Time) error {
	const op = "ValidateJWT"
	fmt.Println("inside ValidateJWT")
	var expiresAt time.Time
	err := r.Pool.QueryRow(ctx, "select expires_at from sessions where token = $1", Token).Scan(&expiresAt)

	if err == pgx.ErrNoRows {
		return errors.New("token not found")
	}

	if err != nil {
		return errors.Wrap(err, op)
	}

	fmt.Println("я тут")

	if expiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}
	return nil
}
