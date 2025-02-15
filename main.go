package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/routes"
)

func main() {
	// Membuat database
	database.CreateDatabase()

	// Koneksi ke database
	database.ConnectDB()
	
	// Migrasi database
	database.Migrate()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Setup routes
	routes.AuthRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
