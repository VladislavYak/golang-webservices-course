package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ auth.AuthRepository = new(AuthRepoPostgres)

type AuthRepoPostgres struct {
	Pool *pgxpool.Pool
}

func NewAuthRepoPostgres(pool *pgxpool.Pool) *AuthRepoPostgres {
	return &AuthRepoPostgres{Pool: pool}
}

func (a *AuthRepoPostgres) AddJWT(ctx context.Context, Token string, UserID string, IssuedAt time.Time, ExpiresAt time.Time) error {
	const op = "AddJWT"
	fmt.Println("before Add JWT")

	fmt.Println("")
	_, err := a.Pool.Exec(ctx,
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

func (a *AuthRepoPostgres) ValidateJWT(ctx context.Context, Token string, ExpiresAt time.Time) error {
	const op = "ValidateJWT"
	fmt.Println("inside ValidateJWT")
	var expiresAt time.Time
	err := a.Pool.QueryRow(ctx, "select expires_at from sessions where token = $1", Token).Scan(&expiresAt)

	if err == pgx.ErrNoRows {
		return errors.New("token not found")
	}

	if err != nil {
		return errors.Wrap(err, op)
	}

	fmt.Println("я тут")

	if expiresAt.Before(time.Now()) {
		return auth.ExpiredTokenError
	}
	return nil
}
