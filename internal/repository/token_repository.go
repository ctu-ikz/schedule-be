package repository

import (
	"context"
	"github.com/ctu-ikz/schedule-be/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		pool: pool,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens (
	    id, user_id, hashed_token, expires_at, device_info, ip_address
	) 
	VALUES ($1, $2, $3, $4, $5, $6)
`
	ip := ""
	if token.IPAddress != nil {
		ip = token.IPAddress.String()
	}

	_, err := r.pool.Exec(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.HashedToken,
		token.ExpiresAt,
		token.DeviceInfo,
		ip,
	)

	return err
}

func (r *TokenRepository) GetRefreshInfoByHashedToken(ctx context.Context, hashedToken string) (bool, time.Time, uuid.UUID, error) {
	query := `SELECT expires_at, revoked, user_id FROM refresh_tokens where hashed_token = $1`

	var expiresAt time.Time
	var revoked bool
	var userID uuid.UUID

	err := r.pool.QueryRow(ctx, query, hashedToken).Scan(&expiresAt, &revoked, &userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, time.Time{}, uuid.UUID{}, nil
		}
		return false, time.Time{}, uuid.UUID{}, err
	}

	return revoked, expiresAt, userID, nil
}

func (r *TokenRepository) RevokeTokenByHashedToken(ctx context.Context, hashedToken string) error {
	query := `
	UPDATE refresh_tokens SET revoked = true WHERE hashed_token = $1`
	_, err := r.pool.Exec(ctx, query, hashedToken)
	return err
}
