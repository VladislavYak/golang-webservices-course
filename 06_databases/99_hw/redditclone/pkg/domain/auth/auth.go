package auth

import (
	"context"
	"time"
)

type AuthRepository interface {
	AddJWT(ctx context.Context, Token string, UserID string, IssuedAt time.Time, ExpiresAt time.Time) error
	ValidateJWT(ctx context.Context, Token string, ExpiresAt time.Time) error
}
