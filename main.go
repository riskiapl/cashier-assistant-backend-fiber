package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/riskiapl/fiber-app/config"
	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/routes"
)

func main() {
	// Load environment variables based on environment
	config.LoadEnv()

	// Log which environment we're running in
	env := config.GetEnvironment()
	log.Printf("Application running in %s mode\n", env)

	// Membuat database
	database.CreateDatabase()

	// Koneksi ke database
	database.ConnectDB()

	// Migrasi database
	database.Migrate()

	app := fiber.New()

	// Get CORS configuration based on environment
	corsConfig := config.GetCorsConfig()

	log.Printf("CORS configuration: %v\n", corsConfig)

	// Setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig["origins"],
		AllowHeaders:     corsConfig["headers"],
		AllowMethods:     corsConfig["methods"],
		AllowCredentials: corsConfig["credentials"] == "true",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Setup routes
	routes.AuthRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
