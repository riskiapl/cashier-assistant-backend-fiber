package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/controllers"
)

// AuthRoutes mengatur route untuk /auth
func AuthRoutes(app fiber.Router) {
	auth := app.Group("/auth") // Buat grup route "/auth"

	auth.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Auth route"})
	})

	authController := controllers.NewAuthController()

	auth.Post("/login", authController.Login)	
}
