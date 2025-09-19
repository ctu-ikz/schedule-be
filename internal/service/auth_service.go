package service

import (
	"context"
	"errors"
	"github.com/ctu-ikz/schedule-be/internal/domain"
	"github.com/ctu-ikz/schedule-be/internal/util"
	"github.com/google/uuid"
	"net"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetRefreshInfoByHashedToken(ctx context.Context, hashedToken string) (bool, time.Time, uuid.UUID, error)
	RevokeTokenByHashedToken(ctx context.Context, hashedToken string) error
}

type AuthService struct {
	repoUser  UserRepository
	repoToken TokenRepository
}

func NewAuthService(repoUser UserRepository, repoToken TokenRepository) *AuthService {
	return &AuthService{repoUser: repoUser, repoToken: repoToken}
}

func (s *AuthService) Register(ctx context.Context, username, plainPassword string) (*domain.User, error) {
	existing, _ := s.repoUser.FindByUsername(ctx, username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()

	user := &domain.User{
		ID:       id,
		Username: username,
		Password: string(hashed),
	}

	if err := s.repoUser.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, plainPassword string, deviceInfo string, IP net.IP) (string, string, error) {
	user, err := s.repoUser.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	refreshRaw, err := util.GenerateRandomString(64)
	if err != nil {
		return "", "", err
	}

	hashedToken := util.HashRefreshToken(refreshRaw)

	accessToken, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refresh := &domain.RefreshToken{
		ID:          uuid.New(),
		UserID:      user.ID,
		HashedToken: hashedToken,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
		Revoked:     false,
		CreatedAt:   time.Now(),
		LastUsedAt:  time.Time{},
		DeviceInfo:  deviceInfo,
		IPAddress:   IP,
	}

	if err := s.repoToken.Create(ctx, refresh); err != nil {
		return "", "", err
	}

	return refreshRaw, accessToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string, deviceInfo string, IP net.IP) (string, string, error) {
	hashedToken := util.HashRefreshToken(refreshToken)

	revoked, expiresAt, userID, err := s.repoToken.GetRefreshInfoByHashedToken(
		ctx,
		hashedToken)
	if err != nil {
		return "", "", err
	}

	if revoked {
		return "", "", errors.New("token is revoked")
	}

	if expiresAt.Before(time.Now()) {
		return "", "", errors.New("token is expired")
	}

	accessToken, err := util.GenerateAccessToken(userID)

	if err != nil {
		return "", "", err
	}

	refreshRaw, err := util.GenerateRandomString(64)
	if err != nil {
		return "", "", err
	}
	hashedRefreshToken := util.HashRefreshToken(refreshRaw)
	refresh := &domain.RefreshToken{
		ID:          uuid.New(),
		UserID:      userID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
		Revoked:     false,
		CreatedAt:   time.Now(),
		LastUsedAt:  time.Time{},
		DeviceInfo:  deviceInfo,
		IPAddress:   IP,
	}

	if err := s.repoToken.Create(ctx, refresh); err != nil {
		return "", "", errors.New("refresh token creation failed")
	}

	err = s.repoToken.RevokeTokenByHashedToken(
		ctx,
		hashedToken,
	)

	if err != nil {
		return "", "", errors.New("failed to revoke refresh token")
	}

	return refreshRaw, accessToken, nil

}
