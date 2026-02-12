package adapters

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

var jwtKey = []byte("super-secret-key")

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (j *JWTService) GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"typ":     "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (j *JWTService) GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func (j *JWTService) ValidateAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

// func ParseToken(tokenStr string) (jwt.MapClaims, error) {
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		return nil, err
// 	}

// 	return token.Claims.(jwt.MapClaims), nil
// }

// func Protected(c *fiber.Ctx) error {
// 	auth := c.Cookies("auth_token")

// 	if auth == "" {
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}

// 	claims, err := h.uc.ValidateAccessToken(auth)

// 	if err != nil {
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}
// 	c.Locals("user_id", claims["user_id"])
// 	return c.Next()
// }
