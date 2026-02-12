package adapters

import (
	"hole/use_cases"

	"github.com/gofiber/fiber/v2"
)

func Protected(ts use_cases.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get token from cookie
		auth := c.Cookies("auth_token")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication token",
			})
		}

		// // 2. Validate using the injected service
		// claims, err := ts.ValidateAccessToken(auth)
		// if err != nil {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "Invalid or expired token",
		// 	})
		// }

		// // 3. Extract and convert user_id to uuid.UUID
		// // JWT claims often store values as strings or float64
		// userIDStr, ok := claims["user_id"].(string)
		// if !ok {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "Invalid token payload",
		// 	})
		// }

		// userUID, err := uuid.Parse(userIDStr)
		// if err != nil {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "Invalid user identification",
		// 	})
		// }

		// 4. Set as actual UUID type for handlers to use
		// c.Locals("user_id", userUID)

		return c.Next()
	}
}
