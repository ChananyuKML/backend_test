package use_cases

import (
	"errors"
	"hole/entities"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
}

type RefreshTokenRepository interface {
	Create(token *entities.RefreshToken) error
	FindByToken(token string) (*entities.RefreshToken, error)
	Revoke(token string) error
}

type TokenService interface {
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateAccessToken(token string) (jwt.MapClaims, error)
}

type AuthUseCase struct {
	repo        UserRepository
	refreshRepo RefreshTokenRepository
	token       TokenService
}

func NewAuthUseCase(r UserRepository, re RefreshTokenRepository, t TokenService) *AuthUseCase {
	return &AuthUseCase{repo: r, refreshRepo: re, token: t}
}

func (u *AuthUseCase) Register(email, password string) error {
	return u.repo.Create(&entities.User{
		Email:    email,
		Password: password,
	})
}

func (uc *AuthUseCase) Login(email, password string) (string, string, error) {
	user, err := uc.repo.FindByEmail(email)

	if err != nil {
		return "", "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, err := uc.token.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// 4️⃣ Generate refresh token
	refresh, err := uc.token.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	if err := uc.refreshRepo.Create(&entities.RefreshToken{
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (uc *AuthUseCase) Refresh(refreshToken string) (string, string, error) {
	rt, err := uc.refreshRepo.FindByToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if rt.Revoked {
		return "", "", errors.New("refresh token revoked")
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = uc.refreshRepo.Revoke(refreshToken)
		return "", "", errors.New("refresh token expired")
	}

	if err := uc.refreshRepo.Revoke(refreshToken); err != nil {
		return "", "", err
	}

	newAccess, err := uc.token.GenerateAccessToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := uc.token.GenerateRefreshToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	if err := uc.refreshRepo.Create(&entities.RefreshToken{
		UserID:    rt.UserID,
		Token:     newRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
