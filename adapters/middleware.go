package adapters

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtService JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		userID, err := jwtService.ValidateAccessToken(tokenStr)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 🔑 make user_id available downstream
		c.Locals("user_id", userID)

		return c.Next()
	}
}
