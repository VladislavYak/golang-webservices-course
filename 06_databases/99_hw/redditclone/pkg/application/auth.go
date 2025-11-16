package application

import (
	"context"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	"github.com/go-faster/errors"
)

type AuthRepo interface {
	ValidateSession(ctx context.Context, Token string, ExpiresAt time.Time)
}

type AuthImpl struct {
	ar auth.AuthRepository
}

func NewAuthImpl(repo auth.AuthRepository) *AuthImpl {
	return &AuthImpl{ar: repo}
}

func (ai *AuthImpl) ValidateSession(ctx context.Context, Token string, ExpiresAt time.Time) error {
	const op = "ValidateSession"
	err := ai.ar.ValidateJWT(ctx, Token, ExpiresAt)

	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
