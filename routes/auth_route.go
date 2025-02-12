package routes

import (
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes mengatur route untuk /auth
func AuthRoutes(app fiber.Router) {
	auth := app.Group("/auth") // Buat grup route "/auth"

	auth.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Auth route"})
	})

	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Login endpoint"})
	})
}
